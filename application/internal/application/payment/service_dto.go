package paymentapp

import "time"

// AuthorizePaymentServiceRequest 是 AuthorizePayment 方法的请求参数。
type AuthorizePaymentServiceRequest struct {
	OrderId  string  `json:"order_id"`
	UserId   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Method   string  `json:"method"`
	Provider string  `json:"provider,omitempty"`
}

// AuthorizePaymentServiceResponse 是 AuthorizePayment 方法的响应结果。
type AuthorizePaymentServiceResponse struct {
	PaymentId string `json:"payment_id"`
	Status    string `json:"status"`
}

// CapturePaymentServiceRequest 是 CapturePayment 方法的请求参数。
type CapturePaymentServiceRequest struct {
	PaymentId string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
}

// CapturePaymentServiceResponse 是 CapturePayment 方法的响应结果。
type CapturePaymentServiceResponse struct {
	PaymentId string     `json:"payment_id"`
	Status    string     `json:"status"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
}

// RefundPaymentServiceRequest 是 RefundPayment 方法的请求参数。
type RefundPaymentServiceRequest struct {
	PaymentId    string  `json:"payment_id"`
	RefundAmount float64 `json:"refund_amount"`
	Reason       string  `json:"reason,omitempty"`
}

// RefundPaymentServiceResponse 是 RefundPayment 方法的响应结果。
type RefundPaymentServiceResponse struct {
	PaymentId  string     `json:"payment_id"`
	Status     string     `json:"status"`
	RefundedAt *time.Time `json:"refunded_at,omitempty"`
}

// CancelPaymentServiceRequest 是 CancelPayment 方法的请求参数。
type CancelPaymentServiceRequest struct {
	PaymentId string `json:"payment_id"`
}

// CancelPaymentServiceResponse 是 CancelPayment 方法的响应结果。
type CancelPaymentServiceResponse struct {
	PaymentId string `json:"payment_id"`
	Status    string `json:"status"`
}
