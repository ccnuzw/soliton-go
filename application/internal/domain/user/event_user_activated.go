package user

import (
	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// UserActivatedEvent 是领域事件。
type UserActivatedEvent struct {
	ddd.BaseDomainEvent
	UserId string `json:"user_id"`
}

// EventName 返回事件名称（主题）。
func (e UserActivatedEvent) EventName() string {
	return "user.activated"
}

// NewUserActivatedEvent 创建一个新的事件实例。
func NewUserActivatedEvent(userId string) UserActivatedEvent {
	return UserActivatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		UserId: userId,
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("user.activated", func() ddd.DomainEvent {
		return &UserActivatedEvent{}
	})
}
