package productapp

import (
	"context"

	"github.com/soliton-go/test-project/internal/domain/product"
)

// GetProductQuery is the query for getting a single Product.
type GetProductQuery struct {
	ID string
}

// GetProductHandler handles GetProductQuery.
type GetProductHandler struct {
	repo product.ProductRepository
}

func NewGetProductHandler(repo product.ProductRepository) *GetProductHandler {
	return &GetProductHandler{repo: repo}
}

func (h *GetProductHandler) Handle(ctx context.Context, query GetProductQuery) (*product.Product, error) {
	return h.repo.Find(ctx, product.ProductID(query.ID))
}

// ListProductsQuery is the query for listing all Products.
type ListProductsQuery struct{}

// ListProductsHandler handles ListProductsQuery.
type ListProductsHandler struct {
	repo product.ProductRepository
}

func NewListProductsHandler(repo product.ProductRepository) *ListProductsHandler {
	return &ListProductsHandler{repo: repo}
}

func (h *ListProductsHandler) Handle(ctx context.Context, query ListProductsQuery) ([]*product.Product, error) {
	return h.repo.FindAll(ctx)
}
