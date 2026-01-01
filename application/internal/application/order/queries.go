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

// ListOrdersQuery is the query for listing all Orders.
type ListOrdersQuery struct{}

// ListOrdersHandler handles ListOrdersQuery.
type ListOrdersHandler struct {
	repo order.OrderRepository
}

func NewListOrdersHandler(repo order.OrderRepository) *ListOrdersHandler {
	return &ListOrdersHandler{repo: repo}
}

func (h *ListOrdersHandler) Handle(ctx context.Context, query ListOrdersQuery) ([]*order.Order, error) {
	return h.repo.FindAll(ctx)
}
