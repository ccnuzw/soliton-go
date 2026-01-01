package userapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/user"
)

// CreateUserRequest is the request body for creating a User.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
}

// UpdateUserRequest is the request body for updating a User.
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
}

// UserResponse is the response body for User data.
type UserResponse struct {
	ID        string    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts entity to response.
func ToUserResponse(e *user.User) UserResponse {
	return UserResponse{
		ID:        string(e.ID),
		Username: e.Username,
		Email: e.Email,
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
