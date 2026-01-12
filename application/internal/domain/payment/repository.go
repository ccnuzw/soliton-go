package payment

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// PaymentRepository 定义 Payment 的持久化接口。
type PaymentRepository interface {
	orm.Repository[*Payment, PaymentID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*Payment, int64, error)
}
