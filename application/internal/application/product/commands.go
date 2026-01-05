package productapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/product"
)

// CreateProductCommand is the command for creating a Product.
type CreateProductCommand struct {
	ID string
	Sku string
	Name string
	Slug string
	Description string
	ShortDescription string
	Brand string
	Category string
	Subcategory string
	Price int64
	OriginalPrice int64
	CostPrice int64
	DiscountPercentage int
	Stock int
	ReservedStock int
	SoldCount int
	ViewCount int
	Rating float64
	ReviewCount int
	Weight float64
	Length float64
	Width float64
	Height float64
	Color string
	Size string
	Material string
	Manufacturer string
	CountryOfOrigin string
	Barcode string
	Status product.ProductStatus
	IsFeatured bool
	IsNew bool
	IsOnSale bool
	IsDigital bool
	RequiresShipping bool
	IsTaxable bool
	TaxRate float64
	MinOrderQuantity int
	MaxOrderQuantity int
	Tags string
	Images string
	VideoUrl string
	PublishedAt *time.Time
	DiscontinuedAt *time.Time
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
	entity := product.NewProduct(cmd.ID, cmd.Sku, cmd.Name, cmd.Slug, cmd.Description, cmd.ShortDescription, cmd.Brand, cmd.Category, cmd.Subcategory, cmd.Price, cmd.OriginalPrice, cmd.CostPrice, cmd.DiscountPercentage, cmd.Stock, cmd.ReservedStock, cmd.SoldCount, cmd.ViewCount, cmd.Rating, cmd.ReviewCount, cmd.Weight, cmd.Length, cmd.Width, cmd.Height, cmd.Color, cmd.Size, cmd.Material, cmd.Manufacturer, cmd.CountryOfOrigin, cmd.Barcode, cmd.Status, cmd.IsFeatured, cmd.IsNew, cmd.IsOnSale, cmd.IsDigital, cmd.RequiresShipping, cmd.IsTaxable, cmd.TaxRate, cmd.MinOrderQuantity, cmd.MaxOrderQuantity, cmd.Tags, cmd.Images, cmd.VideoUrl, cmd.PublishedAt, cmd.DiscontinuedAt)
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
	Sku *string
	Name *string
	Slug *string
	Description *string
	ShortDescription *string
	Brand *string
	Category *string
	Subcategory *string
	Price *int64
	OriginalPrice *int64
	CostPrice *int64
	DiscountPercentage *int
	Stock *int
	ReservedStock *int
	SoldCount *int
	ViewCount *int
	Rating *float64
	ReviewCount *int
	Weight *float64
	Length *float64
	Width *float64
	Height *float64
	Color *string
	Size *string
	Material *string
	Manufacturer *string
	CountryOfOrigin *string
	Barcode *string
	Status *product.ProductStatus
	IsFeatured *bool
	IsNew *bool
	IsOnSale *bool
	IsDigital *bool
	RequiresShipping *bool
	IsTaxable *bool
	TaxRate *float64
	MinOrderQuantity *int
	MaxOrderQuantity *int
	Tags *string
	Images *string
	VideoUrl *string
	PublishedAt *time.Time
	DiscontinuedAt *time.Time
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
	entity.Update(cmd.Sku, cmd.Name, cmd.Slug, cmd.Description, cmd.ShortDescription, cmd.Brand, cmd.Category, cmd.Subcategory, cmd.Price, cmd.OriginalPrice, cmd.CostPrice, cmd.DiscountPercentage, cmd.Stock, cmd.ReservedStock, cmd.SoldCount, cmd.ViewCount, cmd.Rating, cmd.ReviewCount, cmd.Weight, cmd.Length, cmd.Width, cmd.Height, cmd.Color, cmd.Size, cmd.Material, cmd.Manufacturer, cmd.CountryOfOrigin, cmd.Barcode, cmd.Status, cmd.IsFeatured, cmd.IsNew, cmd.IsOnSale, cmd.IsDigital, cmd.RequiresShipping, cmd.IsTaxable, cmd.TaxRate, cmd.MinOrderQuantity, cmd.MaxOrderQuantity, cmd.Tags, cmd.Images, cmd.VideoUrl, cmd.PublishedAt, cmd.DiscontinuedAt)
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
