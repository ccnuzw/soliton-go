package product

import (
	"time"

	"github.com/soliton-go/framework/ddd"
)

// ProductID is a strong typed ID.
type ProductID string

func (id ProductID) String() string {
	return string(id)
}

// Product is the aggregate root.
type Product struct {
	ddd.BaseAggregateRoot
	ID ProductID `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
	Price int64 `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (Product) TableName() string {
	return "products"
}

// NewProduct creates a new Product.
func NewProduct(id string, name string, price int64) *Product {
	e := &Product{
		ID: ProductID(id),
		Name: name,
		Price: price,
	}
	e.AddDomainEvent(NewProductCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *Product) Update(name *string, price *int64) {
	if name != nil {
		e.Name = *name
	}
	if price != nil {
		e.Price = *price
	}
	e.AddDomainEvent(NewProductUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *Product) GetID() ddd.ID {
	return e.ID
}
