package event

import (
	"context"
	"encoding/json"
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

// WatermillEventBus implements EventBus using Watermill.
type WatermillEventBus struct {
	publisher  message.Publisher
	subscriber message.Subscriber
}

// NewLocalEventBus creates a WatermillEventBus using Go channels (in-memory).
func NewLocalEventBus() *WatermillEventBus {
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	return &WatermillEventBus{
		publisher:  pubsub,
		subscriber: pubsub,
	}
}

// NewWatermillEventBus creates a WatermillEventBus with provided publisher/subscriber (e.g. Redis).
func NewWatermillEventBus(pub message.Publisher, sub message.Subscriber) *WatermillEventBus {
	return &WatermillEventBus{
		publisher:  pub,
		subscriber: sub,
	}
}

func (b *WatermillEventBus) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)
		msg.Metadata.Set("occurred_on", event.OccurredOn().Format(time.RFC3339))
		msg.Metadata.Set("event_name", event.EventName())

		if err := b.publisher.Publish(event.EventName(), msg); err != nil {
			return err
		}
	}
	return nil
}

func (b *WatermillEventBus) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	messages, err := b.subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			// A real implementation would need to unmarshal to the correct specific event type.
			// Currently this is a simplified version where we might need a registry of event types.
			// For now, we will just ack the message.
			// In a real system, we'd use a registry or raw JSON in handler.

			// FIXME: This implementation is incomplete regarding unmarshalling specific types.
			// It requires a TypeRegistry or similar mechanism.

			msg.Ack()
		}
	}()

	return nil
}
