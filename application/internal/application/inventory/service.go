package inventoryapp

import (
	"context"
	"fmt"
	"time"

	"github.com/soliton-go/application/internal/domain/inventory"
)

// InventoryService 处理跨领域的业务逻辑编排。
type InventoryService struct {
	repo inventory.InventoryRepository
}

// NewInventoryService 创建 InventoryService 实例。
func NewInventoryService(
	repo inventory.InventoryRepository,
) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

// AdjustStock 实现 AdjustStock 用例。
func (s *InventoryService) AdjustStock(ctx context.Context, req AdjustStockServiceRequest) (*AdjustStockServiceResponse, error) {
	if req.InventoryId == "" {
		return nil, fmt.Errorf("inventory_id is required")
	}
	if req.Delta == 0 {
		return nil, fmt.Errorf("delta must not be zero")
	}

	entity, err := s.repo.Find(ctx, inventory.InventoryID(req.InventoryId))
	if err != nil {
		return nil, err
	}

	newStock := entity.Stock + req.Delta
	if newStock < 0 {
		return nil, fmt.Errorf("stock cannot be negative")
	}
	newAvailable := newStock - entity.ReservedStock
	if newAvailable < 0 {
		return nil, fmt.Errorf("available stock cannot be negative")
	}

	entity.Update(nil, nil, nil, &newStock, nil, &newAvailable, nil, nil, nil, nil, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &AdjustStockServiceResponse{
		Success:        true,
		Message:        "adjusted",
		InventoryId:    string(entity.ID),
		Stock:          entity.Stock,
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// ReserveStock 实现 ReserveStock 用例。
func (s *InventoryService) ReserveStock(ctx context.Context, req ReserveStockServiceRequest) (*ReserveStockServiceResponse, error) {
	if req.InventoryId == "" {
		return nil, fmt.Errorf("inventory_id is required")
	}
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	entity, err := s.repo.Find(ctx, inventory.InventoryID(req.InventoryId))
	if err != nil {
		return nil, err
	}
	if entity.AvailableStock < req.Quantity {
		return nil, fmt.Errorf("insufficient available stock")
	}

	newReserved := entity.ReservedStock + req.Quantity
	newAvailable := entity.AvailableStock - req.Quantity
	entity.Update(nil, nil, nil, nil, &newReserved, &newAvailable, nil, nil, nil, nil, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &ReserveStockServiceResponse{
		Success:        true,
		Message:        "reserved",
		InventoryId:    string(entity.ID),
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// ReleaseStock 实现 ReleaseStock 用例。
func (s *InventoryService) ReleaseStock(ctx context.Context, req ReleaseStockServiceRequest) (*ReleaseStockServiceResponse, error) {
	if req.InventoryId == "" {
		return nil, fmt.Errorf("inventory_id is required")
	}
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	entity, err := s.repo.Find(ctx, inventory.InventoryID(req.InventoryId))
	if err != nil {
		return nil, err
	}
	if entity.ReservedStock < req.Quantity {
		return nil, fmt.Errorf("reserved stock is insufficient")
	}

	newReserved := entity.ReservedStock - req.Quantity
	newAvailable := entity.AvailableStock + req.Quantity
	entity.Update(nil, nil, nil, nil, &newReserved, &newAvailable, nil, nil, nil, nil, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &ReleaseStockServiceResponse{
		Success:        true,
		Message:        "released",
		InventoryId:    string(entity.ID),
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// StockIn 实现 StockIn 用例。
func (s *InventoryService) StockIn(ctx context.Context, req StockInServiceRequest) (*StockInServiceResponse, error) {
	if req.InventoryId == "" {
		return nil, fmt.Errorf("inventory_id is required")
	}
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	entity, err := s.repo.Find(ctx, inventory.InventoryID(req.InventoryId))
	if err != nil {
		return nil, err
	}

	newStock := entity.Stock + req.Quantity
	newAvailable := entity.AvailableStock + req.Quantity
	now := time.Now()
	entity.Update(nil, nil, nil, &newStock, nil, &newAvailable, nil, nil, nil, &now, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &StockInServiceResponse{
		Success:        true,
		Message:        "stocked in",
		InventoryId:    string(entity.ID),
		Stock:          entity.Stock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// StockOut 实现 StockOut 用例。
func (s *InventoryService) StockOut(ctx context.Context, req StockOutServiceRequest) (*StockOutServiceResponse, error) {
	if req.InventoryId == "" {
		return nil, fmt.Errorf("inventory_id is required")
	}
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	entity, err := s.repo.Find(ctx, inventory.InventoryID(req.InventoryId))
	if err != nil {
		return nil, err
	}
	if entity.AvailableStock < req.Quantity {
		return nil, fmt.Errorf("insufficient available stock")
	}

	newStock := entity.Stock - req.Quantity
	newAvailable := entity.AvailableStock - req.Quantity
	entity.Update(nil, nil, nil, &newStock, nil, &newAvailable, nil, nil, nil, nil, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &StockOutServiceResponse{
		Success:        true,
		Message:        "stocked out",
		InventoryId:    string(entity.ID),
		Stock:          entity.Stock,
		AvailableStock: entity.AvailableStock,
	}, nil
}
