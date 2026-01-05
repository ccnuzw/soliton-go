package order

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)

// OrderID 是强类型的实体标识符。
type OrderID string

func (id OrderID) String() string {
	return string(id)
}

// OrderPaymentmethod 表示 Paymentmethod 字段的枚举类型。
type OrderPaymentmethod string

const (
	OrderPaymentmethodCreditCard OrderPaymentmethod = "credit_card"
	OrderPaymentmethodDebitCard OrderPaymentmethod = "debit_card"
	OrderPaymentmethodPaypal OrderPaymentmethod = "paypal"
	OrderPaymentmethodAlipay OrderPaymentmethod = "alipay"
	OrderPaymentmethodWechat OrderPaymentmethod = "wechat"
	OrderPaymentmethodCash OrderPaymentmethod = "cash"
)

// OrderPaymentstatus 表示 Paymentstatus 字段的枚举类型。
type OrderPaymentstatus string

const (
	OrderPaymentstatusPending OrderPaymentstatus = "pending"
	OrderPaymentstatusPaid OrderPaymentstatus = "paid"
	OrderPaymentstatusFailed OrderPaymentstatus = "failed"
	OrderPaymentstatusRefunded OrderPaymentstatus = "refunded"
)

// OrderOrderstatus 表示 Orderstatus 字段的枚举类型。
type OrderOrderstatus string

const (
	OrderOrderstatusPending OrderOrderstatus = "pending"
	OrderOrderstatusConfirmed OrderOrderstatus = "confirmed"
	OrderOrderstatusProcessing OrderOrderstatus = "processing"
	OrderOrderstatusShipped OrderOrderstatus = "shipped"
	OrderOrderstatusDelivered OrderOrderstatus = "delivered"
	OrderOrderstatusCancelled OrderOrderstatus = "cancelled"
	OrderOrderstatusReturned OrderOrderstatus = "returned"
)

// OrderShippingmethod 表示 Shippingmethod 字段的枚举类型。
type OrderShippingmethod string

const (
	OrderShippingmethodStandard OrderShippingmethod = "standard"
	OrderShippingmethodExpress OrderShippingmethod = "express"
	OrderShippingmethodOvernight OrderShippingmethod = "overnight"
)

