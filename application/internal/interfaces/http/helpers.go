package http

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// EnumPtr 是一个辅助函数，用于将 *string 转换为枚举类型的 *T。
// 适用于处理更新请求中的可选枚举字段。
func EnumPtr[T any](v *string, parse func(string) T) *T {
	if v == nil {
		return nil
	}
	parsed := parse(*v)
	return &parsed
}

// ServiceError 将业务错误映射为标准 API 响应。
func ServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	message := err.Error()
	switch {
	case strings.Contains(message, "not found"):
		NotFound(c, message)
	case strings.Contains(message, "conflict"),
		strings.Contains(message, "cannot"),
		strings.Contains(message, "already"):
		Conflict(c, message)
	case strings.Contains(message, "required"),
		strings.Contains(message, "must"),
		strings.Contains(message, "insufficient"),
		strings.Contains(message, "invalid"):
		ValidationError(c, message)
	default:
		BadRequest(c, message)
	}
}
