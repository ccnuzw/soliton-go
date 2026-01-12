package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	promotionapp "github.com/soliton-go/application/internal/application/promotion"
	"github.com/soliton-go/application/internal/domain/promotion"
)

// PromotionHandler 处理 Promotion 相关的 HTTP 请求。
type PromotionHandler struct {
	createHandler *promotionapp.CreatePromotionHandler
	updateHandler *promotionapp.UpdatePromotionHandler
	deleteHandler *promotionapp.DeletePromotionHandler
	getHandler    *promotionapp.GetPromotionHandler
	listHandler   *promotionapp.ListPromotionsHandler
}

// NewPromotionHandler 创建 PromotionHandler 实例。
func NewPromotionHandler(
	createHandler *promotionapp.CreatePromotionHandler,
	updateHandler *promotionapp.UpdatePromotionHandler,
	deleteHandler *promotionapp.DeletePromotionHandler,
	getHandler *promotionapp.GetPromotionHandler,
	listHandler *promotionapp.ListPromotionsHandler,
) *PromotionHandler {
	return &PromotionHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes 注册 Promotion 相关路由。
func (h *PromotionHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/promotions")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/promotions
func (h *PromotionHandler) Create(c *gin.Context) {
	var req promotionapp.CreatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := promotionapp.CreatePromotionCommand{
		ID: uuid.New().String(),
		Code: req.Code,
		Name: req.Name,
		Description: req.Description,
		DiscountType: promotion.PromotionDiscountType(req.DiscountType),
		DiscountValue: req.DiscountValue,
		Currency: req.Currency,
		MinOrderAmount: req.MinOrderAmount,
		MaxDiscountAmount: req.MaxDiscountAmount,
		UsageLimit: req.UsageLimit,
		UsedCount: req.UsedCount,
		PerUserLimit: req.PerUserLimit,
		StartsAt: req.StartsAt,
		EndsAt: req.EndsAt,
		Status: promotion.PromotionStatus(req.Status),
		Metadata: req.Metadata,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, promotionapp.ToPromotionResponse(entity))
}

// Get 处理 GET /api/promotions/:id
func (h *PromotionHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), promotionapp.GetPromotionQuery{ID: id})
	if err != nil {
		NotFound(c, "promotion not found")
		return
	}

	Success(c, promotionapp.ToPromotionResponse(entity))
}

// List 处理 GET /api/promotions?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *PromotionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), promotionapp.ListPromotionsQuery{
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
		"items":       promotionapp.ToPromotionResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/promotions/:id
func (h *PromotionHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req promotionapp.UpdatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := promotionapp.UpdatePromotionCommand{
		ID: id,
		Code: req.Code,
		Name: req.Name,
		Description: req.Description,
		DiscountType: EnumPtr(req.DiscountType, func(v string) promotion.PromotionDiscountType { return promotion.PromotionDiscountType(v) }),
		DiscountValue: req.DiscountValue,
		Currency: req.Currency,
		MinOrderAmount: req.MinOrderAmount,
		MaxDiscountAmount: req.MaxDiscountAmount,
		UsageLimit: req.UsageLimit,
		UsedCount: req.UsedCount,
		PerUserLimit: req.PerUserLimit,
		StartsAt: req.StartsAt,
		EndsAt: req.EndsAt,
		Status: EnumPtr(req.Status, func(v string) promotion.PromotionStatus { return promotion.PromotionStatus(v) }),
		Metadata: req.Metadata,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, promotionapp.ToPromotionResponse(entity))
}

// Delete 处理 DELETE /api/promotions/:id
func (h *PromotionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := promotionapp.DeletePromotionCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
