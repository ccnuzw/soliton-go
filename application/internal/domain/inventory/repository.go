package inventory

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// InventoryRepository 定义 Inventory 的持久化接口。
type InventoryRepository interface {
	orm.Repository[*Inventory, InventoryID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*Inventory, int64, error)
}
