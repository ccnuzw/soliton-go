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
)

// RegisterMigration registers the Product table migration.
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateProduct(db)
}
