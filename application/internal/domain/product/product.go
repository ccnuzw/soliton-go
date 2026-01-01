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
	ID   ProductID `gorm:"primaryKey"`
	Name string            `gorm:"size:255"`
	// TODO: Add more fields here
}

// TableName returns the table name for GORM.
func (Product) TableName() string {
	return "products"
}

// NewProduct creates a new Product.
func NewProduct(id, name string) *Product {
	e := &Product{
		ID:   ProductID(id),
		Name: name,
	}
	e.AddDomainEvent(NewProductCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *Product) Update(name string) {
	e.Name = name
	e.AddDomainEvent(NewProductUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *Product) GetID() ddd.ID {
	return e.ID
}
