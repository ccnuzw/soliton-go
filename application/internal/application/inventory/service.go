package inventoryapp

import (
	"context"
	"errors"

	"github.com/soliton-go/application/internal/domain/inventory"
)

// ServiceRemark: 库存服务

// InventoryService 处理跨领域的业务逻辑编排。
type InventoryService struct {
	repo inventory.InventoryRepository
}

// NewInventoryService 创建 InventoryService 实例。
func NewInventoryService(repo inventory.InventoryRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (s *InventoryService) loadInventory(ctx context.Context, id string) (*inventory.Inventory, error) {
	if id == "" {
		return nil, errors.New("inventory_id is required")
	}
	return s.repo.Find(ctx, inventory.InventoryID(id))
}

// AdjustStock 实现 AdjustStock 用例。
func (s *InventoryService) AdjustStock(ctx context.Context, req AdjustStockServiceRequest) (*AdjustStockServiceResponse, error) {
	entity, err := s.loadInventory(ctx, req.InventoryId)
	if err != nil {
		return nil, err
	}
	if req.Delta == 0 {
		return &AdjustStockServiceResponse{
			InventoryId:    entity.ID.String(),
			Stock:          entity.Stock,
			ReservedStock:  entity.ReservedStock,
			AvailableStock: entity.AvailableStock,
		}, nil
	}
	newStock := entity.Stock + req.Delta
	if newStock < 0 {
		return nil, errors.New("insufficient stock")
	}
	if entity.ReservedStock > newStock {
		return nil, errors.New("reserved stock exceeds total stock")
	}
	entity.Stock = newStock
	entity.AvailableStock = entity.Stock - entity.ReservedStock
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &AdjustStockServiceResponse{
		InventoryId:    entity.ID.String(),
		Stock:          entity.Stock,
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// ReserveStock 实现 ReserveStock 用例。
func (s *InventoryService) ReserveStock(ctx context.Context, req ReserveStockServiceRequest) (*ReserveStockServiceResponse, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	entity, err := s.loadInventory(ctx, req.InventoryId)
	if err != nil {
		return nil, err
	}
	reserved := entity.ReservedStock + req.Quantity
	if reserved > entity.Stock {
		return nil, errors.New("insufficient stock")
	}
	entity.ReservedStock = reserved
	entity.AvailableStock = entity.Stock - entity.ReservedStock
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &ReserveStockServiceResponse{
		InventoryId:    entity.ID.String(),
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// ReleaseStock 实现 ReleaseStock 用例。
func (s *InventoryService) ReleaseStock(ctx context.Context, req ReleaseStockServiceRequest) (*ReleaseStockServiceResponse, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	entity, err := s.loadInventory(ctx, req.InventoryId)
	if err != nil {
		return nil, err
	}
	if req.Quantity > entity.ReservedStock {
		return nil, errors.New("release quantity exceeds reserved stock")
	}
	entity.ReservedStock -= req.Quantity
	entity.AvailableStock = entity.Stock - entity.ReservedStock
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &ReleaseStockServiceResponse{
		InventoryId:    entity.ID.String(),
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// StockIn 实现 StockIn 用例。
func (s *InventoryService) StockIn(ctx context.Context, req StockInServiceRequest) (*StockInServiceResponse, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	entity, err := s.loadInventory(ctx, req.InventoryId)
	if err != nil {
		return nil, err
	}
	entity.Stock += req.Quantity
	if entity.ReservedStock > entity.Stock {
		return nil, errors.New("reserved stock exceeds total stock")
	}
	entity.AvailableStock = entity.Stock - entity.ReservedStock
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &StockInServiceResponse{
		InventoryId:    entity.ID.String(),
		Stock:          entity.Stock,
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}

// StockOut 实现 StockOut 用例。
func (s *InventoryService) StockOut(ctx context.Context, req StockOutServiceRequest) (*StockOutServiceResponse, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	entity, err := s.loadInventory(ctx, req.InventoryId)
	if err != nil {
		return nil, err
	}
	newStock := entity.Stock - req.Quantity
	if newStock < 0 {
		return nil, errors.New("insufficient stock")
	}
	if entity.ReservedStock > newStock {
		return nil, errors.New("reserved stock exceeds total stock")
	}
	entity.Stock = newStock
	entity.AvailableStock = entity.Stock - entity.ReservedStock
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &StockOutServiceResponse{
		InventoryId:    entity.ID.String(),
		Stock:          entity.Stock,
		ReservedStock:  entity.ReservedStock,
		AvailableStock: entity.AvailableStock,
	}, nil
}
