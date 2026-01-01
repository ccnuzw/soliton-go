package http

import "github.com/gin-gonic/gin"

// Response is the standard API response structure.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success returns a successful response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage returns a successful response with custom message.
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error returns an error response.
func Error(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData returns an error response with additional data.
func ErrorWithData(c *gin.Context, httpCode int, code int, message string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest returns a 400 error response.
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, 400, message)
}

// NotFound returns a 404 error response.
func NotFound(c *gin.Context, message string) {
	Error(c, 404, 404, message)
}

// InternalError returns a 500 error response.
func InternalError(c *gin.Context, message string) {
	Error(c, 500, 500, message)
}

// Unauthorized returns a 401 error response.
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, 401, message)
}

// Forbidden returns a 403 error response.
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, 403, message)
}
