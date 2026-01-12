package paymentapp

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soliton-go/application/internal/domain/payment"
)

// ServiceRemark: 支付服务

// PaymentService 处理跨领域的业务逻辑编排。
type PaymentService struct {
	repo payment.PaymentRepository
}

// NewPaymentService 创建 PaymentService 实例。
func NewPaymentService(repo payment.PaymentRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) loadPayment(ctx context.Context, id string) (*payment.Payment, error) {
	if id == "" {
		return nil, errors.New("payment_id is required")
	}
	return s.repo.Find(ctx, payment.PaymentID(id))
}

func parsePaymentMethod(method string) (payment.PaymentMethod, error) {
	switch strings.ToLower(method) {
	case string(payment.PaymentMethodCreditCard):
		return payment.PaymentMethodCreditCard, nil
	case string(payment.PaymentMethodDebitCard):
		return payment.PaymentMethodDebitCard, nil
	case string(payment.PaymentMethodPaypal):
		return payment.PaymentMethodPaypal, nil
	case string(payment.PaymentMethodAlipay):
		return payment.PaymentMethodAlipay, nil
	case string(payment.PaymentMethodWechat):
		return payment.PaymentMethodWechat, nil
	case string(payment.PaymentMethodCash):
		return payment.PaymentMethodCash, nil
	case string(payment.PaymentMethodBankTransfer):
		return payment.PaymentMethodBankTransfer, nil
	default:
		return "", errors.New("unsupported payment method")
	}
}

// AuthorizePayment 实现 AuthorizePayment 用例。
// MethodRemark: AuthorizePayment 支付授权
func (s *PaymentService) AuthorizePayment(ctx context.Context, req AuthorizePaymentServiceRequest) (*AuthorizePaymentServiceResponse, error) {
	if req.OrderId == "" || req.UserId == "" {
		return nil, errors.New("order_id and user_id are required")
	}
	method, err := parsePaymentMethod(req.Method)
	if err != nil {
		return nil, err
	}
	paymentID := uuid.NewString()
	entity := payment.NewPayment(
		paymentID,
		req.OrderId,
		req.UserId,
		req.Amount,
		req.Currency,
		method,
		payment.PaymentStatusAuthorized,
		req.Provider,
		"",
		nil,
		nil,
		"",
		nil,
	)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &AuthorizePaymentServiceResponse{
		PaymentId: paymentID,
		Status:    string(entity.Status),
	}, nil
}

// CapturePayment 实现 CapturePayment 用例。
// MethodRemark: CapturePayment 支付扣款
func (s *PaymentService) CapturePayment(ctx context.Context, req CapturePaymentServiceRequest) (*CapturePaymentServiceResponse, error) {
	entity, err := s.loadPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}
	entity.Status = payment.PaymentStatusPaid
	now := time.Now()
	entity.PaidAt = &now
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &CapturePaymentServiceResponse{
		PaymentId: entity.ID.String(),
		Status:    string(entity.Status),
		PaidAt:    entity.PaidAt,
	}, nil
}

// RefundPayment 实现 RefundPayment 用例。
// MethodRemark: RefundPayment 退款处理
func (s *PaymentService) RefundPayment(ctx context.Context, req RefundPaymentServiceRequest) (*RefundPaymentServiceResponse, error) {
	entity, err := s.loadPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}
	if entity.Status != payment.PaymentStatusPaid {
		return nil, errors.New("payment is not in paid status")
	}
	entity.Status = payment.PaymentStatusRefunded
	now := time.Now()
	entity.RefundedAt = &now
	if req.Reason != "" {
		entity.FailureReason = req.Reason
	}
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &RefundPaymentServiceResponse{
		PaymentId:  entity.ID.String(),
		Status:     string(entity.Status),
		RefundedAt: entity.RefundedAt,
	}, nil
}

// CancelPayment 实现 CancelPayment 用例。
// MethodRemark: CancelPayment 取消支付
func (s *PaymentService) CancelPayment(ctx context.Context, req CancelPaymentServiceRequest) (*CancelPaymentServiceResponse, error) {
	entity, err := s.loadPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}
	if entity.Status == payment.PaymentStatusRefunded {
		return nil, errors.New("refunded payment cannot be cancelled")
	}
	entity.Status = payment.PaymentStatusCancelled
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return &CancelPaymentServiceResponse{
		PaymentId: entity.ID.String(),
		Status:    string(entity.Status),
	}, nil
}
