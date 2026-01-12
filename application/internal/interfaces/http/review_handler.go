package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	reviewapp "github.com/soliton-go/application/internal/application/review"
	"github.com/soliton-go/application/internal/domain/review"
)

// ReviewHandler 处理 Review 相关的 HTTP 请求。
type ReviewHandler struct {
	createHandler *reviewapp.CreateReviewHandler
	updateHandler *reviewapp.UpdateReviewHandler
	deleteHandler *reviewapp.DeleteReviewHandler
	getHandler    *reviewapp.GetReviewHandler
	listHandler   *reviewapp.ListReviewsHandler
}

// NewReviewHandler 创建 ReviewHandler 实例。
func NewReviewHandler(
	createHandler *reviewapp.CreateReviewHandler,
	updateHandler *reviewapp.UpdateReviewHandler,
	deleteHandler *reviewapp.DeleteReviewHandler,
	getHandler *reviewapp.GetReviewHandler,
	listHandler *reviewapp.ListReviewsHandler,
) *ReviewHandler {
	return &ReviewHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes 注册 Review 相关路由。
func (h *ReviewHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/reviews")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/reviews
func (h *ReviewHandler) Create(c *gin.Context) {
	var req reviewapp.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := reviewapp.CreateReviewCommand{
		ID: uuid.New().String(),
		ProductId: req.ProductId,
		UserId: req.UserId,
		OrderId: req.OrderId,
		Rating: req.Rating,
		Title: req.Title,
		Content: req.Content,
		Status: review.ReviewStatus(req.Status),
		IsAnonymous: req.IsAnonymous,
		HelpfulCount: req.HelpfulCount,
		Reply: req.Reply,
		Images: req.Images,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, reviewapp.ToReviewResponse(entity))
}

// Get 处理 GET /api/reviews/:id
func (h *ReviewHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), reviewapp.GetReviewQuery{ID: id})
	if err != nil {
		NotFound(c, "review not found")
		return
	}

	Success(c, reviewapp.ToReviewResponse(entity))
}

// List 处理 GET /api/reviews?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *ReviewHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), reviewapp.ListReviewsQuery{
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
		"items":       reviewapp.ToReviewResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/reviews/:id
func (h *ReviewHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req reviewapp.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := reviewapp.UpdateReviewCommand{
		ID: id,
		ProductId: req.ProductId,
		UserId: req.UserId,
		OrderId: req.OrderId,
		Rating: req.Rating,
		Title: req.Title,
		Content: req.Content,
		Status: EnumPtr(req.Status, func(v string) review.ReviewStatus { return review.ReviewStatus(v) }),
		IsAnonymous: req.IsAnonymous,
		HelpfulCount: req.HelpfulCount,
		Reply: req.Reply,
		Images: req.Images,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, reviewapp.ToReviewResponse(entity))
}

// Delete 处理 DELETE /api/reviews/:id
func (h *ReviewHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := reviewapp.DeleteReviewCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
