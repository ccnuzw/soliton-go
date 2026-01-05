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

func NewGetProductHandler(repo product.ProductRepository) *GetProductHandler {
	return &GetProductHandler{repo: repo}
}

func (h *GetProductHandler) Handle(ctx context.Context, query GetProductQuery) (*product.Product, error) {
	return h.repo.Find(ctx, product.ProductID(query.ID))
}

// ListProductsQuery is the query for listing Products with pagination.
type ListProductsQuery struct {
	Page     int // Page number (1-based)
	PageSize int // Items per page (default: 20, max: 100)
}

// ListProductsResult is the paginated result for ListProductsQuery.
type ListProductsResult struct {
	Items      []*product.Product
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListProductsHandler handles ListProductsQuery.
type ListProductsHandler struct {
	repo product.ProductRepository
}

func NewListProductsHandler(repo product.ProductRepository) *ListProductsHandler {
	return &ListProductsHandler{repo: repo}
}

func (h *ListProductsHandler) Handle(ctx context.Context, query ListProductsQuery) (*ListProductsResult, error) {
	// Normalize pagination parameters
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get total count and items
	items, total, err := h.repo.FindPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &ListProductsResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
