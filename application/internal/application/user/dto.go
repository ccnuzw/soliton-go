package userapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/user"
)

// CreateUserRequest is the request body for creating a User.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Phone string `json:"phone"`
	Nickname string `json:"nickname"`
	Role string `json:"role"`
	Status string `json:"status"`
}

// UpdateUserRequest is the request body for updating a User.
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Phone string `json:"phone"`
	Nickname string `json:"nickname"`
	Role string `json:"role"`
	Status string `json:"status"`
}

// UserResponse is the response body for User data.
type UserResponse struct {
	ID        string    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Phone string `json:"phone"`
	Nickname string `json:"nickname"`
	Role string `json:"role"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts entity to response.
func ToUserResponse(e *user.User) UserResponse {
	return UserResponse{
		ID:        string(e.ID),
		Username: e.Username,
		Email: e.Email,
		PasswordHash: e.PasswordHash,
		Phone: e.Phone,
		Nickname: e.Nickname,
		Role: string(e.Role),
		Status: string(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToUserResponseList converts entities to response list.
func ToUserResponseList(entities []*user.User) []UserResponse {
	result := make([]UserResponse, len(entities))
	for i, e := range entities {
		result[i] = ToUserResponse(e)
	}
	return result
}
