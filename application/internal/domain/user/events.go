package user

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// UserCreatedEvent is published when a new User is created.
type UserCreatedEvent struct {
	ddd.BaseDomainEvent
	UserID string `json:"user_id"`
}

func (e UserCreatedEvent) EventName() string {
	return "user.created"
}

func NewUserCreatedEvent(id string) UserCreatedEvent {
	return UserCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		UserID: id,
	}
}

// UserUpdatedEvent is published when a User is updated.
type UserUpdatedEvent struct {
	ddd.BaseDomainEvent
	UserID string `json:"user_id"`
}

func (e UserUpdatedEvent) EventName() string {
	return "user.updated"
}

func NewUserUpdatedEvent(id string) UserUpdatedEvent {
	return UserUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		UserID: id,
	}
}

// UserDeletedEvent is published when a User is deleted.
type UserDeletedEvent struct {
	ddd.BaseDomainEvent
	UserID string    `json:"user_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e UserDeletedEvent) EventName() string {
	return "user.deleted"
}

func NewUserDeletedEvent(id string) UserDeletedEvent {
	return UserDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		UserID: id,
		DeletedAt: time.Now(),
	}
}

// init registers events with the global registry.
func init() {
	event.RegisterEvent("user.created", func() ddd.DomainEvent {
		return &UserCreatedEvent{}
	})
	event.RegisterEvent("user.updated", func() ddd.DomainEvent {
		return &UserUpdatedEvent{}
	})
	event.RegisterEvent("user.deleted", func() ddd.DomainEvent {
		return &UserDeletedEvent{}
	})
}
