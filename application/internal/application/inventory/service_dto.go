package inventoryapp

// AdjustStockServiceRequest 是 AdjustStock 方法的请求参数。
type AdjustStockServiceRequest struct {
	InventoryId string `json:"inventory_id"`           // 库存ID
	Delta       int    `json:"delta"`                  // 调整数量（正入库，负出库）
	Reason      string `json:"reason,omitempty"`       // 调整原因
	OperatorId  string `json:"operator_id,omitempty"`  // 操作人ID
	ReferenceNo string `json:"reference_no,omitempty"` // 关联单号
}

// AdjustStockServiceResponse 是 AdjustStock 方法的响应结果。
type AdjustStockServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	InventoryId    string `json:"inventory_id,omitempty"`    // 库存ID
	Stock          int    `json:"stock,omitempty"`           // 当前库存
	ReservedStock  int    `json:"reserved_stock,omitempty"`  // 预占库存
	AvailableStock int    `json:"available_stock,omitempty"` // 可用库存
}

// ReserveStockServiceRequest 是 ReserveStock 方法的请求参数。
type ReserveStockServiceRequest struct {
	InventoryId string `json:"inventory_id"`          // 库存ID
	Quantity    int    `json:"quantity"`              // 预占数量
	OrderId     string `json:"order_id,omitempty"`    // 关联订单ID
	Reason      string `json:"reason,omitempty"`      // 预占原因
	OperatorId  string `json:"operator_id,omitempty"` // 操作人ID
}

// ReserveStockServiceResponse 是 ReserveStock 方法的响应结果。
type ReserveStockServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	InventoryId    string `json:"inventory_id,omitempty"`    // 库存ID
	ReservedStock  int    `json:"reserved_stock,omitempty"`  // 预占库存
	AvailableStock int    `json:"available_stock,omitempty"` // 可用库存
}

// ReleaseStockServiceRequest 是 ReleaseStock 方法的请求参数。
type ReleaseStockServiceRequest struct {
	InventoryId string `json:"inventory_id"`          // 库存ID
	Quantity    int    `json:"quantity"`              // 释放数量
	OrderId     string `json:"order_id,omitempty"`    // 关联订单ID
	Reason      string `json:"reason,omitempty"`      // 释放原因
	OperatorId  string `json:"operator_id,omitempty"` // 操作人ID
}

// ReleaseStockServiceResponse 是 ReleaseStock 方法的响应结果。
type ReleaseStockServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	InventoryId    string `json:"inventory_id,omitempty"`    // 库存ID
	ReservedStock  int    `json:"reserved_stock,omitempty"`  // 预占库存
	AvailableStock int    `json:"available_stock,omitempty"` // 可用库存
}

// StockInServiceRequest 是 StockIn 方法的请求参数。
type StockInServiceRequest struct {
	InventoryId string `json:"inventory_id"`           // 库存ID
	Quantity    int    `json:"quantity"`               // 入库数量
	ReferenceNo string `json:"reference_no,omitempty"` // 入库单号
	OperatorId  string `json:"operator_id,omitempty"`  // 操作人ID
	Note        string `json:"note,omitempty"`         // 备注
}

// StockInServiceResponse 是 StockIn 方法的响应结果。
type StockInServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	InventoryId    string `json:"inventory_id,omitempty"`    // 库存ID
	Stock          int    `json:"stock,omitempty"`           // 当前库存
	AvailableStock int    `json:"available_stock,omitempty"` // 可用库存
}

// StockOutServiceRequest 是 StockOut 方法的请求参数。
type StockOutServiceRequest struct {
	InventoryId string `json:"inventory_id"`          // 库存ID
	Quantity    int    `json:"quantity"`              // 出库数量
	OrderId     string `json:"order_id,omitempty"`    // 关联订单ID
	OperatorId  string `json:"operator_id,omitempty"` // 操作人ID
	Note        string `json:"note,omitempty"`        // 备注
}

// StockOutServiceResponse 是 StockOut 方法的响应结果。
type StockOutServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	InventoryId    string `json:"inventory_id,omitempty"`    // 库存ID
	Stock          int    `json:"stock,omitempty"`           // 当前库存
	AvailableStock int    `json:"available_stock,omitempty"` // 可用库存
}
