package handler

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserAPIKeyHandler struct {
	svc *service.APIKeyService
}

func NewUserAPIKeyHandler(svc *service.APIKeyService) *UserAPIKeyHandler {
	return &UserAPIKeyHandler{svc: svc}
}

func (h *UserAPIKeyHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		// For user key authenticated requests, extract from user_key
		if k, exists := c.Get("user_key"); exists {
			userKey := k.(*model.UserAPIKey)
			userID = userKey.UserID
		}
	}
	if userID == "" {
		userID = c.GetString("admin_id")
	}

	var req service.CreateKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	key, fullKey, err := h.svc.Create(userID, &req)
	if err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	data := map[string]interface{}{
		"id":         key.ID,
		"name":       key.Name,
		"key":        fullKey,
		"key_prefix": key.KeyPrefix,
		"rate_limit": key.RateLimit,
		"quota_limit": key.QuotaLimit,
		"created_at": key.CreatedAt,
	}
	if key.ExpiresAt != nil {
		data["expires_at"] = key.ExpiresAt
	}

	response.Created(c, data)
}

func (h *UserAPIKeyHandler) List(c *gin.Context) {
	// Admin sees all keys
	if adminID := c.GetString("admin_id"); adminID != "" {
		keys, err := h.svc.ListAll()
		if err != nil {
			response.InternalError(c)
			return
		}
		response.Success(c, keys)
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		if k, exists := c.Get("user_key"); exists {
			userKey := k.(*model.UserAPIKey)
			userID = userKey.UserID
		}
	}

	keys, err := h.svc.ListByUser(userID)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, keys)
}

func (h *UserAPIKeyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid key id")
		return
	}

	// Admin can delete any key
	if adminID := c.GetString("admin_id"); adminID != "" {
		if err := h.svc.AdminDelete(id); err != nil {
			response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
			return
		}
		response.Success(c, nil)
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		if k, exists := c.Get("user_key"); exists {
			userKey := k.(*model.UserAPIKey)
			userID = userKey.UserID
		}
	}

	if err := h.svc.Delete(id, userID); err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	response.Success(c, nil)
}
