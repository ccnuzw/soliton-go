package reviewapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/review"
)

// GetReviewQuery 是获取单个 Review 的查询。
type GetReviewQuery struct {
	ID string
}

// GetReviewHandler 处理 GetReviewQuery。
type GetReviewHandler struct {
	repo review.ReviewRepository
}

func NewGetReviewHandler(repo review.ReviewRepository) *GetReviewHandler {
	return &GetReviewHandler{repo: repo}
}

func (h *GetReviewHandler) Handle(ctx context.Context, query GetReviewQuery) (*review.Review, error) {
	return h.repo.Find(ctx, review.ReviewID(query.ID))
}

// ListReviewsQuery 是分页列表查询。
type ListReviewsQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListReviewsResult 是分页查询结果。
type ListReviewsResult struct {
	Items      []*review.Review
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListReviewsHandler 处理 ListReviewsQuery。
type ListReviewsHandler struct {
	repo review.ReviewRepository
}

func NewListReviewsHandler(repo review.ReviewRepository) *ListReviewsHandler {
	return &ListReviewsHandler{repo: repo}
}

func (h *ListReviewsHandler) Handle(ctx context.Context, query ListReviewsQuery) (*ListReviewsResult, error) {
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
		"product_id": {},
		"user_id": {},
		"order_id": {},
		"rating": {},
		"title": {},
		"content": {},
		"status": {},
		"is_anonymous": {},
		"helpful_count": {},
		"reply": {},
		"images": {},
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

	return &ListReviewsResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
