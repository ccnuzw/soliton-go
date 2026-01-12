package inventory

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// InventoryID 是强类型的实体标识符。
type InventoryID string

func (id InventoryID) String() string {
	return string(id)
}

// InventoryStatus 表示 Status 字段的枚举类型。
type InventoryStatus string

const (
	InventoryStatusActive    InventoryStatus = "active"
	InventoryStatusInactive  InventoryStatus = "inactive"
	InventoryStatusSuspended InventoryStatus = "suspended"
)

// Inventory 是聚合根实体。
type Inventory struct {
	ddd.BaseAggregateRoot
	ID             InventoryID     `gorm:"primaryKey"`
	ProductId      string          `gorm:"size:36;index"`            // 商品ID
	WarehouseId    string          `gorm:"size:36;index"`            // 仓库ID
	LocationCode   string          `gorm:"size:255"`                 // 库位编码
	Stock          int             `gorm:"not null;default:0"`       // 当前库存
	ReservedStock  int             `gorm:"not null;default:0"`       // 预占库存
	AvailableStock int             `gorm:"not null;default:0"`       // 可用库存
	SafetyStock    int             `gorm:"not null;default:0"`       // 安全库存
	RestockLevel   int             `gorm:"not null;default:0"`       // 补货阈值
	Status         InventoryStatus `gorm:"size:50;default:'active'"` // 库存状态
	LastStockedAt  *time.Time      // 最近入库时间
	LastCheckedAt  *time.Time      // 最近盘点时间
	Notes          string          `gorm:"type:text"` // 备注
	Metadata       datatypes.JSON  // 扩展信息
	CreatedAt      time.Time       `gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt  `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Inventory) TableName() string {
	return "inventories"
}

// NewInventory 创建一个新的 Inventory 实体。
func NewInventory(id string, productId string, warehouseId string, locationCode string, stock int, reservedStock int, availableStock int, safetyStock int, restockLevel int, status InventoryStatus, lastStockedAt *time.Time, lastCheckedAt *time.Time, notes string, metadata datatypes.JSON) *Inventory {
	e := &Inventory{
		ID:             InventoryID(id),
		ProductId:      productId,
		WarehouseId:    warehouseId,
		LocationCode:   locationCode,
		Stock:          stock,
		ReservedStock:  reservedStock,
		AvailableStock: availableStock,
		SafetyStock:    safetyStock,
		RestockLevel:   restockLevel,
		Status:         status,
		LastStockedAt:  lastStockedAt,
		LastCheckedAt:  lastCheckedAt,
		Notes:          notes,
		Metadata:       metadata,
	}
	e.AddDomainEvent(NewInventoryCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Inventory) Update(productId *string, warehouseId *string, locationCode *string, stock *int, reservedStock *int, availableStock *int, safetyStock *int, restockLevel *int, status *InventoryStatus, lastStockedAt *time.Time, lastCheckedAt *time.Time, notes *string, metadata *datatypes.JSON) {
	if productId != nil {
		e.ProductId = *productId
	}
	if warehouseId != nil {
		e.WarehouseId = *warehouseId
	}
	if locationCode != nil {
		e.LocationCode = *locationCode
	}
	if stock != nil {
		e.Stock = *stock
	}
	if reservedStock != nil {
		e.ReservedStock = *reservedStock
	}
	if availableStock != nil {
		e.AvailableStock = *availableStock
	}
	if safetyStock != nil {
		e.SafetyStock = *safetyStock
	}
	if restockLevel != nil {
		e.RestockLevel = *restockLevel
	}
	if status != nil {
		e.Status = *status
	}
	if lastStockedAt != nil {
		e.LastStockedAt = lastStockedAt
	}
	if lastCheckedAt != nil {
		e.LastCheckedAt = lastCheckedAt
	}
	if notes != nil {
		e.Notes = *notes
	}
	if metadata != nil {
		e.Metadata = *metadata
	}
	e.AddDomainEvent(NewInventoryUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Inventory) GetID() ddd.ID {
	return e.ID
}
