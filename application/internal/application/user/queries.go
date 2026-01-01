package userapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/user"
)

// GetUserQuery is the query for getting a single user.
type GetUserQuery struct {
	ID string
}

// GetUserHandler handles GetUserQuery.
type GetUserHandler struct {
	repo user.UserRepository
}

// NewGetUserHandler creates a new GetUserHandler.
func NewGetUserHandler(repo user.UserRepository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

// Handle processes the GetUserQuery.
func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*user.User, error) {
	return h.repo.Find(ctx, user.UserID(query.ID))
}

// ListUsersQuery is the query for listing all users.
type ListUsersQuery struct{}

// ListUsersHandler handles ListUsersQuery.
type ListUsersHandler struct {
	repo user.UserRepository
}

// NewListUsersHandler creates a new ListUsersHandler.
func NewListUsersHandler(repo user.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{repo: repo}
}

// Handle processes the ListUsersQuery.
func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) ([]*user.User, error) {
	return h.repo.FindAll(ctx)
}
