package orderapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/order"
)

// CreateOrderCommand is the command for creating a Order.
type CreateOrderCommand struct {
	ID string
	UserId string
	OrderNo string
	TotalAmount int64
	Status order.OrderStatus
	ReceiverName string
	ReceiverPhone string
	ReceiverAddress string
}

// CreateOrderHandler handles CreateOrderCommand.
type CreateOrderHandler struct {
	repo order.OrderRepository
}

func NewCreateOrderHandler(repo order.OrderRepository) *CreateOrderHandler {
	return &CreateOrderHandler{repo: repo}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*order.Order, error) {
	entity := order.NewOrder(cmd.ID, cmd.UserId, cmd.OrderNo, cmd.TotalAmount, cmd.Status, cmd.ReceiverName, cmd.ReceiverPhone, cmd.ReceiverAddress)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// UpdateOrderCommand is the command for updating a Order.
type UpdateOrderCommand struct {
	ID string
	UserId string
	OrderNo string
	TotalAmount int64
	Status order.OrderStatus
	ReceiverName string
	ReceiverPhone string
	ReceiverAddress string
}

// UpdateOrderHandler handles UpdateOrderCommand.
type UpdateOrderHandler struct {
	repo order.OrderRepository
}

func NewUpdateOrderHandler(repo order.OrderRepository) *UpdateOrderHandler {
	return &UpdateOrderHandler{repo: repo}
}

func (h *UpdateOrderHandler) Handle(ctx context.Context, cmd UpdateOrderCommand) (*order.Order, error) {
	entity, err := h.repo.Find(ctx, order.OrderID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.UserId, cmd.OrderNo, cmd.TotalAmount, cmd.Status, cmd.ReceiverName, cmd.ReceiverPhone, cmd.ReceiverAddress)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteOrderCommand is the command for deleting a Order.
type DeleteOrderCommand struct {
	ID string
}

// DeleteOrderHandler handles DeleteOrderCommand.
type DeleteOrderHandler struct {
	repo order.OrderRepository
}

func NewDeleteOrderHandler(repo order.OrderRepository) *DeleteOrderHandler {
	return &DeleteOrderHandler{repo: repo}
}

func (h *DeleteOrderHandler) Handle(ctx context.Context, cmd DeleteOrderCommand) error {
	return h.repo.Delete(ctx, order.OrderID(cmd.ID))
}
