package productapp

import (
	"time"

	"github.com/soliton-go/test-project/internal/domain/product"
)

// CreateProductRequest is the request body for creating a Product.
type CreateProductRequest struct {
	Name string `json:"name" binding:"required"`
	Price int64 `json:"price"`
}

// UpdateProductRequest is the request body for updating a Product.
type UpdateProductRequest struct {
	Name *string `json:"name,omitempty"`
	Price *int64 `json:"price,omitempty"`
}

// ProductResponse is the response body for Product data.
type ProductResponse struct {
	ID        string    `json:"id"`
	Name string `json:"name"`
	Price int64 `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToProductResponse converts entity to response.
func ToProductResponse(e *product.Product) ProductResponse {
	return ProductResponse{
		ID:        string(e.ID),
		Name: e.Name,
		Price: e.Price,
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
