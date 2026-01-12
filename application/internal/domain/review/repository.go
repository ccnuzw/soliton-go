package review

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// ReviewRepository 定义 Review 的持久化接口。
type ReviewRepository interface {
	orm.Repository[*Review, ReviewID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*Review, int64, error)
}
