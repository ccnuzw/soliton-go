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
	Shortdescription string `json:"shortdescription" binding:"required"`
	Brand string `json:"brand" binding:"required"`
	Category string `json:"category" binding:"required"`
	Subcategory string `json:"subcategory" binding:"required"`
	Price int64 `json:"price"`
	Originalprice int64 `json:"originalprice"`
	Costprice int64 `json:"costprice"`
	Discountpercentage int `json:"discountpercentage"`
	Stock int `json:"stock"`
	Reservedstock int `json:"reservedstock"`
	Soldcount int `json:"soldcount"`
	Viewcount int `json:"viewcount"`
	Rating float64 `json:"rating"`
	Reviewcount int `json:"reviewcount"`
	Weight float64 `json:"weight"`
	Length float64 `json:"length"`
	Width float64 `json:"width"`
	Height float64 `json:"height"`
	Color string `json:"color" binding:"required"`
	Size string `json:"size" binding:"required"`
	Material string `json:"material" binding:"required"`
	Manufacturer string `json:"manufacturer" binding:"required"`
	Countryoforigin string `json:"countryoforigin" binding:"required"`
	Barcode string `json:"barcode" binding:"required"`
	Status string `json:"status" binding:"required,oneof=draft active inactive out_of_stock discontinued"`
	Isfeatured bool `json:"isfeatured"`
	Isnew bool `json:"isnew"`
	Isonsale bool `json:"isonsale"`
	Isdigital bool `json:"isdigital"`
	Requiresshipping bool `json:"requiresshipping"`
	Istaxable bool `json:"istaxable"`
	Taxrate float64 `json:"taxrate"`
	Minorderquantity int `json:"minorderquantity"`
	Maxorderquantity int `json:"maxorderquantity"`
	Tags string `json:"tags" binding:"required"`
	Images string `json:"images" binding:"required"`
	Videourl string `json:"videourl" binding:"required"`
	Publishedat time.Time `json:"publishedat"`
	Discontinuedat time.Time `json:"discontinuedat"`
}

// UpdateProductRequest 是更新 Product 的请求体。
type UpdateProductRequest struct {
	Sku *string `json:"sku,omitempty"`
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
	Shortdescription *string `json:"shortdescription,omitempty"`
	Brand *string `json:"brand,omitempty"`
	Category *string `json:"category,omitempty"`
	Subcategory *string `json:"subcategory,omitempty"`
	Price *int64 `json:"price,omitempty"`
	Originalprice *int64 `json:"originalprice,omitempty"`
	Costprice *int64 `json:"costprice,omitempty"`
	Discountpercentage *int `json:"discountpercentage,omitempty"`
	Stock *int `json:"stock,omitempty"`
	Reservedstock *int `json:"reservedstock,omitempty"`
	Soldcount *int `json:"soldcount,omitempty"`
	Viewcount *int `json:"viewcount,omitempty"`
	Rating *float64 `json:"rating,omitempty"`
	Reviewcount *int `json:"reviewcount,omitempty"`
	Weight *float64 `json:"weight,omitempty"`
	Length *float64 `json:"length,omitempty"`
	Width *float64 `json:"width,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Color *string `json:"color,omitempty"`
	Size *string `json:"size,omitempty"`
	Material *string `json:"material,omitempty"`
	Manufacturer *string `json:"manufacturer,omitempty"`
	Countryoforigin *string `json:"countryoforigin,omitempty"`
	Barcode *string `json:"barcode,omitempty"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=draft active inactive out_of_stock discontinued"`
	Isfeatured *bool `json:"isfeatured,omitempty"`
	Isnew *bool `json:"isnew,omitempty"`
	Isonsale *bool `json:"isonsale,omitempty"`
	Isdigital *bool `json:"isdigital,omitempty"`
	Requiresshipping *bool `json:"requiresshipping,omitempty"`
	Istaxable *bool `json:"istaxable,omitempty"`
	Taxrate *float64 `json:"taxrate,omitempty"`
	Minorderquantity *int `json:"minorderquantity,omitempty"`
	Maxorderquantity *int `json:"maxorderquantity,omitempty"`
	Tags *string `json:"tags,omitempty"`
	Images *string `json:"images,omitempty"`
	Videourl *string `json:"videourl,omitempty"`
	Publishedat *time.Time `json:"publishedat,omitempty"`
	Discontinuedat *time.Time `json:"discontinuedat,omitempty"`
}

// ProductResponse 是 Product 的响应体。
type ProductResponse struct {
	ID        string    `json:"id"`
	Sku string `json:"sku"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	Shortdescription string `json:"shortdescription"`
	Brand string `json:"brand"`
	Category string `json:"category"`
	Subcategory string `json:"subcategory"`
	Price int64 `json:"price"`
	Originalprice int64 `json:"originalprice"`
	Costprice int64 `json:"costprice"`
	Discountpercentage int `json:"discountpercentage"`
	Stock int `json:"stock"`
	Reservedstock int `json:"reservedstock"`
	Soldcount int `json:"soldcount"`
	Viewcount int `json:"viewcount"`
	Rating float64 `json:"rating"`
	Reviewcount int `json:"reviewcount"`
	Weight float64 `json:"weight"`
	Length float64 `json:"length"`
	Width float64 `json:"width"`
	Height float64 `json:"height"`
	Color string `json:"color"`
	Size string `json:"size"`
	Material string `json:"material"`
	Manufacturer string `json:"manufacturer"`
	Countryoforigin string `json:"countryoforigin"`
	Barcode string `json:"barcode"`
	Status string `json:"status"`
	Isfeatured bool `json:"isfeatured"`
	Isnew bool `json:"isnew"`
	Isonsale bool `json:"isonsale"`
	Isdigital bool `json:"isdigital"`
	Requiresshipping bool `json:"requiresshipping"`
	Istaxable bool `json:"istaxable"`
	Taxrate float64 `json:"taxrate"`
	Minorderquantity int `json:"minorderquantity"`
	Maxorderquantity int `json:"maxorderquantity"`
	Tags string `json:"tags"`
	Images string `json:"images"`
	Videourl string `json:"videourl"`
	Publishedat time.Time `json:"publishedat"`
	Discontinuedat time.Time `json:"discontinuedat"`
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
		Shortdescription: e.Shortdescription,
		Brand: e.Brand,
		Category: e.Category,
		Subcategory: e.Subcategory,
		Price: e.Price,
		Originalprice: e.Originalprice,
		Costprice: e.Costprice,
		Discountpercentage: e.Discountpercentage,
		Stock: e.Stock,
		Reservedstock: e.Reservedstock,
		Soldcount: e.Soldcount,
		Viewcount: e.Viewcount,
		Rating: e.Rating,
		Reviewcount: e.Reviewcount,
		Weight: e.Weight,
		Length: e.Length,
		Width: e.Width,
		Height: e.Height,
		Color: e.Color,
		Size: e.Size,
		Material: e.Material,
		Manufacturer: e.Manufacturer,
		Countryoforigin: e.Countryoforigin,
		Barcode: e.Barcode,
		Status: string(e.Status),
		Isfeatured: e.Isfeatured,
		Isnew: e.Isnew,
		Isonsale: e.Isonsale,
		Isdigital: e.Isdigital,
		Requiresshipping: e.Requiresshipping,
		Istaxable: e.Istaxable,
		Taxrate: e.Taxrate,
		Minorderquantity: e.Minorderquantity,
		Maxorderquantity: e.Maxorderquantity,
		Tags: e.Tags,
		Images: e.Images,
		Videourl: e.Videourl,
		Publishedat: e.Publishedat,
		Discontinuedat: e.Discontinuedat,
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
