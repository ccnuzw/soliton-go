package inventoryapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/inventory"
)

// GetInventoryQuery 是获取单个 Inventory 的查询。
type GetInventoryQuery struct {
	ID string
}

// GetInventoryHandler 处理 GetInventoryQuery。
type GetInventoryHandler struct {
	repo inventory.InventoryRepository
}

func NewGetInventoryHandler(repo inventory.InventoryRepository) *GetInventoryHandler {
	return &GetInventoryHandler{repo: repo}
}

func (h *GetInventoryHandler) Handle(ctx context.Context, query GetInventoryQuery) (*inventory.Inventory, error) {
	return h.repo.Find(ctx, inventory.InventoryID(query.ID))
}

// ListInventorysQuery 是分页列表查询。
type ListInventorysQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListInventorysResult 是分页查询结果。
type ListInventorysResult struct {
	Items      []*inventory.Inventory
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListInventorysHandler 处理 ListInventorysQuery。
type ListInventorysHandler struct {
	repo inventory.InventoryRepository
}

func NewListInventorysHandler(repo inventory.InventoryRepository) *ListInventorysHandler {
	return &ListInventorysHandler{repo: repo}
}

func (h *ListInventorysHandler) Handle(ctx context.Context, query ListInventorysQuery) (*ListInventorysResult, error) {
	// 规范化分页参数
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 排序字段白名单
	sortBy := strings.ToLower(strings.TrimSpace(query.SortBy))
	sortOrder := strings.ToLower(strings.TrimSpace(query.SortOrder))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	allowedSorts := map[string]struct{}{
		"id":         {},
		"created_at": {},
		"updated_at": {},
		"product_id": {},
		"warehouse_id": {},
		"location_code": {},
		"stock": {},
		"reserved_stock": {},
		"available_stock": {},
		"safety_stock": {},
		"restock_level": {},
		"status": {},
		"last_stocked_at": {},
		"last_checked_at": {},
		"notes": {},
		"metadata": {},
	}
	if _, ok := allowedSorts[sortBy]; !ok {
		sortBy = "id"
	}

	// 获取总数和分页数据
	items, total, err := h.repo.FindPaginated(ctx, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &ListInventorysResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
