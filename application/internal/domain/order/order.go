package order

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)
// DomainRemark: 订单领域

// OrderID 是强类型的实体标识符。
type OrderID string

func (id OrderID) String() string {
	return string(id)
}

// OrderPaymentMethod 表示 PaymentMethod 字段的枚举类型。
type OrderPaymentMethod string

const (
	OrderPaymentMethodCreditCard OrderPaymentMethod = "credit_card"
	OrderPaymentMethodDebitCard OrderPaymentMethod = "debit_card"
	OrderPaymentMethodPaypal OrderPaymentMethod = "paypal"
	OrderPaymentMethodAlipay OrderPaymentMethod = "alipay"
	OrderPaymentMethodWechat OrderPaymentMethod = "wechat"
	OrderPaymentMethodCash OrderPaymentMethod = "cash"
)

// OrderPaymentStatus 表示 PaymentStatus 字段的枚举类型。
type OrderPaymentStatus string

const (
	OrderPaymentStatusPending OrderPaymentStatus = "pending"
	OrderPaymentStatusPaid OrderPaymentStatus = "paid"
	OrderPaymentStatusFailed OrderPaymentStatus = "failed"
	OrderPaymentStatusRefunded OrderPaymentStatus = "refunded"
)

// OrderOrderStatus 表示 OrderStatus 字段的枚举类型。
type OrderOrderStatus string

const (
	OrderOrderStatusPending OrderOrderStatus = "pending"
	OrderOrderStatusConfirmed OrderOrderStatus = "confirmed"
	OrderOrderStatusProcessing OrderOrderStatus = "processing"
	OrderOrderStatusShipped OrderOrderStatus = "shipped"
	OrderOrderStatusDelivered OrderOrderStatus = "delivered"
	OrderOrderStatusCancelled OrderOrderStatus = "cancelled"
	OrderOrderStatusReturned OrderOrderStatus = "returned"
)

// OrderShippingMethod 表示 ShippingMethod 字段的枚举类型。
type OrderShippingMethod string

const (
	OrderShippingMethodStandard OrderShippingMethod = "standard"
	OrderShippingMethodExpress OrderShippingMethod = "express"
	OrderShippingMethodOvernight OrderShippingMethod = "overnight"
)

