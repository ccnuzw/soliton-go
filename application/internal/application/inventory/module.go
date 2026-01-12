package inventoryapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/inventory"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 Inventory 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) inventory.InventoryRepository {
		return persistence.NewInventoryRepository(db)
	}),

	// Domain Services
	fx.Provide(inventory.NewInventoryDomainService),

	// Command Handlers
	fx.Provide(NewCreateInventoryHandler),
	fx.Provide(NewUpdateInventoryHandler),
	fx.Provide(NewDeleteInventoryHandler),

	// Query Handlers
	fx.Provide(NewGetInventoryHandler),
	fx.Provide(NewListInventorysHandler),
	
	fx.Provide(NewInventoryService),
	// soliton-gen:services
	// soliton-gen:event-handlers

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreateInventoryHandler,
	//     updateHandler *UpdateInventoryHandler,
	//     deleteHandler *DeleteInventoryHandler,
	//     getHandler *GetInventoryHandler,
	//     listHandler *ListInventorysHandler) {
	//     cmdBus.Register(CreateInventoryCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdateInventoryCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeleteInventoryCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetInventoryQuery{}, getHandler.Handle)
	//     queryBus.Register(ListInventorysQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration 注册 Inventory 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateInventory(db)
}
