package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	productapp "github.com/soliton-go/application/internal/application/product"
	"github.com/soliton-go/application/internal/domain/product"
)

// ProductHandler 处理 Product 相关的 HTTP 请求。
type ProductHandler struct {
	createHandler *productapp.CreateProductHandler
	updateHandler *productapp.UpdateProductHandler
	deleteHandler *productapp.DeleteProductHandler
	getHandler    *productapp.GetProductHandler
	listHandler   *productapp.ListProductsHandler
}

// NewProductHandler 创建 ProductHandler 实例。
func NewProductHandler(
	createHandler *productapp.CreateProductHandler,
	updateHandler *productapp.UpdateProductHandler,
	deleteHandler *productapp.DeleteProductHandler,
	getHandler *productapp.GetProductHandler,
	listHandler *productapp.ListProductsHandler,
) *ProductHandler {
	return &ProductHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes 注册 Product 相关路由。
func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/products")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req productapp.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := productapp.CreateProductCommand{
		ID: uuid.New().String(),
		Sku: req.Sku,
		Name: req.Name,
		Slug: req.Slug,
		Description: req.Description,
		ShortDescription: req.ShortDescription,
		Brand: req.Brand,
		Category: req.Category,
		Subcategory: req.Subcategory,
		Price: req.Price,
		OriginalPrice: req.OriginalPrice,
		CostPrice: req.CostPrice,
		DiscountPercentage: req.DiscountPercentage,
		Stock: req.Stock,
		ReservedStock: req.ReservedStock,
		SoldCount: req.SoldCount,
		ViewCount: req.ViewCount,
		Rating: req.Rating,
		ReviewCount: req.ReviewCount,
		Weight: req.Weight,
		Length: req.Length,
		Width: req.Width,
		Height: req.Height,
		Color: req.Color,
		Size: req.Size,
		Material: req.Material,
		Manufacturer: req.Manufacturer,
		CountryOfOrigin: req.CountryOfOrigin,
		Barcode: req.Barcode,
		Status: product.ProductStatus(req.Status),
		IsFeatured: req.IsFeatured,
		IsNew: req.IsNew,
		IsOnSale: req.IsOnSale,
		IsDigital: req.IsDigital,
		RequiresShipping: req.RequiresShipping,
		IsTaxable: req.IsTaxable,
		TaxRate: req.TaxRate,
		MinOrderQuantity: req.MinOrderQuantity,
		MaxOrderQuantity: req.MaxOrderQuantity,
		Tags: req.Tags,
		Images: req.Images,
		VideoUrl: req.VideoUrl,
		PublishedAt: req.PublishedAt,
		DiscontinuedAt: req.DiscontinuedAt,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, productapp.ToProductResponse(entity))
}

// Get 处理 GET /api/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), productapp.GetProductQuery{ID: id})
	if err != nil {
		NotFound(c, "product not found")
		return
	}

	Success(c, productapp.ToProductResponse(entity))
}

// List 处理 GET /api/products?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), productapp.ListProductsQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{
		"items":       productapp.ToProductResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req productapp.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := productapp.UpdateProductCommand{
		ID: id,
		Sku: req.Sku,
		Name: req.Name,
		Slug: req.Slug,
		Description: req.Description,
		ShortDescription: req.ShortDescription,
		Brand: req.Brand,
		Category: req.Category,
		Subcategory: req.Subcategory,
		Price: req.Price,
		OriginalPrice: req.OriginalPrice,
		CostPrice: req.CostPrice,
		DiscountPercentage: req.DiscountPercentage,
		Stock: req.Stock,
		ReservedStock: req.ReservedStock,
		SoldCount: req.SoldCount,
		ViewCount: req.ViewCount,
		Rating: req.Rating,
		ReviewCount: req.ReviewCount,
		Weight: req.Weight,
		Length: req.Length,
		Width: req.Width,
		Height: req.Height,
		Color: req.Color,
		Size: req.Size,
		Material: req.Material,
		Manufacturer: req.Manufacturer,
		CountryOfOrigin: req.CountryOfOrigin,
		Barcode: req.Barcode,
		Status: EnumPtr(req.Status, func(v string) product.ProductStatus { return product.ProductStatus(v) }),
		IsFeatured: req.IsFeatured,
		IsNew: req.IsNew,
		IsOnSale: req.IsOnSale,
		IsDigital: req.IsDigital,
		RequiresShipping: req.RequiresShipping,
		IsTaxable: req.IsTaxable,
		TaxRate: req.TaxRate,
		MinOrderQuantity: req.MinOrderQuantity,
		MaxOrderQuantity: req.MaxOrderQuantity,
		Tags: req.Tags,
		Images: req.Images,
		VideoUrl: req.VideoUrl,
		PublishedAt: req.PublishedAt,
		DiscontinuedAt: req.DiscontinuedAt,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, productapp.ToProductResponse(entity))
}

// Delete 处理 DELETE /api/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := productapp.DeleteProductCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
