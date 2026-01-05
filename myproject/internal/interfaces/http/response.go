package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error codes
const (
	CodeSuccess      = 0     // Success
	CodeBadRequest   = 400   // Bad request (validation error)
	CodeUnauthorized = 401   // Unauthorized
	CodeForbidden    = 403   // Forbidden
	CodeNotFound     = 404   // Resource not found
	CodeInternal     = 500   // Internal server error

	// Business error codes (1000+)
	CodeValidation   = 1001  // Validation failed
	CodeDuplicate    = 1002  // Duplicate entry
	CodeConflict     = 1003  // Business conflict
)

// Response is the standard API response.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success returns a successful response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// BadRequest returns a 400 error response.
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeBadRequest,
		Message: message,
	})
}

// NotFound returns a 404 error response.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// InternalError returns a 500 error response.
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeInternal,
		Message: message,
	})
}

// ValidationError returns a validation error response.
func ValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeValidation,
		Message: message,
	})
}
