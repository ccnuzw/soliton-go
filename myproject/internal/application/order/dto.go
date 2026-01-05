package orderapp

// CreateOrderRequest is the request for CreateOrder.
type CreateOrderRequest struct {
	// Add your request fields here.
	// Common patterns:
	ID string `json:"id,omitempty"` // Entity ID (for Get/Update/Delete operations)
	// Data   any    `json:"data,omitempty"` // Payload for Create/Update operations
}

// CreateOrderResponse is the response for CreateOrder.
type CreateOrderResponse struct {
	Success bool   `json:"success"`           // Operation success flag
	Message string `json:"message,omitempty"` // Human-readable message
	Data    any    `json:"data,omitempty"`    // Response payload
}

// GetOrderRequest is the request for GetOrder.
type GetOrderRequest struct {
	// Add your request fields here.
	// Common patterns:
	ID string `json:"id,omitempty"` // Entity ID (for Get/Update/Delete operations)
	// Data   any    `json:"data,omitempty"` // Payload for Create/Update operations
}

// GetOrderResponse is the response for GetOrder.
type GetOrderResponse struct {
	Success bool   `json:"success"`           // Operation success flag
	Message string `json:"message,omitempty"` // Human-readable message
	Data    any    `json:"data,omitempty"`    // Response payload
}

// ListOrdersRequest is the request for ListOrders.
type ListOrdersRequest struct {
	// Add your request fields here.
	// Common patterns:
	ID string `json:"id,omitempty"` // Entity ID (for Get/Update/Delete operations)
	// Data   any    `json:"data,omitempty"` // Payload for Create/Update operations
}

// ListOrdersResponse is the response for ListOrders.
type ListOrdersResponse struct {
	Success bool   `json:"success"`           // Operation success flag
	Message string `json:"message,omitempty"` // Human-readable message
	Data    any    `json:"data,omitempty"`    // Response payload
}
