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
