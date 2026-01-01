package userapp

import "github.com/soliton-go/application/internal/domain/user"

// CreateUserRequest is the request body for creating a User.
type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateUserRequest is the request body for updating a User.
type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// UserResponse is the response body for User data.
type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToUserResponse converts entity to response.
func ToUserResponse(e *user.User) UserResponse {
	return UserResponse{
		ID:   string(e.ID),
		Name: e.Name,
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
