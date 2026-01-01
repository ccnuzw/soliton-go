package order

import "github.com/soliton-go/framework/ddd"

// OrderID is a strong typed ID.
type OrderID string

func (id OrderID) String() string {
	return string(id)
}

// Order is the aggregate root.
type Order struct {
	ddd.BaseAggregateRoot
	ID   OrderID `gorm:"primaryKey"`
	Name string            `gorm:"size:255"`
	// TODO: Add more fields here
}

// TableName returns the table name for GORM.
func (Order) TableName() string {
	return "orders"
}

// NewOrder creates a new Order.
func NewOrder(id, name string) *Order {
	e := &Order{
		ID:   OrderID(id),
		Name: name,
	}
	e.AddDomainEvent(NewOrderCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *Order) Update(name string) {
	e.Name = name
	e.AddDomainEvent(NewOrderUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *Order) GetID() ddd.ID {
	return e.ID
}
