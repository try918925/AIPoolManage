package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// User represents a registered user
type User struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"type:varchar(100);not null;uniqueIndex"`
	Email        string    `json:"email,omitempty" gorm:"type:varchar(200);uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	Role         string    `json:"role" gorm:"type:varchar(20);default:'user'"`
	Enabled      bool      `json:"enabled" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

type Provider struct {
	ID              int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string         `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
	Type            string         `json:"type" gorm:"type:varchar(50);not null;index"`
	BaseURL         string         `json:"base_url" gorm:"type:varchar(500);not null"`
	APIKeyEncrypted string         `json:"-" gorm:"column:api_key_encrypted;type:text;not null"`
	OrgID           string         `json:"org_id,omitempty" gorm:"type:varchar(200)"`
	Enabled         bool           `json:"enabled" gorm:"default:true;index"`
	Config          datatypes.JSON `json:"config,omitempty" gorm:"type:jsonb;default:'{}'"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Provider) TableName() string { return "providers" }

type ProviderModel struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ProviderID       int64          `json:"provider_id" gorm:"not null;index"`
	ModelName        string         `json:"model_name" gorm:"type:varchar(100);not null;uniqueIndex"`
	ModelID          string         `json:"model_id" gorm:"type:varchar(100);not null"`
	ModelType        string         `json:"model_type" gorm:"type:varchar(20);default:'chat'"`
	Enabled          bool           `json:"enabled" gorm:"default:true;index"`
	Weight           int            `json:"weight" gorm:"default:1"`
	Priority         int            `json:"priority" gorm:"default:0"`
	MaxContextTokens *int           `json:"max_context_tokens,omitempty"`
	InputPrice       *float64       `json:"input_price,omitempty" gorm:"type:decimal(10,6)"`
	OutputPrice      *float64       `json:"output_price,omitempty" gorm:"type:decimal(10,6)"`
	Config           datatypes.JSON `json:"config,omitempty" gorm:"type:jsonb;default:'{}'"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`

	Provider *Provider `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
}

func (ProviderModel) TableName() string { return "provider_models" }

type UserAPIKey struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        string         `json:"user_id" gorm:"type:varchar(100);not null;index"`
	Name          string         `json:"name,omitempty" gorm:"type:varchar(100)"`
	KeyHash       string         `json:"-" gorm:"type:varchar(64);not null;uniqueIndex"`
	KeyPrefix     string         `json:"key_prefix" gorm:"type:varchar(12);not null"`
	Enabled       bool           `json:"enabled" gorm:"default:true"`
	RateLimit     int            `json:"rate_limit" gorm:"default:60"`
	QuotaLimit    int64          `json:"quota_limit" gorm:"default:0"`
	QuotaUsed     int64          `json:"quota_used" gorm:"default:0"`
	AllowedModels pq.StringArray `json:"allowed_models,omitempty" gorm:"type:text[]"`
	ExpiresAt     *time.Time     `json:"expires_at,omitempty"`
	LastUsedAt    *time.Time     `json:"last_used_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (UserAPIKey) TableName() string { return "user_api_keys" }

func (k *UserAPIKey) HasModelAccess(modelName string) bool {
	if len(k.AllowedModels) == 0 {
		return true
	}
	for _, m := range k.AllowedModels {
		if m == modelName {
			return true
		}
	}
	return false
}

type UsageLog struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserKeyID        int64     `json:"user_key_id" gorm:"not null;index"`
	UserID           string    `json:"user_id" gorm:"type:varchar(100);not null;index:idx_usage_user"`
	ProviderID       int64     `json:"provider_id" gorm:"not null"`
	ModelName        string    `json:"model_name" gorm:"type:varchar(100);not null;index:idx_usage_model"`
	PromptTokens     int       `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int       `json:"completion_tokens" gorm:"default:0"`
	TotalTokens      int       `json:"total_tokens" gorm:"default:0"`
	LatencyMs        int       `json:"latency_ms"`
	Status           string    `json:"status" gorm:"type:varchar(20);not null;default:'success'"`
	ErrorMessage     string    `json:"error_message,omitempty" gorm:"type:text"`
	RequestIP        string    `json:"request_ip,omitempty" gorm:"type:varchar(45)"`
	ChannelID        int64     `json:"channel_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"index:idx_usage_created"`
}

func (UsageLog) TableName() string { return "usage_logs" }

// Chat request/response types (OpenAI compatible)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model" binding:"required"`
	Messages    []ChatMessage `json:"messages" binding:"required"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream"`
	TopP        *float64      `json:"top_p,omitempty"`
	N           *int          `json:"n,omitempty"`
	Stop        interface{}   `json:"stop,omitempty"`
}

type ChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []ChatChoice   `json:"choices"`
	Usage   *ChatUsage     `json:"usage,omitempty"`
	Route   *RouteInfo     `json:"-"`
}

type ChatChoice struct {
	Index        int          `json:"index"`
	Message      *ChatMessage `json:"message,omitempty"`
	Delta        *ChatMessage `json:"delta,omitempty"`
	FinishReason *string      `json:"finish_reason"`
}

type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type RouteInfo struct {
	ChannelID    int64  `json:"channel_id"`
	ProviderType string `json:"provider_type"`
	ProviderName string `json:"provider_name"`
}

type RouteHint struct {
	ChannelID         int64
	PreferredProvider string
}
