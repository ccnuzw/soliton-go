package shippingapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/shipping"
)

// CreateShippingRequest 是创建 Shipping 的请求体。
type CreateShippingRequest struct {
	OrderId string `json:"order_id" binding:"required"`
	Carrier string `json:"carrier" binding:"required"`
	ShippingMethod string `json:"shipping_method" binding:"required,oneof=standard express overnight"`
	TrackingNumber string `json:"tracking_number" binding:"required"`
	Status string `json:"status" binding:"required,oneof=pending label_created in_transit delivered returned cancelled"`
	ShippedAt *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	ReceiverName string `json:"receiver_name" binding:"required"`
	ReceiverPhone string `json:"receiver_phone" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	ReceiverCity string `json:"receiver_city" binding:"required"`
	ReceiverState string `json:"receiver_state" binding:"required"`
	ReceiverCountry string `json:"receiver_country" binding:"required"`
	ReceiverPostalCode string `json:"receiver_postal_code" binding:"required"`
	Notes string `json:"notes" binding:"required"`
}

// UpdateShippingRequest 是更新 Shipping 的请求体。
type UpdateShippingRequest struct {
	OrderId *string `json:"order_id,omitempty"`
	Carrier *string `json:"carrier,omitempty"`
	ShippingMethod *string `json:"shipping_method,omitempty" binding:"omitempty,oneof=standard express overnight"`
	TrackingNumber *string `json:"tracking_number,omitempty"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=pending label_created in_transit delivered returned cancelled"`
	ShippedAt *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	ReceiverName *string `json:"receiver_name,omitempty"`
	ReceiverPhone *string `json:"receiver_phone,omitempty"`
	ReceiverAddress *string `json:"receiver_address,omitempty"`
	ReceiverCity *string `json:"receiver_city,omitempty"`
	ReceiverState *string `json:"receiver_state,omitempty"`
	ReceiverCountry *string `json:"receiver_country,omitempty"`
	ReceiverPostalCode *string `json:"receiver_postal_code,omitempty"`
	Notes *string `json:"notes,omitempty"`
}

// ShippingResponse 是 Shipping 的响应体。
type ShippingResponse struct {
	ID        string    `json:"id"`
	OrderId string `json:"order_id"`
	Carrier string `json:"carrier"`
	ShippingMethod string `json:"shipping_method"`
	TrackingNumber string `json:"tracking_number"`
	Status string `json:"status"`
	ShippedAt *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	ReceiverName string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	ReceiverAddress string `json:"receiver_address"`
	ReceiverCity string `json:"receiver_city"`
	ReceiverState string `json:"receiver_state"`
	ReceiverCountry string `json:"receiver_country"`
	ReceiverPostalCode string `json:"receiver_postal_code"`
	Notes string `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToShippingResponse 将实体转换为响应体。
func ToShippingResponse(e *shipping.Shipping) ShippingResponse {
	return ShippingResponse{
		ID:        string(e.ID),
		OrderId: e.OrderId,
		Carrier: e.Carrier,
		ShippingMethod: string(e.ShippingMethod),
		TrackingNumber: e.TrackingNumber,
		Status: string(e.Status),
		ShippedAt: e.ShippedAt,
		DeliveredAt: e.DeliveredAt,
		ReceiverName: e.ReceiverName,
		ReceiverPhone: e.ReceiverPhone,
		ReceiverAddress: e.ReceiverAddress,
		ReceiverCity: e.ReceiverCity,
		ReceiverState: e.ReceiverState,
		ReceiverCountry: e.ReceiverCountry,
		ReceiverPostalCode: e.ReceiverPostalCode,
		Notes: e.Notes,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToShippingResponseList 将实体列表转换为响应体列表。
func ToShippingResponseList(entities []*shipping.Shipping) []ShippingResponse {
	result := make([]ShippingResponse, len(entities))
	for i, e := range entities {
		result[i] = ToShippingResponse(e)
	}
	return result
}
