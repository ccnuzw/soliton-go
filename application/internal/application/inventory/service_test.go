package inventoryapp

import (
	"context"
	"testing"
	"time"

	"github.com/soliton-go/application/internal/domain/inventory"
)

type inventoryRepoStub struct {
	items map[string]*inventory.Inventory
}

func newInventoryRepoStub() *inventoryRepoStub {
	return &inventoryRepoStub{items: map[string]*inventory.Inventory{}}
}

func (r *inventoryRepoStub) Find(ctx context.Context, id inventory.InventoryID) (*inventory.Inventory, error) {
	item, ok := r.items[id.String()]
	if !ok {
		return nil, errNotFound("inventory not found")
	}
	return item, nil
}

func (r *inventoryRepoStub) FindAll(ctx context.Context) ([]*inventory.Inventory, error) {
	result := make([]*inventory.Inventory, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	return result, nil
}

func (r *inventoryRepoStub) Save(ctx context.Context, entity *inventory.Inventory) error {
	r.items[entity.ID.String()] = entity
	return nil
}

func (r *inventoryRepoStub) Delete(ctx context.Context, id inventory.InventoryID) error {
	delete(r.items, id.String())
	return nil
}

func (r *inventoryRepoStub) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*inventory.Inventory, int64, error) {
	items, _ := r.FindAll(ctx)
	return items, int64(len(items)), nil
}

func TestInventoryServiceReserveAndRelease(t *testing.T) {
	repo := newInventoryRepoStub()
	service := NewInventoryService(repo)

	entity := inventory.NewInventory(
		"inv-1",
		"prod-1",
		"wh-1",
		"A1",
		10,
		2,
		8,
		0,
		0,
		inventory.InventoryStatusActive,
		nil,
		nil,
		"",
		nil,
	)
	if err := repo.Save(context.Background(), entity); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	reserveResp, err := service.ReserveStock(context.Background(), ReserveStockServiceRequest{
		InventoryId: "inv-1",
		Quantity:    3,
	})
	if err != nil {
		t.Fatalf("reserve failed: %v", err)
	}
	if reserveResp.ReservedStock != 5 || reserveResp.AvailableStock != 5 {
		t.Fatalf("unexpected stock after reserve: %+v", reserveResp)
	}

	releaseResp, err := service.ReleaseStock(context.Background(), ReleaseStockServiceRequest{
		InventoryId: "inv-1",
		Quantity:    2,
	})
	if err != nil {
		t.Fatalf("release failed: %v", err)
	}
	if releaseResp.ReservedStock != 3 || releaseResp.AvailableStock != 7 {
		t.Fatalf("unexpected stock after release: %+v", releaseResp)
	}
}

func TestInventoryServiceStockInOut(t *testing.T) {
	repo := newInventoryRepoStub()
	service := NewInventoryService(repo)

	now := time.Now()
	entity := inventory.NewInventory(
		"inv-2",
		"prod-2",
		"wh-1",
		"B1",
		5,
		0,
		5,
		0,
		0,
		inventory.InventoryStatusActive,
		&now,
		nil,
		"",
		nil,
	)
	if err := repo.Save(context.Background(), entity); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	stockOutResp, err := service.StockOut(context.Background(), StockOutServiceRequest{
		InventoryId: "inv-2",
		Quantity:    3,
	})
	if err != nil {
		t.Fatalf("stock out failed: %v", err)
	}
	if stockOutResp.Stock != 2 || stockOutResp.AvailableStock != 2 {
		t.Fatalf("unexpected stock after stock out: %+v", stockOutResp)
	}

	stockInResp, err := service.StockIn(context.Background(), StockInServiceRequest{
		InventoryId: "inv-2",
		Quantity:    4,
	})
	if err != nil {
		t.Fatalf("stock in failed: %v", err)
	}
	if stockInResp.Stock != 6 || stockInResp.AvailableStock != 6 {
		t.Fatalf("unexpected stock after stock in: %+v", stockInResp)
	}
}

type errNotFound string

func (e errNotFound) Error() string {
	return string(e)
}
