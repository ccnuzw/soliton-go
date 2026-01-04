package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	productapp "github.com/soliton-go/test-project/internal/application/product"
)

// ProductHandler handles HTTP requests for Product operations.
type ProductHandler struct {
	createHandler *productapp.CreateProductHandler
	updateHandler *productapp.UpdateProductHandler
	deleteHandler *productapp.DeleteProductHandler
	getHandler    *productapp.GetProductHandler
	listHandler   *productapp.ListProductsHandler
}

// NewProductHandler creates a new ProductHandler.
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

// RegisterRoutes registers Product routes.
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

// Create handles POST /api/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req productapp.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := productapp.CreateProductCommand{
		ID: uuid.New().String(),
		Name: req.Name,
		Price: req.Price,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, productapp.ToProductResponse(entity))
}

// Get handles GET /api/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), productapp.GetProductQuery{ID: id})
	if err != nil {
		NotFound(c, "product not found")
		return
	}

	Success(c, productapp.ToProductResponse(entity))
}

// List handles GET /api/products?page=1&page_size=20
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.listHandler.Handle(c.Request.Context(), productapp.ListProductsQuery{
		Page:     page,
		PageSize: pageSize,
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

// Update handles PUT /api/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req productapp.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := productapp.UpdateProductCommand{
		ID: id,
		Name: req.Name,
		Price: req.Price,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, productapp.ToProductResponse(entity))
}

// Delete handles DELETE /api/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := productapp.DeleteProductCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
