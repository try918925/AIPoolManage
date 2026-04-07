package handler

import (
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminModelHandler struct {
	svc *service.ModelService
}

func NewAdminModelHandler(svc *service.ModelService) *AdminModelHandler {
	return &AdminModelHandler{svc: svc}
}

func (h *AdminModelHandler) Create(c *gin.Context) {
	providerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	var req service.CreateModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	m, err := h.svc.Create(providerID, &req)
	if err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	response.Created(c, m)
}

func (h *AdminModelHandler) ListByProvider(c *gin.Context) {
	providerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	models, err := h.svc.ListByProvider(providerID)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, models)
}

func (h *AdminModelHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	var req service.UpdateModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	m, err := h.svc.Update(id, &req)
	if err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	response.Success(c, m)
}

func (h *AdminModelHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	response.Success(c, nil)
}
