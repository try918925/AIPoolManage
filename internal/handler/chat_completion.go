package handler

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/pkg/response"
	"awesomeProject/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatCompletionHandler struct {
	chatSvc *service.ChatService
}

func NewChatCompletionHandler(chatSvc *service.ChatService) *ChatCompletionHandler {
	return &ChatCompletionHandler{chatSvc: chatSvc}
}

func (h *ChatCompletionHandler) Handle(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userKey := c.MustGet("user_key").(*model.UserAPIKey)

	// Parse route hints from headers
	var hint *model.RouteHint
	channelIDStr := c.GetHeader("X-Channel-Id")
	preferredProvider := c.GetHeader("X-Preferred-Provider")

	if channelIDStr != "" || preferredProvider != "" {
		hint = &model.RouteHint{
			PreferredProvider: preferredProvider,
		}
		if channelIDStr != "" {
			if id, err := strconv.ParseInt(channelIDStr, 10, 64); err == nil {
				hint.ChannelID = id
			}
		}
	}

	if req.Stream {
		h.handleStream(c, &req, userKey, hint)
		return
	}

	h.handleNonStream(c, &req, userKey, hint)
}

func (h *ChatCompletionHandler) handleNonStream(c *gin.Context, req *model.ChatRequest, userKey *model.UserAPIKey, hint *model.RouteHint) {
	result, err := h.chatSvc.ChatCompletion(c.Request.Context(), req, userKey, hint)
	if err != nil {
		handleChatError(c, err)
		return
	}

	// Set route info headers
	if result.RouteInfo != nil {
		c.Header("X-Channel-Id", strconv.FormatInt(result.RouteInfo.ChannelID, 10))
		c.Header("X-Provider-Type", result.RouteInfo.ProviderType)
		c.Header("X-Provider-Name", result.RouteInfo.ProviderName)
	}

	c.JSON(http.StatusOK, result.Response)
}

func (h *ChatCompletionHandler) handleStream(c *gin.Context, req *model.ChatRequest, userKey *model.UserAPIKey, hint *model.RouteHint) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.InternalError(c)
		return
	}

	routeInfo, err := h.chatSvc.ChatCompletionStream(c.Request.Context(), req, userKey, hint, c.Writer, flusher)
	if err != nil {
		handleChatError(c, err)
		return
	}

	if routeInfo != nil {
		c.Header("X-Channel-Id", strconv.FormatInt(routeInfo.ChannelID, 10))
		c.Header("X-Provider-Type", routeInfo.ProviderType)
		c.Header("X-Provider-Name", routeInfo.ProviderName)
	}
}

func handleChatError(c *gin.Context, err error) {
	msg := err.Error()
	switch msg {
	case "model_not_found":
		response.NotFound(c, "the model does not exist or you do not have access to it")
	case "model_forbidden":
		response.Forbidden(c, "you do not have access to this model")
	case "quota_exceeded":
		response.QuotaExceeded(c)
	default:
		if len(msg) > 200 {
			msg = msg[:200]
		}
		response.ServiceUnavailable(c, msg)
	}
}
