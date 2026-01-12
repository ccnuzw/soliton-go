package promotionapp


// ApplyPromotionServiceRequest 是 ApplyPromotion 方法的请求参数。
type ApplyPromotionServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// ApplyPromotionServiceResponse 是 ApplyPromotion 方法的响应结果。
type ApplyPromotionServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// ValidatePromotionServiceRequest 是 ValidatePromotion 方法的请求参数。
type ValidatePromotionServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// ValidatePromotionServiceResponse 是 ValidatePromotion 方法的响应结果。
type ValidatePromotionServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// RevokePromotionServiceRequest 是 RevokePromotion 方法的请求参数。
type RevokePromotionServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// RevokePromotionServiceResponse 是 RevokePromotion 方法的响应结果。
type RevokePromotionServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// EvaluatePromotionServiceRequest 是 EvaluatePromotion 方法的请求参数。
type EvaluatePromotionServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// EvaluatePromotionServiceResponse 是 EvaluatePromotion 方法的响应结果。
type EvaluatePromotionServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// FindByCodeServiceRequest 是 FindByCode 方法的请求参数。
type FindByCodeServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// FindByCodeServiceResponse 是 FindByCode 方法的响应结果。
type FindByCodeServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

