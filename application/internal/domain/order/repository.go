package order

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// OrderRepository is the interface for Order persistence.
type OrderRepository interface {
	orm.Repository[*Order, OrderID]
	// FindPaginated returns a page of entities with total count.
	FindPaginated(ctx context.Context, page, pageSize int) ([]*Order, int64, error)
}
