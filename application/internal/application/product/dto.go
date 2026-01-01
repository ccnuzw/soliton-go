package productapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/product"
)

// CreateProductRequest is the request body for creating a Product.
type CreateProductRequest struct {
	Name string `json:"name"`
	Sku string `json:"sku"`
	Description string `json:"description"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"original_price"`
	Stock int `json:"stock"`
	CategoryId string `json:"category_id"`
	Status string `json:"status"`
}

// UpdateProductRequest is the request body for updating a Product.
type UpdateProductRequest struct {
	Name string `json:"name"`
	Sku string `json:"sku"`
	Description string `json:"description"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"original_price"`
	Stock int `json:"stock"`
	CategoryId string `json:"category_id"`
	Status string `json:"status"`
}

// ProductResponse is the response body for Product data.
type ProductResponse struct {
	ID        string    `json:"id"`
	Name string `json:"name"`
	Sku string `json:"sku"`
	Description string `json:"description"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"original_price"`
	Stock int `json:"stock"`
	CategoryId string `json:"category_id"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToProductResponse converts entity to response.
func ToProductResponse(e *product.Product) ProductResponse {
	return ProductResponse{
		ID:        string(e.ID),
		Name: e.Name,
		Sku: e.Sku,
		Description: e.Description,
		Price: e.Price,
		OriginalPrice: e.OriginalPrice,
		Stock: e.Stock,
		CategoryId: e.CategoryId,
		Status: string(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToProductResponseList converts entities to response list.
func ToProductResponseList(entities []*product.Product) []ProductResponse {
	result := make([]ProductResponse, len(entities))
	for i, e := range entities {
		result[i] = ToProductResponse(e)
	}
	return result
}
