package orderapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/order"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module provides all Order dependencies for Fx.
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

	// Optional: Register with CQRS bus
	// Uncomment to enable CQRS pattern:
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

// RegisterMigration registers the Order table migration.
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateOrder(db	fx.Provide(NewOrderService),
)
}
