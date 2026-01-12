package shippingapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/shipping"
)

// GetShippingQuery 是获取单个 Shipping 的查询。
type GetShippingQuery struct {
	ID string
}

// GetShippingHandler 处理 GetShippingQuery。
type GetShippingHandler struct {
	repo shipping.ShippingRepository
}

func NewGetShippingHandler(repo shipping.ShippingRepository) *GetShippingHandler {
	return &GetShippingHandler{repo: repo}
}

func (h *GetShippingHandler) Handle(ctx context.Context, query GetShippingQuery) (*shipping.Shipping, error) {
	return h.repo.Find(ctx, shipping.ShippingID(query.ID))
}

// ListShippingsQuery 是分页列表查询。
type ListShippingsQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListShippingsResult 是分页查询结果。
type ListShippingsResult struct {
	Items      []*shipping.Shipping
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListShippingsHandler 处理 ListShippingsQuery。
type ListShippingsHandler struct {
	repo shipping.ShippingRepository
}

func NewListShippingsHandler(repo shipping.ShippingRepository) *ListShippingsHandler {
	return &ListShippingsHandler{repo: repo}
}

func (h *ListShippingsHandler) Handle(ctx context.Context, query ListShippingsQuery) (*ListShippingsResult, error) {
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
		"order_id": {},
		"carrier": {},
		"shipping_method": {},
		"tracking_number": {},
		"status": {},
		"shipped_at": {},
		"delivered_at": {},
		"receiver_name": {},
		"receiver_phone": {},
		"receiver_address": {},
		"receiver_city": {},
		"receiver_state": {},
		"receiver_country": {},
		"receiver_postal_code": {},
		"notes": {},
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

	return &ListShippingsResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
