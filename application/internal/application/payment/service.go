package paymentapp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/soliton-go/application/internal/domain/payment"
)

// PaymentService 处理跨领域的业务逻辑编排。
type PaymentService struct {
	repo payment.PaymentRepository
}

// NewPaymentService 创建 PaymentService 实例。
func NewPaymentService(
	repo payment.PaymentRepository,
) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

// AuthorizePayment 实现 AuthorizePayment 用例。
func (s *PaymentService) AuthorizePayment(ctx context.Context, req AuthorizePaymentServiceRequest) (*AuthorizePaymentServiceResponse, error) {
	if req.OrderId == "" {
		return nil, fmt.Errorf("order_id is required")
	}
	if req.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if req.Method == "" {
		return nil, fmt.Errorf("method is required")
	}

	entity := payment.NewPayment(
		uuid.New().String(),
		req.OrderId,
		req.UserId,
		req.Amount,
		req.Currency,
		payment.PaymentMethod(req.Method),
		payment.PaymentStatusAuthorized,
		req.Provider,
		req.ProviderTxnId,
		nil,
		nil,
		"",
		req.Metadata,
	)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	authorizedAt := time.Now()
	return &AuthorizePaymentServiceResponse{
		Success:      true,
		Message:      "authorized",
		PaymentId:    string(entity.ID),
		Status:       string(entity.Status),
		AuthorizedAt: &authorizedAt,
	}, nil
}

// CapturePayment 实现 CapturePayment 用例。
func (s *PaymentService) CapturePayment(ctx context.Context, req CapturePaymentServiceRequest) (*CapturePaymentServiceResponse, error) {
	if req.PaymentId == "" {
		return nil, fmt.Errorf("payment_id is required")
	}

	entity, err := s.repo.Find(ctx, payment.PaymentID(req.PaymentId))
	if err != nil {
		return nil, err
	}
	if entity.Status != payment.PaymentStatusAuthorized && entity.Status != payment.PaymentStatusPending {
		return nil, fmt.Errorf("payment status must be authorized or pending")
	}

	var amount *float64
	if req.Amount > 0 {
		amount = &req.Amount
	}
	status := payment.PaymentStatusPaid
	paidAt := req.PaidAt
	if paidAt == nil {
		now := time.Now()
		paidAt = &now
	}
	var providerTxnId *string
	if req.ProviderTxnId != "" {
		providerTxnId = &req.ProviderTxnId
	}

	entity.Update(nil, nil, amount, nil, nil, &status, nil, providerTxnId, paidAt, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &CapturePaymentServiceResponse{
		Success:   true,
		Message:   "captured",
		PaymentId: string(entity.ID),
		Status:    string(entity.Status),
		PaidAt:    entity.PaidAt,
	}, nil
}

// RefundPayment 实现 RefundPayment 用例。
func (s *PaymentService) RefundPayment(ctx context.Context, req RefundPaymentServiceRequest) (*RefundPaymentServiceResponse, error) {
	if req.PaymentId == "" {
		return nil, fmt.Errorf("payment_id is required")
	}
	if req.RefundAmount <= 0 {
		return nil, fmt.Errorf("refund_amount must be greater than 0")
	}

	entity, err := s.repo.Find(ctx, payment.PaymentID(req.PaymentId))
	if err != nil {
		return nil, err
	}
	if entity.Status != payment.PaymentStatusPaid {
		return nil, fmt.Errorf("payment status must be paid")
	}

	status := payment.PaymentStatusRefunded
	refundedAt := req.RefundedAt
	if refundedAt == nil {
		now := time.Now()
		refundedAt = &now
	}
	var reason *string
	if req.Reason != "" {
		reason = &req.Reason
	}
	var providerTxnId *string
	if req.ProviderTxnId != "" {
		providerTxnId = &req.ProviderTxnId
	}

	entity.Update(nil, nil, nil, nil, nil, &status, nil, providerTxnId, nil, refundedAt, reason, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &RefundPaymentServiceResponse{
		Success:    true,
		Message:    "refunded",
		PaymentId:  string(entity.ID),
		Status:     string(entity.Status),
		RefundedAt: entity.RefundedAt,
	}, nil
}

// CancelPayment 实现 CancelPayment 用例。
func (s *PaymentService) CancelPayment(ctx context.Context, req CancelPaymentServiceRequest) (*CancelPaymentServiceResponse, error) {
	if req.PaymentId == "" {
		return nil, fmt.Errorf("payment_id is required")
	}

	entity, err := s.repo.Find(ctx, payment.PaymentID(req.PaymentId))
	if err != nil {
		return nil, err
	}
	if entity.Status == payment.PaymentStatusPaid || entity.Status == payment.PaymentStatusRefunded {
		return nil, fmt.Errorf("payment cannot be cancelled in current status")
	}
	if entity.Status == payment.PaymentStatusCancelled {
		return &CancelPaymentServiceResponse{
			Success:   true,
			Message:   "already cancelled",
			PaymentId: string(entity.ID),
			Status:    string(entity.Status),
		}, nil
	}

	status := payment.PaymentStatusCancelled
	var reason *string
	if req.Reason != "" {
		reason = &req.Reason
	}

	entity.Update(nil, nil, nil, nil, nil, &status, nil, nil, nil, nil, reason, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &CancelPaymentServiceResponse{
		Success:   true,
		Message:   "cancelled",
		PaymentId: string(entity.ID),
		Status:    string(entity.Status),
	}, nil
}
