package productapp

import (
	"context"

	"github.com/soliton-go/test-project/internal/domain/product"
)

// CreateProductCommand is the command for creating a Product.
type CreateProductCommand struct {
	ID string
	Name string
	Price int64
}

// CreateProductHandler handles CreateProductCommand.
type CreateProductHandler struct {
	repo product.ProductRepository
	// Optional: Add event bus for domain event publishing
	// eventBus event.EventBus
}

func NewCreateProductHandler(repo product.ProductRepository) *CreateProductHandler {
	return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*product.Product, error) {
	entity := product.NewProduct(cmd.ID, cmd.Name, cmd.Price)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// Optional: Publish domain events
	// Uncomment to enable event publishing:
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateProductCommand is the command for updating a Product.
type UpdateProductCommand struct {
	ID string
	Name *string
	Price *int64
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
	entity.Update(cmd.Name, cmd.Price)
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
