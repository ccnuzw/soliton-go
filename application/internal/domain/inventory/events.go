package inventory

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// InventoryCreatedEvent 在创建 Inventory 时发布。
type InventoryCreatedEvent struct {
	ddd.BaseDomainEvent
	InventoryID string `json:"inventory_id"`
}

func (e InventoryCreatedEvent) EventName() string {
	return "inventory.created"
}

func NewInventoryCreatedEvent(id string) InventoryCreatedEvent {
	return InventoryCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		InventoryID: id,
	}
}

// InventoryUpdatedEvent 在更新 Inventory 时发布。
type InventoryUpdatedEvent struct {
	ddd.BaseDomainEvent
	InventoryID string `json:"inventory_id"`
}

func (e InventoryUpdatedEvent) EventName() string {
	return "inventory.updated"
}

func NewInventoryUpdatedEvent(id string) InventoryUpdatedEvent {
	return InventoryUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		InventoryID: id,
	}
}

// InventoryDeletedEvent 在删除 Inventory 时发布。
type InventoryDeletedEvent struct {
	ddd.BaseDomainEvent
	InventoryID string    `json:"inventory_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e InventoryDeletedEvent) EventName() string {
	return "inventory.deleted"
}

func NewInventoryDeletedEvent(id string) InventoryDeletedEvent {
	return InventoryDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		InventoryID: id,
		DeletedAt: time.Now(),
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("inventory.created", func() ddd.DomainEvent {
		return &InventoryCreatedEvent{}
	})
	event.RegisterEvent("inventory.updated", func() ddd.DomainEvent {
		return &InventoryUpdatedEvent{}
	})
	event.RegisterEvent("inventory.deleted", func() ddd.DomainEvent {
		return &InventoryDeletedEvent{}
	})
}
