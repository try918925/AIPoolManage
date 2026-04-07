package handler

import (
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminProviderHandler struct {
	svc *service.ProviderService
}

func NewAdminProviderHandler(svc *service.ProviderService) *AdminProviderHandler {
	return &AdminProviderHandler{svc: svc}
}

func (h *AdminProviderHandler) Create(c *gin.Context) {
	var req service.CreateProviderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	provider, err := h.svc.Create(&req)
	if err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	response.Created(c, provider)
}

func (h *AdminProviderHandler) List(c *gin.Context) {
	providers, err := h.svc.List()
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, providers)
}

func (h *AdminProviderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	var req service.UpdateProviderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	provider, err := h.svc.Update(id, &req)
	if err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	response.Success(c, provider)
}

func (h *AdminProviderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		response.ErrorResponse(c, 500, "api_error", "internal_error", err.Error())
		return
	}

	response.Success(c, nil)
}
