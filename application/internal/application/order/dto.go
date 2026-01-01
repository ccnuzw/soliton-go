package orderapp

import "github.com/soliton-go/application/internal/domain/order"

// CreateOrderRequest is the request body for creating a Order.
type CreateOrderRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateOrderRequest is the request body for updating a Order.
type UpdateOrderRequest struct {
	Name string `json:"name" binding:"required"`
}

// OrderResponse is the response body for Order data.
type OrderResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToOrderResponse converts entity to response.
func ToOrderResponse(e *order.Order) OrderResponse {
	return OrderResponse{
		ID:   string(e.ID),
		Name: e.Name,
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
