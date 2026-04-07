package handler

import (
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	svc *service.UserService
}

func NewAdminUserHandler(svc *service.UserService) *AdminUserHandler {
	return &AdminUserHandler{svc: svc}
}

func (h *AdminUserHandler) List(c *gin.Context) {
	var req service.AdminListUsersReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	users, total, err := h.svc.AdminListUsers(&req)
	if err != nil {
		response.InternalError(c)
		return
	}

	type userItem struct {
		ID        int64  `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		Enabled   bool   `json:"enabled"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	items := make([]userItem, 0, len(users))
	for _, u := range users {
		items = append(items, userItem{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role,
			Enabled:   u.Enabled,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.Success(c, gin.H{
		"users":     items,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

func (h *AdminUserHandler) Create(c *gin.Context) {
	var req service.AdminCreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.svc.AdminCreateUser(&req)
	if err != nil {
		response.ErrorResponse(c, 400, "invalid_request_error", "create_user_failed", err.Error())
		return
	}

	response.Created(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"enabled":  user.Enabled,
	})
}

func (h *AdminUserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req service.AdminUpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.svc.AdminUpdateUser(id, &req)
	if err != nil {
		response.ErrorResponse(c, 400, "invalid_request_error", "update_user_failed", err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"enabled":  user.Enabled,
	})
}

func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.svc.AdminDeleteUser(id); err != nil {
		response.ErrorResponse(c, 400, "invalid_request_error", "delete_user_failed", err.Error())
		return
	}

	response.Success(c, nil)
}
