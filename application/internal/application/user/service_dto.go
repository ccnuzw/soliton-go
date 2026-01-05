package userapp


// CreateUserServiceRequest 是 CreateUser 方法的请求参数。
type CreateUserServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// CreateUserServiceResponse 是 CreateUser 方法的响应结果。
type CreateUserServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// GetUserServiceRequest 是 GetUser 方法的请求参数。
type GetUserServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// GetUserServiceResponse 是 GetUser 方法的响应结果。
type GetUserServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

// ListUsersServiceRequest 是 ListUsers 方法的请求参数。
type ListUsersServiceRequest struct {
	// 在此添加请求字段：
	ID string `json:"id,omitempty"` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    `json:"data,omitempty"` // 请求数据（用于 Create/Update 操作）
}

// ListUsersServiceResponse 是 ListUsers 方法的响应结果。
type ListUsersServiceResponse struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 提示消息
	Data    any    `json:"data,omitempty"`    // 响应数据
}

