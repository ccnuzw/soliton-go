package productapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/product"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module provides all Product dependencies for Fx.
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) product.ProductRepository {
		return persistence.NewProductRepository(db)
	}),

	// Command Handlers
	fx.Provide(NewCreateProductHandler),
	fx.Provide(NewUpdateProductHandler),
	fx.Provide(NewDeleteProductHandler),

	// Query Handlers
	fx.Provide(NewGetProductHandler),
	fx.Provide(NewListProductsHandler),

	// Optional: Register with CQRS bus
	// Uncomment to enable CQRS pattern:
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

// RegisterMigration registers the Product table migration.
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateProduct(db	fx.Provide(NewProductService),
)
}
