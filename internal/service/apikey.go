package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/pkg/hash"
	"awesomeProject/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type APIKeyService struct {
	repo  *repository.APIKeyRepo
	redis *redis.Client
}

func NewAPIKeyService(repo *repository.APIKeyRepo, rdb *redis.Client) *APIKeyService {
	return &APIKeyService{repo: repo, redis: rdb}
}

type CreateKeyReq struct {
	Name          string   `json:"name"`
	RateLimit     int      `json:"rate_limit"`
	QuotaLimit    int64    `json:"quota_limit"`
	AllowedModels []string `json:"allowed_models"`
	ExpiresAt     *string  `json:"expires_at"`
}

func (s *APIKeyService) Create(userID string, req *CreateKeyReq) (*model.UserAPIKey, string, error) {
	fullKey, prefix, keyHash := hash.GenerateAPIKey()

	rateLimit := 60
	if req.RateLimit > 0 {
		rateLimit = req.RateLimit
	}

	key := &model.UserAPIKey{
		UserID:        userID,
		Name:          req.Name,
		KeyHash:       keyHash,
		KeyPrefix:     prefix,
		Enabled:       true,
		RateLimit:     rateLimit,
		QuotaLimit:    req.QuotaLimit,
		AllowedModels: req.AllowedModels,
	}

	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err == nil {
			key.ExpiresAt = &t
		}
	}

	if err := s.repo.Create(key); err != nil {
		return nil, "", err
	}

	return key, fullKey, nil
}

func (s *APIKeyService) ListByUser(userID string) ([]model.UserAPIKey, error) {
	return s.repo.FindByUserID(userID)
}

func (s *APIKeyService) ListAll() ([]model.UserAPIKey, error) {
	return s.repo.FindAll()
}

func (s *APIKeyService) AdminDelete(id int64) error {
	key, err := s.repo.FindByID(id)
	if err == nil {
		ctx := context.Background()
		s.redis.Del(ctx, "apikey:"+key.KeyHash)
	}
	return s.repo.AdminDelete(id)
}

func (s *APIKeyService) Delete(id int64, userID string) error {
	// Invalidate cache
	key, err := s.repo.FindByID(id)
	if err == nil {
		ctx := context.Background()
		s.redis.Del(ctx, "apikey:"+key.KeyHash)
	}
	return s.repo.Delete(id, userID)
}

func (s *APIKeyService) ValidateKey(apiKey string) (*model.UserAPIKey, error) {
	keyHash := hash.SHA256Hex(apiKey)
	ctx := context.Background()

	// Check Redis cache
	cached, err := s.redis.Get(ctx, "apikey:"+keyHash).Result()
	if err == nil {
		var key model.UserAPIKey
		if err := json.Unmarshal([]byte(cached), &key); err == nil {
			return &key, nil
		}
	}

	// Cache miss, query DB
	key, err := s.repo.FindByHash(keyHash)
	if err != nil {
		return nil, fmt.Errorf("invalid api key")
	}

	if !key.Enabled {
		return nil, fmt.Errorf("api key disabled")
	}

	// Write to cache (TTL 5 min)
	data, _ := json.Marshal(key)
	s.redis.Set(ctx, "apikey:"+keyHash, string(data), 5*time.Minute)

	return key, nil
}

func (s *APIKeyService) UpdateLastUsed(id int64) {
	s.repo.UpdateLastUsed(id)
}

func (s *APIKeyService) IncrementQuota(id int64, tokens int64) error {
	return s.repo.IncrementQuotaUsed(id, tokens)
}

func (s *APIKeyService) InvalidateCache(keyHash string) {
	ctx := context.Background()
	s.redis.Del(ctx, "apikey:"+keyHash)
}
