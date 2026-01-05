package http

// EnumPtr 是一个辅助函数，用于将 *string 转换为枚举类型的 *T。
// 适用于处理更新请求中的可选枚举字段。
func EnumPtr[T any](v *string, parse func(string) T) *T {
	if v == nil {
		return nil
	}
	parsed := parse(*v)
	return &parsed
}
