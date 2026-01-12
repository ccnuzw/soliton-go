package persistence

import (
	"context"
	"fmt"

	"github.com/soliton-go/application/internal/domain/shipping"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type ShippingRepoImpl struct {
	*orm.GormRepository[*shipping.Shipping, shipping.ShippingID]
	db *gorm.DB
}

func NewShippingRepository(db *gorm.DB) shipping.ShippingRepository {
	return &ShippingRepoImpl{
		GormRepository: orm.NewGormRepository[*shipping.Shipping, shipping.ShippingID](db),
		db:             db,
	}
}

// FindPaginated 返回分页数据和总数。
func (r *ShippingRepoImpl) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*shipping.Shipping, int64, error) {
	var entities []*shipping.Shipping
	var total int64

	// 查询总数
	baseQuery := r.db.WithContext(ctx).Model(&shipping.Shipping{})
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := r.db.WithContext(ctx).Offset(offset).Limit(pageSize)
	if sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateShipping 创建数据库表（如不存在）。
func MigrateShipping(db *gorm.DB) error {
	return db.AutoMigrate(&shipping.Shipping{})
}
