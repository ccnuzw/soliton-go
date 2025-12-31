package ddd

import "time"

// DomainEvent is the interface that all domain events should implement.
type DomainEvent interface {
	EventName() string
	OccurredOn() time.Time
}

// BaseDomainEvent is a struct that can be embedded in domain events to provide common behavior.
type BaseDomainEvent struct {
	occurredOn time.Time
}

func NewBaseDomainEvent() BaseDomainEvent {
	return BaseDomainEvent{occurredOn: time.Now()}
}

func (e BaseDomainEvent) OccurredOn() time.Time {
	return e.occurredOn
}
