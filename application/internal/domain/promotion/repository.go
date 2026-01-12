package promotion

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// PromotionRepository 定义 Promotion 的持久化接口。
type PromotionRepository interface {
	orm.Repository[*Promotion, PromotionID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*Promotion, int64, error)
}
