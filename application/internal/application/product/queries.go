package productapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/product"
)

// GetProductQuery 是获取单个 Product 的查询。
type GetProductQuery struct {
	ID string
}

// GetProductHandler 处理 GetProductQuery。
type GetProductHandler struct {
	repo product.ProductRepository
}

func NewGetProductHandler(repo product.ProductRepository) *GetProductHandler {
	return &GetProductHandler{repo: repo}
}

func (h *GetProductHandler) Handle(ctx context.Context, query GetProductQuery) (*product.Product, error) {
	return h.repo.Find(ctx, product.ProductID(query.ID))
}

// ListProductsQuery 是分页列表查询。
type ListProductsQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListProductsResult 是分页查询结果。
type ListProductsResult struct {
	Items      []*product.Product
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListProductsHandler 处理 ListProductsQuery。
type ListProductsHandler struct {
	repo product.ProductRepository
}

func NewListProductsHandler(repo product.ProductRepository) *ListProductsHandler {
	return &ListProductsHandler{repo: repo}
}

func (h *ListProductsHandler) Handle(ctx context.Context, query ListProductsQuery) (*ListProductsResult, error) {
	// 规范化分页参数
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

	// 排序字段白名单
	sortBy := strings.ToLower(strings.TrimSpace(query.SortBy))
	sortOrder := strings.ToLower(strings.TrimSpace(query.SortOrder))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	allowedSorts := map[string]struct{}{
		"id":         {},
		"created_at": {},
		"updated_at": {},
		"sku": {},
		"name": {},
		"slug": {},
		"description": {},
		"short_description": {},
		"brand": {},
		"category": {},
		"subcategory": {},
		"price": {},
		"original_price": {},
		"cost_price": {},
		"discount_percentage": {},
		"stock": {},
		"reserved_stock": {},
		"sold_count": {},
		"view_count": {},
		"rating": {},
		"review_count": {},
		"weight": {},
		"length": {},
		"width": {},
		"height": {},
		"color": {},
		"size": {},
		"material": {},
		"manufacturer": {},
		"country_of_origin": {},
		"barcode": {},
		"status": {},
		"is_featured": {},
		"is_new": {},
		"is_on_sale": {},
		"is_digital": {},
		"requires_shipping": {},
		"is_taxable": {},
		"tax_rate": {},
		"min_order_quantity": {},
		"max_order_quantity": {},
		"tags": {},
		"images": {},
		"video_url": {},
		"published_at": {},
		"discontinued_at": {},
	}
	if _, ok := allowedSorts[sortBy]; !ok {
		sortBy = "id"
	}

	// 获取总数和分页数据
	items, total, err := h.repo.FindPaginated(ctx, page, pageSize, sortBy, sortOrder)
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
