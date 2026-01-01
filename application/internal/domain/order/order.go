package order

import (
	"time"

	"github.com/soliton-go/framework/ddd"
)

// OrderID is a strong typed ID.
type OrderID string

func (id OrderID) String() string {
	return string(id)
}

// OrderStatus represents the Status enum.
type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusPaid OrderStatus = "paid"
	OrderStatusShipped OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order is the aggregate root.
type Order struct {
	ddd.BaseAggregateRoot
	ID OrderID `gorm:"primaryKey"`
	UserId string `gorm:"size:36;index"`
	OrderNo string `gorm:"size:255"`
	TotalAmount int64 `gorm:"not null;default:0"`
	Status OrderStatus `gorm:"size:50;default:'pending'"`
	ReceiverName string `gorm:"size:255"`
	ReceiverPhone string `gorm:"size:255"`
	ReceiverAddress string `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (Order) TableName() string {
	return "orders"
}

// NewOrder creates a new Order.
func NewOrder(id string, userId string, orderNo string, totalAmount int64, status OrderStatus, receiverName string, receiverPhone string, receiverAddress string) *Order {
	e := &Order{
		ID: OrderID(id),
		UserId: userId,
		OrderNo: orderNo,
		TotalAmount: totalAmount,
		Status: status,
		ReceiverName: receiverName,
		ReceiverPhone: receiverPhone,
		ReceiverAddress: receiverAddress,
	}
	e.AddDomainEvent(NewOrderCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *Order) Update(userId string, orderNo string, totalAmount int64, status OrderStatus, receiverName string, receiverPhone string, receiverAddress string) {
	e.UserId = userId
	e.OrderNo = orderNo
	e.TotalAmount = totalAmount
	e.Status = status
	e.ReceiverName = receiverName
	e.ReceiverPhone = receiverPhone
	e.ReceiverAddress = receiverAddress
	e.AddDomainEvent(NewOrderUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *Order) GetID() ddd.ID {
	return e.ID
}
