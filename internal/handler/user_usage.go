package handler

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserUsageHandler struct {
	usageRepo *repository.UsageRepo
}

func NewUserUsageHandler(usageRepo *repository.UsageRepo) *UserUsageHandler {
	return &UserUsageHandler{usageRepo: usageRepo}
}

func (h *UserUsageHandler) GetUsage(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		if k, exists := c.Get("user_key"); exists {
			userKey := k.(*model.UserAPIKey)
			userID = userKey.UserID
		}
	}
	if userID == "" {
		response.Unauthorized(c, "missing user identity")
		return
	}

	startDate, endDate := parseDateRange(c)
	modelFilter := c.Query("model")
	groupBy := c.DefaultQuery("group_by", "day")

	summary, err := h.usageRepo.GetUsageSummary(userID, startDate, endDate, modelFilter)
	if err != nil {
		response.InternalError(c)
		return
	}

	breakdown, err := h.usageRepo.GetUsageBreakdown(userID, startDate, endDate, modelFilter, groupBy)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, gin.H{
		"total_requests":          summary.TotalRequests,
		"total_tokens":            summary.TotalTokens,
		"total_prompt_tokens":     summary.TotalPromptTokens,
		"total_completion_tokens": summary.TotalCompletionTokens,
		"breakdown":               breakdown,
	})
}

func (h *UserUsageHandler) GetUsageDetails(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		if k, exists := c.Get("user_key"); exists {
			userKey := k.(*model.UserAPIKey)
			userID = userKey.UserID
		}
	}
	if userID == "" {
		response.Unauthorized(c, "missing user identity")
		return
	}

	startDate, endDate := parseDateRange(c)
	modelFilter := c.Query("model")
	statusFilter := c.Query("status")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.usageRepo.GetUsageDetails(userID, startDate, endDate, modelFilter, statusFilter, page, pageSize)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, response.PagedData{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    logs,
	})
}

func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now

	if s := c.Query("start_date"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			startDate = t
		}
	}
	if s := c.Query("end_date"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			endDate = t.Add(24*time.Hour - time.Second)
		}
	}

	return startDate, endDate
}
