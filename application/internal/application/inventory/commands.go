package inventoryapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/inventory"
)

// CreateInventoryCommand 是创建 Inventory 的命令。
type CreateInventoryCommand struct {
	ID string
	ProductId string
	WarehouseId string
	LocationCode string
	Stock int
	ReservedStock int
	AvailableStock int
	SafetyStock int
	RestockLevel int
	Status inventory.InventoryStatus
	LastStockedAt *time.Time
	LastCheckedAt *time.Time
	Notes string
	Metadata datatypes.JSON
}

// CreateInventoryHandler 处理 CreateInventoryCommand。
type CreateInventoryHandler struct {
	repo inventory.InventoryRepository
	service *inventory.InventoryDomainService
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreateInventoryHandler(repo inventory.InventoryRepository, service *inventory.InventoryDomainService) *CreateInventoryHandler {
	return &CreateInventoryHandler{repo: repo, service: service}
}

func (h *CreateInventoryHandler) Handle(ctx context.Context, cmd CreateInventoryCommand) (*inventory.Inventory, error) {
	entity := inventory.NewInventory(cmd.ID, cmd.ProductId, cmd.WarehouseId, cmd.LocationCode, cmd.Stock, cmd.ReservedStock, cmd.AvailableStock, cmd.SafetyStock, cmd.RestockLevel, cmd.Status, cmd.LastStockedAt, cmd.LastCheckedAt, cmd.Notes, cmd.Metadata)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// 可选：发布领域事件
	// 取消注释以启用事件发布：
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateInventoryCommand 是更新 Inventory 的命令。
type UpdateInventoryCommand struct {
	ID string
	ProductId *string
	WarehouseId *string
	LocationCode *string
	Stock *int
	ReservedStock *int
	AvailableStock *int
	SafetyStock *int
	RestockLevel *int
	Status *inventory.InventoryStatus
	LastStockedAt *time.Time
	LastCheckedAt *time.Time
	Notes *string
	Metadata *datatypes.JSON
}

// UpdateInventoryHandler 处理 UpdateInventoryCommand。
type UpdateInventoryHandler struct {
	repo inventory.InventoryRepository
	service *inventory.InventoryDomainService
}

func NewUpdateInventoryHandler(repo inventory.InventoryRepository, service *inventory.InventoryDomainService) *UpdateInventoryHandler {
	return &UpdateInventoryHandler{repo: repo, service: service}
}

func (h *UpdateInventoryHandler) Handle(ctx context.Context, cmd UpdateInventoryCommand) (*inventory.Inventory, error) {
	entity, err := h.repo.Find(ctx, inventory.InventoryID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.ProductId, cmd.WarehouseId, cmd.LocationCode, cmd.Stock, cmd.ReservedStock, cmd.AvailableStock, cmd.SafetyStock, cmd.RestockLevel, cmd.Status, cmd.LastStockedAt, cmd.LastCheckedAt, cmd.Notes, cmd.Metadata)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteInventoryCommand 是删除 Inventory 的命令。
type DeleteInventoryCommand struct {
	ID string
}

// DeleteInventoryHandler 处理 DeleteInventoryCommand。
type DeleteInventoryHandler struct {
	repo inventory.InventoryRepository
	service *inventory.InventoryDomainService
}

func NewDeleteInventoryHandler(repo inventory.InventoryRepository, service *inventory.InventoryDomainService) *DeleteInventoryHandler {
	return &DeleteInventoryHandler{repo: repo, service: service}
}

func (h *DeleteInventoryHandler) Handle(ctx context.Context, cmd DeleteInventoryCommand) error {
	return h.repo.Delete(ctx, inventory.InventoryID(cmd.ID))
}
