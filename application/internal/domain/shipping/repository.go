package shipping

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// ShippingRepository 定义 Shipping 的持久化接口。
type ShippingRepository interface {
	orm.Repository[*Shipping, ShippingID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*Shipping, int64, error)
}
