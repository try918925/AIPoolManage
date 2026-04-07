package middleware

import (
	"awesomeProject/internal/pkg/hash"
	"awesomeProject/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"awesomeProject/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type AdminClaims struct {
	AdminID string `json:"admin_id"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func AdminAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": gin.H{"message": "missing admin token", "type": "authentication_error", "code": "invalid_api_key"},
			})
			return
		}

		claims := &AdminClaims{}
		parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(401, gin.H{
				"error": gin.H{"message": "invalid admin token", "type": "authentication_error", "code": "invalid_api_key"},
			})
			return
		}

		if claims.Role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{
				"error": gin.H{"message": "admin access required", "type": "permission_error", "code": "forbidden"},
			})
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Next()
	}
}

func APIKeyAuthMiddleware(repo *repository.APIKeyRepo, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := extractBearerToken(c)
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": gin.H{"message": "missing api key", "type": "authentication_error", "code": "invalid_api_key"},
			})
			return
		}

		keyHash := hash.SHA256Hex(apiKey)
		ctx := context.Background()

		// Check Redis cache
		var userKey *model.UserAPIKey
		cached, err := rdb.Get(ctx, "apikey:"+keyHash).Result()
		if err == nil {
			var k model.UserAPIKey
			if err := json.Unmarshal([]byte(cached), &k); err == nil {
				userKey = &k
			}
		}

		if userKey == nil {
			// Cache miss
			k, err := repo.FindByHash(keyHash)
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{
					"error": gin.H{"message": "invalid api key", "type": "authentication_error", "code": "invalid_api_key"},
				})
				return
			}

			if !k.Enabled {
				c.AbortWithStatusJSON(401, gin.H{
					"error": gin.H{"message": "api key disabled", "type": "authentication_error", "code": "invalid_api_key"},
				})
				return
			}

			userKey = k
			data, _ := json.Marshal(k)
			rdb.Set(ctx, "apikey:"+keyHash, string(data), 5*time.Minute)
		}

		// Check expiration
		if userKey.ExpiresAt != nil && time.Now().After(*userKey.ExpiresAt) {
			c.AbortWithStatusJSON(401, gin.H{
				"error": gin.H{"message": "api key expired", "type": "authentication_error", "code": "invalid_api_key"},
			})
			return
		}

		// Check quota
		if userKey.QuotaLimit > 0 && userKey.QuotaUsed >= userKey.QuotaLimit {
			c.AbortWithStatusJSON(429, gin.H{
				"error": gin.H{"message": "quota exceeded", "type": "rate_limit_error", "code": "quota_exceeded"},
			})
			return
		}

		c.Set("user_key", userKey)
		c.Set("user_id", userKey.UserID)
		c.Next()
	}
}

// UserClaims for user JWT tokens
type UserClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// UserAuthMiddleware validates user JWT tokens
func UserAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": gin.H{"message": "missing user token", "type": "authentication_error", "code": "invalid_token"},
			})
			return
		}

		claims := &UserClaims{}
		parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(401, gin.H{
				"error": gin.H{"message": "invalid or expired token", "type": "authentication_error", "code": "invalid_token"},
			})
			return
		}

		c.Set("user_id", fmt.Sprintf("%d", claims.UserID))
		c.Set("user_id_int", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

// GenerateAdminToken creates a JWT token for admin authentication
func GenerateAdminToken(adminID, jwtSecret string, expire time.Duration) (string, error) {
	claims := AdminClaims{
		AdminID: adminID,
		Role:    "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func extractBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
