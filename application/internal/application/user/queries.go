package userapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/user"
)

// GetUserQuery 是获取单个 User 的查询。
type GetUserQuery struct {
	ID string
}

// GetUserHandler 处理 GetUserQuery。
type GetUserHandler struct {
	repo user.UserRepository
}

func NewGetUserHandler(repo user.UserRepository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*user.User, error) {
	return h.repo.Find(ctx, user.UserID(query.ID))
}

// ListUsersQuery 是分页列表查询。
type ListUsersQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListUsersResult 是分页查询结果。
type ListUsersResult struct {
	Items      []*user.User
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListUsersHandler 处理 ListUsersQuery。
type ListUsersHandler struct {
	repo user.UserRepository
}

func NewListUsersHandler(repo user.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{repo: repo}
}

func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) (*ListUsersResult, error) {
	// 规范化分页参数
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

	// 排序字段白名单
	sortBy := strings.ToLower(strings.TrimSpace(query.SortBy))
	sortOrder := strings.ToLower(strings.TrimSpace(query.SortOrder))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	allowedSorts := map[string]struct{}{
		"id":         {},
		"created_at": {},
		"updated_at": {},
		"username": {},
		"email": {},
	}
	if _, ok := allowedSorts[sortBy]; !ok {
		sortBy = "id"
	}

	// 获取总数和分页数据
	items, total, err := h.repo.FindPaginated(ctx, page, pageSize, sortBy, sortOrder)
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
