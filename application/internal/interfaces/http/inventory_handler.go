package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	inventoryapp "github.com/soliton-go/application/internal/application/inventory"
	"github.com/soliton-go/application/internal/domain/inventory"
)

// InventoryHandler 处理 Inventory 相关的 HTTP 请求。
type InventoryHandler struct {
	createHandler *inventoryapp.CreateInventoryHandler
	updateHandler *inventoryapp.UpdateInventoryHandler
	deleteHandler *inventoryapp.DeleteInventoryHandler
	getHandler    *inventoryapp.GetInventoryHandler
	listHandler   *inventoryapp.ListInventorysHandler
}

// NewInventoryHandler 创建 InventoryHandler 实例。
func NewInventoryHandler(
	createHandler *inventoryapp.CreateInventoryHandler,
	updateHandler *inventoryapp.UpdateInventoryHandler,
	deleteHandler *inventoryapp.DeleteInventoryHandler,
	getHandler *inventoryapp.GetInventoryHandler,
	listHandler *inventoryapp.ListInventorysHandler,
) *InventoryHandler {
	return &InventoryHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes 注册 Inventory 相关路由。
func (h *InventoryHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/inventories")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/inventorys
func (h *InventoryHandler) Create(c *gin.Context) {
	var req inventoryapp.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := inventoryapp.CreateInventoryCommand{
		ID: uuid.New().String(),
		ProductId: req.ProductId,
		WarehouseId: req.WarehouseId,
		LocationCode: req.LocationCode,
		Stock: req.Stock,
		ReservedStock: req.ReservedStock,
		AvailableStock: req.AvailableStock,
		SafetyStock: req.SafetyStock,
		RestockLevel: req.RestockLevel,
		Status: inventory.InventoryStatus(req.Status),
		LastStockedAt: req.LastStockedAt,
		LastCheckedAt: req.LastCheckedAt,
		Notes: req.Notes,
		Metadata: req.Metadata,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, inventoryapp.ToInventoryResponse(entity))
}

// Get 处理 GET /api/inventorys/:id
func (h *InventoryHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), inventoryapp.GetInventoryQuery{ID: id})
	if err != nil {
		NotFound(c, "inventory not found")
		return
	}

	Success(c, inventoryapp.ToInventoryResponse(entity))
}

// List 处理 GET /api/inventorys?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *InventoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), inventoryapp.ListInventorysQuery{
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
		"items":       inventoryapp.ToInventoryResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/inventorys/:id
func (h *InventoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req inventoryapp.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := inventoryapp.UpdateInventoryCommand{
		ID: id,
		ProductId: req.ProductId,
		WarehouseId: req.WarehouseId,
		LocationCode: req.LocationCode,
		Stock: req.Stock,
		ReservedStock: req.ReservedStock,
		AvailableStock: req.AvailableStock,
		SafetyStock: req.SafetyStock,
		RestockLevel: req.RestockLevel,
		Status: EnumPtr(req.Status, func(v string) inventory.InventoryStatus { return inventory.InventoryStatus(v) }),
		LastStockedAt: req.LastStockedAt,
		LastCheckedAt: req.LastCheckedAt,
		Notes: req.Notes,
		Metadata: req.Metadata,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, inventoryapp.ToInventoryResponse(entity))
}

// Delete 处理 DELETE /api/inventorys/:id
func (h *InventoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := inventoryapp.DeleteInventoryCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
