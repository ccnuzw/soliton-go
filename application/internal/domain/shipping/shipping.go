package shipping

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)
// DomainRemark: 物流领域

// ShippingID 是强类型的实体标识符。
type ShippingID string

func (id ShippingID) String() string {
	return string(id)
}

// ShippingShippingMethod 表示 ShippingMethod 字段的枚举类型。
type ShippingShippingMethod string

const (
	ShippingShippingMethodStandard ShippingShippingMethod = "standard"
	ShippingShippingMethodExpress ShippingShippingMethod = "express"
	ShippingShippingMethodOvernight ShippingShippingMethod = "overnight"
)

// ShippingStatus 表示 Status 字段的枚举类型。
type ShippingStatus string

const (
	ShippingStatusPending ShippingStatus = "pending"
	ShippingStatusLabelCreated ShippingStatus = "label_created"
	ShippingStatusInTransit ShippingStatus = "in_transit"
	ShippingStatusDelivered ShippingStatus = "delivered"
	ShippingStatusReturned ShippingStatus = "returned"
	ShippingStatusCancelled ShippingStatus = "cancelled"
)

// Shipping 是聚合根实体。
type Shipping struct {
	ddd.BaseAggregateRoot
	ID ShippingID `gorm:"primaryKey"`
	OrderId string `gorm:"size:255"` // 订单ID
	Carrier string `gorm:"size:255"` // 物流承运商
	ShippingMethod ShippingShippingMethod `gorm:"size:50;default:'standard'"` // 配送方式
	TrackingNumber string `gorm:"size:255"` // 物流单号
	Status ShippingStatus `gorm:"size:50;default:'pending'"` // 物流状态
	ShippedAt *time.Time  // 发货时间
	DeliveredAt *time.Time  // 签收时间
	ReceiverName string `gorm:"size:255"` // 收件人姓名
	ReceiverPhone string `gorm:"size:255"` // 收件人电话
	ReceiverAddress string `gorm:"size:255"` // 收件人地址
	ReceiverCity string `gorm:"size:255"` // 收件人城市
	ReceiverState string `gorm:"size:255"` // 收件人省/州
	ReceiverCountry string `gorm:"size:255"` // 收件人国家
	ReceiverPostalCode string `gorm:"size:255"` // 邮编
	Notes string `gorm:"size:255"` // 备注
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Shipping) TableName() string {
	return "shippings"
}

// NewShipping 创建一个新的 Shipping 实体。
func NewShipping(id string, orderId string, carrier string, shippingMethod ShippingShippingMethod, trackingNumber string, status ShippingStatus, shippedAt *time.Time, deliveredAt *time.Time, receiverName string, receiverPhone string, receiverAddress string, receiverCity string, receiverState string, receiverCountry string, receiverPostalCode string, notes string) *Shipping {
	e := &Shipping{
		ID: ShippingID(id),
		OrderId: orderId,
		Carrier: carrier,
		ShippingMethod: shippingMethod,
		TrackingNumber: trackingNumber,
		Status: status,
		ShippedAt: shippedAt,
		DeliveredAt: deliveredAt,
		ReceiverName: receiverName,
		ReceiverPhone: receiverPhone,
		ReceiverAddress: receiverAddress,
		ReceiverCity: receiverCity,
		ReceiverState: receiverState,
		ReceiverCountry: receiverCountry,
		ReceiverPostalCode: receiverPostalCode,
		Notes: notes,
	}
	e.AddDomainEvent(NewShippingCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Shipping) Update(orderId *string, carrier *string, shippingMethod *ShippingShippingMethod, trackingNumber *string, status *ShippingStatus, shippedAt *time.Time, deliveredAt *time.Time, receiverName *string, receiverPhone *string, receiverAddress *string, receiverCity *string, receiverState *string, receiverCountry *string, receiverPostalCode *string, notes *string) {
	if orderId != nil {
		e.OrderId = *orderId
	}
	if carrier != nil {
		e.Carrier = *carrier
	}
	if shippingMethod != nil {
		e.ShippingMethod = *shippingMethod
	}
	if trackingNumber != nil {
		e.TrackingNumber = *trackingNumber
	}
	if status != nil {
		e.Status = *status
	}
	if shippedAt != nil {
		e.ShippedAt = shippedAt
	}
	if deliveredAt != nil {
		e.DeliveredAt = deliveredAt
	}
	if receiverName != nil {
		e.ReceiverName = *receiverName
	}
	if receiverPhone != nil {
		e.ReceiverPhone = *receiverPhone
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
	e.AddDomainEvent(NewShippingUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Shipping) GetID() ddd.ID {
	return e.ID
}
