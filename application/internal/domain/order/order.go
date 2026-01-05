package order

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)

// OrderID is a strong typed ID.
type OrderID string

func (id OrderID) String() string {
	return string(id)
}

// OrderPaymentMethod represents the PaymentMethod enum.
type OrderPaymentMethod string

const (
	OrderPaymentMethodCreditCard OrderPaymentMethod = "credit_card"
	OrderPaymentMethodDebitCard OrderPaymentMethod = "debit_card"
	OrderPaymentMethodPaypal OrderPaymentMethod = "paypal"
	OrderPaymentMethodAlipay OrderPaymentMethod = "alipay"
	OrderPaymentMethodWechat OrderPaymentMethod = "wechat"
	OrderPaymentMethodCash OrderPaymentMethod = "cash"
)

// OrderPaymentStatus represents the PaymentStatus enum.
type OrderPaymentStatus string

const (
	OrderPaymentStatusPending OrderPaymentStatus = "pending"
	OrderPaymentStatusPaid OrderPaymentStatus = "paid"
	OrderPaymentStatusFailed OrderPaymentStatus = "failed"
	OrderPaymentStatusRefunded OrderPaymentStatus = "refunded"
)

// OrderOrderStatus represents the OrderStatus enum.
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

// OrderShippingMethod represents the ShippingMethod enum.
type OrderShippingMethod string

const (
	OrderShippingMethodStandard OrderShippingMethod = "standard"
	OrderShippingMethodExpress OrderShippingMethod = "express"
	OrderShippingMethodOvernight OrderShippingMethod = "overnight"
)

// Order is the aggregate root.
type Order struct {
	ddd.BaseAggregateRoot
	ID OrderID `gorm:"primaryKey"`
	UserId string `gorm:"size:255"` // 用户ID
	OrderNo string `gorm:"size:255"` // 订单号
	TotalAmount int64 `gorm:"not null;default:0"` // 总金额
	DiscountAmount int64 `gorm:"not null;default:0"` // 折扣金额
	TaxAmount int64 `gorm:"not null;default:0"` // 税费
	ShippingFee int64 `gorm:"not null;default:0"` // 运费
	FinalAmount int64 `gorm:"not null;default:0"` // 最终金额
	Currency string `gorm:"size:255"` // 货币
	PaymentMethod OrderPaymentMethod `gorm:"size:50;default:'credit_card'"` // 支付方式
	PaymentStatus OrderPaymentStatus `gorm:"size:50;default:'pending'"` // 支付状态
	OrderStatus OrderOrderStatus `gorm:"size:50;default:'pending'"` // 订单状态
	ShippingMethod OrderShippingMethod `gorm:"size:50;default:'standard'"` // 配送方式
	TrackingNumber string `gorm:"size:255"` // 物流单号
	ReceiverName string `gorm:"size:255"` // 收货人姓名
	ReceiverPhone string `gorm:"size:255"` // 收货人电话
	ReceiverEmail string `gorm:"size:255"` // 收货人邮箱
	ReceiverAddress string `gorm:"size:255"` // 收货地址
	ReceiverCity string `gorm:"size:255"` // 城市
	ReceiverState string `gorm:"size:255"` // 省份
	ReceiverCountry string `gorm:"size:255"` // 国家
	ReceiverPostalCode string `gorm:"size:255"` // 邮编
	Notes string `gorm:"size:255"` // 订单备注
	PaidAt time.Time `gorm:"type:timestamp"` // 支付时间
	ShippedAt time.Time `gorm:"type:timestamp"` // 发货时间
	DeliveredAt time.Time `gorm:"type:timestamp"` // 送达时间
	CancelledAt time.Time `gorm:"type:timestamp"` // 取消时间
	RefundAmount int64 `gorm:"not null;default:0"` // 退款金额
	RefundReason string `gorm:"size:255"` // 退款原因
	ItemCount int `gorm:"not null;default:0"` // 商品数量
	Weight float64 `gorm:"default:0"` // 重量
	IsGift bool `gorm:"default:false"` // 是否礼物
	GiftMessage string `gorm:"size:255"` // 礼物留言
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for GORM.
func (Order) TableName() string {
	return "orders"
}

// NewOrder creates a new Order.
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

// Update updates the entity fields.
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

// GetID returns the entity ID.
func (e *Order) GetID() ddd.ID {
	return e.ID
}
