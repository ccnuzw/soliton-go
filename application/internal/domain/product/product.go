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
	Shortdescription string `gorm:"size:255"`
	Brand string `gorm:"size:255"`
	Category string `gorm:"size:255"`
	Subcategory string `gorm:"size:255"`
	Price int64 `gorm:"not null;default:0"`
	Originalprice int64 `gorm:"not null;default:0"`
	Costprice int64 `gorm:"not null;default:0"`
	Discountpercentage int `gorm:"not null;default:0"`
	Stock int `gorm:"not null;default:0"`
	Reservedstock int `gorm:"not null;default:0"`
	Soldcount int `gorm:"not null;default:0"`
	Viewcount int `gorm:"not null;default:0"`
	Rating float64 `gorm:"default:0"`
	Reviewcount int `gorm:"not null;default:0"`
	Weight float64 `gorm:"default:0"`
	Length float64 `gorm:"default:0"`
	Width float64 `gorm:"default:0"`
	Height float64 `gorm:"default:0"`
	Color string `gorm:"size:255"`
	Size string `gorm:"size:255"`
	Material string `gorm:"size:255"`
	Manufacturer string `gorm:"size:255"`
	Countryoforigin string `gorm:"size:255"`
	Barcode string `gorm:"size:255"`
	Status ProductStatus `gorm:"size:50;default:'draft'"`
	Isfeatured bool `gorm:"default:false"`
	Isnew bool `gorm:"default:false"`
	Isonsale bool `gorm:"default:false"`
	Isdigital bool `gorm:"default:false"`
	Requiresshipping bool `gorm:"default:false"`
	Istaxable bool `gorm:"default:false"`
	Taxrate float64 `gorm:"default:0"`
	Minorderquantity int `gorm:"not null;default:0"`
	Maxorderquantity int `gorm:"not null;default:0"`
	Tags string `gorm:"size:255"`
	Images string `gorm:"size:255"`
	Videourl string `gorm:"size:255"`
	Publishedat time.Time `gorm:"type:timestamp"`
	Discontinuedat time.Time `gorm:"type:timestamp"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Product) TableName() string {
	return "products"
}

// NewProduct 创建一个新的 Product 实体。
func NewProduct(id string, sku string, name string, slug string, description string, shortdescription string, brand string, category string, subcategory string, price int64, originalprice int64, costprice int64, discountpercentage int, stock int, reservedstock int, soldcount int, viewcount int, rating float64, reviewcount int, weight float64, length float64, width float64, height float64, color string, size string, material string, manufacturer string, countryoforigin string, barcode string, status ProductStatus, isfeatured bool, isnew bool, isonsale bool, isdigital bool, requiresshipping bool, istaxable bool, taxrate float64, minorderquantity int, maxorderquantity int, tags string, images string, videourl string, publishedat time.Time, discontinuedat time.Time) *Product {
	e := &Product{
		ID: ProductID(id),
		Sku: sku,
		Name: name,
		Slug: slug,
		Description: description,
		Shortdescription: shortdescription,
		Brand: brand,
		Category: category,
		Subcategory: subcategory,
		Price: price,
		Originalprice: originalprice,
		Costprice: costprice,
		Discountpercentage: discountpercentage,
		Stock: stock,
		Reservedstock: reservedstock,
		Soldcount: soldcount,
		Viewcount: viewcount,
		Rating: rating,
		Reviewcount: reviewcount,
		Weight: weight,
		Length: length,
		Width: width,
		Height: height,
		Color: color,
		Size: size,
		Material: material,
		Manufacturer: manufacturer,
		Countryoforigin: countryoforigin,
		Barcode: barcode,
		Status: status,
		Isfeatured: isfeatured,
		Isnew: isnew,
		Isonsale: isonsale,
		Isdigital: isdigital,
		Requiresshipping: requiresshipping,
		Istaxable: istaxable,
		Taxrate: taxrate,
		Minorderquantity: minorderquantity,
		Maxorderquantity: maxorderquantity,
		Tags: tags,
		Images: images,
		Videourl: videourl,
		Publishedat: publishedat,
		Discontinuedat: discontinuedat,
	}
	e.AddDomainEvent(NewProductCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Product) Update(sku *string, name *string, slug *string, description *string, shortdescription *string, brand *string, category *string, subcategory *string, price *int64, originalprice *int64, costprice *int64, discountpercentage *int, stock *int, reservedstock *int, soldcount *int, viewcount *int, rating *float64, reviewcount *int, weight *float64, length *float64, width *float64, height *float64, color *string, size *string, material *string, manufacturer *string, countryoforigin *string, barcode *string, status *ProductStatus, isfeatured *bool, isnew *bool, isonsale *bool, isdigital *bool, requiresshipping *bool, istaxable *bool, taxrate *float64, minorderquantity *int, maxorderquantity *int, tags *string, images *string, videourl *string, publishedat *time.Time, discontinuedat *time.Time) {
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
	if shortdescription != nil {
		e.Shortdescription = *shortdescription
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
	if originalprice != nil {
		e.Originalprice = *originalprice
	}
	if costprice != nil {
		e.Costprice = *costprice
	}
	if discountpercentage != nil {
		e.Discountpercentage = *discountpercentage
	}
	if stock != nil {
		e.Stock = *stock
	}
	if reservedstock != nil {
		e.Reservedstock = *reservedstock
	}
	if soldcount != nil {
		e.Soldcount = *soldcount
	}
	if viewcount != nil {
		e.Viewcount = *viewcount
	}
	if rating != nil {
		e.Rating = *rating
	}
	if reviewcount != nil {
		e.Reviewcount = *reviewcount
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
	if countryoforigin != nil {
		e.Countryoforigin = *countryoforigin
	}
	if barcode != nil {
		e.Barcode = *barcode
	}
	if status != nil {
		e.Status = *status
	}
	if isfeatured != nil {
		e.Isfeatured = *isfeatured
	}
	if isnew != nil {
		e.Isnew = *isnew
	}
	if isonsale != nil {
		e.Isonsale = *isonsale
	}
	if isdigital != nil {
		e.Isdigital = *isdigital
	}
	if requiresshipping != nil {
		e.Requiresshipping = *requiresshipping
	}
	if istaxable != nil {
		e.Istaxable = *istaxable
	}
	if taxrate != nil {
		e.Taxrate = *taxrate
	}
	if minorderquantity != nil {
		e.Minorderquantity = *minorderquantity
	}
	if maxorderquantity != nil {
		e.Maxorderquantity = *maxorderquantity
	}
	if tags != nil {
		e.Tags = *tags
	}
	if images != nil {
		e.Images = *images
	}
	if videourl != nil {
		e.Videourl = *videourl
	}
	if publishedat != nil {
		e.Publishedat = *publishedat
	}
	if discontinuedat != nil {
		e.Discontinuedat = *discontinuedat
	}
	e.AddDomainEvent(NewProductUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Product) GetID() ddd.ID {
	return e.ID
}
