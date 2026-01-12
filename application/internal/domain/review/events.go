package review

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// ReviewCreatedEvent 在创建 Review 时发布。
type ReviewCreatedEvent struct {
	ddd.BaseDomainEvent
	ReviewID string `json:"review_id"`
}

func (e ReviewCreatedEvent) EventName() string {
	return "review.created"
}

func NewReviewCreatedEvent(id string) ReviewCreatedEvent {
	return ReviewCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ReviewID: id,
	}
}

// ReviewUpdatedEvent 在更新 Review 时发布。
type ReviewUpdatedEvent struct {
	ddd.BaseDomainEvent
	ReviewID string `json:"review_id"`
}

func (e ReviewUpdatedEvent) EventName() string {
	return "review.updated"
}

func NewReviewUpdatedEvent(id string) ReviewUpdatedEvent {
	return ReviewUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ReviewID: id,
	}
}

// ReviewDeletedEvent 在删除 Review 时发布。
type ReviewDeletedEvent struct {
	ddd.BaseDomainEvent
	ReviewID string    `json:"review_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e ReviewDeletedEvent) EventName() string {
	return "review.deleted"
}

func NewReviewDeletedEvent(id string) ReviewDeletedEvent {
	return ReviewDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		ReviewID: id,
		DeletedAt: time.Now(),
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("review.created", func() ddd.DomainEvent {
		return &ReviewCreatedEvent{}
	})
	event.RegisterEvent("review.updated", func() ddd.DomainEvent {
		return &ReviewUpdatedEvent{}
	})
	event.RegisterEvent("review.deleted", func() ddd.DomainEvent {
		return &ReviewDeletedEvent{}
	})
}
