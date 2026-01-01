package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	userapp "github.com/soliton-go/application/internal/application/user"
	"github.com/soliton-go/application/internal/domain/user"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	repo          user.UserRepository
	createHandler *userapp.CreateUserHandler
	getHandler    *userapp.GetUserHandler
	listHandler   *userapp.ListUsersHandler
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(
	repo user.UserRepository,
	createHandler *userapp.CreateUserHandler,
	getHandler *userapp.GetUserHandler,
	listHandler *userapp.ListUsersHandler,
) *UserHandler {
	return &UserHandler{
		repo:          repo,
		createHandler: createHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes registers user routes on the Gin engine.
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("", h.CreateUser)
			users.GET("", h.ListUsers)
			users.GET("/:id", h.GetUser)
		}
	}
}

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UserResponse is the response body for user data.
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// Generate a new ID
	id := uuid.New().String()

	cmd := userapp.CreateUserCommand{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.createHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, UserResponse{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
}

// GetUser handles GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	query := userapp.GetUserQuery{ID: id}
	u, err := h.getHandler.Handle(c.Request.Context(), query)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, UserResponse{
		ID:    string(u.ID),
		Name:  u.Name,
		Email: u.Email,
	})
}

// ListUsers handles GET /api/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	query := userapp.ListUsersQuery{}
	users, err := h.listHandler.Handle(c.Request.Context(), query)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			ID:    string(u.ID),
			Name:  u.Name,
			Email: u.Email,
		})
	}

	Success(c, response)
}
