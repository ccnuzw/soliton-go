package product

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)

// ProductID 是强类型的实体标识符。
type ProductID string

func (id ProductID) String() string {
	return string(id)
}

// ProductStatus 表示 Status 字段的枚举类型。
type ProductStatus string

const (
	ProductStatusDraft ProductStatus = "draft"
	ProductStatusActive ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// Product 是聚合根实体。
type Product struct {
	ddd.BaseAggregateRoot
	ID ProductID `gorm:"primaryKey"`
	Sku string `gorm:"size:255"`
	Name string `gorm:"size:255"`
	Slug string `gorm:"size:255"`
	Description string `gorm:"size:255"`
	ShortDescription string `gorm:"size:255"`
	Brand string `gorm:"size:255"`
	Category string `gorm:"size:255"`
	Subcategory string `gorm:"size:255"`
	Price int64 `gorm:"not null;default:0"`
	OriginalPrice int64 `gorm:"not null;default:0"`
	CostPrice int64 `gorm:"not null;default:0"`
	DiscountPercentage int `gorm:"not null;default:0"`
	Stock int `gorm:"not null;default:0"`
	ReservedStock int `gorm:"not null;default:0"`
	SoldCount int `gorm:"not null;default:0"`
	ViewCount int `gorm:"not null;default:0"`
	Rating float64 `gorm:"default:0"`
	ReviewCount int `gorm:"not null;default:0"`
	Weight float64 `gorm:"default:0"`
	Length float64 `gorm:"default:0"`
	Width float64 `gorm:"default:0"`
	Height float64 `gorm:"default:0"`
	Color string `gorm:"size:255"`
	Size string `gorm:"size:255"`
	Material string `gorm:"size:255"`
	Manufacturer string `gorm:"size:255"`
	CountryOfOrigin string `gorm:"size:255"`
	Barcode string `gorm:"size:255"`
	Status ProductStatus `gorm:"size:50;default:'draft'"`
	IsFeatured bool `gorm:"default:false"`
	IsNew bool `gorm:"default:false"`
	IsOnSale bool `gorm:"default:false"`
	IsDigital bool `gorm:"default:false"`
	RequiresShipping bool `gorm:"default:false"`
	IsTaxable bool `gorm:"default:false"`
	TaxRate float64 `gorm:"default:0"`
	MinOrderQuantity int `gorm:"not null;default:0"`
	MaxOrderQuantity int `gorm:"not null;default:0"`
	Tags string `gorm:"size:255"`
	Images string `gorm:"size:255"`
	VideoUrl string `gorm:"size:255"`
	PublishedAt time.Time `gorm:"type:timestamp"`
	DiscontinuedAt time.Time `gorm:"type:timestamp"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Product) TableName() string {
	return "products"
}

// NewProduct 创建一个新的 Product 实体。
func NewProduct(id string, sku string, name string, slug string, description string, shortDescription string, brand string, category string, subcategory string, price int64, originalPrice int64, costPrice int64, discountPercentage int, stock int, reservedStock int, soldCount int, viewCount int, rating float64, reviewCount int, weight float64, length float64, width float64, height float64, color string, size string, material string, manufacturer string, countryOfOrigin string, barcode string, status ProductStatus, isFeatured bool, isNew bool, isOnSale bool, isDigital bool, requiresShipping bool, isTaxable bool, taxRate float64, minOrderQuantity int, maxOrderQuantity int, tags string, images string, videoUrl string, publishedAt time.Time, discontinuedAt time.Time) *Product {
	e := &Product{
		ID: ProductID(id),
		Sku: sku,
		Name: name,
		Slug: slug,
		Description: description,
		ShortDescription: shortDescription,
		Brand: brand,
		Category: category,
		Subcategory: subcategory,
		Price: price,
		OriginalPrice: originalPrice,
		CostPrice: costPrice,
		DiscountPercentage: discountPercentage,
		Stock: stock,
		ReservedStock: reservedStock,
		SoldCount: soldCount,
		ViewCount: viewCount,
		Rating: rating,
		ReviewCount: reviewCount,
		Weight: weight,
		Length: length,
		Width: width,
		Height: height,
		Color: color,
		Size: size,
		Material: material,
		Manufacturer: manufacturer,
		CountryOfOrigin: countryOfOrigin,
		Barcode: barcode,
		Status: status,
		IsFeatured: isFeatured,
		IsNew: isNew,
		IsOnSale: isOnSale,
		IsDigital: isDigital,
		RequiresShipping: requiresShipping,
		IsTaxable: isTaxable,
		TaxRate: taxRate,
		MinOrderQuantity: minOrderQuantity,
		MaxOrderQuantity: maxOrderQuantity,
		Tags: tags,
		Images: images,
		VideoUrl: videoUrl,
		PublishedAt: publishedAt,
		DiscontinuedAt: discontinuedAt,
	}
	e.AddDomainEvent(NewProductCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Product) Update(sku *string, name *string, slug *string, description *string, shortDescription *string, brand *string, category *string, subcategory *string, price *int64, originalPrice *int64, costPrice *int64, discountPercentage *int, stock *int, reservedStock *int, soldCount *int, viewCount *int, rating *float64, reviewCount *int, weight *float64, length *float64, width *float64, height *float64, color *string, size *string, material *string, manufacturer *string, countryOfOrigin *string, barcode *string, status *ProductStatus, isFeatured *bool, isNew *bool, isOnSale *bool, isDigital *bool, requiresShipping *bool, isTaxable *bool, taxRate *float64, minOrderQuantity *int, maxOrderQuantity *int, tags *string, images *string, videoUrl *string, publishedAt *time.Time, discontinuedAt *time.Time) {
	if sku != nil {
		e.Sku = *sku
	}
	if name != nil {
		e.Name = *name
	}
	if slug != nil {
		e.Slug = *slug
	}
	if description != nil {
		e.Description = *description
	}
	if shortDescription != nil {
		e.ShortDescription = *shortDescription
	}
	if brand != nil {
		e.Brand = *brand
	}
	if category != nil {
		e.Category = *category
	}
	if subcategory != nil {
		e.Subcategory = *subcategory
	}
	if price != nil {
		e.Price = *price
	}
	if originalPrice != nil {
		e.OriginalPrice = *originalPrice
	}
	if costPrice != nil {
		e.CostPrice = *costPrice
	}
	if discountPercentage != nil {
		e.DiscountPercentage = *discountPercentage
	}
	if stock != nil {
		e.Stock = *stock
	}
	if reservedStock != nil {
		e.ReservedStock = *reservedStock
	}
	if soldCount != nil {
		e.SoldCount = *soldCount
	}
	if viewCount != nil {
		e.ViewCount = *viewCount
	}
	if rating != nil {
		e.Rating = *rating
	}
	if reviewCount != nil {
		e.ReviewCount = *reviewCount
	}
	if weight != nil {
		e.Weight = *weight
	}
	if length != nil {
		e.Length = *length
	}
	if width != nil {
		e.Width = *width
	}
	if height != nil {
		e.Height = *height
	}
	if color != nil {
		e.Color = *color
	}
	if size != nil {
		e.Size = *size
	}
	if material != nil {
		e.Material = *material
	}
	if manufacturer != nil {
		e.Manufacturer = *manufacturer
	}
	if countryOfOrigin != nil {
		e.CountryOfOrigin = *countryOfOrigin
	}
	if barcode != nil {
		e.Barcode = *barcode
	}
	if status != nil {
		e.Status = *status
	}
	if isFeatured != nil {
		e.IsFeatured = *isFeatured
	}
	if isNew != nil {
		e.IsNew = *isNew
	}
	if isOnSale != nil {
		e.IsOnSale = *isOnSale
	}
	if isDigital != nil {
		e.IsDigital = *isDigital
	}
	if requiresShipping != nil {
		e.RequiresShipping = *requiresShipping
	}
	if isTaxable != nil {
		e.IsTaxable = *isTaxable
	}
	if taxRate != nil {
		e.TaxRate = *taxRate
	}
	if minOrderQuantity != nil {
		e.MinOrderQuantity = *minOrderQuantity
	}
	if maxOrderQuantity != nil {
		e.MaxOrderQuantity = *maxOrderQuantity
	}
	if tags != nil {
		e.Tags = *tags
	}
	if images != nil {
		e.Images = *images
	}
	if videoUrl != nil {
		e.VideoUrl = *videoUrl
	}
	if publishedAt != nil {
		e.PublishedAt = *publishedAt
	}
	if discontinuedAt != nil {
		e.DiscontinuedAt = *discontinuedAt
	}
	e.AddDomainEvent(NewProductUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Product) GetID() ddd.ID {
	return e.ID
}
