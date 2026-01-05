package product

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// ProductRepository 定义 Product 的持久化接口。
type ProductRepository interface {
	orm.Repository[*Product, ProductID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int) ([]*Product, int64, error)
}
