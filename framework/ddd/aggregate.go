package ddd

// ID is the interface for entity identifiers.
type ID interface {
	String() string
}

// Entity is the marker interface for entities.
type Entity interface {
	GetID() ID
}

// AggregateRoot is the marker interface for aggregate roots.
type AggregateRoot interface {
	Entity
	AddDomainEvent(event DomainEvent)
	PullDomainEvents() []DomainEvent
}

// BaseAggregateRoot is a base struct for aggregates that handles domain events.
type BaseAggregateRoot struct {
	events []DomainEvent
}

func (b *BaseAggregateRoot) AddDomainEvent(event DomainEvent) {
	b.events = append(b.events, event)
}

func (b *BaseAggregateRoot) PullDomainEvents() []DomainEvent {
	events := b.events
	b.events = nil
	return events
}
