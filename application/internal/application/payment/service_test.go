package paymentapp

import (
	"context"
	"testing"

	"github.com/soliton-go/application/internal/domain/payment"
)

type paymentRepoStub struct {
	items map[string]*payment.Payment
}

func newPaymentRepoStub() *paymentRepoStub {
	return &paymentRepoStub{items: map[string]*payment.Payment{}}
}

func (r *paymentRepoStub) Find(ctx context.Context, id payment.PaymentID) (*payment.Payment, error) {
	item, ok := r.items[id.String()]
	if !ok {
		return nil, errNotFound("payment not found")
	}
	return item, nil
}

func (r *paymentRepoStub) FindAll(ctx context.Context) ([]*payment.Payment, error) {
	result := make([]*payment.Payment, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	return result, nil
}

func (r *paymentRepoStub) Save(ctx context.Context, entity *payment.Payment) error {
	r.items[entity.ID.String()] = entity
	return nil
}

func (r *paymentRepoStub) Delete(ctx context.Context, id payment.PaymentID) error {
	delete(r.items, id.String())
	return nil
}

func (r *paymentRepoStub) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*payment.Payment, int64, error) {
	items, _ := r.FindAll(ctx)
	return items, int64(len(items)), nil
}

func TestPaymentServiceAuthorizeAndCapture(t *testing.T) {
	repo := newPaymentRepoStub()
	service := NewPaymentService(repo)

	authResp, err := service.AuthorizePayment(context.Background(), AuthorizePaymentServiceRequest{
		OrderId:  "order-1",
		UserId:   "user-1",
		Amount:   100,
		Currency: "CNY",
		Method:   "alipay",
	})
	if err != nil {
		t.Fatalf("authorize failed: %v", err)
	}
	if authResp.PaymentId == "" {
		t.Fatalf("expected payment id")
	}

	captureResp, err := service.CapturePayment(context.Background(), CapturePaymentServiceRequest{
		PaymentId: authResp.PaymentId,
		Amount:    100,
	})
	if err != nil {
		t.Fatalf("capture failed: %v", err)
	}
	if captureResp.Status != string(payment.PaymentStatusPaid) {
		t.Fatalf("expected status paid, got %s", captureResp.Status)
	}
	if captureResp.PaidAt == nil {
		t.Fatalf("expected paid_at to be set")
	}
}

func TestPaymentServiceRefundAndCancel(t *testing.T) {
	repo := newPaymentRepoStub()
	service := NewPaymentService(repo)

	entity := payment.NewPayment(
		"pay-1",
		"order-1",
		"user-1",
		100,
		"CNY",
		payment.PaymentMethodWechat,
		payment.PaymentStatusPaid,
		"",
		"",
		nil,
		nil,
		"",
		nil,
	)
	if err := repo.Save(context.Background(), entity); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	refundResp, err := service.RefundPayment(context.Background(), RefundPaymentServiceRequest{
		PaymentId:    "pay-1",
		RefundAmount: 50,
		Reason:       "test",
	})
	if err != nil {
		t.Fatalf("refund failed: %v", err)
	}
	if refundResp.Status != string(payment.PaymentStatusRefunded) {
		t.Fatalf("expected status refunded, got %s", refundResp.Status)
	}

	_, err = service.CancelPayment(context.Background(), CancelPaymentServiceRequest{
		PaymentId: "pay-1",
	})
	if err == nil {
		t.Fatalf("expected cancel to fail for refunded payment")
	}
}

type errNotFound string

func (e errNotFound) Error() string {
	return string(e)
}
