package productapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/product"
)

// CreateProductCommand is the command for creating a Product.
type CreateProductCommand struct {
	ID string
	// Add command fields here
}

// CreateProductHandler handles CreateProductCommand.
type CreateProductHandler struct {
	repo product.ProductRepository
}

// NewCreateProductHandler creates a new handler.
func NewCreateProductHandler(repo product.ProductRepository) *CreateProductHandler {
	return &CreateProductHandler{repo: repo}
}

// Handle processes the command.
func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) error {
	entity := product.NewProduct(cmd.ID)
	return h.repo.Save(ctx, entity)
}

// UpdateProductCommand is the command for updating a Product.
type UpdateProductCommand struct {
	ID string
	// Add update fields here
}

// UpdateProductHandler handles UpdateProductCommand.
type UpdateProductHandler struct {
	repo product.ProductRepository
}

// NewUpdateProductHandler creates a new handler.
func NewUpdateProductHandler(repo product.ProductRepository) *UpdateProductHandler {
	return &UpdateProductHandler{repo: repo}
}

// Handle processes the command.
func (h *UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProductCommand) error {
	entity, err := h.repo.Find(ctx, product.ProductID(cmd.ID))
	if err != nil {
		return err
	}
	// Update entity fields here
	entity.AddDomainEvent(product.NewProductUpdatedEvent(cmd.ID))
	return h.repo.Save(ctx, entity)
}

// DeleteProductCommand is the command for deleting a Product.
type DeleteProductCommand struct {
	ID string
}

// DeleteProductHandler handles DeleteProductCommand.
type DeleteProductHandler struct {
	repo product.ProductRepository
}

// NewDeleteProductHandler creates a new handler.
func NewDeleteProductHandler(repo product.ProductRepository) *DeleteProductHandler {
	return &DeleteProductHandler{repo: repo}
}

// Handle processes the command.
func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	return h.repo.Delete(ctx, product.ProductID(cmd.ID))
}
