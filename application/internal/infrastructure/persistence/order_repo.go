package persistence

import (
	"context"

	"github.com/soliton-go/application/internal/domain/order"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type OrderRepoImpl struct {
	*orm.GormRepository[*order.Order, order.OrderID]
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) order.OrderRepository {
	return &OrderRepoImpl{
		GormRepository: orm.NewGormRepository[*order.Order, order.OrderID](db),
		db:             db,
	}
}

// FindPaginated 返回分页数据和总数。
func (r *OrderRepoImpl) FindPaginated(ctx context.Context, page, pageSize int) ([]*order.Order, int64, error) {
	var entities []*order.Order
	var total int64

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&order.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateOrder 创建数据库表（如不存在）。
func MigrateOrder(db *gorm.DB) error {
	return db.AutoMigrate(&order.Order{})
}
