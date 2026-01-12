package shippingapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/shipping"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 Shipping 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) shipping.ShippingRepository {
		return persistence.NewShippingRepository(db)
	}),

	// Domain Services
	fx.Provide(shipping.NewShippingDomainService),

	// Command Handlers
	fx.Provide(NewCreateShippingHandler),
	fx.Provide(NewUpdateShippingHandler),
	fx.Provide(NewDeleteShippingHandler),

	// Query Handlers
	fx.Provide(NewGetShippingHandler),
	fx.Provide(NewListShippingsHandler),

	// Application Services
	fx.Provide(NewShippingService),

	// soliton-gen:services
	// soliton-gen:event-handlers

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreateShippingHandler,
	//     updateHandler *UpdateShippingHandler,
	//     deleteHandler *DeleteShippingHandler,
	//     getHandler *GetShippingHandler,
	//     listHandler *ListShippingsHandler) {
	//     cmdBus.Register(CreateShippingCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdateShippingCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeleteShippingCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetShippingQuery{}, getHandler.Handle)
	//     queryBus.Register(ListShippingsQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration 注册 Shipping 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateShipping(db)
}
