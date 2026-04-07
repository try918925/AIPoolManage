package handler

import (
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminChannelHandler struct {
	lb        *service.LoadBalancer
	modelRepo *repository.ModelRepo
	usageRepo *repository.UsageRepo
}

func NewAdminChannelHandler(lb *service.LoadBalancer, modelRepo *repository.ModelRepo, usageRepo *repository.UsageRepo) *AdminChannelHandler {
	return &AdminChannelHandler{lb: lb, modelRepo: modelRepo, usageRepo: usageRepo}
}

func (h *AdminChannelHandler) GetHealth(c *gin.Context) {
	models, err := h.modelRepo.FindAllEnabled()
	if err != nil {
		response.InternalError(c)
		return
	}

	var result []map[string]interface{}
	for _, m := range models {
		health := h.lb.GetChannelHealth(m.ID)

		entry := map[string]interface{}{
			"channel_id":           m.ID,
			"model_name":           m.ModelName,
			"status":               health["status"],
			"consecutive_failures": health["consecutive_failures"],
			"weight":               m.Weight,
			"priority":             m.Priority,
		}

		if m.Provider != nil {
			entry["provider_name"] = m.Provider.Name
		}
		if v, ok := health["last_failure_at"]; ok && v != "" {
			entry["last_failure_at"] = v
		}
		if v, ok := health["recovery_at"]; ok && v != "" {
			entry["recovery_at"] = v
		}

		result = append(result, entry)
	}

	response.Success(c, result)
}

func (h *AdminChannelHandler) ResetHealth(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid channel id")
		return
	}

	h.lb.ResetChannelHealth(id)
	response.Success(c, nil)
}
