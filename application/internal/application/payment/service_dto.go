package paymentapp


// ProcessPaymentServiceRequest 是 ProcessPayment 方法的请求参数。
type ProcessPaymentServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// ProcessPaymentServiceResponse 是 ProcessPayment 方法的响应结果。
type ProcessPaymentServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// RefundPaymentServiceRequest 是 RefundPayment 方法的请求参数。
type RefundPaymentServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// RefundPaymentServiceResponse 是 RefundPayment 方法的响应结果。
type RefundPaymentServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

