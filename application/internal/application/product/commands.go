package productapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/product"
)

// CreateProductCommand 是创建 Product 的命令。
type CreateProductCommand struct {
	ID string
	Sku string
	Name string
	Slug string
	Description string
	Shortdescription string
	Brand string
	Category string
	Subcategory string
	Price int64
	Originalprice int64
	Costprice int64
	Discountpercentage int
	Stock int
	Reservedstock int
	Soldcount int
	Viewcount int
	Rating float64
	Reviewcount int
	Weight float64
	Length float64
	Width float64
	Height float64
	Color string
	Size string
	Material string
	Manufacturer string
	Countryoforigin string
	Barcode string
	Status product.ProductStatus
	Isfeatured bool
	Isnew bool
	Isonsale bool
	Isdigital bool
	Requiresshipping bool
	Istaxable bool
	Taxrate float64
	Minorderquantity int
	Maxorderquantity int
	Tags string
	Images string
	Videourl string
	Publishedat time.Time
	Discontinuedat time.Time
}

// CreateProductHandler 处理 CreateProductCommand。
type CreateProductHandler struct {
	repo product.ProductRepository
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreateProductHandler(repo product.ProductRepository) *CreateProductHandler {
	return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*product.Product, error) {
	entity := product.NewProduct(cmd.ID, cmd.Sku, cmd.Name, cmd.Slug, cmd.Description, cmd.Shortdescription, cmd.Brand, cmd.Category, cmd.Subcategory, cmd.Price, cmd.Originalprice, cmd.Costprice, cmd.Discountpercentage, cmd.Stock, cmd.Reservedstock, cmd.Soldcount, cmd.Viewcount, cmd.Rating, cmd.Reviewcount, cmd.Weight, cmd.Length, cmd.Width, cmd.Height, cmd.Color, cmd.Size, cmd.Material, cmd.Manufacturer, cmd.Countryoforigin, cmd.Barcode, cmd.Status, cmd.Isfeatured, cmd.Isnew, cmd.Isonsale, cmd.Isdigital, cmd.Requiresshipping, cmd.Istaxable, cmd.Taxrate, cmd.Minorderquantity, cmd.Maxorderquantity, cmd.Tags, cmd.Images, cmd.Videourl, cmd.Publishedat, cmd.Discontinuedat)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// 可选：发布领域事件
	// 取消注释以启用事件发布：
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateProductCommand 是更新 Product 的命令。
type UpdateProductCommand struct {
	ID string
	Sku *string
	Name *string
	Slug *string
	Description *string
	Shortdescription *string
	Brand *string
	Category *string
	Subcategory *string
	Price *int64
	Originalprice *int64
	Costprice *int64
	Discountpercentage *int
	Stock *int
	Reservedstock *int
	Soldcount *int
	Viewcount *int
	Rating *float64
	Reviewcount *int
	Weight *float64
	Length *float64
	Width *float64
	Height *float64
	Color *string
	Size *string
	Material *string
	Manufacturer *string
	Countryoforigin *string
	Barcode *string
	Status *product.ProductStatus
	Isfeatured *bool
	Isnew *bool
	Isonsale *bool
	Isdigital *bool
	Requiresshipping *bool
	Istaxable *bool
	Taxrate *float64
	Minorderquantity *int
	Maxorderquantity *int
	Tags *string
	Images *string
	Videourl *string
	Publishedat *time.Time
	Discontinuedat *time.Time
}

// UpdateProductHandler 处理 UpdateProductCommand。
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
	entity.Update(cmd.Sku, cmd.Name, cmd.Slug, cmd.Description, cmd.Shortdescription, cmd.Brand, cmd.Category, cmd.Subcategory, cmd.Price, cmd.Originalprice, cmd.Costprice, cmd.Discountpercentage, cmd.Stock, cmd.Reservedstock, cmd.Soldcount, cmd.Viewcount, cmd.Rating, cmd.Reviewcount, cmd.Weight, cmd.Length, cmd.Width, cmd.Height, cmd.Color, cmd.Size, cmd.Material, cmd.Manufacturer, cmd.Countryoforigin, cmd.Barcode, cmd.Status, cmd.Isfeatured, cmd.Isnew, cmd.Isonsale, cmd.Isdigital, cmd.Requiresshipping, cmd.Istaxable, cmd.Taxrate, cmd.Minorderquantity, cmd.Maxorderquantity, cmd.Tags, cmd.Images, cmd.Videourl, cmd.Publishedat, cmd.Discontinuedat)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteProductCommand 是删除 Product 的命令。
type DeleteProductCommand struct {
	ID string
}

// DeleteProductHandler 处理 DeleteProductCommand。
type DeleteProductHandler struct {
	repo product.ProductRepository
}

func NewDeleteProductHandler(repo product.ProductRepository) *DeleteProductHandler {
	return &DeleteProductHandler{repo: repo}
}

func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	return h.repo.Delete(ctx, product.ProductID(cmd.ID))
}
