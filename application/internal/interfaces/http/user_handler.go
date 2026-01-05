package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	userapp "github.com/soliton-go/application/internal/application/user"
	"github.com/soliton-go/application/internal/domain/user"
)

// UserHandler 处理 User 相关的 HTTP 请求。
type UserHandler struct {
	createHandler *userapp.CreateUserHandler
	updateHandler *userapp.UpdateUserHandler
	deleteHandler *userapp.DeleteUserHandler
	getHandler    *userapp.GetUserHandler
	listHandler   *userapp.ListUsersHandler
}

// NewUserHandler 创建 UserHandler 实例。
func NewUserHandler(
	createHandler *userapp.CreateUserHandler,
	updateHandler *userapp.UpdateUserHandler,
	deleteHandler *userapp.DeleteUserHandler,
	getHandler *userapp.GetUserHandler,
	listHandler *userapp.ListUsersHandler,
) *UserHandler {
	return &UserHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes 注册 User 相关路由。
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/users")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create 处理 POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req userapp.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := userapp.CreateUserCommand{
		ID: uuid.New().String(),
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone: req.Phone,
		Avatar: req.Avatar,
		Bio: req.Bio,
		BirthDate: req.BirthDate,
		Gender: user.UserGender(req.Gender),
		Role: user.UserRole(req.Role),
		Status: user.UserStatus(req.Status),
		EmailVerified: req.EmailVerified,
		PhoneVerified: req.PhoneVerified,
		LastLoginAt: req.LastLoginAt,
		LoginCount: req.LoginCount,
		FailedLoginCount: req.FailedLoginCount,
		Balance: req.Balance,
		Points: req.Points,
		VipLevel: req.VipLevel,
		Preferences: req.Preferences,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, userapp.ToUserResponse(entity))
}

// Get 处理 GET /api/users/:id
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), userapp.GetUserQuery{ID: id})
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, userapp.ToUserResponse(entity))
}

// List 处理 GET /api/users?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), userapp.ListUsersQuery{
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
		"items":       userapp.ToUserResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// Update 处理 PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req userapp.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := userapp.UpdateUserCommand{
		ID: id,
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone: req.Phone,
		Avatar: req.Avatar,
		Bio: req.Bio,
		BirthDate: req.BirthDate,
		Gender: EnumPtr(req.Gender, func(v string) user.UserGender { return user.UserGender(v) }),
		Role: EnumPtr(req.Role, func(v string) user.UserRole { return user.UserRole(v) }),
		Status: EnumPtr(req.Status, func(v string) user.UserStatus { return user.UserStatus(v) }),
		EmailVerified: req.EmailVerified,
		PhoneVerified: req.PhoneVerified,
		LastLoginAt: req.LastLoginAt,
		LoginCount: req.LoginCount,
		FailedLoginCount: req.FailedLoginCount,
		Balance: req.Balance,
		Points: req.Points,
		VipLevel: req.VipLevel,
		Preferences: req.Preferences,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, userapp.ToUserResponse(entity))
}

// Delete 处理 DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := userapp.DeleteUserCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
