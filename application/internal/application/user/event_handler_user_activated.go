package userapp

import (
	"context"
	"fmt"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/user"
)

// UserActivatedHandler 处理 UserActivatedEvent 事件。
type UserActivatedHandler struct{}

// NewUserActivatedHandler 创建 UserActivatedHandler 实例。
func NewUserActivatedHandler() *UserActivatedHandler {
	return &UserActivatedHandler{}
}

// Handle 处理领域事件。
func (h *UserActivatedHandler) Handle(ctx context.Context, evt ddd.DomainEvent) error {
	e, ok := evt.(*user.UserActivatedEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", evt)
	}
	_ = e
	// TODO: 在此实现事件处理逻辑
	return nil
}

// RegisterUserActivatedHandler 注册事件处理器。
func RegisterUserActivatedHandler(lc fx.Lifecycle, bus event.EventBus, handler *UserActivatedHandler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return bus.Subscribe(ctx, "user.activated", handler.Handle)
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
