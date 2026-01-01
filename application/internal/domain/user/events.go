package user

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// UserCreatedEvent is published when a new user is created.
type UserCreatedEvent struct {
	ddd.BaseDomainEvent
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// EventName returns the name of the event.
func (e UserCreatedEvent) EventName() string {
	return "user.created"
}

// NewUserCreatedEvent creates a new UserCreatedEvent.
func NewUserCreatedEvent(userID, name, email string) UserCreatedEvent {
	return UserCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		UserID:          userID,
		Name:            name,
		Email:           email,
	}
}

// UserUpdatedEvent is published when a user is updated.
type UserUpdatedEvent struct {
	ddd.BaseDomainEvent
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// EventName returns the name of the event.
func (e UserUpdatedEvent) EventName() string {
	return "user.updated"
}

// NewUserUpdatedEvent creates a new UserUpdatedEvent.
func NewUserUpdatedEvent(userID, name, email string) UserUpdatedEvent {
	return UserUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		UserID:          userID,
		Name:            name,
		Email:           email,
	}
}

// UserDeletedEvent is published when a user is deleted.
type UserDeletedEvent struct {
	ddd.BaseDomainEvent
	UserID    string    `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}

// EventName returns the name of the event.
func (e UserDeletedEvent) EventName() string {
	return "user.deleted"
}

// NewUserDeletedEvent creates a new UserDeletedEvent.
func NewUserDeletedEvent(userID string) UserDeletedEvent {
	return UserDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		UserID:          userID,
		DeletedAt:       time.Now(),
	}
}

// init registers user events with the global event registry.
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
