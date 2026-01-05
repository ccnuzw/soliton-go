package orderapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/order"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 Order 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) order.OrderRepository {
		return persistence.NewOrderRepository(db)
	}),

	// Command Handlers
	fx.Provide(NewCreateOrderHandler),
	fx.Provide(NewUpdateOrderHandler),
	fx.Provide(NewDeleteOrderHandler),

	// Query Handlers
	fx.Provide(NewGetOrderHandler),
	fx.Provide(NewListOrdersHandler),
	
	fx.Provide(NewOrderService),
	// soliton-gen:services

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreateOrderHandler,
	//     updateHandler *UpdateOrderHandler,
	//     deleteHandler *DeleteOrderHandler,
	//     getHandler *GetOrderHandler,
	//     listHandler *ListOrdersHandler) {
	//     cmdBus.Register(CreateOrderCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdateOrderCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeleteOrderCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetOrderQuery{}, getHandler.Handle)
	//     queryBus.Register(ListOrdersQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration 注册 Order 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateOrder(db)
}
