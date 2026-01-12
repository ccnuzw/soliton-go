package inventoryapp

// AdjustStockServiceRequest 是 AdjustStock 方法的请求参数。
type AdjustStockServiceRequest struct {
	InventoryId string `json:"inventory_id"`
	Delta       int    `json:"delta"`
}

// AdjustStockServiceResponse 是 AdjustStock 方法的响应结果。
type AdjustStockServiceResponse struct {
	InventoryId    string `json:"inventory_id"`
	Stock          int    `json:"stock"`
	ReservedStock  int    `json:"reserved_stock"`
	AvailableStock int    `json:"available_stock"`
}

// ReserveStockServiceRequest 是 ReserveStock 方法的请求参数。
type ReserveStockServiceRequest struct {
	InventoryId string `json:"inventory_id"`
	Quantity    int    `json:"quantity"`
}

// ReserveStockServiceResponse 是 ReserveStock 方法的响应结果。
type ReserveStockServiceResponse struct {
	InventoryId    string `json:"inventory_id"`
	ReservedStock  int    `json:"reserved_stock"`
	AvailableStock int    `json:"available_stock"`
}

// ReleaseStockServiceRequest 是 ReleaseStock 方法的请求参数。
type ReleaseStockServiceRequest struct {
	InventoryId string `json:"inventory_id"`
	Quantity    int    `json:"quantity"`
}

// ReleaseStockServiceResponse 是 ReleaseStock 方法的响应结果。
type ReleaseStockServiceResponse struct {
	InventoryId    string `json:"inventory_id"`
	ReservedStock  int    `json:"reserved_stock"`
	AvailableStock int    `json:"available_stock"`
}

// StockInServiceRequest 是 StockIn 方法的请求参数。
type StockInServiceRequest struct {
	InventoryId string `json:"inventory_id"`
	Quantity    int    `json:"quantity"`
}

// StockInServiceResponse 是 StockIn 方法的响应结果。
type StockInServiceResponse struct {
	InventoryId    string `json:"inventory_id"`
	Stock          int    `json:"stock"`
	ReservedStock  int    `json:"reserved_stock"`
	AvailableStock int    `json:"available_stock"`
}

// StockOutServiceRequest 是 StockOut 方法的请求参数。
type StockOutServiceRequest struct {
	InventoryId string `json:"inventory_id"`
	Quantity    int    `json:"quantity"`
}

// StockOutServiceResponse 是 StockOut 方法的响应结果。
type StockOutServiceResponse struct {
	InventoryId    string `json:"inventory_id"`
	Stock          int    `json:"stock"`
	ReservedStock  int    `json:"reserved_stock"`
	AvailableStock int    `json:"available_stock"`
}