// Order 是聚合根实体。
type Order struct {
	ddd.BaseAggregateRoot
	ID OrderID `gorm:"primaryKey"`
	UserId string `gorm:"size:255"`
	OrderNo string `gorm:"size:255"`
	TotalAmount int64 `gorm:"not null;default:0"`
	DiscountAmount int64 `gorm:"not null;default:0"`
	TaxAmount int64 `gorm:"not null;default:0"`
	ShippingFee int64 `gorm:"not null;default:0"`
	FinalAmount int64 `gorm:"not null;default:0"`
	Currency string `gorm:"size:255"`
	PaymentMethod OrderPaymentMethod `gorm:"size:50;default:'credit_card'"`
	PaymentStatus OrderPaymentStatus `gorm:"size:50;default:'pending'"`
	OrderStatus OrderOrderStatus `gorm:"size:50;default:'pending'"`
	ShippingMethod OrderShippingMethod `gorm:"size:50;default:'standard'"`
	TrackingNumber string `gorm:"size:255"`
	ReceiverName string `gorm:"size:255"`
	ReceiverPhone string `gorm:"size:255"`
	ReceiverEmail string `gorm:"size:255"`
	ReceiverAddress string `gorm:"size:255"`
	ReceiverCity string `gorm:"size:255"`
	ReceiverState string `gorm:"size:255"`
	ReceiverCountry string `gorm:"size:255"`
	ReceiverPostalCode string `gorm:"size:255"`
	Notes string `gorm:"size:255"`
	PaidAt time.Time `gorm:"type:timestamp"`
	ShippedAt time.Time `gorm:"type:timestamp"`
	DeliveredAt time.Time `gorm:"type:timestamp"`
	CancelledAt time.Time `gorm:"type:timestamp"`
	RefundAmount int64 `gorm:"not null;default:0"`
	RefundReason string `gorm:"size:255"`
	ItemCount int `gorm:"not null;default:0"`
	Weight float64 `gorm:"default:0"`
	IsGift bool `gorm:"default:false"`
	GiftMessage string `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Order) TableName() string {
	return "orders"
}

// NewOrder 创建一个新的 Order 实体。
func NewOrder(id string, userId string, orderNo string, totalAmount int64, discountAmount int64, taxAmount int64, shippingFee int64, finalAmount int64, currency string, paymentMethod OrderPaymentMethod, paymentStatus OrderPaymentStatus, orderStatus OrderOrderStatus, shippingMethod OrderShippingMethod, trackingNumber string, receiverName string, receiverPhone string, receiverEmail string, receiverAddress string, receiverCity string, receiverState string, receiverCountry string, receiverPostalCode string, notes string, paidAt time.Time, shippedAt time.Time, deliveredAt time.Time, cancelledAt time.Time, refundAmount int64, refundReason string, itemCount int, weight float64, isGift bool, giftMessage string) *Order {
	e := &Order{
		ID: OrderID(id),
		UserId: userId,
		OrderNo: orderNo,
		TotalAmount: totalAmount,
		DiscountAmount: discountAmount,
		TaxAmount: taxAmount,
		ShippingFee: shippingFee,
		FinalAmount: finalAmount,
		Currency: currency,
		PaymentMethod: paymentMethod,
		PaymentStatus: paymentStatus,
		OrderStatus: orderStatus,
		ShippingMethod: shippingMethod,
		TrackingNumber: trackingNumber,
		ReceiverName: receiverName,
		ReceiverPhone: receiverPhone,
		ReceiverEmail: receiverEmail,
		ReceiverAddress: receiverAddress,
		ReceiverCity: receiverCity,
		ReceiverState: receiverState,
		ReceiverCountry: receiverCountry,
		ReceiverPostalCode: receiverPostalCode,
		Notes: notes,
		PaidAt: paidAt,
		ShippedAt: shippedAt,
		DeliveredAt: deliveredAt,
		CancelledAt: cancelledAt,
		RefundAmount: refundAmount,
		RefundReason: refundReason,
		ItemCount: itemCount,
		Weight: weight,
		IsGift: isGift,
		GiftMessage: giftMessage,
	}
	e.AddDomainEvent(NewOrderCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Order) Update(userId *string, orderNo *string, totalAmount *int64, discountAmount *int64, taxAmount *int64, shippingFee *int64, finalAmount *int64, currency *string, paymentMethod *OrderPaymentMethod, paymentStatus *OrderPaymentStatus, orderStatus *OrderOrderStatus, shippingMethod *OrderShippingMethod, trackingNumber *string, receiverName *string, receiverPhone *string, receiverEmail *string, receiverAddress *string, receiverCity *string, receiverState *string, receiverCountry *string, receiverPostalCode *string, notes *string, paidAt *time.Time, shippedAt *time.Time, deliveredAt *time.Time, cancelledAt *time.Time, refundAmount *int64, refundReason *string, itemCount *int, weight *float64, isGift *bool, giftMessage *string) {
	if userId != nil {
		e.UserId = *userId
	}
	if orderNo != nil {
		e.OrderNo = *orderNo
	}
	if totalAmount != nil {
		e.TotalAmount = *totalAmount
	}
	if discountAmount != nil {
		e.DiscountAmount = *discountAmount
	}
	if taxAmount != nil {
		e.TaxAmount = *taxAmount
	}
	if shippingFee != nil {
		e.ShippingFee = *shippingFee
	}
	if finalAmount != nil {
		e.FinalAmount = *finalAmount
	}
	if currency != nil {
		e.Currency = *currency
	}
	if paymentMethod != nil {
		e.PaymentMethod = *paymentMethod
	}
	if paymentStatus != nil {
		e.PaymentStatus = *paymentStatus
	}
	if orderStatus != nil {
		e.OrderStatus = *orderStatus
	}
	if shippingMethod != nil {
		e.ShippingMethod = *shippingMethod
	}
	if trackingNumber != nil {
		e.TrackingNumber = *trackingNumber
	}
	if receiverName != nil {
		e.ReceiverName = *receiverName
	}
	if receiverPhone != nil {
		e.ReceiverPhone = *receiverPhone
	}
	if receiverEmail != nil {
		e.ReceiverEmail = *receiverEmail
	}
	if receiverAddress != nil {
		e.ReceiverAddress = *receiverAddress
	}
	if receiverCity != nil {
		e.ReceiverCity = *receiverCity
	}
	if receiverState != nil {
		e.ReceiverState = *receiverState
	}
	if receiverCountry != nil {
		e.ReceiverCountry = *receiverCountry
	}
	if receiverPostalCode != nil {
		e.ReceiverPostalCode = *receiverPostalCode
	}
	if notes != nil {
		e.Notes = *notes
	}
	if paidAt != nil {
		e.PaidAt = *paidAt
	}
	if shippedAt != nil {
		e.ShippedAt = *shippedAt
	}
	if deliveredAt != nil {
		e.DeliveredAt = *deliveredAt
	}
	if cancelledAt != nil {
		e.CancelledAt = *cancelledAt
	}
	if refundAmount != nil {
		e.RefundAmount = *refundAmount
	}
	if refundReason != nil {
		e.RefundReason = *refundReason
	}
	if itemCount != nil {
		e.ItemCount = *itemCount
	}
	if weight != nil {
		e.Weight = *weight
	}
	if isGift != nil {
		e.IsGift = *isGift
	}
	if giftMessage != nil {
		e.GiftMessage = *giftMessage
	}
	e.AddDomainEvent(NewOrderUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Order) GetID() ddd.ID {
	return e.ID
}
