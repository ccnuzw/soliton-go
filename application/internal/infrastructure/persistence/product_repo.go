package persistence

import (
	"context"

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

// FindPaginated returns a page of entities with total count.
func (r *ProductRepoImpl) FindPaginated(ctx context.Context, page, pageSize int) ([]*product.Product, int64, error) {
	var entities []*product.Product
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&product.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get page
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateProduct creates the table if it doesn't exist.
func MigrateProduct(db *gorm.DB) error {
	return db.AutoMigrate(&product.Product{})
}
