package reviewapp


// CreateReviewServiceRequest 是 CreateReview 方法的请求参数。
type CreateReviewServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// CreateReviewServiceResponse 是 CreateReview 方法的响应结果。
type CreateReviewServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// ModerateReviewServiceRequest 是 ModerateReview 方法的请求参数。
type ModerateReviewServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// ModerateReviewServiceResponse 是 ModerateReview 方法的响应结果。
type ModerateReviewServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// ReplyReviewServiceRequest 是 ReplyReview 方法的请求参数。
type ReplyReviewServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// ReplyReviewServiceResponse 是 ReplyReview 方法的响应结果。
type ReplyReviewServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

