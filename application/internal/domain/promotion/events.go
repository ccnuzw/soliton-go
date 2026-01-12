package promotion

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// PromotionCreatedEvent 在创建 Promotion 时发布。
type PromotionCreatedEvent struct {
	ddd.BaseDomainEvent
	PromotionID string `json:"promotion_id"`
}

func (e PromotionCreatedEvent) EventName() string {
	return "promotion.created"
}

func NewPromotionCreatedEvent(id string) PromotionCreatedEvent {
	return PromotionCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		PromotionID: id,
	}
}

// PromotionUpdatedEvent 在更新 Promotion 时发布。
type PromotionUpdatedEvent struct {
	ddd.BaseDomainEvent
	PromotionID string `json:"promotion_id"`
}

func (e PromotionUpdatedEvent) EventName() string {
	return "promotion.updated"
}

func NewPromotionUpdatedEvent(id string) PromotionUpdatedEvent {
	return PromotionUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		PromotionID: id,
	}
}

// PromotionDeletedEvent 在删除 Promotion 时发布。
type PromotionDeletedEvent struct {
	ddd.BaseDomainEvent
	PromotionID string    `json:"promotion_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e PromotionDeletedEvent) EventName() string {
	return "promotion.deleted"
}

func NewPromotionDeletedEvent(id string) PromotionDeletedEvent {
	return PromotionDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		PromotionID: id,
		DeletedAt: time.Now(),
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("promotion.created", func() ddd.DomainEvent {
		return &PromotionCreatedEvent{}
	})
	event.RegisterEvent("promotion.updated", func() ddd.DomainEvent {
		return &PromotionUpdatedEvent{}
	})
	event.RegisterEvent("promotion.deleted", func() ddd.DomainEvent {
		return &PromotionDeletedEvent{}
	})
}
