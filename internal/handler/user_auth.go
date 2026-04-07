package handler

import (
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/service"

	"github.com/gin-gonic/gin"
)

type UserAuthHandler struct {
	svc *service.UserService
}

func NewUserAuthHandler(svc *service.UserService) *UserAuthHandler {
	return &UserAuthHandler{svc: svc}
}

func (h *UserAuthHandler) Register(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		response.ErrorResponse(c, 400, "invalid_request_error", "register_failed", err.Error())
		return
	}

	response.Created(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *UserAuthHandler) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, token, err := h.svc.Login(&req)
	if err != nil {
		response.ErrorResponse(c, 401, "authentication_error", "login_failed", err.Error())
		return
	}

	response.Success(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *UserAuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id_int")
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid token")
		return
	}

	user, err := h.svc.GetByID(uid)
	if err != nil {
		response.ErrorResponse(c, 404, "not_found", "user_not_found", "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}
