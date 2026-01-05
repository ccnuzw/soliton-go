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
		Userid: req.Userid,
		Orderno: req.Orderno,
		Totalamount: req.Totalamount,
		Discountamount: req.Discountamount,
		Taxamount: req.Taxamount,
		Shippingfee: req.Shippingfee,
		Finalamount: req.Finalamount,
		Currency: req.Currency,
		Paymentmethod: order.OrderPaymentmethod(req.Paymentmethod),
		Paymentstatus: order.OrderPaymentstatus(req.Paymentstatus),
		Orderstatus: order.OrderOrderstatus(req.Orderstatus),
		Shippingmethod: order.OrderShippingmethod(req.Shippingmethod),
		Trackingnumber: req.Trackingnumber,
		Receivername: req.Receivername,
		Receiverphone: req.Receiverphone,
		Receiveremail: req.Receiveremail,
		Receiveraddress: req.Receiveraddress,
		Receivercity: req.Receivercity,
		Receiverstate: req.Receiverstate,
		Receivercountry: req.Receivercountry,
		Receiverpostalcode: req.Receiverpostalcode,
		Notes: req.Notes,
		Paidat: req.Paidat,
		Shippedat: req.Shippedat,
		Deliveredat: req.Deliveredat,
		Cancelledat: req.Cancelledat,
		Refundamount: req.Refundamount,
		Refundreason: req.Refundreason,
		Itemcount: req.Itemcount,
		Weight: req.Weight,
		Isgift: req.Isgift,
		Giftmessage: req.Giftmessage,
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
		Userid: req.Userid,
		Orderno: req.Orderno,
		Totalamount: req.Totalamount,
		Discountamount: req.Discountamount,
		Taxamount: req.Taxamount,
		Shippingfee: req.Shippingfee,
		Finalamount: req.Finalamount,
		Currency: req.Currency,
		Paymentmethod: EnumPtr(req.Paymentmethod, func(v string) order.OrderPaymentmethod { return order.OrderPaymentmethod(v) }),
		Paymentstatus: EnumPtr(req.Paymentstatus, func(v string) order.OrderPaymentstatus { return order.OrderPaymentstatus(v) }),
		Orderstatus: EnumPtr(req.Orderstatus, func(v string) order.OrderOrderstatus { return order.OrderOrderstatus(v) }),
		Shippingmethod: EnumPtr(req.Shippingmethod, func(v string) order.OrderShippingmethod { return order.OrderShippingmethod(v) }),
		Trackingnumber: req.Trackingnumber,
		Receivername: req.Receivername,
		Receiverphone: req.Receiverphone,
		Receiveremail: req.Receiveremail,
		Receiveraddress: req.Receiveraddress,
		Receivercity: req.Receivercity,
		Receiverstate: req.Receiverstate,
		Receivercountry: req.Receivercountry,
		Receiverpostalcode: req.Receiverpostalcode,
		Notes: req.Notes,
		Paidat: req.Paidat,
		Shippedat: req.Shippedat,
		Deliveredat: req.Deliveredat,
		Cancelledat: req.Cancelledat,
		Refundamount: req.Refundamount,
		Refundreason: req.Refundreason,
		Itemcount: req.Itemcount,
		Weight: req.Weight,
		Isgift: req.Isgift,
		Giftmessage: req.Giftmessage,
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
