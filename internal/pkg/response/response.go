package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PagedData struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Items    interface{} `json:"items"`
}

type APIError struct {
	HTTPStatus int    `json:"-"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Code       string `json:"code"`
	Param      string `json:"param,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "success", Data: data})
}

func ErrorResponse(c *gin.Context, status int, errType, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
			"type":    errType,
			"code":    code,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, 400, "invalid_request_error", "invalid_request", message)
}

func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, 401, "authentication_error", "invalid_api_key", message)
}

func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, 403, "permission_error", "model_not_allowed", message)
}

func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, 404, "invalid_request_error", "model_not_found", message)
}

func RateLimited(c *gin.Context, message string) {
	ErrorResponse(c, 429, "rate_limit_error", "rate_limit_exceeded", message)
}

func QuotaExceeded(c *gin.Context) {
	ErrorResponse(c, 429, "rate_limit_error", "quota_exceeded", "token quota exceeded")
}

func InternalError(c *gin.Context) {
	ErrorResponse(c, 500, "api_error", "internal_error", "internal server error")
}

func UpstreamError(c *gin.Context, message string) {
	ErrorResponse(c, 502, "api_error", "upstream_error", message)
}

func ServiceUnavailable(c *gin.Context, message string) {
	ErrorResponse(c, 503, "api_error", "provider_unavailable", message)
}
