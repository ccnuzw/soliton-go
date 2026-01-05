package product

import (
	"time"

	"github.com/soliton-go/framework/ddd"
)

// ProductID is a strong typed ID.
type ProductID string

func (id ProductID) String() string {
	return string(id)
}

// ProductStatus represents the Status enum.
type ProductStatus string

const (
	ProductStatusDraft ProductStatus = "draft"
	ProductStatusActive ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// Product is the aggregate root.
type Product struct {
	ddd.BaseAggregateRoot
	ID ProductID `gorm:"primaryKey"`
	Sku string `gorm:"size:255"` // SKU编号
	Name string `gorm:"size:255"` // 商品名称
	Slug string `gorm:"size:255"` // URL别名
	Description string `gorm:"type:text"` // 详细描述
	ShortDescription string `gorm:"type:text"` // 简短描述
	Brand string `gorm:"size:255"` // 品牌
	Category string `gorm:"size:255"` // 分类
	Subcategory string `gorm:"size:255"` // 子分类
	Price int64 `gorm:"not null;default:0"` // 售价
	OriginalPrice int64 `gorm:"not null;default:0"` // 原价
	CostPrice int64 `gorm:"not null;default:0"` // 成本价
	DiscountPercentage int `gorm:"not null;default:0"` // 折扣百分比
	Stock int `gorm:"not null;default:0"` // 库存
	ReservedStock int `gorm:"not null;default:0"` // 预留库存
	SoldCount int `gorm:"not null;default:0"` // 已售数量
	ViewCount int `gorm:"not null;default:0"` // 浏览次数
	Rating float64 `gorm:"default:0"` // 评分
	ReviewCount int `gorm:"not null;default:0"` // 评论数
	Weight float64 `gorm:"default:0"` // 重量
	Length float64 `gorm:"default:0"` // 长度
	Width float64 `gorm:"default:0"` // 宽度
	Height float64 `gorm:"default:0"` // 高度
	Color string `gorm:"size:255"` // 颜色
	Size string `gorm:"size:255"` // 尺寸
	Material string `gorm:"size:255"` // 材质
	Manufacturer string `gorm:"size:255"` // 制造商
	CountryOfOrigin string `gorm:"size:255"` // 原产国
	Barcode string `gorm:"size:255"` // 条形码
	Status ProductStatus `gorm:"size:50;default:'draft'"` // 状态
	IsFeatured bool `gorm:"default:false"` // 是否精选
	IsNew bool `gorm:"default:false"` // 是否新品
	IsOnSale bool `gorm:"default:false"` // 是否促销
	IsDigital bool `gorm:"default:false"` // 是否数字商品
	RequiresShipping bool `gorm:"default:false"` // 是否需要配送
	IsTaxable bool `gorm:"default:false"` // 是否需要税费
	TaxRate float64 `gorm:"default:0"` // 税率
	MinOrderQuantity int `gorm:"not null;default:0"` // 最小订购量
	MaxOrderQuantity int `gorm:"not null;default:0"` // 最大订购量
	Tags string `gorm:"type:text"` // 标签
	Images string `gorm:"type:text"` // 图片列表
	VideoUrl string `gorm:"size:255"` // 视频URL
	PublishedAt *time.Time  // 发布时间
	DiscontinuedAt *time.Time  // 停产时间
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (Product) TableName() string {
	return "products"
}

// NewProduct creates a new Product.
func NewProduct(id string, sku string, name string, slug string, description string, shortDescription string, brand string, category string, subcategory string, price int64, originalPrice int64, costPrice int64, discountPercentage int, stock int, reservedStock int, soldCount int, viewCount int, rating float64, reviewCount int, weight float64, length float64, width float64, height float64, color string, size string, material string, manufacturer string, countryOfOrigin string, barcode string, status ProductStatus, isFeatured bool, isNew bool, isOnSale bool, isDigital bool, requiresShipping bool, isTaxable bool, taxRate float64, minOrderQuantity int, maxOrderQuantity int, tags string, images string, videoUrl string, publishedAt *time.Time, discontinuedAt *time.Time) *Product {
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

// Update updates the entity fields.
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
		e.PublishedAt = publishedAt
	}
	if discontinuedAt != nil {
		e.DiscontinuedAt = discontinuedAt
	}
	e.AddDomainEvent(NewProductUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *Product) GetID() ddd.ID {
	return e.ID
}
