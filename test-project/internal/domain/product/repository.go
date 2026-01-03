package product

import (
	"github.com/soliton-go/framework/orm"
)

// ProductRepository is the interface for Product persistence.
type ProductRepository interface {
	orm.Repository[*Product, ProductID]
	// TODO: Add custom query methods here
}
