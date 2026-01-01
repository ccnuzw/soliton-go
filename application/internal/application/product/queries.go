package productapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/product"
)

// GetProductQuery is the query for getting a single Product.
type GetProductQuery struct {
	ID string
}

// GetProductHandler handles GetProductQuery.
type GetProductHandler struct {
	repo product.ProductRepository
}

// NewGetProductHandler creates a new handler.
func NewGetProductHandler(repo product.ProductRepository) *GetProductHandler {
	return &GetProductHandler{repo: repo}
}

// Handle processes the query.
func (h *GetProductHandler) Handle(ctx context.Context, query GetProductQuery) (*product.Product, error) {
	return h.repo.Find(ctx, product.ProductID(query.ID))
}

// ListProductsQuery is the query for listing all Products.
type ListProductsQuery struct{}

// ListProductsHandler handles ListProductsQuery.
type ListProductsHandler struct {
	repo product.ProductRepository
}

// NewListProductsHandler creates a new handler.
func NewListProductsHandler(repo product.ProductRepository) *ListProductsHandler {
	return &ListProductsHandler{repo: repo}
}

// Handle processes the query.
func (h *ListProductsHandler) Handle(ctx context.Context, query ListProductsQuery) ([]*product.Product, error) {
	return h.repo.FindAll(ctx)
}
