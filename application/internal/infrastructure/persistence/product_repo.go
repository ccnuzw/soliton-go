package persistence

import (
	"github.com/soliton-go/application/internal/domain/product"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type ProductRepoImpl struct {
	*orm.GormRepository[*product.Product, product.ProductID]
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.ProductRepository {
	return &ProductRepoImpl{
		GormRepository: orm.NewGormRepository[*product.Product, product.ProductID](db),
		db:             db,
	}
}

// Migrate creates the table if it doesn't exist.
func MigrateProduct(db *gorm.DB) error {
	return db.AutoMigrate(&product.Product{})
}
