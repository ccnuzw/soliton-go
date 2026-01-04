package user

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// UserRepository is the interface for User persistence.
type UserRepository interface {
	orm.Repository[*User, UserID]
	// FindPaginated returns a page of entities with total count.
	FindPaginated(ctx context.Context, page, pageSize int) ([]*User, int64, error)
}
