package event

import (
	"fmt"
	"sync"

	"github.com/soliton-go/framework/ddd"
)

// EventFactory is a function that creates a new instance of a domain event.
type EventFactory func() ddd.DomainEvent

// EventRegistry manages the mapping between event names and their factories.
// This enables proper deserialization of events from JSON to their concrete types.
type EventRegistry interface {
	// Register associates an event name with a factory function.
	Register(eventName string, factory EventFactory)
	// Create returns a new instance of the event for the given name.
	Create(eventName string) (ddd.DomainEvent, error)
	// Has checks if an event type is registered.
	Has(eventName string) bool
}

// DefaultEventRegistry is the default in-memory implementation of EventRegistry.
type DefaultEventRegistry struct {
	mu        sync.RWMutex
	factories map[string]EventFactory
}

// NewEventRegistry creates a new DefaultEventRegistry.
func NewEventRegistry() *DefaultEventRegistry {
	return &DefaultEventRegistry{
		factories: make(map[string]EventFactory),
	}
}

// Register associates an event name with its factory.
func (r *DefaultEventRegistry) Register(eventName string, factory EventFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[eventName] = factory
}

// Create returns a new instance of the event for the given name.
func (r *DefaultEventRegistry) Create(eventName string) (ddd.DomainEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.factories[eventName]
	if !ok {
		return nil, fmt.Errorf("event type not registered: %s", eventName)
	}
	return factory(), nil
}

// Has checks if an event type is registered.
func (r *DefaultEventRegistry) Has(eventName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.factories[eventName]
	return ok
}

// globalRegistry is the default global event registry.
var (
	globalRegistry     *DefaultEventRegistry
	globalRegistryOnce sync.Once
)

// GlobalRegistry returns the global event registry singleton.
// Use this for convenience when you don't need multiple registries.
func GlobalRegistry() *DefaultEventRegistry {
	globalRegistryOnce.Do(func() {
		globalRegistry = NewEventRegistry()
	})
	return globalRegistry
}

// RegisterEvent is a convenience function to register an event in the global registry.
func RegisterEvent(eventName string, factory EventFactory) {
	GlobalRegistry().Register(eventName, factory)
}
