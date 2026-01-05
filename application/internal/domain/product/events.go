package product

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// ProductCreatedEvent 在创建 Product 时发布。
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

// ProductUpdatedEvent 在更新 Product 时发布。
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

// ProductDeletedEvent 在删除 Product 时发布。
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

// init 将事件注册到全局注册表。
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
