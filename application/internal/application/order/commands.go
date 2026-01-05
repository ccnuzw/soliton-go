package orderapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/order"
)

// CreateOrderCommand 是创建 Order 的命令。
type CreateOrderCommand struct {
	ID string
	Userid string
	Orderno string
	Totalamount int64
	Discountamount int64
	Taxamount int64
	Shippingfee int64
	Finalamount int64
	Currency string
	Paymentmethod order.OrderPaymentmethod
	Paymentstatus order.OrderPaymentstatus
	Orderstatus order.OrderOrderstatus
	Shippingmethod order.OrderShippingmethod
	Trackingnumber string
	Receivername string
	Receiverphone string
	Receiveremail string
	Receiveraddress string
	Receivercity string
	Receiverstate string
	Receivercountry string
	Receiverpostalcode string
	Notes string
	Paidat time.Time
	Shippedat time.Time
	Deliveredat time.Time
	Cancelledat time.Time
	Refundamount int64
	Refundreason string
	Itemcount int
	Weight float64
	Isgift bool
	Giftmessage string
}

// CreateOrderHandler 处理 CreateOrderCommand。
type CreateOrderHandler struct {
	repo order.OrderRepository
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreateOrderHandler(repo order.OrderRepository) *CreateOrderHandler {
	return &CreateOrderHandler{repo: repo}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*order.Order, error) {
	entity := order.NewOrder(cmd.ID, cmd.Userid, cmd.Orderno, cmd.Totalamount, cmd.Discountamount, cmd.Taxamount, cmd.Shippingfee, cmd.Finalamount, cmd.Currency, cmd.Paymentmethod, cmd.Paymentstatus, cmd.Orderstatus, cmd.Shippingmethod, cmd.Trackingnumber, cmd.Receivername, cmd.Receiverphone, cmd.Receiveremail, cmd.Receiveraddress, cmd.Receivercity, cmd.Receiverstate, cmd.Receivercountry, cmd.Receiverpostalcode, cmd.Notes, cmd.Paidat, cmd.Shippedat, cmd.Deliveredat, cmd.Cancelledat, cmd.Refundamount, cmd.Refundreason, cmd.Itemcount, cmd.Weight, cmd.Isgift, cmd.Giftmessage)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// 可选：发布领域事件
	// 取消注释以启用事件发布：
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateOrderCommand 是更新 Order 的命令。
type UpdateOrderCommand struct {
	ID string
	Userid *string
	Orderno *string
	Totalamount *int64
	Discountamount *int64
	Taxamount *int64
	Shippingfee *int64
	Finalamount *int64
	Currency *string
	Paymentmethod *order.OrderPaymentmethod
	Paymentstatus *order.OrderPaymentstatus
	Orderstatus *order.OrderOrderstatus
	Shippingmethod *order.OrderShippingmethod
	Trackingnumber *string
	Receivername *string
	Receiverphone *string
	Receiveremail *string
	Receiveraddress *string
	Receivercity *string
	Receiverstate *string
	Receivercountry *string
	Receiverpostalcode *string
	Notes *string
	Paidat *time.Time
	Shippedat *time.Time
	Deliveredat *time.Time
	Cancelledat *time.Time
	Refundamount *int64
	Refundreason *string
	Itemcount *int
	Weight *float64
	Isgift *bool
	Giftmessage *string
}

// UpdateOrderHandler 处理 UpdateOrderCommand。
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
	entity.Update(cmd.Userid, cmd.Orderno, cmd.Totalamount, cmd.Discountamount, cmd.Taxamount, cmd.Shippingfee, cmd.Finalamount, cmd.Currency, cmd.Paymentmethod, cmd.Paymentstatus, cmd.Orderstatus, cmd.Shippingmethod, cmd.Trackingnumber, cmd.Receivername, cmd.Receiverphone, cmd.Receiveremail, cmd.Receiveraddress, cmd.Receivercity, cmd.Receiverstate, cmd.Receivercountry, cmd.Receiverpostalcode, cmd.Notes, cmd.Paidat, cmd.Shippedat, cmd.Deliveredat, cmd.Cancelledat, cmd.Refundamount, cmd.Refundreason, cmd.Itemcount, cmd.Weight, cmd.Isgift, cmd.Giftmessage)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteOrderCommand 是删除 Order 的命令。
type DeleteOrderCommand struct {
	ID string
}

// DeleteOrderHandler 处理 DeleteOrderCommand。
type DeleteOrderHandler struct {
	repo order.OrderRepository
}

func NewDeleteOrderHandler(repo order.OrderRepository) *DeleteOrderHandler {
	return &DeleteOrderHandler{repo: repo}
}

func (h *DeleteOrderHandler) Handle(ctx context.Context, cmd DeleteOrderCommand) error {
	return h.repo.Delete(ctx, order.OrderID(cmd.ID))
}
