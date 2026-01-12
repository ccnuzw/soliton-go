package shipping

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// ShippingCreatedEvent 在创建 Shipping 时发布。
type ShippingCreatedEvent struct {
	ddd.BaseDomainEvent
	ShippingID string `json:"shipping_id"`
}

func (e ShippingCreatedEvent) EventName() string {
	return "shipping.created"
}

func NewShippingCreatedEvent(id string) ShippingCreatedEvent {
	return ShippingCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ShippingID: id,
	}
}

// ShippingUpdatedEvent 在更新 Shipping 时发布。
type ShippingUpdatedEvent struct {
	ddd.BaseDomainEvent
	ShippingID string `json:"shipping_id"`
}

func (e ShippingUpdatedEvent) EventName() string {
	return "shipping.updated"
}

func NewShippingUpdatedEvent(id string) ShippingUpdatedEvent {
	return ShippingUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ShippingID: id,
	}
}

// ShippingDeletedEvent 在删除 Shipping 时发布。
type ShippingDeletedEvent struct {
	ddd.BaseDomainEvent
	ShippingID string    `json:"shipping_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e ShippingDeletedEvent) EventName() string {
	return "shipping.deleted"
}

func NewShippingDeletedEvent(id string) ShippingDeletedEvent {
	return ShippingDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ShippingID: id,
		DeletedAt: time.Now(),
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("shipping.created", func() ddd.DomainEvent {
		return &ShippingCreatedEvent{}
	})
	event.RegisterEvent("shipping.updated", func() ddd.DomainEvent {
		return &ShippingUpdatedEvent{}
	})
	event.RegisterEvent("shipping.deleted", func() ddd.DomainEvent {
		return &ShippingDeletedEvent{}
	})
}
