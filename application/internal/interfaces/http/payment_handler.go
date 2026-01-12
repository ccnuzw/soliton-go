package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	paymentapp "github.com/soliton-go/application/internal/application/payment"
	"github.com/soliton-go/application/internal/domain/payment"
)

// PaymentHandler 处理 Payment 相关的 HTTP 请求。
type PaymentHandler struct {
	createHandler *paymentapp.CreatePaymentHandler
	updateHandler *paymentapp.UpdatePaymentHandler
	deleteHandler *paymentapp.DeletePaymentHandler
	getHandler    *paymentapp.GetPaymentHandler
	listHandler   *paymentapp.ListPaymentsHandler
	service       *paymentapp.PaymentService
}

// NewPaymentHandler 创建 PaymentHandler 实例。
func NewPaymentHandler(
	createHandler *paymentapp.CreatePaymentHandler,
	updateHandler *paymentapp.UpdatePaymentHandler,
	deleteHandler *paymentapp.DeletePaymentHandler,
	getHandler *paymentapp.GetPaymentHandler,
	listHandler *paymentapp.ListPaymentsHandler,
	service *paymentapp.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
		service:       service,
	}
}

// RegisterRoutes 注册 Payment 相关路由。
func (h *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/payments")
	{
		api.POST("", h.Create)
		api.POST("/authorize", h.Authorize)
		api.POST("/:id/capture", h.Capture)
		api.POST("/:id/refund", h.Refund)
		api.POST("/:id/cancel", h.Cancel)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/payments
func (h *PaymentHandler) Create(c *gin.Context) {
	var req paymentapp.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := paymentapp.CreatePaymentCommand{
		ID:            uuid.New().String(),
		OrderId:       req.OrderId,
		UserId:        req.UserId,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Method:        payment.PaymentMethod(req.Method),
		Status:        payment.PaymentStatus(req.Status),
		Provider:      req.Provider,
		ProviderTxnId: req.ProviderTxnId,
		PaidAt:        req.PaidAt,
		RefundedAt:    req.RefundedAt,
		FailureReason: req.FailureReason,
		Metadata:      req.Metadata,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, paymentapp.ToPaymentResponse(entity))
}

// Get 处理 GET /api/payments/:id
func (h *PaymentHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), paymentapp.GetPaymentQuery{ID: id})
	if err != nil {
		NotFound(c, "payment not found")
		return
	}

	Success(c, paymentapp.ToPaymentResponse(entity))
}

// List 处理 GET /api/payments?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *PaymentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), paymentapp.ListPaymentsQuery{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{
		"items":       paymentapp.ToPaymentResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/payments/:id
func (h *PaymentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req paymentapp.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := paymentapp.UpdatePaymentCommand{
		ID:            id,
		OrderId:       req.OrderId,
		UserId:        req.UserId,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Method:        EnumPtr(req.Method, func(v string) payment.PaymentMethod { return payment.PaymentMethod(v) }),
		Status:        EnumPtr(req.Status, func(v string) payment.PaymentStatus { return payment.PaymentStatus(v) }),
		Provider:      req.Provider,
		ProviderTxnId: req.ProviderTxnId,
		PaidAt:        req.PaidAt,
		RefundedAt:    req.RefundedAt,
		FailureReason: req.FailureReason,
		Metadata:      req.Metadata,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, paymentapp.ToPaymentResponse(entity))
}

// Delete 处理 DELETE /api/payments/:id
func (h *PaymentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := paymentapp.DeletePaymentCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}

// Authorize 处理 POST /api/payments/authorize
func (h *PaymentHandler) Authorize(c *gin.Context) {
	var req paymentapp.AuthorizePaymentServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	resp, err := h.service.AuthorizePayment(c.Request.Context(), req)
	if err != nil {
		ServiceError(c, err)
		return
	}

	Success(c, resp)
}

// Capture 处理 POST /api/payments/:id/capture
func (h *PaymentHandler) Capture(c *gin.Context) {
	id := c.Param("id")

	var req paymentapp.CapturePaymentServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	if req.PaymentId == "" {
		req.PaymentId = id
	} else if req.PaymentId != id {
		BadRequest(c, "payment_id mismatch")
		return
	}

	resp, err := h.service.CapturePayment(c.Request.Context(), req)
	if err != nil {
		ServiceError(c, err)
		return
	}

	Success(c, resp)
}

// Refund 处理 POST /api/payments/:id/refund
func (h *PaymentHandler) Refund(c *gin.Context) {
	id := c.Param("id")

	var req paymentapp.RefundPaymentServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	if req.PaymentId == "" {
		req.PaymentId = id
	} else if req.PaymentId != id {
		BadRequest(c, "payment_id mismatch")
		return
	}

	resp, err := h.service.RefundPayment(c.Request.Context(), req)
	if err != nil {
		ServiceError(c, err)
		return
	}

	Success(c, resp)
}

// Cancel 处理 POST /api/payments/:id/cancel
func (h *PaymentHandler) Cancel(c *gin.Context) {
	id := c.Param("id")

	var req paymentapp.CancelPaymentServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	if req.PaymentId == "" {
		req.PaymentId = id
	} else if req.PaymentId != id {
		BadRequest(c, "payment_id mismatch")
		return
	}

	resp, err := h.service.CancelPayment(c.Request.Context(), req)
	if err != nil {
		ServiceError(c, err)
		return
	}

	Success(c, resp)
}
