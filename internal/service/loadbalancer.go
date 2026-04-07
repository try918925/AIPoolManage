package service

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultFailureThreshold = 5
	DefaultRecoveryTimeout  = 30 * time.Second
	DefaultHalfOpenMax      = 3
)

type LoadBalancer struct {
	mu               sync.Mutex
	rdb              *redis.Client
	counters         map[int64]int64
	failureThreshold int
	recoveryTimeout  time.Duration
	halfOpenMax      int
}

func NewLoadBalancer(rdb *redis.Client, failureThreshold int, recoveryTimeout time.Duration, halfOpenMax int) *LoadBalancer {
	if failureThreshold <= 0 {
		failureThreshold = DefaultFailureThreshold
	}
	if recoveryTimeout <= 0 {
		recoveryTimeout = DefaultRecoveryTimeout
	}
	if halfOpenMax <= 0 {
		halfOpenMax = DefaultHalfOpenMax
	}
	return &LoadBalancer{
		rdb:              rdb,
		counters:         make(map[int64]int64),
		failureThreshold: failureThreshold,
		recoveryTimeout:  recoveryTimeout,
		halfOpenMax:      halfOpenMax,
	}
}

func (lb *LoadBalancer) SelectChannels(channels []*model.ProviderModel) []*model.ProviderModel {
	// Group by priority
	groups := make(map[int][]*model.ProviderModel)
	for _, ch := range channels {
		groups[ch.Priority] = append(groups[ch.Priority], ch)
	}

	priorities := make([]int, 0, len(groups))
	for p := range groups {
		priorities = append(priorities, p)
	}
	sort.Ints(priorities)

	var result []*model.ProviderModel
	for _, p := range priorities {
		group := groups[p]
		healthy := lb.filterHealthy(group)
		if len(healthy) == 0 {
			continue
		}
		sorted := lb.weightedRoundRobin(healthy)
		result = append(result, sorted...)
	}

	return result
}

func (lb *LoadBalancer) filterHealthy(channels []*model.ProviderModel) []*model.ProviderModel {
	var healthy []*model.ProviderModel
	for _, ch := range channels {
		if lb.IsHealthy(ch.ID) {
			healthy = append(healthy, ch)
		}
	}
	return healthy
}

func (lb *LoadBalancer) weightedRoundRobin(channels []*model.ProviderModel) []*model.ProviderModel {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	totalWeight := 0
	for _, ch := range channels {
		totalWeight += ch.Weight
	}

	type weightedChannel struct {
		channel       *model.ProviderModel
		currentWeight int
	}

	items := make([]weightedChannel, len(channels))
	for i, ch := range channels {
		items[i] = weightedChannel{
			channel:       ch,
			currentWeight: int(lb.counters[ch.ID]) + ch.Weight,
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].currentWeight > items[j].currentWeight
	})

	if len(items) > 0 {
		items[0].currentWeight -= totalWeight
	}

	for _, item := range items {
		lb.counters[item.channel.ID] = int64(item.currentWeight)
	}

	result := make([]*model.ProviderModel, len(items))
	for i, item := range items {
		result[i] = item.channel
	}
	return result
}

func (lb *LoadBalancer) IsHealthy(channelID int64) bool {
	ctx := context.Background()
	key := fmt.Sprintf("channel:health:%d", channelID)
	data, err := lb.rdb.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return true
	}

	status := data["status"]
	switch status {
	case "closed", "":
		return true
	case "open":
		recoveryAt, err := time.Parse(time.RFC3339, data["recovery_at"])
		if err != nil {
			return true
		}
		if time.Now().After(recoveryAt) {
			lb.rdb.HSet(ctx, key, "status", "half_open", "probe_count", "0")
			return true
		}
		return false
	case "half_open":
		probeCount, _ := strconv.Atoi(data["probe_count"])
		if probeCount < lb.halfOpenMax {
			lb.rdb.HIncrBy(ctx, key, "probe_count", 1)
			return true
		}
		return false
	}
	return true
}

func (lb *LoadBalancer) RecordFailure(channelID int64) {
	ctx := context.Background()
	key := fmt.Sprintf("channel:health:%d", channelID)

	failures, _ := lb.rdb.HIncrBy(ctx, key, "consecutive_failures", 1).Result()
	lb.rdb.HSet(ctx, key, "last_failure_at", time.Now().Format(time.RFC3339))

	if failures >= int64(lb.failureThreshold) {
		lb.rdb.HSet(ctx, key,
			"status", "open",
			"recovery_at", time.Now().Add(lb.recoveryTimeout).Format(time.RFC3339),
		)
		lb.rdb.Expire(ctx, key, 10*time.Minute)
	}
}

func (lb *LoadBalancer) RecordSuccess(channelID int64) {
	ctx := context.Background()
	key := fmt.Sprintf("channel:health:%d", channelID)
	lb.rdb.HSet(ctx, key,
		"status", "closed",
		"consecutive_failures", "0",
	)
	lb.rdb.Expire(ctx, key, 5*time.Minute)
}

func (lb *LoadBalancer) GetChannelHealth(channelID int64) map[string]string {
	ctx := context.Background()
	key := fmt.Sprintf("channel:health:%d", channelID)
	data, err := lb.rdb.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return map[string]string{
			"status":               "closed",
			"consecutive_failures": "0",
		}
	}
	if data["status"] == "" {
		data["status"] = "closed"
	}
	return data
}

func (lb *LoadBalancer) ResetChannelHealth(channelID int64) {
	ctx := context.Background()
	key := fmt.Sprintf("channel:health:%d", channelID)
	lb.rdb.Del(ctx, key)
}
