package orderapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/order"
)

// GetOrderQuery is the query for getting a single Order.
type GetOrderQuery struct {
	ID string
}

// GetOrderHandler handles GetOrderQuery.
type GetOrderHandler struct {
	repo order.OrderRepository
}

func NewGetOrderHandler(repo order.OrderRepository) *GetOrderHandler {
	return &GetOrderHandler{repo: repo}
}

func (h *GetOrderHandler) Handle(ctx context.Context, query GetOrderQuery) (*order.Order, error) {
	return h.repo.Find(ctx, order.OrderID(query.ID))
}

// ListOrdersQuery is the query for listing Orders with pagination.
type ListOrdersQuery struct {
	Page     int // Page number (1-based)
	PageSize int // Items per page (default: 20, max: 100)
}

// ListOrdersResult is the paginated result for ListOrdersQuery.
type ListOrdersResult struct {
	Items      []*order.Order
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListOrdersHandler handles ListOrdersQuery.
type ListOrdersHandler struct {
	repo order.OrderRepository
}

func NewListOrdersHandler(repo order.OrderRepository) *ListOrdersHandler {
	return &ListOrdersHandler{repo: repo}
}

func (h *ListOrdersHandler) Handle(ctx context.Context, query ListOrdersQuery) (*ListOrdersResult, error) {
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

	return &ListOrdersResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
