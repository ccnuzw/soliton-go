package persistence

import (
	"context"
	"fmt"

	"github.com/soliton-go/application/internal/domain/inventory"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type InventoryRepoImpl struct {
	*orm.GormRepository[*inventory.Inventory, inventory.InventoryID]
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) inventory.InventoryRepository {
	return &InventoryRepoImpl{
		GormRepository: orm.NewGormRepository[*inventory.Inventory, inventory.InventoryID](db),
		db:             db,
	}
}

// FindPaginated 返回分页数据和总数。
func (r *InventoryRepoImpl) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*inventory.Inventory, int64, error) {
	var entities []*inventory.Inventory
	var total int64

	// 查询总数
	baseQuery := r.db.WithContext(ctx).Model(&inventory.Inventory{})
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

// MigrateInventory 创建数据库表（如不存在）。
func MigrateInventory(db *gorm.DB) error {
	return db.AutoMigrate(&inventory.Inventory{})
}
