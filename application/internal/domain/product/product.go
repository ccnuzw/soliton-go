package product

import "github.com/soliton-go/framework/ddd"

// ProductID is a strong typed ID.
type ProductID string

func (id ProductID) String() string {
	return string(id)
}

// Product is the aggregate root.
type Product struct {
	ddd.BaseAggregateRoot
	ID ProductID `gorm:"primaryKey"`
	// Add your fields here
}

// NewProduct creates a new Product.
func NewProduct(id string) *Product {
	e := &Product{
		ID: ProductID(id),
	}
	e.AddDomainEvent(NewProductCreatedEvent(id))
	return e
}

// GetID returns the entity ID.
func (e *Product) GetID() ddd.ID {
	return e.ID
}
