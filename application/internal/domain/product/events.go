package product

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// ProductCreatedEvent is published when a new Product is created.
type ProductCreatedEvent struct {
	ddd.BaseDomainEvent
	ProductID string `json:"product_id"`
}

func (e ProductCreatedEvent) EventName() string {
	return "product.created"
}

func NewProductCreatedEvent(id string) ProductCreatedEvent {
	return ProductCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ProductID: id,
	}
}

// ProductUpdatedEvent is published when a Product is updated.
type ProductUpdatedEvent struct {
	ddd.BaseDomainEvent
	ProductID string `json:"product_id"`
}

func (e ProductUpdatedEvent) EventName() string {
	return "product.updated"
}

func NewProductUpdatedEvent(id string) ProductUpdatedEvent {
	return ProductUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ProductID: id,
	}
}

// ProductDeletedEvent is published when a Product is deleted.
type ProductDeletedEvent struct {
	ddd.BaseDomainEvent
	ProductID string    `json:"product_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e ProductDeletedEvent) EventName() string {
	return "product.deleted"
}

func NewProductDeletedEvent(id string) ProductDeletedEvent {
	return ProductDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ProductID: id,
		DeletedAt: time.Now(),
	}
}

// init registers events with the global registry.
func init() {
	event.RegisterEvent("product.created", func() ddd.DomainEvent {
		return &ProductCreatedEvent{}
	})
	event.RegisterEvent("product.updated", func() ddd.DomainEvent {
		return &ProductUpdatedEvent{}
	})
	event.RegisterEvent("product.deleted", func() ddd.DomainEvent {
		return &ProductDeletedEvent{}
	})
}
