package productapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/product"
)

// CreateProductCommand is the command for creating a Product.
type CreateProductCommand struct {
	ID   string
	Name string
}

// CreateProductHandler handles CreateProductCommand.
type CreateProductHandler struct {
	repo product.ProductRepository
}

func NewCreateProductHandler(repo product.ProductRepository) *CreateProductHandler {
	return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*product.Product, error) {
	entity := product.NewProduct(cmd.ID, cmd.Name)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// UpdateProductCommand is the command for updating a Product.
type UpdateProductCommand struct {
	ID   string
	Name string
}

// UpdateProductHandler handles UpdateProductCommand.
type UpdateProductHandler struct {
	repo product.ProductRepository
}

func NewUpdateProductHandler(repo product.ProductRepository) *UpdateProductHandler {
	return &UpdateProductHandler{repo: repo}
}

func (h *UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProductCommand) (*product.Product, error) {
	entity, err := h.repo.Find(ctx, product.ProductID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.Name)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteProductCommand is the command for deleting a Product.
type DeleteProductCommand struct {
	ID string
}

// DeleteProductHandler handles DeleteProductCommand.
type DeleteProductHandler struct {
	repo product.ProductRepository
}

func NewDeleteProductHandler(repo product.ProductRepository) *DeleteProductHandler {
	return &DeleteProductHandler{repo: repo}
}

func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	return h.repo.Delete(ctx, product.ProductID(cmd.ID))
}
