package handler

import (
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"

	"github.com/gin-gonic/gin"
)

type UserModelHandler struct {
	modelSvc  *service.ModelService
	lb        *service.LoadBalancer
	usageRepo *repository.UsageRepo
}

func NewUserModelHandler(modelSvc *service.ModelService, lb *service.LoadBalancer, usageRepo *repository.UsageRepo) *UserModelHandler {
	return &UserModelHandler{modelSvc: modelSvc, lb: lb, usageRepo: usageRepo}
}

func (h *UserModelHandler) ListModels(c *gin.Context) {
	models, err := h.modelSvc.ListAvailableModels()
	if err != nil {
		response.InternalError(c)
		return
	}

	c.JSON(200, gin.H{
		"object": "list",
		"data":   models,
	})
}

// ListModelsPortal returns models in the standard {code, data} format for the user portal
func (h *UserModelHandler) ListModelsPortal(c *gin.Context) {
	models, err := h.modelSvc.ListAvailableModels()
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, models)
}

func (h *UserModelHandler) GetModelDetail(c *gin.Context) {
	modelName := c.Param("model")
	detail, err := h.modelSvc.GetModelDetail(modelName)
	if err != nil || detail == nil {
		response.NotFound(c, "model not found: "+modelName)
		return
	}

	response.Success(c, detail)
}

func (h *UserModelHandler) GetModelChannels(c *gin.Context) {
	modelName := c.Param("model")
	channels, err := h.modelSvc.GetModelChannels(modelName)
	if err != nil || len(channels) == 0 {
		response.NotFound(c, "model not found: "+modelName)
		return
	}

	var channelList []map[string]interface{}
	for _, ch := range channels {
		health := h.lb.GetChannelHealth(ch.ID)
		avgLatency, _ := h.usageRepo.GetAvgLatency(modelName, ch.ID, 24)

		entry := map[string]interface{}{
			"channel_id":     ch.ID,
			"priority":       ch.Priority,
			"status":         health["status"],
			"avg_latency_ms": avgLatency,
		}

		if ch.Provider != nil {
			entry["provider_type"] = ch.Provider.Type
			entry["provider_name"] = ch.Provider.Name
		}

		channelList = append(channelList, entry)
	}

	response.Success(c, gin.H{
		"model_name": modelName,
		"channels":   channelList,
	})
}
