package orderapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/order"
)

// CreateOrderRequest 是创建 Order 的请求体。
type CreateOrderRequest struct {
	Userid string `json:"userid" binding:"required"`
	Orderno string `json:"orderno" binding:"required"`
	Totalamount int64 `json:"totalamount"`
	Discountamount int64 `json:"discountamount"`
	Taxamount int64 `json:"taxamount"`
	Shippingfee int64 `json:"shippingfee"`
	Finalamount int64 `json:"finalamount"`
	Currency string `json:"currency" binding:"required"`
	Paymentmethod string `json:"paymentmethod" binding:"required,oneof=credit_card debit_card paypal alipay wechat cash"`
	Paymentstatus string `json:"paymentstatus" binding:"required,oneof=pending paid failed refunded"`
	Orderstatus string `json:"orderstatus" binding:"required,oneof=pending confirmed processing shipped delivered cancelled returned"`
	Shippingmethod string `json:"shippingmethod" binding:"required,oneof=standard express overnight"`
	Trackingnumber string `json:"trackingnumber" binding:"required"`
	Receivername string `json:"receivername" binding:"required"`
	Receiverphone string `json:"receiverphone" binding:"required"`
	Receiveremail string `json:"receiveremail" binding:"required"`
	Receiveraddress string `json:"receiveraddress" binding:"required"`
	Receivercity string `json:"receivercity" binding:"required"`
	Receiverstate string `json:"receiverstate" binding:"required"`
	Receivercountry string `json:"receivercountry" binding:"required"`
	Receiverpostalcode string `json:"receiverpostalcode" binding:"required"`
	Notes string `json:"notes" binding:"required"`
	Paidat time.Time `json:"paidat"`
	Shippedat time.Time `json:"shippedat"`
	Deliveredat time.Time `json:"deliveredat"`
	Cancelledat time.Time `json:"cancelledat"`
	Refundamount int64 `json:"refundamount"`
	Refundreason string `json:"refundreason" binding:"required"`
	Itemcount int `json:"itemcount"`
	Weight float64 `json:"weight"`
	Isgift bool `json:"isgift"`
	Giftmessage string `json:"giftmessage" binding:"required"`
}

