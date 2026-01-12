package paymentapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/payment"
)

// CreatePaymentRequest 是创建 Payment 的请求体。
type CreatePaymentRequest struct {
	OrderId string `json:"order_id" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency" binding:"required"`
	Method string `json:"method" binding:"required,oneof=credit_card debit_card paypal alipay wechat cash bank_transfer"`
	Status string `json:"status" binding:"required,oneof=pending authorized paid failed refunded cancelled"`
	Provider string `json:"provider" binding:"required"`
	ProviderTxnId string `json:"provider_txn_id" binding:"required"`
	PaidAt *time.Time `json:"paid_at"`
	RefundedAt *time.Time `json:"refunded_at"`
	FailureReason string `json:"failure_reason" binding:"required"`
	Metadata datatypes.JSON `json:"metadata"`
}

// UpdatePaymentRequest 是更新 Payment 的请求体。
type UpdatePaymentRequest struct {
	OrderId *string `json:"order_id,omitempty"`
	UserId *string `json:"user_id,omitempty"`
	Amount *float64 `json:"amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Method *string `json:"method,omitempty" binding:"omitempty,oneof=credit_card debit_card paypal alipay wechat cash bank_transfer"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=pending authorized paid failed refunded cancelled"`
	Provider *string `json:"provider,omitempty"`
	ProviderTxnId *string `json:"provider_txn_id,omitempty"`
	PaidAt *time.Time `json:"paid_at,omitempty"`
	RefundedAt *time.Time `json:"refunded_at,omitempty"`
	FailureReason *string `json:"failure_reason,omitempty"`
	Metadata *datatypes.JSON `json:"metadata,omitempty"`
}

// PaymentResponse 是 Payment 的响应体。
type PaymentResponse struct {
	ID        string    `json:"id"`
	OrderId string `json:"order_id"`
	UserId string `json:"user_id"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	Method string `json:"method"`
	Status string `json:"status"`
	Provider string `json:"provider"`
	ProviderTxnId string `json:"provider_txn_id"`
	PaidAt *time.Time `json:"paid_at"`
	RefundedAt *time.Time `json:"refunded_at"`
	FailureReason string `json:"failure_reason"`
	Metadata datatypes.JSON `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToPaymentResponse 将实体转换为响应体。
func ToPaymentResponse(e *payment.Payment) PaymentResponse {
	return PaymentResponse{
		ID:        string(e.ID),
		OrderId: e.OrderId,
		UserId: e.UserId,
		Amount: e.Amount,
		Currency: e.Currency,
		Method: string(e.Method),
		Status: string(e.Status),
		Provider: e.Provider,
		ProviderTxnId: e.ProviderTxnId,
		PaidAt: e.PaidAt,
		RefundedAt: e.RefundedAt,
		FailureReason: e.FailureReason,
		Metadata: e.Metadata,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToPaymentResponseList 将实体列表转换为响应体列表。
func ToPaymentResponseList(entities []*payment.Payment) []PaymentResponse {
	result := make([]PaymentResponse, len(entities))
	for i, e := range entities {
		result[i] = ToPaymentResponse(e)
	}
	return result
}
