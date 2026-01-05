package orderapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/order"
)

// CreateOrderRequest 是创建 Order 的请求体。
type CreateOrderRequest struct {
	UserId string `json:"user_id" binding:"required"`
	OrderNo string `json:"order_no" binding:"required"`
	TotalAmount int64 `json:"total_amount"`
	DiscountAmount int64 `json:"discount_amount"`
	TaxAmount int64 `json:"tax_amount"`
	ShippingFee int64 `json:"shipping_fee"`
	FinalAmount int64 `json:"final_amount"`
	Currency string `json:"currency" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=credit_card debit_card paypal alipay wechat cash"`
	PaymentStatus string `json:"payment_status" binding:"required,oneof=pending paid failed refunded"`
	OrderStatus string `json:"order_status" binding:"required,oneof=pending confirmed processing shipped delivered cancelled returned"`
	ShippingMethod string `json:"shipping_method" binding:"required,oneof=standard express overnight"`
	TrackingNumber string `json:"tracking_number" binding:"required"`
	ReceiverName string `json:"receiver_name" binding:"required"`
	ReceiverPhone string `json:"receiver_phone" binding:"required"`
	ReceiverEmail string `json:"receiver_email" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	ReceiverCity string `json:"receiver_city" binding:"required"`
	ReceiverState string `json:"receiver_state" binding:"required"`
	ReceiverCountry string `json:"receiver_country" binding:"required"`
	ReceiverPostalCode string `json:"receiver_postal_code" binding:"required"`
	Notes string `json:"notes" binding:"required"`
	PaidAt *time.Time `json:"paid_at"`
	ShippedAt *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
	RefundAmount int64 `json:"refund_amount"`
	RefundReason string `json:"refund_reason" binding:"required"`
	ItemCount int `json:"item_count"`
	Weight float64 `json:"weight"`
	IsGift bool `json:"is_gift"`
	GiftMessage string `json:"gift_message" binding:"required"`
}

// UpdateOrderRequest 是更新 Order 的请求体。
type UpdateOrderRequest struct {
	UserId *string `json:"user_id,omitempty"`
	OrderNo *string `json:"order_no,omitempty"`
	TotalAmount *int64 `json:"total_amount,omitempty"`
	DiscountAmount *int64 `json:"discount_amount,omitempty"`
	TaxAmount *int64 `json:"tax_amount,omitempty"`
	ShippingFee *int64 `json:"shipping_fee,omitempty"`
	FinalAmount *int64 `json:"final_amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
	PaymentMethod *string `json:"payment_method,omitempty" binding:"omitempty,oneof=credit_card debit_card paypal alipay wechat cash"`
	PaymentStatus *string `json:"payment_status,omitempty" binding:"omitempty,oneof=pending paid failed refunded"`
	OrderStatus *string `json:"order_status,omitempty" binding:"omitempty,oneof=pending confirmed processing shipped delivered cancelled returned"`
	ShippingMethod *string `json:"shipping_method,omitempty" binding:"omitempty,oneof=standard express overnight"`
	TrackingNumber *string `json:"tracking_number,omitempty"`
	ReceiverName *string `json:"receiver_name,omitempty"`
	ReceiverPhone *string `json:"receiver_phone,omitempty"`
	ReceiverEmail *string `json:"receiver_email,omitempty"`
	ReceiverAddress *string `json:"receiver_address,omitempty"`
	ReceiverCity *string `json:"receiver_city,omitempty"`
	ReceiverState *string `json:"receiver_state,omitempty"`
	ReceiverCountry *string `json:"receiver_country,omitempty"`
	ReceiverPostalCode *string `json:"receiver_postal_code,omitempty"`
	Notes *string `json:"notes,omitempty"`
	PaidAt *time.Time `json:"paid_at,omitempty"`
	ShippedAt *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	RefundAmount *int64 `json:"refund_amount,omitempty"`
	RefundReason *string `json:"refund_reason,omitempty"`
	ItemCount *int `json:"item_count,omitempty"`
	Weight *float64 `json:"weight,omitempty"`
	IsGift *bool `json:"is_gift,omitempty"`
	GiftMessage *string `json:"gift_message,omitempty"`
}

// OrderResponse 是 Order 的响应体。
type OrderResponse struct {
	ID        string    `json:"id"`
	UserId string `json:"user_id"`
	OrderNo string `json:"order_no"`
	TotalAmount int64 `json:"total_amount"`
	DiscountAmount int64 `json:"discount_amount"`
	TaxAmount int64 `json:"tax_amount"`
	ShippingFee int64 `json:"shipping_fee"`
	FinalAmount int64 `json:"final_amount"`
	Currency string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	PaymentStatus string `json:"payment_status"`
	OrderStatus string `json:"order_status"`
	ShippingMethod string `json:"shipping_method"`
	TrackingNumber string `json:"tracking_number"`
	ReceiverName string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	ReceiverEmail string `json:"receiver_email"`
	ReceiverAddress string `json:"receiver_address"`
	ReceiverCity string `json:"receiver_city"`
	ReceiverState string `json:"receiver_state"`
	ReceiverCountry string `json:"receiver_country"`
	ReceiverPostalCode string `json:"receiver_postal_code"`
	Notes string `json:"notes"`
	PaidAt *time.Time `json:"paid_at"`
	ShippedAt *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
	RefundAmount int64 `json:"refund_amount"`
	RefundReason string `json:"refund_reason"`
	ItemCount int `json:"item_count"`
	Weight float64 `json:"weight"`
	IsGift bool `json:"is_gift"`
	GiftMessage string `json:"gift_message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToOrderResponse 将实体转换为响应体。
func ToOrderResponse(e *order.Order) OrderResponse {
	return OrderResponse{
		ID:        string(e.ID),
		UserId: e.UserId,
		OrderNo: e.OrderNo,
		TotalAmount: e.TotalAmount,
		DiscountAmount: e.DiscountAmount,
		TaxAmount: e.TaxAmount,
		ShippingFee: e.ShippingFee,
		FinalAmount: e.FinalAmount,
		Currency: e.Currency,
		PaymentMethod: string(e.PaymentMethod),
		PaymentStatus: string(e.PaymentStatus),
		OrderStatus: string(e.OrderStatus),
		ShippingMethod: string(e.ShippingMethod),
		TrackingNumber: e.TrackingNumber,
		ReceiverName: e.ReceiverName,
		ReceiverPhone: e.ReceiverPhone,
		ReceiverEmail: e.ReceiverEmail,
		ReceiverAddress: e.ReceiverAddress,
		ReceiverCity: e.ReceiverCity,
		ReceiverState: e.ReceiverState,
		ReceiverCountry: e.ReceiverCountry,
		ReceiverPostalCode: e.ReceiverPostalCode,
		Notes: e.Notes,
		PaidAt: e.PaidAt,
		ShippedAt: e.ShippedAt,
		DeliveredAt: e.DeliveredAt,
		CancelledAt: e.CancelledAt,
		RefundAmount: e.RefundAmount,
		RefundReason: e.RefundReason,
		ItemCount: e.ItemCount,
		Weight: e.Weight,
		IsGift: e.IsGift,
		GiftMessage: e.GiftMessage,
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