// Order 是聚合根实体。
type Order struct {
	ddd.BaseAggregateRoot
	ID OrderID `gorm:"primaryKey"`
	Userid string `gorm:"size:255"`
	Orderno string `gorm:"size:255"`
	Totalamount int64 `gorm:"not null;default:0"`
	Discountamount int64 `gorm:"not null;default:0"`
	Taxamount int64 `gorm:"not null;default:0"`
	Shippingfee int64 `gorm:"not null;default:0"`
	Finalamount int64 `gorm:"not null;default:0"`
	Currency string `gorm:"size:255"`
	Paymentmethod OrderPaymentmethod `gorm:"size:50;default:'credit_card'"`
	Paymentstatus OrderPaymentstatus `gorm:"size:50;default:'pending'"`
	Orderstatus OrderOrderstatus `gorm:"size:50;default:'pending'"`
	Shippingmethod OrderShippingmethod `gorm:"size:50;default:'standard'"`
	Trackingnumber string `gorm:"size:255"`
	Receivername string `gorm:"size:255"`
	Receiverphone string `gorm:"size:255"`
	Receiveremail string `gorm:"size:255"`
	Receiveraddress string `gorm:"size:255"`
	Receivercity string `gorm:"size:255"`
	Receiverstate string `gorm:"size:255"`
	Receivercountry string `gorm:"size:255"`
	Receiverpostalcode string `gorm:"size:255"`
	Notes string `gorm:"size:255"`
	Paidat time.Time `gorm:"type:timestamp"`
	Shippedat time.Time `gorm:"type:timestamp"`
	Deliveredat time.Time `gorm:"type:timestamp"`
	Cancelledat time.Time `gorm:"type:timestamp"`
	Refundamount int64 `gorm:"not null;default:0"`
	Refundreason string `gorm:"size:255"`
	Itemcount int `gorm:"not null;default:0"`
	Weight float64 `gorm:"default:0"`
	Isgift bool `gorm:"default:false"`
	Giftmessage string `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Order) TableName() string {
	return "orders"
}

// NewOrder 创建一个新的 Order 实体。
func NewOrder(id string, userid string, orderno string, totalamount int64, discountamount int64, taxamount int64, shippingfee int64, finalamount int64, currency string, paymentmethod OrderPaymentmethod, paymentstatus OrderPaymentstatus, orderstatus OrderOrderstatus, shippingmethod OrderShippingmethod, trackingnumber string, receivername string, receiverphone string, receiveremail string, receiveraddress string, receivercity string, receiverstate string, receivercountry string, receiverpostalcode string, notes string, paidat time.Time, shippedat time.Time, deliveredat time.Time, cancelledat time.Time, refundamount int64, refundreason string, itemcount int, weight float64, isgift bool, giftmessage string) *Order {
	e := &Order{
		ID: OrderID(id),
		Userid: userid,
		Orderno: orderno,
		Totalamount: totalamount,
		Discountamount: discountamount,
		Taxamount: taxamount,
		Shippingfee: shippingfee,
		Finalamount: finalamount,
		Currency: currency,
		Paymentmethod: paymentmethod,
		Paymentstatus: paymentstatus,
		Orderstatus: orderstatus,
		Shippingmethod: shippingmethod,
		Trackingnumber: trackingnumber,
		Receivername: receivername,
		Receiverphone: receiverphone,
		Receiveremail: receiveremail,
		Receiveraddress: receiveraddress,
		Receivercity: receivercity,
		Receiverstate: receiverstate,
		Receivercountry: receivercountry,
		Receiverpostalcode: receiverpostalcode,
		Notes: notes,
		Paidat: paidat,
		Shippedat: shippedat,
		Deliveredat: deliveredat,
		Cancelledat: cancelledat,
		Refundamount: refundamount,
		Refundreason: refundreason,
		Itemcount: itemcount,
		Weight: weight,
		Isgift: isgift,
		Giftmessage: giftmessage,
	}
	e.AddDomainEvent(NewOrderCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Order) Update(userid *string, orderno *string, totalamount *int64, discountamount *int64, taxamount *int64, shippingfee *int64, finalamount *int64, currency *string, paymentmethod *OrderPaymentmethod, paymentstatus *OrderPaymentstatus, orderstatus *OrderOrderstatus, shippingmethod *OrderShippingmethod, trackingnumber *string, receivername *string, receiverphone *string, receiveremail *string, receiveraddress *string, receivercity *string, receiverstate *string, receivercountry *string, receiverpostalcode *string, notes *string, paidat *time.Time, shippedat *time.Time, deliveredat *time.Time, cancelledat *time.Time, refundamount *int64, refundreason *string, itemcount *int, weight *float64, isgift *bool, giftmessage *string) {
	if userid != nil {
		e.Userid = *userid
	}
	if orderno != nil {
		e.Orderno = *orderno
	}
	if totalamount != nil {
		e.Totalamount = *totalamount
	}
	if discountamount != nil {
		e.Discountamount = *discountamount
	}
	if taxamount != nil {
		e.Taxamount = *taxamount
	}
	if shippingfee != nil {
		e.Shippingfee = *shippingfee
	}
	if finalamount != nil {
		e.Finalamount = *finalamount
	}
	if currency != nil {
		e.Currency = *currency
	}
	if paymentmethod != nil {
		e.Paymentmethod = *paymentmethod
	}
	if paymentstatus != nil {
		e.Paymentstatus = *paymentstatus
	}
	if orderstatus != nil {
		e.Orderstatus = *orderstatus
	}
	if shippingmethod != nil {
		e.Shippingmethod = *shippingmethod
	}
	if trackingnumber != nil {
		e.Trackingnumber = *trackingnumber
	}
	if receivername != nil {
		e.Receivername = *receivername
	}
	if receiverphone != nil {
		e.Receiverphone = *receiverphone
	}
	if receiveremail != nil {
		e.Receiveremail = *receiveremail
	}
	if receiveraddress != nil {
		e.Receiveraddress = *receiveraddress
	}
	if receivercity != nil {
		e.Receivercity = *receivercity
	}
	if receiverstate != nil {
		e.Receiverstate = *receiverstate
	}
	if receivercountry != nil {
		e.Receivercountry = *receivercountry
	}
	if receiverpostalcode != nil {
		e.Receiverpostalcode = *receiverpostalcode
	}
	if notes != nil {
		e.Notes = *notes
	}
	if paidat != nil {
		e.Paidat = *paidat
	}
	if shippedat != nil {
		e.Shippedat = *shippedat
	}
	if deliveredat != nil {
		e.Deliveredat = *deliveredat
	}
	if cancelledat != nil {
		e.Cancelledat = *cancelledat
	}
	if refundamount != nil {
		e.Refundamount = *refundamount
	}
	if refundreason != nil {
		e.Refundreason = *refundreason
	}
	if itemcount != nil {
		e.Itemcount = *itemcount
	}
	if weight != nil {
		e.Weight = *weight
	}
	if isgift != nil {
		e.Isgift = *isgift
	}
	if giftmessage != nil {
		e.Giftmessage = *giftmessage
	}
	e.AddDomainEvent(NewOrderUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Order) GetID() ddd.ID {
	return e.ID
}
