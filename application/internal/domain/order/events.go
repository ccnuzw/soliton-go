package order

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// OrderCreatedEvent 在创建 Order 时发布。
type OrderCreatedEvent struct {
	ddd.BaseDomainEvent
	OrderID string `json:"order_id"`
}

func (e OrderCreatedEvent) EventName() string {
	return "order.created"
}

func NewOrderCreatedEvent(id string) OrderCreatedEvent {
	return OrderCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		OrderID: id,
	}
}

// OrderUpdatedEvent 在更新 Order 时发布。
type OrderUpdatedEvent struct {
	ddd.BaseDomainEvent
	OrderID string `json:"order_id"`
}

func (e OrderUpdatedEvent) EventName() string {
	return "order.updated"
}

func NewOrderUpdatedEvent(id string) OrderUpdatedEvent {
	return OrderUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		OrderID: id,
	}
}

// OrderDeletedEvent 在删除 Order 时发布。
type OrderDeletedEvent struct {
	ddd.BaseDomainEvent
	OrderID string    `json:"order_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e OrderDeletedEvent) EventName() string {
	return "order.deleted"
}

func NewOrderDeletedEvent(id string) OrderDeletedEvent {
	return OrderDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		OrderID: id,
		DeletedAt: time.Now(),
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("order.created", func() ddd.DomainEvent {
		return &OrderCreatedEvent{}
	})
	event.RegisterEvent("order.updated", func() ddd.DomainEvent {
		return &OrderUpdatedEvent{}
	})
	event.RegisterEvent("order.deleted", func() ddd.DomainEvent {
		return &OrderDeletedEvent{}
	})
}
