package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	shippingapp "github.com/soliton-go/application/internal/application/shipping"
	"github.com/soliton-go/application/internal/domain/shipping"
)

// ShippingHandler 处理 Shipping 相关的 HTTP 请求。
type ShippingHandler struct {
	createHandler *shippingapp.CreateShippingHandler
	updateHandler *shippingapp.UpdateShippingHandler
	deleteHandler *shippingapp.DeleteShippingHandler
	getHandler    *shippingapp.GetShippingHandler
	listHandler   *shippingapp.ListShippingsHandler
}

// NewShippingHandler 创建 ShippingHandler 实例。
func NewShippingHandler(
	createHandler *shippingapp.CreateShippingHandler,
	updateHandler *shippingapp.UpdateShippingHandler,
	deleteHandler *shippingapp.DeleteShippingHandler,
	getHandler *shippingapp.GetShippingHandler,
	listHandler *shippingapp.ListShippingsHandler,
) *ShippingHandler {
	return &ShippingHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes 注册 Shipping 相关路由。
func (h *ShippingHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/shippings")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/shippings
func (h *ShippingHandler) Create(c *gin.Context) {
	var req shippingapp.CreateShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := shippingapp.CreateShippingCommand{
		ID: uuid.New().String(),
		OrderId: req.OrderId,
		Carrier: req.Carrier,
		ShippingMethod: shipping.ShippingShippingMethod(req.ShippingMethod),
		TrackingNumber: req.TrackingNumber,
		Status: shipping.ShippingStatus(req.Status),
		ShippedAt: req.ShippedAt,
		DeliveredAt: req.DeliveredAt,
		ReceiverName: req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		ReceiverAddress: req.ReceiverAddress,
		ReceiverCity: req.ReceiverCity,
		ReceiverState: req.ReceiverState,
		ReceiverCountry: req.ReceiverCountry,
		ReceiverPostalCode: req.ReceiverPostalCode,
		Notes: req.Notes,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, shippingapp.ToShippingResponse(entity))
}

// Get 处理 GET /api/shippings/:id
func (h *ShippingHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), shippingapp.GetShippingQuery{ID: id})
	if err != nil {
		NotFound(c, "shipping not found")
		return
	}

	Success(c, shippingapp.ToShippingResponse(entity))
}

// List 处理 GET /api/shippings?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *ShippingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), shippingapp.ListShippingsQuery{
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
		"items":       shippingapp.ToShippingResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/shippings/:id
func (h *ShippingHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req shippingapp.UpdateShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := shippingapp.UpdateShippingCommand{
		ID: id,
		OrderId: req.OrderId,
		Carrier: req.Carrier,
		ShippingMethod: EnumPtr(req.ShippingMethod, func(v string) shipping.ShippingShippingMethod { return shipping.ShippingShippingMethod(v) }),
		TrackingNumber: req.TrackingNumber,
		Status: EnumPtr(req.Status, func(v string) shipping.ShippingStatus { return shipping.ShippingStatus(v) }),
		ShippedAt: req.ShippedAt,
		DeliveredAt: req.DeliveredAt,
		ReceiverName: req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		ReceiverAddress: req.ReceiverAddress,
		ReceiverCity: req.ReceiverCity,
		ReceiverState: req.ReceiverState,
		ReceiverCountry: req.ReceiverCountry,
		ReceiverPostalCode: req.ReceiverPostalCode,
		Notes: req.Notes,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, shippingapp.ToShippingResponse(entity))
}

// Delete 处理 DELETE /api/shippings/:id
func (h *ShippingHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := shippingapp.DeleteShippingCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
