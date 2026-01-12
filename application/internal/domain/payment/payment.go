package payment

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)
// DomainRemark: 支付领域

// PaymentID 是强类型的实体标识符。
type PaymentID string

func (id PaymentID) String() string {
	return string(id)
}

// PaymentMethod 表示 Method 字段的枚举类型。
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard PaymentMethod = "debit_card"
	PaymentMethodPaypal PaymentMethod = "paypal"
	PaymentMethodAlipay PaymentMethod = "alipay"
	PaymentMethodWechat PaymentMethod = "wechat"
	PaymentMethodCash PaymentMethod = "cash"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

// PaymentStatus 表示 Status 字段的枚举类型。
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusPaid PaymentStatus = "paid"
	PaymentStatusFailed PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Payment 是聚合根实体。
type Payment struct {
	ddd.BaseAggregateRoot
	ID PaymentID `gorm:"primaryKey"`
	OrderId string `gorm:"size:255"` // 订单ID
	UserId string `gorm:"size:255"` // 用户ID
	Amount float64 `gorm:"default:0"` // 支付金额
	Currency string `gorm:"size:255"` // 币种
	Method PaymentMethod `gorm:"size:50;default:'credit_card'"` // 支付方式
	Status PaymentStatus `gorm:"size:50;default:'pending'"` // 支付状态
	Provider string `gorm:"size:255"` // 支付渠道
	ProviderTxnId string `gorm:"size:255"` // 渠道交易号
	PaidAt *time.Time  // 支付完成时间
	RefundedAt *time.Time  // 退款完成时间
	FailureReason string `gorm:"size:255"` // 失败原因
	Metadata datatypes.JSON  // 扩展信息
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Payment) TableName() string {
	return "payments"
}

// NewPayment 创建一个新的 Payment 实体。
func NewPayment(id string, orderId string, userId string, amount float64, currency string, method PaymentMethod, status PaymentStatus, provider string, providerTxnId string, paidAt *time.Time, refundedAt *time.Time, failureReason string, metadata datatypes.JSON) *Payment {
	e := &Payment{
		ID: PaymentID(id),
		OrderId: orderId,
		UserId: userId,
		Amount: amount,
		Currency: currency,
		Method: method,
		Status: status,
		Provider: provider,
		ProviderTxnId: providerTxnId,
		PaidAt: paidAt,
		RefundedAt: refundedAt,
		FailureReason: failureReason,
		Metadata: metadata,
	}
	e.AddDomainEvent(NewPaymentCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Payment) Update(orderId *string, userId *string, amount *float64, currency *string, method *PaymentMethod, status *PaymentStatus, provider *string, providerTxnId *string, paidAt *time.Time, refundedAt *time.Time, failureReason *string, metadata *datatypes.JSON) {
	if orderId != nil {
		e.OrderId = *orderId
	}
	if userId != nil {
		e.UserId = *userId
	}
	if amount != nil {
		e.Amount = *amount
	}
	if currency != nil {
		e.Currency = *currency
	}
	if method != nil {
		e.Method = *method
	}
	if status != nil {
		e.Status = *status
	}
	if provider != nil {
		e.Provider = *provider
	}
	if providerTxnId != nil {
		e.ProviderTxnId = *providerTxnId
	}
	if paidAt != nil {
		e.PaidAt = paidAt
	}
	if refundedAt != nil {
		e.RefundedAt = refundedAt
	}
	if failureReason != nil {
		e.FailureReason = *failureReason
	}
	if metadata != nil {
		e.Metadata = *metadata
	}
	e.AddDomainEvent(NewPaymentUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Payment) GetID() ddd.ID {
	return e.ID
}
