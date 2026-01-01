package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	orderapp "github.com/soliton-go/application/internal/application/order"
)

// OrderHandler handles HTTP requests for Order operations.
type OrderHandler struct {
	createHandler *orderapp.CreateOrderHandler
	updateHandler *orderapp.UpdateOrderHandler
	deleteHandler *orderapp.DeleteOrderHandler
	getHandler    *orderapp.GetOrderHandler
	listHandler   *orderapp.ListOrdersHandler
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(
	createHandler *orderapp.CreateOrderHandler,
	updateHandler *orderapp.UpdateOrderHandler,
	deleteHandler *orderapp.DeleteOrderHandler,
	getHandler *orderapp.GetOrderHandler,
	listHandler *orderapp.ListOrdersHandler,
) *OrderHandler {
	return &OrderHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes registers Order routes.
func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/orders")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /api/orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req orderapp.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := orderapp.CreateOrderCommand{
		ID:   uuid.New().String(),
		Name: req.Name,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, orderapp.ToOrderResponse(entity))
}

// Get handles GET /api/orders/:id
func (h *OrderHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), orderapp.GetOrderQuery{ID: id})
	if err != nil {
		NotFound(c, "order not found")
		return
	}

	Success(c, orderapp.ToOrderResponse(entity))
}

// List handles GET /api/orders
func (h *OrderHandler) List(c *gin.Context) {
	entities, err := h.listHandler.Handle(c.Request.Context(), orderapp.ListOrdersQuery{})
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, orderapp.ToOrderResponseList(entities))
}

// Update handles PUT /api/orders/:id
func (h *OrderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req orderapp.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := orderapp.UpdateOrderCommand{
		ID:   id,
		Name: req.Name,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, orderapp.ToOrderResponse(entity))
}

// Delete handles DELETE /api/orders/:id
func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := orderapp.DeleteOrderCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
