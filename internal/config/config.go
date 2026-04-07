package config

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Redis        RedisConfig        `mapstructure:"redis"`
	Security     SecurityConfig     `mapstructure:"security"`
	Proxy        ProxyConfig        `mapstructure:"proxy"`
	LoadBalancer LoadBalancerConfig `mapstructure:"loadbalancer"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type SecurityConfig struct {
	MasterKey string        `mapstructure:"master_key"`
	JWTSecret string        `mapstructure:"jwt_secret"`
	JWTExpire time.Duration `mapstructure:"jwt_expire"`
}

type ProxyConfig struct {
	Timeout         time.Duration `mapstructure:"timeout"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	IdleConnTimeout time.Duration `mapstructure:"idle_conn_timeout"`
}

type LoadBalancerConfig struct {
	Strategy       string               `mapstructure:"strategy"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
}

type CircuitBreakerConfig struct {
	FailureThreshold int           `mapstructure:"failure_threshold"`
	RecoveryTimeout  time.Duration `mapstructure:"recovery_timeout"`
	HalfOpenMax      int           `mapstructure:"half_open_max_probes"`
}

func (c *Config) GetMasterKeyBytes() ([]byte, error) {
	return hex.DecodeString(c.Security.MasterKey)
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "120s")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "gateway")
	v.SetDefault("database.dbname", "ai_gateway")
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 100)
	v.SetDefault("security.jwt_expire", "24h")
	v.SetDefault("proxy.timeout", "120s")
	v.SetDefault("proxy.max_idle_conns", 200)
	v.SetDefault("proxy.idle_conn_timeout", "90s")
	v.SetDefault("loadbalancer.strategy", "weighted_round_robin")
	v.SetDefault("loadbalancer.circuit_breaker.failure_threshold", 5)
	v.SetDefault("loadbalancer.circuit_breaker.recovery_timeout", "30s")
	v.SetDefault("loadbalancer.circuit_breaker.half_open_max_probes", 3)

	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Map environment variables
	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.dbname", "DB_NAME")
	_ = v.BindEnv("redis.addr", "REDIS_ADDR")
	_ = v.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = v.BindEnv("security.master_key", "MASTER_ENCRYPTION_KEY")
	_ = v.BindEnv("security.jwt_secret", "JWT_SECRET")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		zap.L().Warn("config file not found, using defaults and env vars")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
