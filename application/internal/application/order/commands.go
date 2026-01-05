package orderapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/order"
)

// CreateOrderCommand is the command for creating a Order.
type CreateOrderCommand struct {
	ID string
	UserId string
	OrderNo string
	TotalAmount int64
	DiscountAmount int64
	TaxAmount int64
	ShippingFee int64
	FinalAmount int64
	Currency string
	PaymentMethod order.OrderPaymentMethod
	PaymentStatus order.OrderPaymentStatus
	OrderStatus order.OrderOrderStatus
	ShippingMethod order.OrderShippingMethod
	TrackingNumber string
	ReceiverName string
	ReceiverPhone string
	ReceiverEmail string
	ReceiverAddress string
	ReceiverCity string
	ReceiverState string
	ReceiverCountry string
	ReceiverPostalCode string
	Notes string
	PaidAt time.Time
	ShippedAt time.Time
	DeliveredAt time.Time
	CancelledAt time.Time
	RefundAmount int64
	RefundReason string
	ItemCount int
	Weight float64
	IsGift bool
	GiftMessage string
}

// CreateOrderHandler handles CreateOrderCommand.
type CreateOrderHandler struct {
	repo order.OrderRepository
	// Optional: Add event bus for domain event publishing
	// eventBus event.EventBus
}

func NewCreateOrderHandler(repo order.OrderRepository) *CreateOrderHandler {
	return &CreateOrderHandler{repo: repo}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*order.Order, error) {
	entity := order.NewOrder(cmd.ID, cmd.UserId, cmd.OrderNo, cmd.TotalAmount, cmd.DiscountAmount, cmd.TaxAmount, cmd.ShippingFee, cmd.FinalAmount, cmd.Currency, cmd.PaymentMethod, cmd.PaymentStatus, cmd.OrderStatus, cmd.ShippingMethod, cmd.TrackingNumber, cmd.ReceiverName, cmd.ReceiverPhone, cmd.ReceiverEmail, cmd.ReceiverAddress, cmd.ReceiverCity, cmd.ReceiverState, cmd.ReceiverCountry, cmd.ReceiverPostalCode, cmd.Notes, cmd.PaidAt, cmd.ShippedAt, cmd.DeliveredAt, cmd.CancelledAt, cmd.RefundAmount, cmd.RefundReason, cmd.ItemCount, cmd.Weight, cmd.IsGift, cmd.GiftMessage)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// Optional: Publish domain events
	// Uncomment to enable event publishing:
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateOrderCommand is the command for updating a Order.
type UpdateOrderCommand struct {
	ID string
	UserId *string
	OrderNo *string
	TotalAmount *int64
	DiscountAmount *int64
	TaxAmount *int64
	ShippingFee *int64
	FinalAmount *int64
	Currency *string
	PaymentMethod *order.OrderPaymentMethod
	PaymentStatus *order.OrderPaymentStatus
	OrderStatus *order.OrderOrderStatus
	ShippingMethod *order.OrderShippingMethod
	TrackingNumber *string
	ReceiverName *string
	ReceiverPhone *string
	ReceiverEmail *string
	ReceiverAddress *string
	ReceiverCity *string
	ReceiverState *string
	ReceiverCountry *string
	ReceiverPostalCode *string
	Notes *string
	PaidAt *time.Time
	ShippedAt *time.Time
	DeliveredAt *time.Time
	CancelledAt *time.Time
	RefundAmount *int64
	RefundReason *string
	ItemCount *int
	Weight *float64
	IsGift *bool
	GiftMessage *string
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
	entity.Update(cmd.UserId, cmd.OrderNo, cmd.TotalAmount, cmd.DiscountAmount, cmd.TaxAmount, cmd.ShippingFee, cmd.FinalAmount, cmd.Currency, cmd.PaymentMethod, cmd.PaymentStatus, cmd.OrderStatus, cmd.ShippingMethod, cmd.TrackingNumber, cmd.ReceiverName, cmd.ReceiverPhone, cmd.ReceiverEmail, cmd.ReceiverAddress, cmd.ReceiverCity, cmd.ReceiverState, cmd.ReceiverCountry, cmd.ReceiverPostalCode, cmd.Notes, cmd.PaidAt, cmd.ShippedAt, cmd.DeliveredAt, cmd.CancelledAt, cmd.RefundAmount, cmd.RefundReason, cmd.ItemCount, cmd.Weight, cmd.IsGift, cmd.GiftMessage)
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
