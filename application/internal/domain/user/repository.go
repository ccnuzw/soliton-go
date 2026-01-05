package user

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// UserRepository 定义 User 的持久化接口。
type UserRepository interface {
	orm.Repository[*User, UserID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*User, int64, error)
}
