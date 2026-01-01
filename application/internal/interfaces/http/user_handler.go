package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	userapp "github.com/soliton-go/application/internal/application/user"
)

// UserHandler handles HTTP requests for User operations.
type UserHandler struct {
	createHandler *userapp.CreateUserHandler
	updateHandler *userapp.UpdateUserHandler
	deleteHandler *userapp.DeleteUserHandler
	getHandler    *userapp.GetUserHandler
	listHandler   *userapp.ListUsersHandler
}

// NewUserHandler creates a new UserHandler.
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

// RegisterRoutes registers User routes.
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/users")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req userapp.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := userapp.CreateUserCommand{
		ID:   uuid.New().String(),
		Name: req.Name,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, userapp.ToUserResponse(entity))
}

// Get handles GET /api/users/:id
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), userapp.GetUserQuery{ID: id})
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, userapp.ToUserResponse(entity))
}

// List handles GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	entities, err := h.listHandler.Handle(c.Request.Context(), userapp.ListUsersQuery{})
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, userapp.ToUserResponseList(entities))
}

// Update handles PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req userapp.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := userapp.UpdateUserCommand{
		ID:   id,
		Name: req.Name,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, userapp.ToUserResponse(entity))
}

// Delete handles DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := userapp.DeleteUserCommand{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
