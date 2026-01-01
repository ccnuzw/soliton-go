package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/soliton-go/framework/ddd"
)

// EventBus dispatches domain events.
type EventBus interface {
	Publish(ctx context.Context, events ...ddd.DomainEvent) error
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
}

// EventHandler handles a domain event.
type EventHandler func(ctx context.Context, event ddd.DomainEvent) error

// RawEventHandler handles raw event data when type is not registered.
type RawEventHandler func(ctx context.Context, eventName string, payload []byte) error

// WatermillEventBus implements EventBus using Watermill.
type WatermillEventBus struct {
	publisher  message.Publisher
	subscriber message.Subscriber
	registry   EventRegistry
	logger     watermill.LoggerAdapter
}

// WatermillEventBusOption is a functional option for WatermillEventBus.
type WatermillEventBusOption func(*WatermillEventBus)

// WithRegistry sets a custom event registry.
func WithRegistry(registry EventRegistry) WatermillEventBusOption {
	return func(b *WatermillEventBus) {
		b.registry = registry
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger watermill.LoggerAdapter) WatermillEventBusOption {
	return func(b *WatermillEventBus) {
		b.logger = logger
	}
}

// NewLocalEventBus creates a WatermillEventBus using Go channels (in-memory).
func NewLocalEventBus(opts ...WatermillEventBusOption) *WatermillEventBus {
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	bus := &WatermillEventBus{
		publisher:  pubsub,
		subscriber: pubsub,
		registry:   GlobalRegistry(),
		logger:     logger,
	}
	for _, opt := range opts {
		opt(bus)
	}
	return bus
}

// NewWatermillEventBus creates a WatermillEventBus with provided publisher/subscriber (e.g. Redis).
func NewWatermillEventBus(pub message.Publisher, sub message.Subscriber, opts ...WatermillEventBusOption) *WatermillEventBus {
	bus := &WatermillEventBus{
		publisher:  pub,
		subscriber: sub,
		registry:   GlobalRegistry(),
		logger:     watermill.NewStdLogger(false, false),
	}
	for _, opt := range opts {
		opt(bus)
	}
	return bus
}

// Registry returns the event registry used by this bus.
func (b *WatermillEventBus) Registry() EventRegistry {
	return b.registry
}

func (b *WatermillEventBus) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event %s: %w", event.EventName(), err)
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)
		msg.Metadata.Set("occurred_on", event.OccurredOn().Format(time.RFC3339))
		msg.Metadata.Set("event_name", event.EventName())

		if err := b.publisher.Publish(event.EventName(), msg); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.EventName(), err)
		}
	}
	return nil
}

// Subscribe registers a handler for events on the given topic.
// The handler receives properly deserialized events based on the registered event types.
func (b *WatermillEventBus) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	messages, err := b.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-messages:
				if !ok {
					return
				}
				b.handleMessage(ctx, msg, handler)
			}
		}
	}()

	return nil
}

// SubscribeRaw registers a handler that receives raw event data without deserialization.
// Use this when you need to handle events dynamically or when type registration is not possible.
func (b *WatermillEventBus) SubscribeRaw(ctx context.Context, topic string, handler RawEventHandler) error {
	messages, err := b.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-messages:
				if !ok {
					return
				}
				eventName := msg.Metadata.Get("event_name")
				if err := handler(ctx, eventName, msg.Payload); err != nil {
					b.logger.Error("Failed to handle raw event", err, watermill.LogFields{
						"event_name": eventName,
						"message_id": msg.UUID,
					})
					msg.Nack()
					continue
				}
				msg.Ack()
			}
		}
	}()

	return nil
}

func (b *WatermillEventBus) handleMessage(ctx context.Context, msg *message.Message, handler EventHandler) {
	eventName := msg.Metadata.Get("event_name")
	if eventName == "" {
		b.logger.Error("Message missing event_name metadata", nil, watermill.LogFields{
			"message_id": msg.UUID,
		})
		msg.Nack()
		return
	}

	// Try to create event instance from registry
	event, err := b.registry.Create(eventName)
	if err != nil {
		b.logger.Error("Failed to create event from registry", err, watermill.LogFields{
			"event_name": eventName,
			"message_id": msg.UUID,
		})
		msg.Nack()
		return
	}

	// Unmarshal the payload into the event
	if err := json.Unmarshal(msg.Payload, event); err != nil {
		b.logger.Error("Failed to unmarshal event", err, watermill.LogFields{
			"event_name": eventName,
			"message_id": msg.UUID,
		})
		msg.Nack()
		return
	}

	// Call the handler
	if err := handler(ctx, event); err != nil {
		b.logger.Error("Event handler failed", err, watermill.LogFields{
			"event_name": eventName,
			"message_id": msg.UUID,
		})
		msg.Nack()
		return
	}

	msg.Ack()
}
