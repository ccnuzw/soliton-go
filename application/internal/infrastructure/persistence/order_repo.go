package persistence

import (
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

// Migrate creates the table if it doesn't exist.
func MigrateOrder(db *gorm.DB) error {
	return db.AutoMigrate(&order.Order{})
}
