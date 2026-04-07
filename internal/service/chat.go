package service

import (
	"awesomeProject/internal/adapter"
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ChatService struct {
	modelRepo       *repository.ModelRepo
	providerRepo    *repository.ProviderRepo
	usageRepo       *repository.UsageRepo
	apiKeyRepo      *repository.APIKeyRepo
	providerService *ProviderService
	loadBalancer    *LoadBalancer
	adapterRegistry *adapter.Registry
	httpClient      *http.Client
	rdb             *redis.Client
}

func NewChatService(
	modelRepo *repository.ModelRepo,
	providerRepo *repository.ProviderRepo,
	usageRepo *repository.UsageRepo,
	apiKeyRepo *repository.APIKeyRepo,
	providerService *ProviderService,
	lb *LoadBalancer,
	registry *adapter.Registry,
	httpClient *http.Client,
	rdb *redis.Client,
) *ChatService {
	return &ChatService{
		modelRepo:       modelRepo,
		providerRepo:    providerRepo,
		usageRepo:       usageRepo,
		apiKeyRepo:      apiKeyRepo,
		providerService: providerService,
		loadBalancer:    lb,
		adapterRegistry: registry,
		httpClient:      httpClient,
		rdb:             rdb,
	}
}

type ChatResult struct {
	Response  *model.ChatResponse
	Usage     *model.ChatUsage
	RouteInfo *model.RouteInfo
}

func (s *ChatService) ChatCompletion(ctx context.Context, req *model.ChatRequest, userKey *model.UserAPIKey, hint *model.RouteHint) (*ChatResult, error) {
	// 1. Find channels for the model
	channels, err := s.modelRepo.FindAllByName(req.Model)
	if err != nil || len(channels) == 0 {
		return nil, fmt.Errorf("model_not_found")
	}

	// 2. Check model permission
	if !userKey.HasModelAccess(req.Model) {
		return nil, fmt.Errorf("model_forbidden")
	}

	// 3. Check quota
	if userKey.QuotaLimit > 0 && userKey.QuotaUsed >= userKey.QuotaLimit {
		return nil, fmt.Errorf("quota_exceeded")
	}

	// 4. Apply route hints and select channels
	channelPtrs := make([]*model.ProviderModel, len(channels))
	for i := range channels {
		channelPtrs[i] = &channels[i]
	}
	sortedChannels := s.applyRouteHint(channelPtrs, hint)

	// 5. Try channels in order (failover)
	var lastErr error
	for _, channel := range sortedChannels {
		if !s.loadBalancer.IsHealthy(channel.ID) {
			continue
		}

		provider, err := s.providerRepo.FindByID(channel.ProviderID)
		if err != nil {
			continue
		}

		providerAPIKey, err := s.providerService.DecryptAPIKey(provider)
		if err != nil {
			continue
		}

		adp, err := s.adapterRegistry.Get(provider.Type)
		if err != nil {
			continue
		}

		start := time.Now()
		upstreamReq, err := adp.ConvertRequest(req, channel, provider, providerAPIKey)
		if err != nil {
			continue
		}
		upstreamReq = upstreamReq.WithContext(ctx)

		upstreamResp, err := s.httpClient.Do(upstreamReq)
		if err != nil {
			s.loadBalancer.RecordFailure(channel.ID)
			lastErr = err
			continue
		}

		if upstreamResp.StatusCode >= 500 {
			s.loadBalancer.RecordFailure(channel.ID)
			lastErr = fmt.Errorf("upstream returned %d", upstreamResp.StatusCode)
			upstreamResp.Body.Close()
			continue
		}

		resp, usage, err := adp.ConvertResponse(upstreamResp, req.Model)
		if err != nil {
			s.loadBalancer.RecordFailure(channel.ID)
			lastErr = err
			continue
		}

		latency := time.Since(start).Milliseconds()
		s.loadBalancer.RecordSuccess(channel.ID)

		routeInfo := &model.RouteInfo{
			ChannelID:    channel.ID,
			ProviderType: provider.Type,
			ProviderName: provider.Name,
		}

		// Async record usage
		go s.recordUsage(userKey, channel, provider, usage, int(latency), "success", "")

		return &ChatResult{
			Response:  resp,
			Usage:     usage,
			RouteInfo: routeInfo,
		}, nil
	}

	return nil, fmt.Errorf("all channels exhausted: %v", lastErr)
}

func (s *ChatService) ChatCompletionStream(ctx context.Context, req *model.ChatRequest, userKey *model.UserAPIKey, hint *model.RouteHint, writer http.ResponseWriter, flusher http.Flusher) (*model.RouteInfo, error) {
	channels, err := s.modelRepo.FindAllByName(req.Model)
	if err != nil || len(channels) == 0 {
		return nil, fmt.Errorf("model_not_found")
	}

	if !userKey.HasModelAccess(req.Model) {
		return nil, fmt.Errorf("model_forbidden")
	}

	if userKey.QuotaLimit > 0 && userKey.QuotaUsed >= userKey.QuotaLimit {
		return nil, fmt.Errorf("quota_exceeded")
	}

	channelPtrs := make([]*model.ProviderModel, len(channels))
	for i := range channels {
		channelPtrs[i] = &channels[i]
	}
	sortedChannels := s.applyRouteHint(channelPtrs, hint)

	var lastErr error
	for _, channel := range sortedChannels {
		if !s.loadBalancer.IsHealthy(channel.ID) {
			continue
		}

		provider, err := s.providerRepo.FindByID(channel.ProviderID)
		if err != nil {
			continue
		}

		providerAPIKey, err := s.providerService.DecryptAPIKey(provider)
		if err != nil {
			continue
		}

		adp, err := s.adapterRegistry.Get(provider.Type)
		if err != nil {
			continue
		}

		start := time.Now()
		upstreamReq, err := adp.ConvertRequest(req, channel, provider, providerAPIKey)
		if err != nil {
			continue
		}
		upstreamReq = upstreamReq.WithContext(ctx)

		upstreamResp, err := s.httpClient.Do(upstreamReq)
		if err != nil {
			s.loadBalancer.RecordFailure(channel.ID)
			lastErr = err
			continue
		}

		if upstreamResp.StatusCode >= 500 {
			s.loadBalancer.RecordFailure(channel.ID)
			lastErr = fmt.Errorf("upstream returned %d", upstreamResp.StatusCode)
			upstreamResp.Body.Close()
			continue
		}

		usage, err := adp.ConvertStreamResponse(upstreamResp, writer, flusher, req.Model)
		latency := time.Since(start).Milliseconds()

		if err != nil {
			s.loadBalancer.RecordFailure(channel.ID)
			lastErr = err
			continue
		}

		s.loadBalancer.RecordSuccess(channel.ID)

		routeInfo := &model.RouteInfo{
			ChannelID:    channel.ID,
			ProviderType: provider.Type,
			ProviderName: provider.Name,
		}

		go s.recordUsage(userKey, channel, provider, usage, int(latency), "success", "")

		return routeInfo, nil
	}

	return nil, fmt.Errorf("all channels exhausted: %v", lastErr)
}

func (s *ChatService) applyRouteHint(channels []*model.ProviderModel, hint *model.RouteHint) []*model.ProviderModel {
	if hint == nil {
		return s.loadBalancer.SelectChannels(channels)
	}

	if hint.ChannelID > 0 {
		for _, ch := range channels {
			if ch.ID == hint.ChannelID {
				result := []*model.ProviderModel{ch}
				for _, other := range channels {
					if other.ID != hint.ChannelID {
						result = append(result, other)
					}
				}
				return result
			}
		}
	}

	if hint.PreferredProvider != "" {
		preferred := make([]*model.ProviderModel, 0)
		fallback := make([]*model.ProviderModel, 0)
		for _, ch := range channels {
			if ch.Provider != nil && ch.Provider.Type == hint.PreferredProvider {
				preferred = append(preferred, ch)
			} else {
				fallback = append(fallback, ch)
			}
		}
		sorted := s.loadBalancer.SelectChannels(preferred)
		sorted = append(sorted, s.loadBalancer.SelectChannels(fallback)...)
		return sorted
	}

	return s.loadBalancer.SelectChannels(channels)
}

func (s *ChatService) recordUsage(userKey *model.UserAPIKey, channel *model.ProviderModel, provider *model.Provider, usage *model.ChatUsage, latencyMs int, status, errMsg string) {
	if usage == nil {
		usage = &model.ChatUsage{}
	}

	log := &model.UsageLog{
		UserKeyID:        userKey.ID,
		UserID:           userKey.UserID,
		ProviderID:       provider.ID,
		ModelName:        channel.ModelName,
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
		LatencyMs:        latencyMs,
		Status:           status,
		ErrorMessage:     errMsg,
		ChannelID:        channel.ID,
	}

	if err := s.usageRepo.Create(log); err != nil {
		zap.L().Error("failed to record usage", zap.Error(err))
	}

	// Update quota
	if usage.TotalTokens > 0 {
		if err := s.apiKeyRepo.IncrementQuotaUsed(userKey.ID, int64(usage.TotalTokens)); err != nil {
			zap.L().Error("failed to increment quota", zap.Error(err))
		}
		// Also update Redis for real-time tracking
		ctx := context.Background()
		s.rdb.IncrBy(ctx, fmt.Sprintf("quota:%d", userKey.ID), int64(usage.TotalTokens))
	}

	// Update last used
	s.apiKeyRepo.UpdateLastUsed(userKey.ID)
}
