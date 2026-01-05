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

// FindPaginated 返回分页数据和总数。
func (r *ProductRepoImpl) FindPaginated(ctx context.Context, page, pageSize int) ([]*product.Product, int64, error) {
	var entities []*product.Product
	var total int64

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&product.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateProduct 创建数据库表（如不存在）。
func MigrateProduct(db *gorm.DB) error {
	return db.AutoMigrate(&product.Product{})
}
