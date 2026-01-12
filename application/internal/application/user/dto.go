package userapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/user"
)

// CreateUserRequest 是创建 User 的请求体。
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// UpdateUserRequest 是更新 User 的请求体。
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email *string `json:"email,omitempty"`
}

// UserResponse 是 User 的响应体。
type UserResponse struct {
	ID        string    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse 将实体转换为响应体。
func ToUserResponse(e *user.User) UserResponse {
	return UserResponse{
		ID:        string(e.ID),
		Username: e.Username,
		Email: e.Email,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToUserResponseList 将实体列表转换为响应体列表。
func ToUserResponseList(entities []*user.User) []UserResponse {
	result := make([]UserResponse, len(entities))
	for i, e := range entities {
		result[i] = ToUserResponse(e)
	}
	return result
}