// UpdateOrderRequest 是更新 Order 的请求体。
type UpdateOrderRequest struct {
	Userid *string `json:"userid,omitempty"`
	Orderno *string `json:"orderno,omitempty"`
	Totalamount *int64 `json:"totalamount,omitempty"`
	Discountamount *int64 `json:"discountamount,omitempty"`
	Taxamount *int64 `json:"taxamount,omitempty"`
	Shippingfee *int64 `json:"shippingfee,omitempty"`
	Finalamount *int64 `json:"finalamount,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Paymentmethod *string `json:"paymentmethod,omitempty" binding:"omitempty,oneof=credit_card debit_card paypal alipay wechat cash"`
	Paymentstatus *string `json:"paymentstatus,omitempty" binding:"omitempty,oneof=pending paid failed refunded"`
	Orderstatus *string `json:"orderstatus,omitempty" binding:"omitempty,oneof=pending confirmed processing shipped delivered cancelled returned"`
	Shippingmethod *string `json:"shippingmethod,omitempty" binding:"omitempty,oneof=standard express overnight"`
	Trackingnumber *string `json:"trackingnumber,omitempty"`
	Receivername *string `json:"receivername,omitempty"`
	Receiverphone *string `json:"receiverphone,omitempty"`
	Receiveremail *string `json:"receiveremail,omitempty"`
	Receiveraddress *string `json:"receiveraddress,omitempty"`
	Receivercity *string `json:"receivercity,omitempty"`
	Receiverstate *string `json:"receiverstate,omitempty"`
	Receivercountry *string `json:"receivercountry,omitempty"`
	Receiverpostalcode *string `json:"receiverpostalcode,omitempty"`
	Notes *string `json:"notes,omitempty"`
	Paidat *time.Time `json:"paidat,omitempty"`
	Shippedat *time.Time `json:"shippedat,omitempty"`
	Deliveredat *time.Time `json:"deliveredat,omitempty"`
	Cancelledat *time.Time `json:"cancelledat,omitempty"`
	Refundamount *int64 `json:"refundamount,omitempty"`
	Refundreason *string `json:"refundreason,omitempty"`
	Itemcount *int `json:"itemcount,omitempty"`
	Weight *float64 `json:"weight,omitempty"`
	Isgift *bool `json:"isgift,omitempty"`
	Giftmessage *string `json:"giftmessage,omitempty"`
}

// OrderResponse 是 Order 的响应体。
type OrderResponse struct {
	ID        string    `json:"id"`
	Userid string `json:"userid"`
	Orderno string `json:"orderno"`
	Totalamount int64 `json:"totalamount"`
	Discountamount int64 `json:"discountamount"`
	Taxamount int64 `json:"taxamount"`
	Shippingfee int64 `json:"shippingfee"`
	Finalamount int64 `json:"finalamount"`
	Currency string `json:"currency"`
	Paymentmethod string `json:"paymentmethod"`
	Paymentstatus string `json:"paymentstatus"`
	Orderstatus string `json:"orderstatus"`
	Shippingmethod string `json:"shippingmethod"`
	Trackingnumber string `json:"trackingnumber"`
	Receivername string `json:"receivername"`
	Receiverphone string `json:"receiverphone"`
	Receiveremail string `json:"receiveremail"`
	Receiveraddress string `json:"receiveraddress"`
	Receivercity string `json:"receivercity"`
	Receiverstate string `json:"receiverstate"`
	Receivercountry string `json:"receivercountry"`
	Receiverpostalcode string `json:"receiverpostalcode"`
	Notes string `json:"notes"`
	Paidat time.Time `json:"paidat"`
	Shippedat time.Time `json:"shippedat"`
	Deliveredat time.Time `json:"deliveredat"`
	Cancelledat time.Time `json:"cancelledat"`
	Refundamount int64 `json:"refundamount"`
	Refundreason string `json:"refundreason"`
	Itemcount int `json:"itemcount"`
	Weight float64 `json:"weight"`
	Isgift bool `json:"isgift"`
	Giftmessage string `json:"giftmessage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToOrderResponse 将实体转换为响应体。
func ToOrderResponse(e *order.Order) OrderResponse {
	return OrderResponse{
		ID:        string(e.ID),
		Userid: e.Userid,
		Orderno: e.Orderno,
		Totalamount: e.Totalamount,
		Discountamount: e.Discountamount,
		Taxamount: e.Taxamount,
		Shippingfee: e.Shippingfee,
		Finalamount: e.Finalamount,
		Currency: e.Currency,
		Paymentmethod: string(e.Paymentmethod),
		Paymentstatus: string(e.Paymentstatus),
		Orderstatus: string(e.Orderstatus),
		Shippingmethod: string(e.Shippingmethod),
		Trackingnumber: e.Trackingnumber,
		Receivername: e.Receivername,
		Receiverphone: e.Receiverphone,
		Receiveremail: e.Receiveremail,
		Receiveraddress: e.Receiveraddress,
		Receivercity: e.Receivercity,
		Receiverstate: e.Receiverstate,
		Receivercountry: e.Receivercountry,
		Receiverpostalcode: e.Receiverpostalcode,
		Notes: e.Notes,
		Paidat: e.Paidat,
		Shippedat: e.Shippedat,
		Deliveredat: e.Deliveredat,
		Cancelledat: e.Cancelledat,
		Refundamount: e.Refundamount,
		Refundreason: e.Refundreason,
		Itemcount: e.Itemcount,
		Weight: e.Weight,
		Isgift: e.Isgift,
		Giftmessage: e.Giftmessage,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToOrderResponseList 将实体列表转换为响应体列表。
func ToOrderResponseList(entities []*order.Order) []OrderResponse {
	result := make([]OrderResponse, len(entities))
	for i, e := range entities {
		result[i] = ToOrderResponse(e)
	}
	return result
}
