package inventoryapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/inventory"
)

// CreateInventoryRequest 是创建 Inventory 的请求体。
type CreateInventoryRequest struct {
	ProductId string `json:"product_id" binding:"required"`
	WarehouseId string `json:"warehouse_id" binding:"required"`
	LocationCode string `json:"location_code" binding:"required"`
	Stock int `json:"stock"`
	ReservedStock int `json:"reserved_stock"`
	AvailableStock int `json:"available_stock"`
	SafetyStock int `json:"safety_stock"`
	RestockLevel int `json:"restock_level"`
	Status string `json:"status" binding:"required,oneof=active inactive suspended"`
	LastStockedAt *time.Time `json:"last_stocked_at"`
	LastCheckedAt *time.Time `json:"last_checked_at"`
	Notes string `json:"notes" binding:"required"`
	Metadata datatypes.JSON `json:"metadata"`
}

// UpdateInventoryRequest 是更新 Inventory 的请求体。
type UpdateInventoryRequest struct {
	ProductId *string `json:"product_id,omitempty"`
	WarehouseId *string `json:"warehouse_id,omitempty"`
	LocationCode *string `json:"location_code,omitempty"`
	Stock *int `json:"stock,omitempty"`
	ReservedStock *int `json:"reserved_stock,omitempty"`
	AvailableStock *int `json:"available_stock,omitempty"`
	SafetyStock *int `json:"safety_stock,omitempty"`
	RestockLevel *int `json:"restock_level,omitempty"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive suspended"`
	LastStockedAt *time.Time `json:"last_stocked_at,omitempty"`
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty"`
	Notes *string `json:"notes,omitempty"`
	Metadata *datatypes.JSON `json:"metadata,omitempty"`
}

// InventoryResponse 是 Inventory 的响应体。
type InventoryResponse struct {
	ID        string    `json:"id"`
	ProductId string `json:"product_id"`
	WarehouseId string `json:"warehouse_id"`
	LocationCode string `json:"location_code"`
	Stock int `json:"stock"`
	ReservedStock int `json:"reserved_stock"`
	AvailableStock int `json:"available_stock"`
	SafetyStock int `json:"safety_stock"`
	RestockLevel int `json:"restock_level"`
	Status string `json:"status"`
	LastStockedAt *time.Time `json:"last_stocked_at"`
	LastCheckedAt *time.Time `json:"last_checked_at"`
	Notes string `json:"notes"`
	Metadata datatypes.JSON `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToInventoryResponse 将实体转换为响应体。
func ToInventoryResponse(e *inventory.Inventory) InventoryResponse {
	return InventoryResponse{
		ID:        string(e.ID),
		ProductId: e.ProductId,
		WarehouseId: e.WarehouseId,
		LocationCode: e.LocationCode,
		Stock: e.Stock,
		ReservedStock: e.ReservedStock,
		AvailableStock: e.AvailableStock,
		SafetyStock: e.SafetyStock,
		RestockLevel: e.RestockLevel,
		Status: string(e.Status),
		LastStockedAt: e.LastStockedAt,
		LastCheckedAt: e.LastCheckedAt,
		Notes: e.Notes,
		Metadata: e.Metadata,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToInventoryResponseList 将实体列表转换为响应体列表。
func ToInventoryResponseList(entities []*inventory.Inventory) []InventoryResponse {
	result := make([]InventoryResponse, len(entities))
	for i, e := range entities {
		result[i] = ToInventoryResponse(e)
	}
	return result
}
