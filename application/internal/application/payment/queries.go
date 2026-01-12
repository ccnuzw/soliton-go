package paymentapp

import (
	"context"
	"strings"

	"github.com/soliton-go/application/internal/domain/payment"
)

// GetPaymentQuery 是获取单个 Payment 的查询。
type GetPaymentQuery struct {
	ID string
}

// GetPaymentHandler 处理 GetPaymentQuery。
type GetPaymentHandler struct {
	repo payment.PaymentRepository
}

func NewGetPaymentHandler(repo payment.PaymentRepository) *GetPaymentHandler {
	return &GetPaymentHandler{repo: repo}
}

func (h *GetPaymentHandler) Handle(ctx context.Context, query GetPaymentQuery) (*payment.Payment, error) {
	return h.repo.Find(ctx, payment.PaymentID(query.ID))
}

// ListPaymentsQuery 是分页列表查询。
type ListPaymentsQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// ListPaymentsResult 是分页查询结果。
type ListPaymentsResult struct {
	Items      []*payment.Payment
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListPaymentsHandler 处理 ListPaymentsQuery。
type ListPaymentsHandler struct {
	repo payment.PaymentRepository
}

func NewListPaymentsHandler(repo payment.PaymentRepository) *ListPaymentsHandler {
	return &ListPaymentsHandler{repo: repo}
}

func (h *ListPaymentsHandler) Handle(ctx context.Context, query ListPaymentsQuery) (*ListPaymentsResult, error) {
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
		"order_id": {},
		"user_id": {},
		"amount": {},
		"currency": {},
		"method": {},
		"status": {},
		"provider": {},
		"provider_txn_id": {},
		"paid_at": {},
		"refunded_at": {},
		"failure_reason": {},
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

	return &ListPaymentsResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
