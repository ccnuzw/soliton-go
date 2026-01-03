package userapp

import (
	"time"

	"github.com/soliton-go/test-project/internal/domain/user"
)

// CreateUserRequest is the request body for creating a User.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Status string `json:"status"`
}

// UpdateUserRequest is the request body for updating a User.
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email *string `json:"email,omitempty"`
	Status *string `json:"status,omitempty"`
}

// UserResponse is the response body for User data.
type UserResponse struct {
	ID        string    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
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
