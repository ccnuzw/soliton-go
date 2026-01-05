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

// FindPaginated returns a page of entities with total count.
func (r *OrderRepoImpl) FindPaginated(ctx context.Context, page, pageSize int) ([]*order.Order, int64, error) {
	var entities []*order.Order
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&order.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get page
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateOrder creates the table if it doesn't exist.
func MigrateOrder(db *gorm.DB) error {
	return db.AutoMigrate(&order.Order{})
}
