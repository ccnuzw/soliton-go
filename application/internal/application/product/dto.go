package productapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/product"
)

// CreateProductRequest 是创建 Product 的请求体。
type CreateProductRequest struct {
	Sku string `json:"sku" binding:"required"`
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
	Description string `json:"description" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Brand string `json:"brand" binding:"required"`
	Category string `json:"category" binding:"required"`
	Subcategory string `json:"subcategory" binding:"required"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"original_price"`
	CostPrice int64 `json:"cost_price"`
	DiscountPercentage int `json:"discount_percentage"`
	Stock int `json:"stock"`
	ReservedStock int `json:"reserved_stock"`
	SoldCount int `json:"sold_count"`
	ViewCount int `json:"view_count"`
	Rating float64 `json:"rating"`
	ReviewCount int `json:"review_count"`
	Weight float64 `json:"weight"`
	Length float64 `json:"length"`
	Width float64 `json:"width"`
	Height float64 `json:"height"`
	Color string `json:"color" binding:"required"`
	Size string `json:"size" binding:"required"`
	Material string `json:"material" binding:"required"`
	Manufacturer string `json:"manufacturer" binding:"required"`
	CountryOfOrigin string `json:"country_of_origin" binding:"required"`
	Barcode string `json:"barcode" binding:"required"`
	Status string `json:"status" binding:"required,oneof=draft active inactive out_of_stock discontinued"`
	IsFeatured bool `json:"is_featured"`
	IsNew bool `json:"is_new"`
	IsOnSale bool `json:"is_on_sale"`
	IsDigital bool `json:"is_digital"`
	RequiresShipping bool `json:"requires_shipping"`
	IsTaxable bool `json:"is_taxable"`
	TaxRate float64 `json:"tax_rate"`
	MinOrderQuantity int `json:"min_order_quantity"`
	MaxOrderQuantity int `json:"max_order_quantity"`
	Tags string `json:"tags" binding:"required"`
	Images string `json:"images" binding:"required"`
	VideoUrl string `json:"video_url" binding:"required"`
	PublishedAt *time.Time `json:"published_at"`
	DiscontinuedAt *time.Time `json:"discontinued_at"`
}

// UpdateProductRequest 是更新 Product 的请求体。
type UpdateProductRequest struct {
	Sku *string `json:"sku,omitempty"`
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
	ShortDescription *string `json:"short_description,omitempty"`
	Brand *string `json:"brand,omitempty"`
	Category *string `json:"category,omitempty"`
	Subcategory *string `json:"subcategory,omitempty"`
	Price *int64 `json:"price,omitempty"`
	OriginalPrice *int64 `json:"original_price,omitempty"`
	CostPrice *int64 `json:"cost_price,omitempty"`
	DiscountPercentage *int `json:"discount_percentage,omitempty"`
	Stock *int `json:"stock,omitempty"`
	ReservedStock *int `json:"reserved_stock,omitempty"`
	SoldCount *int `json:"sold_count,omitempty"`
	ViewCount *int `json:"view_count,omitempty"`
	Rating *float64 `json:"rating,omitempty"`
	ReviewCount *int `json:"review_count,omitempty"`
	Weight *float64 `json:"weight,omitempty"`
	Length *float64 `json:"length,omitempty"`
	Width *float64 `json:"width,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Color *string `json:"color,omitempty"`
	Size *string `json:"size,omitempty"`
	Material *string `json:"material,omitempty"`
	Manufacturer *string `json:"manufacturer,omitempty"`
	CountryOfOrigin *string `json:"country_of_origin,omitempty"`
	Barcode *string `json:"barcode,omitempty"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=draft active inactive out_of_stock discontinued"`
	IsFeatured *bool `json:"is_featured,omitempty"`
	IsNew *bool `json:"is_new,omitempty"`
	IsOnSale *bool `json:"is_on_sale,omitempty"`
	IsDigital *bool `json:"is_digital,omitempty"`
	RequiresShipping *bool `json:"requires_shipping,omitempty"`
	IsTaxable *bool `json:"is_taxable,omitempty"`
	TaxRate *float64 `json:"tax_rate,omitempty"`
	MinOrderQuantity *int `json:"min_order_quantity,omitempty"`
	MaxOrderQuantity *int `json:"max_order_quantity,omitempty"`
	Tags *string `json:"tags,omitempty"`
	Images *string `json:"images,omitempty"`
	VideoUrl *string `json:"video_url,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	DiscontinuedAt *time.Time `json:"discontinued_at,omitempty"`
}

// ProductResponse 是 Product 的响应体。
type ProductResponse struct {
	ID        string    `json:"id"`
	Sku string `json:"sku"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	ShortDescription string `json:"short_description"`
	Brand string `json:"brand"`
	Category string `json:"category"`
	Subcategory string `json:"subcategory"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"original_price"`
	CostPrice int64 `json:"cost_price"`
	DiscountPercentage int `json:"discount_percentage"`
	Stock int `json:"stock"`
	ReservedStock int `json:"reserved_stock"`
	SoldCount int `json:"sold_count"`
	ViewCount int `json:"view_count"`
	Rating float64 `json:"rating"`
	ReviewCount int `json:"review_count"`
	Weight float64 `json:"weight"`
	Length float64 `json:"length"`
	Width float64 `json:"width"`
	Height float64 `json:"height"`
	Color string `json:"color"`
	Size string `json:"size"`
	Material string `json:"material"`
	Manufacturer string `json:"manufacturer"`
	CountryOfOrigin string `json:"country_of_origin"`
	Barcode string `json:"barcode"`
	Status string `json:"status"`
	IsFeatured bool `json:"is_featured"`
	IsNew bool `json:"is_new"`
	IsOnSale bool `json:"is_on_sale"`
	IsDigital bool `json:"is_digital"`
	RequiresShipping bool `json:"requires_shipping"`
	IsTaxable bool `json:"is_taxable"`
	TaxRate float64 `json:"tax_rate"`
	MinOrderQuantity int `json:"min_order_quantity"`
	MaxOrderQuantity int `json:"max_order_quantity"`
	Tags string `json:"tags"`
	Images string `json:"images"`
	VideoUrl string `json:"video_url"`
	PublishedAt *time.Time `json:"published_at"`
	DiscontinuedAt *time.Time `json:"discontinued_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToProductResponse 将实体转换为响应体。
func ToProductResponse(e *product.Product) ProductResponse {
	return ProductResponse{
		ID:        string(e.ID),
		Sku: e.Sku,
		Name: e.Name,
		Slug: e.Slug,
		Description: e.Description,
		ShortDescription: e.ShortDescription,
		Brand: e.Brand,
		Category: e.Category,
		Subcategory: e.Subcategory,
		Price: e.Price,
		OriginalPrice: e.OriginalPrice,
		CostPrice: e.CostPrice,
		DiscountPercentage: e.DiscountPercentage,
		Stock: e.Stock,
		ReservedStock: e.ReservedStock,
		SoldCount: e.SoldCount,
		ViewCount: e.ViewCount,
		Rating: e.Rating,
		ReviewCount: e.ReviewCount,
		Weight: e.Weight,
		Length: e.Length,
		Width: e.Width,
		Height: e.Height,
		Color: e.Color,
		Size: e.Size,
		Material: e.Material,
		Manufacturer: e.Manufacturer,
		CountryOfOrigin: e.CountryOfOrigin,
		Barcode: e.Barcode,
		Status: string(e.Status),
		IsFeatured: e.IsFeatured,
		IsNew: e.IsNew,
		IsOnSale: e.IsOnSale,
		IsDigital: e.IsDigital,
		RequiresShipping: e.RequiresShipping,
		IsTaxable: e.IsTaxable,
		TaxRate: e.TaxRate,
		MinOrderQuantity: e.MinOrderQuantity,
		MaxOrderQuantity: e.MaxOrderQuantity,
		Tags: e.Tags,
		Images: e.Images,
		VideoUrl: e.VideoUrl,
		PublishedAt: e.PublishedAt,
		DiscontinuedAt: e.DiscontinuedAt,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToProductResponseList 将实体列表转换为响应体列表。
func ToProductResponseList(entities []*product.Product) []ProductResponse {
	result := make([]ProductResponse, len(entities))
	for i, e := range entities {
		result[i] = ToProductResponse(e)
	}
	return result
}
