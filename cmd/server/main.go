package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"awesomeProject/internal/adapter"
	"awesomeProject/internal/config"
	"awesomeProject/internal/handler"
	"awesomeProject/internal/middleware"
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// Generate master key if not set
	if cfg.Security.MasterKey == "" {
		key := make([]byte, 32)
		rand.Read(key)
		cfg.Security.MasterKey = hex.EncodeToString(key)
		logger.Warn("No master encryption key set, generated a random one. Set MASTER_ENCRYPTION_KEY for persistence.")
	}

	// Generate JWT secret if not set
	if cfg.Security.JWTSecret == "" {
		key := make([]byte, 32)
		rand.Read(key)
		cfg.Security.JWTSecret = hex.EncodeToString(key)
		logger.Warn("No JWT secret set, generated a random one. Set JWT_SECRET for persistence.")
	}

	masterKey, err := cfg.GetMasterKeyBytes()
	if err != nil {
		log.Fatalf("Invalid master key: %v", err)
	}

	// Init database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	// Auto migrate
	if err := db.AutoMigrate(&model.User{}, &model.Provider{}, &model.ProviderModel{}, &model.UserAPIKey{}, &model.UsageLog{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Init Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// Init repositories
	userRepo := repository.NewUserRepo(db)
	providerRepo := repository.NewProviderRepo(db)
	modelRepo := repository.NewModelRepo(db)
	apiKeyRepo := repository.NewAPIKeyRepo(db)
	usageRepo := repository.NewUsageRepo(db)

	// Init services
	userSvc := service.NewUserService(userRepo, cfg.Security.JWTSecret, cfg.Security.JWTExpire)
	providerSvc := service.NewProviderService(providerRepo, masterKey)
	modelSvc := service.NewModelService(modelRepo, providerRepo)
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo, rdb)
	adapterRegistry := adapter.NewRegistry()

	lb := service.NewLoadBalancer(rdb,
		cfg.LoadBalancer.CircuitBreaker.FailureThreshold,
		cfg.LoadBalancer.CircuitBreaker.RecoveryTimeout,
		cfg.LoadBalancer.CircuitBreaker.HalfOpenMax,
	)

	httpClient := &http.Client{
		Timeout: cfg.Proxy.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:    cfg.Proxy.MaxIdleConns,
			IdleConnTimeout: cfg.Proxy.IdleConnTimeout,
		},
	}

	chatSvc := service.NewChatService(modelRepo, providerRepo, usageRepo, apiKeyRepo, providerSvc, lb, adapterRegistry, httpClient, rdb)

	// Init handlers
	adminProviderH := handler.NewAdminProviderHandler(providerSvc)
	adminModelH := handler.NewAdminModelHandler(modelSvc)
	adminChannelH := handler.NewAdminChannelHandler(lb, modelRepo, usageRepo)
	adminUserH := handler.NewAdminUserHandler(userSvc)
	userKeyH := handler.NewUserAPIKeyHandler(apiKeySvc)
	userModelH := handler.NewUserModelHandler(modelSvc, lb, usageRepo)
	userUsageH := handler.NewUserUsageHandler(usageRepo)
	chatH := handler.NewChatCompletionHandler(chatSvc)
	userAuthH := handler.NewUserAuthHandler(userSvc)

	// Setup router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.CORSMiddleware())

	// Serve Web UI - Admin backend
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, webUIHTML)
	})

	// Serve User Portal
	r.GET("/portal", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, userPortalHTML)
	})

	// User auth routes (public)
	r.POST("/user/register", userAuthH.Register)
	r.POST("/user/login", userAuthH.Login)

	// Admin auth & token generation endpoint
	r.POST("/admin/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": gin.H{"message": err.Error()}})
			return
		}
		// Simple admin auth - you can enhance this
		if req.Username == "admin" && req.Password == cfg.Security.JWTSecret {
			token, err := middleware.GenerateAdminToken("admin", cfg.Security.JWTSecret, cfg.Security.JWTExpire)
			if err != nil {
				c.JSON(500, gin.H{"error": gin.H{"message": "failed to generate token"}})
				return
			}
			c.JSON(200, gin.H{"code": 0, "data": gin.H{"token": token, "expires_in": cfg.Security.JWTExpire.Seconds()}})
		} else {
			c.JSON(401, gin.H{"error": gin.H{"message": "invalid credentials"}})
		}
	})

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware(cfg.Security.JWTSecret))
	{
		admin.POST("/providers", adminProviderH.Create)
		admin.GET("/providers", adminProviderH.List)
		admin.PUT("/providers/:id", adminProviderH.Update)
		admin.DELETE("/providers/:id", adminProviderH.Delete)

		admin.POST("/providers/:id/models", adminModelH.Create)
		admin.GET("/providers/:id/models", adminModelH.ListByProvider)
		admin.PUT("/models/:id", adminModelH.Update)
		admin.DELETE("/models/:id", adminModelH.Delete)

		admin.GET("/channels/health", adminChannelH.GetHealth)
		admin.POST("/channels/:id/reset", adminChannelH.ResetHealth)

		// User management
		admin.GET("/users", adminUserH.List)
		admin.POST("/users", adminUserH.Create)
		admin.PUT("/users/:id", adminUserH.Update)
		admin.DELETE("/users/:id", adminUserH.Delete)

		// Admin can also manage API keys
		admin.POST("/api/keys", userKeyH.Create)
		admin.GET("/api/keys", userKeyH.List)
		admin.DELETE("/api/keys/:id", userKeyH.Delete)
	}

	// User JWT authenticated routes (for portal)
	userAuth := middleware.UserAuthMiddleware(cfg.Security.JWTSecret)
	userAPI := r.Group("/user/api")
	userAPI.Use(userAuth)
	{
		userAPI.GET("/profile", userAuthH.GetProfile)
		userAPI.POST("/keys", userKeyH.Create)
		userAPI.GET("/keys", userKeyH.List)
		userAPI.DELETE("/keys/:id", userKeyH.Delete)
		userAPI.GET("/usage", userUsageH.GetUsage)
		userAPI.GET("/usage/details", userUsageH.GetUsageDetails)
		userAPI.GET("/models", userModelH.ListModelsPortal)
	}

	// API key authenticated routes
	apiKeyAuth := middleware.APIKeyAuthMiddleware(apiKeyRepo, rdb)
	rateLimiter := middleware.RateLimitMiddleware(rdb)

	// API Key management (also accessible via admin)
	api := r.Group("/api")
	api.Use(apiKeyAuth, rateLimiter)
	{
		api.POST("/keys", userKeyH.Create)
		api.GET("/keys", userKeyH.List)
		api.DELETE("/keys/:id", userKeyH.Delete)
		api.GET("/usage", userUsageH.GetUsage)
		api.GET("/usage/details", userUsageH.GetUsageDetails)
	}

	// OpenAI-compatible v1 routes
	v1 := r.Group("/v1")
	v1.Use(apiKeyAuth, rateLimiter)
	{
		v1.GET("/models", userModelH.ListModels)
		v1.GET("/models/:model", userModelH.GetModelDetail)
		v1.GET("/models/:model/channels", userModelH.GetModelChannels)
		v1.POST("/chat/completions", chatH.Handle)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Starting AI Gateway", zap.String("addr", addr))
	fmt.Printf("\n🚀 AI Gateway is running at http://localhost%s\n", addr)
	fmt.Printf("📊 Admin Dashboard: http://localhost%s/\n", addr)
	fmt.Printf("👤 User Portal:     http://localhost%s/portal\n", addr)
	fmt.Printf("🔑 Admin login: POST /admin/login with {\"username\": \"admin\", \"password\": \"%s\"}\n\n", cfg.Security.JWTSecret)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
