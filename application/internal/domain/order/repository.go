package order

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// OrderRepository 定义 Order 的持久化接口。
type OrderRepository interface {
	orm.Repository[*Order, OrderID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int) ([]*Order, int64, error)
}
