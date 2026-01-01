package order

import (
	"github.com/soliton-go/framework/orm"
)

// OrderRepository is the interface for Order persistence.
type OrderRepository interface {
	orm.Repository[*Order, OrderID]
	// TODO: Add custom query methods here
}
