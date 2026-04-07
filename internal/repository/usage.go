package repository

import (
	"awesomeProject/internal/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UsageRepo struct {
	db *gorm.DB
}

func NewUsageRepo(db *gorm.DB) *UsageRepo {
	return &UsageRepo{db: db}
}

func (r *UsageRepo) Create(log *model.UsageLog) error {
	return r.db.Create(log).Error
}

type UsageSummary struct {
	TotalRequests        int64   `json:"total_requests"`
	TotalTokens          int64   `json:"total_tokens"`
	TotalPromptTokens    int64   `json:"total_prompt_tokens"`
	TotalCompletionTokens int64  `json:"total_completion_tokens"`
}

type UsageBreakdown struct {
	Model            string  `json:"model,omitempty"`
	Day              string  `json:"day,omitempty"`
	Requests         int64   `json:"requests"`
	PromptTokens     int64   `json:"prompt_tokens"`
	CompletionTokens int64   `json:"completion_tokens"`
	TotalTokens      int64   `json:"total_tokens"`
	AvgLatencyMs     float64 `json:"avg_latency_ms"`
}

func (r *UsageRepo) GetUsageSummary(userID string, startDate, endDate time.Time, modelFilter string) (*UsageSummary, error) {
	query := r.db.Model(&model.UsageLog{}).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate)

	if modelFilter != "" {
		query = query.Where("model_name = ?", modelFilter)
	}

	var summary UsageSummary
	err := query.Select(
		"COUNT(*) as total_requests, "+
			"COALESCE(SUM(total_tokens), 0) as total_tokens, "+
			"COALESCE(SUM(prompt_tokens), 0) as total_prompt_tokens, "+
			"COALESCE(SUM(completion_tokens), 0) as total_completion_tokens",
	).Scan(&summary).Error

	return &summary, err
}

func (r *UsageRepo) GetUsageBreakdown(userID string, startDate, endDate time.Time, modelFilter, groupBy string) ([]UsageBreakdown, error) {
	query := r.db.Model(&model.UsageLog{}).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate)

	if modelFilter != "" {
		query = query.Where("model_name = ?", modelFilter)
	}

	selectFields := "COUNT(*) as requests, " +
		"COALESCE(SUM(prompt_tokens), 0) as prompt_tokens, " +
		"COALESCE(SUM(completion_tokens), 0) as completion_tokens, " +
		"COALESCE(SUM(total_tokens), 0) as total_tokens, " +
		"COALESCE(AVG(latency_ms), 0) as avg_latency_ms"

	var groupFields []string
	switch groupBy {
	case "model":
		groupFields = []string{"model_name"}
		selectFields = "model_name as model, " + selectFields
	case "day,model":
		groupFields = []string{"DATE(created_at)", "model_name"}
		selectFields = fmt.Sprintf("TO_CHAR(created_at, 'YYYY-MM-DD') as day, model_name as model, %s", selectFields)
	default: // "day"
		groupFields = []string{"DATE(created_at)"}
		selectFields = fmt.Sprintf("TO_CHAR(created_at, 'YYYY-MM-DD') as day, %s", selectFields)
	}

	var breakdowns []UsageBreakdown
	q := query.Select(selectFields)
	for _, g := range groupFields {
		q = q.Group(g)
	}
	err := q.Order("1 ASC").Scan(&breakdowns).Error

	return breakdowns, err
}

func (r *UsageRepo) GetUsageDetails(userID string, startDate, endDate time.Time, modelFilter, statusFilter string, page, pageSize int) ([]model.UsageLog, int64, error) {
	query := r.db.Model(&model.UsageLog{}).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate)

	if modelFilter != "" {
		query = query.Where("model_name = ?", modelFilter)
	}
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	var total int64
	query.Count(&total)

	var logs []model.UsageLog
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

func (r *UsageRepo) GetAvgLatency(modelName string, channelID int64, hours int) (float64, error) {
	var result struct {
		AvgLatency float64
	}
	err := r.db.Model(&model.UsageLog{}).
		Where("model_name = ? AND channel_id = ? AND status = 'success' AND created_at >= ?",
			modelName, channelID, time.Now().Add(-time.Duration(hours)*time.Hour)).
		Select("COALESCE(AVG(latency_ms), 0) as avg_latency").
		Scan(&result).Error
	return result.AvgLatency, err
}
