package shippingapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/shipping"
)

// CreateShippingCommand 是创建 Shipping 的命令。
type CreateShippingCommand struct {
	ID string
	OrderId string
	Carrier string
	ShippingMethod shipping.ShippingShippingMethod
	TrackingNumber string
	Status shipping.ShippingStatus
	ShippedAt *time.Time
	DeliveredAt *time.Time
	ReceiverName string
	ReceiverPhone string
	ReceiverAddress string
	ReceiverCity string
	ReceiverState string
	ReceiverCountry string
	ReceiverPostalCode string
	Notes string
}

// CreateShippingHandler 处理 CreateShippingCommand。
type CreateShippingHandler struct {
	repo shipping.ShippingRepository
	service *shipping.ShippingDomainService
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreateShippingHandler(repo shipping.ShippingRepository, service *shipping.ShippingDomainService) *CreateShippingHandler {
	return &CreateShippingHandler{repo: repo, service: service}
}

func (h *CreateShippingHandler) Handle(ctx context.Context, cmd CreateShippingCommand) (*shipping.Shipping, error) {
	entity := shipping.NewShipping(cmd.ID, cmd.OrderId, cmd.Carrier, cmd.ShippingMethod, cmd.TrackingNumber, cmd.Status, cmd.ShippedAt, cmd.DeliveredAt, cmd.ReceiverName, cmd.ReceiverPhone, cmd.ReceiverAddress, cmd.ReceiverCity, cmd.ReceiverState, cmd.ReceiverCountry, cmd.ReceiverPostalCode, cmd.Notes)
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

// UpdateShippingCommand 是更新 Shipping 的命令。
type UpdateShippingCommand struct {
	ID string
	OrderId *string
	Carrier *string
	ShippingMethod *shipping.ShippingShippingMethod
	TrackingNumber *string
	Status *shipping.ShippingStatus
	ShippedAt *time.Time
	DeliveredAt *time.Time
	ReceiverName *string
	ReceiverPhone *string
	ReceiverAddress *string
	ReceiverCity *string
	ReceiverState *string
	ReceiverCountry *string
	ReceiverPostalCode *string
	Notes *string
}

// UpdateShippingHandler 处理 UpdateShippingCommand。
type UpdateShippingHandler struct {
	repo shipping.ShippingRepository
	service *shipping.ShippingDomainService
}

func NewUpdateShippingHandler(repo shipping.ShippingRepository, service *shipping.ShippingDomainService) *UpdateShippingHandler {
	return &UpdateShippingHandler{repo: repo, service: service}
}

func (h *UpdateShippingHandler) Handle(ctx context.Context, cmd UpdateShippingCommand) (*shipping.Shipping, error) {
	entity, err := h.repo.Find(ctx, shipping.ShippingID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.OrderId, cmd.Carrier, cmd.ShippingMethod, cmd.TrackingNumber, cmd.Status, cmd.ShippedAt, cmd.DeliveredAt, cmd.ReceiverName, cmd.ReceiverPhone, cmd.ReceiverAddress, cmd.ReceiverCity, cmd.ReceiverState, cmd.ReceiverCountry, cmd.ReceiverPostalCode, cmd.Notes)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteShippingCommand 是删除 Shipping 的命令。
type DeleteShippingCommand struct {
	ID string
}

// DeleteShippingHandler 处理 DeleteShippingCommand。
type DeleteShippingHandler struct {
	repo shipping.ShippingRepository
	service *shipping.ShippingDomainService
}

func NewDeleteShippingHandler(repo shipping.ShippingRepository, service *shipping.ShippingDomainService) *DeleteShippingHandler {
	return &DeleteShippingHandler{repo: repo, service: service}
}

func (h *DeleteShippingHandler) Handle(ctx context.Context, cmd DeleteShippingCommand) error {
	return h.repo.Delete(ctx, shipping.ShippingID(cmd.ID))
}
