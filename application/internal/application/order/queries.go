package orderapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/order"
)

// GetOrderQuery 是获取单个 Order 的查询。
type GetOrderQuery struct {
	ID string
}

// GetOrderHandler 处理 GetOrderQuery。
type GetOrderHandler struct {
	repo order.OrderRepository
}

func NewGetOrderHandler(repo order.OrderRepository) *GetOrderHandler {
	return &GetOrderHandler{repo: repo}
}

func (h *GetOrderHandler) Handle(ctx context.Context, query GetOrderQuery) (*order.Order, error) {
	return h.repo.Find(ctx, order.OrderID(query.ID))
}

// ListOrdersQuery 是分页列表查询。
type ListOrdersQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListOrdersResult 是分页查询结果。
type ListOrdersResult struct {
	Items      []*order.Order
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListOrdersHandler 处理 ListOrdersQuery。
type ListOrdersHandler struct {
	repo order.OrderRepository
}

func NewListOrdersHandler(repo order.OrderRepository) *ListOrdersHandler {
	return &ListOrdersHandler{repo: repo}
}

func (h *ListOrdersHandler) Handle(ctx context.Context, query ListOrdersQuery) (*ListOrdersResult, error) {
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
		"user_id": {},
		"order_no": {},
		"total_amount": {},
		"discount_amount": {},
		"tax_amount": {},
		"shipping_fee": {},
		"final_amount": {},
		"currency": {},
		"payment_method": {},
		"payment_status": {},
		"order_status": {},
		"shipping_method": {},
		"tracking_number": {},
		"receiver_name": {},
		"receiver_phone": {},
		"receiver_email": {},
		"receiver_address": {},
		"receiver_city": {},
		"receiver_state": {},
		"receiver_country": {},
		"receiver_postal_code": {},
		"notes": {},
		"paid_at": {},
		"shipped_at": {},
		"delivered_at": {},
		"cancelled_at": {},
		"refund_amount": {},
		"refund_reason": {},
		"item_count": {},
		"weight": {},
		"is_gift": {},
		"gift_message": {},
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

	return &ListOrdersResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
