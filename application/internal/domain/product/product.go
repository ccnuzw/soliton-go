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

// ProductStatus represents the Status enum.
type ProductStatus string

const (
	ProductStatusDraft ProductStatus = "draft"
	ProductStatusActive ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
)

// Product is the aggregate root.
type Product struct {
	ddd.BaseAggregateRoot
	ID ProductID `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
	Sku string `gorm:"size:255"`
	Description string `gorm:"type:text"`
	Price int64 `gorm:"not null;default:0"`
	OriginalPrice int64 `gorm:"not null;default:0"`
	Stock int `gorm:"not null;default:0"`
	CategoryId string `gorm:"size:36;index"`
	Status ProductStatus `gorm:"size:50;default:'draft'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (Product) TableName() string {
	return "products"
}

// NewProduct creates a new Product.
func NewProduct(id string, name string, sku string, description string, price int64, originalPrice int64, stock int, categoryId string, status ProductStatus) *Product {
	e := &Product{
		ID: ProductID(id),
		Name: name,
		Sku: sku,
		Description: description,
		Price: price,
		OriginalPrice: originalPrice,
		Stock: stock,
		CategoryId: categoryId,
		Status: status,
	}
	e.AddDomainEvent(NewProductCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *Product) Update(name string, sku string, description string, price int64, originalPrice int64, stock int, categoryId string, status ProductStatus) {
	e.Name = name
	e.Sku = sku
	e.Description = description
	e.Price = price
	e.OriginalPrice = originalPrice
	e.Stock = stock
	e.CategoryId = categoryId
	e.Status = status
	e.AddDomainEvent(NewProductUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *Product) GetID() ddd.ID {
	return e.ID
}
