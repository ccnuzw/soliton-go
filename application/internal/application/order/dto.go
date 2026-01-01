package orderapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/order"
)

// CreateOrderRequest is the request body for creating a Order.
type CreateOrderRequest struct {
	UserId string `json:"user_id"`
	OrderNo string `json:"order_no"`
	TotalAmount int64 `json:"total_amount"`
	Status string `json:"status"`
	ReceiverName string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	ReceiverAddress string `json:"receiver_address"`
}

// UpdateOrderRequest is the request body for updating a Order.
type UpdateOrderRequest struct {
	UserId string `json:"user_id"`
	OrderNo string `json:"order_no"`
	TotalAmount int64 `json:"total_amount"`
	Status string `json:"status"`
	ReceiverName string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	ReceiverAddress string `json:"receiver_address"`
}

// OrderResponse is the response body for Order data.
type OrderResponse struct {
	ID        string    `json:"id"`
	UserId string `json:"user_id"`
	OrderNo string `json:"order_no"`
	TotalAmount int64 `json:"total_amount"`
	Status string `json:"status"`
	ReceiverName string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	ReceiverAddress string `json:"receiver_address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToOrderResponse converts entity to response.
func ToOrderResponse(e *order.Order) OrderResponse {
	return OrderResponse{
		ID:        string(e.ID),
		UserId: e.UserId,
		OrderNo: e.OrderNo,
		TotalAmount: e.TotalAmount,
		Status: string(e.Status),
		ReceiverName: e.ReceiverName,
		ReceiverPhone: e.ReceiverPhone,
		ReceiverAddress: e.ReceiverAddress,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToOrderResponseList converts entities to response list.
func ToOrderResponseList(entities []*order.Order) []OrderResponse {
	result := make([]OrderResponse, len(entities))
	for i, e := range entities {
		result[i] = ToOrderResponse(e)
	}
	return result
}
