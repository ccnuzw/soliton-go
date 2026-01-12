package promotionapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/promotion"
)

// GetPromotionQuery 是获取单个 Promotion 的查询。
type GetPromotionQuery struct {
	ID string
}

// GetPromotionHandler 处理 GetPromotionQuery。
type GetPromotionHandler struct {
	repo promotion.PromotionRepository
}

func NewGetPromotionHandler(repo promotion.PromotionRepository) *GetPromotionHandler {
	return &GetPromotionHandler{repo: repo}
}

func (h *GetPromotionHandler) Handle(ctx context.Context, query GetPromotionQuery) (*promotion.Promotion, error) {
	return h.repo.Find(ctx, promotion.PromotionID(query.ID))
}

// ListPromotionsQuery 是分页列表查询。
type ListPromotionsQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListPromotionsResult 是分页查询结果。
type ListPromotionsResult struct {
	Items      []*promotion.Promotion
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListPromotionsHandler 处理 ListPromotionsQuery。
type ListPromotionsHandler struct {
	repo promotion.PromotionRepository
}

func NewListPromotionsHandler(repo promotion.PromotionRepository) *ListPromotionsHandler {
	return &ListPromotionsHandler{repo: repo}
}

func (h *ListPromotionsHandler) Handle(ctx context.Context, query ListPromotionsQuery) (*ListPromotionsResult, error) {
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
		"code": {},
		"name": {},
		"description": {},
		"discount_type": {},
		"discount_value": {},
		"currency": {},
		"min_order_amount": {},
		"max_discount_amount": {},
		"usage_limit": {},
		"used_count": {},
		"per_user_limit": {},
		"starts_at": {},
		"ends_at": {},
		"status": {},
		"metadata": {},
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

	return &ListPromotionsResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
