package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	orderapp "github.com/soliton-go/application/internal/application/order"
	"github.com/soliton-go/application/internal/domain/order"
)

// OrderHandler 处理 Order 相关的 HTTP 请求。
type OrderHandler struct {
	createHandler *orderapp.CreateOrderHandler
	updateHandler *orderapp.UpdateOrderHandler
	deleteHandler *orderapp.DeleteOrderHandler
	getHandler    *orderapp.GetOrderHandler
	listHandler   *orderapp.ListOrdersHandler
}

// NewOrderHandler 创建 OrderHandler 实例。
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

// RegisterRoutes 注册 Order 相关路由。
func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/orders")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req orderapp.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := orderapp.CreateOrderCommand{
		ID: uuid.New().String(),
		UserId: req.UserId,
		OrderNo: req.OrderNo,
		TotalAmount: req.TotalAmount,
		DiscountAmount: req.DiscountAmount,
		TaxAmount: req.TaxAmount,
		ShippingFee: req.ShippingFee,
		FinalAmount: req.FinalAmount,
		Currency: req.Currency,
		PaymentMethod: order.OrderPaymentMethod(req.PaymentMethod),
		PaymentStatus: order.OrderPaymentStatus(req.PaymentStatus),
		OrderStatus: order.OrderOrderStatus(req.OrderStatus),
		ShippingMethod: order.OrderShippingMethod(req.ShippingMethod),
		TrackingNumber: req.TrackingNumber,
		ReceiverName: req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		ReceiverEmail: req.ReceiverEmail,
		ReceiverAddress: req.ReceiverAddress,
		ReceiverCity: req.ReceiverCity,
		ReceiverState: req.ReceiverState,
		ReceiverCountry: req.ReceiverCountry,
		ReceiverPostalCode: req.ReceiverPostalCode,
		Notes: req.Notes,
		PaidAt: req.PaidAt,
		ShippedAt: req.ShippedAt,
		DeliveredAt: req.DeliveredAt,
		CancelledAt: req.CancelledAt,
		RefundAmount: req.RefundAmount,
		RefundReason: req.RefundReason,
		ItemCount: req.ItemCount,
		Weight: req.Weight,
		IsGift: req.IsGift,
		GiftMessage: req.GiftMessage,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, orderapp.ToOrderResponse(entity))
}

// Get 处理 GET /api/orders/:id
func (h *OrderHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), orderapp.GetOrderQuery{ID: id})
	if err != nil {
		NotFound(c, "order not found")
		return
	}

	Success(c, orderapp.ToOrderResponse(entity))
}

// List 处理 GET /api/orders?page=1&page_size=20
func (h *OrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.listHandler.Handle(c.Request.Context(), orderapp.ListOrdersQuery{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{
		"items":       orderapp.ToOrderResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/orders/:id
func (h *OrderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req orderapp.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := orderapp.UpdateOrderCommand{
		ID: id,
		UserId: req.UserId,
		OrderNo: req.OrderNo,
		TotalAmount: req.TotalAmount,
		DiscountAmount: req.DiscountAmount,
		TaxAmount: req.TaxAmount,
		ShippingFee: req.ShippingFee,
		FinalAmount: req.FinalAmount,
		Currency: req.Currency,
		PaymentMethod: EnumPtr(req.PaymentMethod, func(v string) order.OrderPaymentMethod { return order.OrderPaymentMethod(v) }),
		PaymentStatus: EnumPtr(req.PaymentStatus, func(v string) order.OrderPaymentStatus { return order.OrderPaymentStatus(v) }),
		OrderStatus: EnumPtr(req.OrderStatus, func(v string) order.OrderOrderStatus { return order.OrderOrderStatus(v) }),
		ShippingMethod: EnumPtr(req.ShippingMethod, func(v string) order.OrderShippingMethod { return order.OrderShippingMethod(v) }),
		TrackingNumber: req.TrackingNumber,
		ReceiverName: req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		ReceiverEmail: req.ReceiverEmail,
		ReceiverAddress: req.ReceiverAddress,
		ReceiverCity: req.ReceiverCity,
		ReceiverState: req.ReceiverState,
		ReceiverCountry: req.ReceiverCountry,
		ReceiverPostalCode: req.ReceiverPostalCode,
		Notes: req.Notes,
		PaidAt: req.PaidAt,
		ShippedAt: req.ShippedAt,
		DeliveredAt: req.DeliveredAt,
		CancelledAt: req.CancelledAt,
		RefundAmount: req.RefundAmount,
		RefundReason: req.RefundReason,
		ItemCount: req.ItemCount,
		Weight: req.Weight,
		IsGift: req.IsGift,
		GiftMessage: req.GiftMessage,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, orderapp.ToOrderResponse(entity))
}

// Delete 处理 DELETE /api/orders/:id
func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := orderapp.DeleteOrderCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
