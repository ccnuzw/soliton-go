package paymentapp


// CreatePaymentServiceRequest 是 CreatePayment 方法的请求参数。
type CreatePaymentServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// CreatePaymentServiceResponse 是 CreatePayment 方法的响应结果。
type CreatePaymentServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// GetPaymentServiceRequest 是 GetPayment 方法的请求参数。
type GetPaymentServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// GetPaymentServiceResponse 是 GetPayment 方法的响应结果。
type GetPaymentServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// ListPaymentsServiceRequest 是 ListPayments 方法的请求参数。
type ListPaymentsServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// ListPaymentsServiceResponse 是 ListPayments 方法的响应结果。
type ListPaymentsServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

