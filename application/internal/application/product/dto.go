package productapp

import "github.com/soliton-go/application/internal/domain/product"

// CreateProductRequest is the request body for creating a Product.
type CreateProductRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateProductRequest is the request body for updating a Product.
type UpdateProductRequest struct {
	Name string `json:"name" binding:"required"`
}

// ProductResponse is the response body for Product data.
type ProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToProductResponse converts entity to response.
func ToProductResponse(e *product.Product) ProductResponse {
	return ProductResponse{
		ID:   string(e.ID),
		Name: e.Name,
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
