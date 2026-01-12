package payment

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// PaymentCreatedEvent 在创建 Payment 时发布。
type PaymentCreatedEvent struct {
	ddd.BaseDomainEvent
	PaymentID string `json:"payment_id"`
}

func (e PaymentCreatedEvent) EventName() string {
	return "payment.created"
}

func NewPaymentCreatedEvent(id string) PaymentCreatedEvent {
	return PaymentCreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		PaymentID: id,
	}
}

// PaymentUpdatedEvent 在更新 Payment 时发布。
type PaymentUpdatedEvent struct {
	ddd.BaseDomainEvent
	PaymentID string `json:"payment_id"`
}

func (e PaymentUpdatedEvent) EventName() string {
	return "payment.updated"
}

func NewPaymentUpdatedEvent(id string) PaymentUpdatedEvent {
	return PaymentUpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		PaymentID: id,
	}
}

// PaymentDeletedEvent 在删除 Payment 时发布。
type PaymentDeletedEvent struct {
	ddd.BaseDomainEvent
	PaymentID string    `json:"payment_id"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (e PaymentDeletedEvent) EventName() string {
	return "payment.deleted"
}

func NewPaymentDeletedEvent(id string) PaymentDeletedEvent {
	return PaymentDeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		PaymentID: id,
		DeletedAt: time.Now(),
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("payment.created", func() ddd.DomainEvent {
		return &PaymentCreatedEvent{}
	})
	event.RegisterEvent("payment.updated", func() ddd.DomainEvent {
		return &PaymentUpdatedEvent{}
	})
	event.RegisterEvent("payment.deleted", func() ddd.DomainEvent {
		return &PaymentDeletedEvent{}
	})
}
