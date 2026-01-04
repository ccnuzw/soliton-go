package product

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// ProductRepository is the interface for Product persistence.
type ProductRepository interface {
	orm.Repository[*Product, ProductID]
	// FindPaginated returns a page of entities with total count.
	FindPaginated(ctx context.Context, page, pageSize int) ([]*Product, int64, error)
}
