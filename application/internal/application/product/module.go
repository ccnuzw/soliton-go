package productapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/product"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 Product 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) product.ProductRepository {
		return persistence.NewProductRepository(db)
	}),

	// Domain Services
	fx.Provide(product.NewProductDomainService),

	// Command Handlers
	fx.Provide(NewCreateProductHandler),
	fx.Provide(NewUpdateProductHandler),
	fx.Provide(NewDeleteProductHandler),

	// Query Handlers
	fx.Provide(NewGetProductHandler),
	fx.Provide(NewListProductsHandler),
	
	// soliton-gen:services
	// soliton-gen:event-handlers

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreateProductHandler,
	//     updateHandler *UpdateProductHandler,
	//     deleteHandler *DeleteProductHandler,
	//     getHandler *GetProductHandler,
	//     listHandler *ListProductsHandler) {
	//     cmdBus.Register(CreateProductCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdateProductCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeleteProductCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetProductQuery{}, getHandler.Handle)
	//     queryBus.Register(ListProductsQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration 注册 Product 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateProduct(db)
}
