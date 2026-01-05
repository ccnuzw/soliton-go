package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误码定义
const (
	CodeSuccess      = 0     // 成功
	CodeBadRequest   = 400   // 请求错误
	CodeUnauthorized = 401   // 未授权
	CodeForbidden    = 403   // 禁止访问
	CodeNotFound     = 404   // 资源不存在
	CodeInternal     = 500   // 服务器内部错误

	// 业务错误码 (1000+)
	CodeValidation   = 1001  // 校验失败
	CodeDuplicate    = 1002  // 重复数据
	CodeConflict     = 1003  // 业务冲突
)

// Response 是标准的 API 响应结构体。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// BadRequest 返回 400 错误响应。
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeBadRequest,
		Message: message,
	})
}

// NotFound 返回 404 错误响应。
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// InternalError 返回 500 错误响应。
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeInternal,
		Message: message,
	})
}

// ValidationError 返回校验错误响应。
func ValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeValidation,
		Message: message,
	})
}
