package paymentapp

import (
	"time"

	"gorm.io/datatypes"
)

// AuthorizePaymentServiceRequest 是 AuthorizePayment 方法的请求参数。
type AuthorizePaymentServiceRequest struct {
	OrderId       string         `json:"order_id"`                  // 订单ID
	UserId        string         `json:"user_id"`                   // 用户ID
	Amount        float64        `json:"amount"`                    // 支付金额
	Currency      string         `json:"currency"`                  // 币种
	Method        string         `json:"method"`                    // 支付方式
	Provider      string         `json:"provider,omitempty"`        // 支付渠道
	ProviderTxnId string         `json:"provider_txn_id,omitempty"` // 渠道交易号
	Metadata      datatypes.JSON `json:"metadata,omitempty"`        // 扩展信息
}

// AuthorizePaymentServiceResponse 是 AuthorizePayment 方法的响应结果。
type AuthorizePaymentServiceResponse struct {
	Success      bool       `json:"success"`                 // 操作是否成功
	Message      string     `json:"message,omitempty"`       // 提示消息
	PaymentId    string     `json:"payment_id,omitempty"`    // 支付ID
	Status       string     `json:"status,omitempty"`        // 支付状态
	AuthorizedAt *time.Time `json:"authorized_at,omitempty"` // 授权时间
}

// CapturePaymentServiceRequest 是 CapturePayment 方法的请求参数。
type CapturePaymentServiceRequest struct {
	PaymentId     string     `json:"payment_id"`                // 支付ID
	Amount        float64    `json:"amount,omitempty"`          // 实际扣款金额
	ProviderTxnId string     `json:"provider_txn_id,omitempty"` // 渠道交易号
	PaidAt        *time.Time `json:"paid_at,omitempty"`         // 支付完成时间
}

// CapturePaymentServiceResponse 是 CapturePayment 方法的响应结果。
type CapturePaymentServiceResponse struct {
	Success   bool       `json:"success"`              // 操作是否成功
	Message   string     `json:"message,omitempty"`    // 提示消息
	PaymentId string     `json:"payment_id,omitempty"` // 支付ID
	Status    string     `json:"status,omitempty"`     // 支付状态
	PaidAt    *time.Time `json:"paid_at,omitempty"`    // 支付完成时间
}

// RefundPaymentServiceRequest 是 RefundPayment 方法的请求参数。
type RefundPaymentServiceRequest struct {
	PaymentId     string     `json:"payment_id"`                // 支付ID
	RefundAmount  float64    `json:"refund_amount"`             // 退款金额
	Reason        string     `json:"reason,omitempty"`          // 退款原因
	ProviderTxnId string     `json:"provider_txn_id,omitempty"` // 渠道交易号
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`     // 退款完成时间
}

// RefundPaymentServiceResponse 是 RefundPayment 方法的响应结果。
type RefundPaymentServiceResponse struct {
	Success    bool       `json:"success"`               // 操作是否成功
	Message    string     `json:"message,omitempty"`     // 提示消息
	PaymentId  string     `json:"payment_id,omitempty"`  // 支付ID
	Status     string     `json:"status,omitempty"`      // 支付状态
	RefundedAt *time.Time `json:"refunded_at,omitempty"` // 退款完成时间
}

// CancelPaymentServiceRequest 是 CancelPayment 方法的请求参数。
type CancelPaymentServiceRequest struct {
	PaymentId   string     `json:"payment_id"`             // 支付ID
	Reason      string     `json:"reason,omitempty"`       // 取消原因
	CancelledAt *time.Time `json:"cancelled_at,omitempty"` // 取消时间
}

// CancelPaymentServiceResponse 是 CancelPayment 方法的响应结果。
type CancelPaymentServiceResponse struct {
	Success   bool   `json:"success"`              // 操作是否成功
	Message   string `json:"message,omitempty"`    // 提示消息
	PaymentId string `json:"payment_id,omitempty"` // 支付ID
	Status    string `json:"status,omitempty"`     // 支付状态
}
