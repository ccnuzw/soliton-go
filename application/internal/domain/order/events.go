package order

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// OrderCreatedEvent is published when a new Order is created.
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

// OrderUpdatedEvent is published when a Order is updated.
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

// OrderDeletedEvent is published when a Order is deleted.
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

// init registers events with the global registry.
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
