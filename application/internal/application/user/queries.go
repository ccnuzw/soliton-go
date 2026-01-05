package userapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/user"
)

// GetUserQuery is the query for getting a single User.
type GetUserQuery struct {
	ID string
}

// GetUserHandler handles GetUserQuery.
type GetUserHandler struct {
	repo user.UserRepository
}

func NewGetUserHandler(repo user.UserRepository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*user.User, error) {
	return h.repo.Find(ctx, user.UserID(query.ID))
}

// ListUsersQuery is the query for listing Users with pagination.
type ListUsersQuery struct {
	Page     int // Page number (1-based)
	PageSize int // Items per page (default: 20, max: 100)
}

// ListUsersResult is the paginated result for ListUsersQuery.
type ListUsersResult struct {
	Items      []*user.User
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListUsersHandler handles ListUsersQuery.
type ListUsersHandler struct {
	repo user.UserRepository
}

func NewListUsersHandler(repo user.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{repo: repo}
}

func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) (*ListUsersResult, error) {
	// Normalize pagination parameters
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get total count and items
	items, total, err := h.repo.FindPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &ListUsersResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
