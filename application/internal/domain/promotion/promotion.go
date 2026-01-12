package promotion

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)
// DomainRemark: 促销领域

// PromotionID 是强类型的实体标识符。
type PromotionID string

func (id PromotionID) String() string {
	return string(id)
}

// PromotionDiscountType 表示 DiscountType 字段的枚举类型。
type PromotionDiscountType string

const (
	PromotionDiscountTypePercentage PromotionDiscountType = "percentage"
	PromotionDiscountTypeFixed PromotionDiscountType = "fixed"
	PromotionDiscountTypeFreeShipping PromotionDiscountType = "free_shipping"
)

// PromotionStatus 表示 Status 字段的枚举类型。
type PromotionStatus string

const (
	PromotionStatusDraft PromotionStatus = "draft"
	PromotionStatusActive PromotionStatus = "active"
	PromotionStatusExpired PromotionStatus = "expired"
	PromotionStatusDisabled PromotionStatus = "disabled"
)

// Promotion 是聚合根实体。
type Promotion struct {
	ddd.BaseAggregateRoot
	ID PromotionID `gorm:"primaryKey"`
	Code string `gorm:"size:255"` // 优惠码
	Name string `gorm:"size:255"` // 活动名称
	Description string `gorm:"size:255"` // 活动说明
	DiscountType PromotionDiscountType `gorm:"size:50;default:'percentage'"` // 优惠类型
	DiscountValue int64 `gorm:"not null;default:0"` // 优惠值
	Currency string `gorm:"size:255"` // 币种
	MinOrderAmount int64 `gorm:"not null;default:0"` // 最低订单金额
	MaxDiscountAmount int64 `gorm:"not null;default:0"` // 最大优惠金额
	UsageLimit int `gorm:"not null;default:0"` // 总使用次数
	UsedCount int `gorm:"not null;default:0"` // 已使用次数
	PerUserLimit int `gorm:"not null;default:0"` // 单用户限次
	StartsAt *time.Time  // 开始时间
	EndsAt *time.Time  // 结束时间
	Status PromotionStatus `gorm:"size:50;default:'draft'"` // 活动状态
	Metadata datatypes.JSON  // 扩展信息
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Promotion) TableName() string {
	return "promotions"
}

// NewPromotion 创建一个新的 Promotion 实体。
func NewPromotion(id string, code string, name string, description string, discountType PromotionDiscountType, discountValue int64, currency string, minOrderAmount int64, maxDiscountAmount int64, usageLimit int, usedCount int, perUserLimit int, startsAt *time.Time, endsAt *time.Time, status PromotionStatus, metadata datatypes.JSON) *Promotion {
	e := &Promotion{
		ID: PromotionID(id),
		Code: code,
		Name: name,
		Description: description,
		DiscountType: discountType,
		DiscountValue: discountValue,
		Currency: currency,
		MinOrderAmount: minOrderAmount,
		MaxDiscountAmount: maxDiscountAmount,
		UsageLimit: usageLimit,
		UsedCount: usedCount,
		PerUserLimit: perUserLimit,
		StartsAt: startsAt,
		EndsAt: endsAt,
		Status: status,
		Metadata: metadata,
	}
	e.AddDomainEvent(NewPromotionCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Promotion) Update(code *string, name *string, description *string, discountType *PromotionDiscountType, discountValue *int64, currency *string, minOrderAmount *int64, maxDiscountAmount *int64, usageLimit *int, usedCount *int, perUserLimit *int, startsAt *time.Time, endsAt *time.Time, status *PromotionStatus, metadata *datatypes.JSON) {
	if code != nil {
		e.Code = *code
	}
	if name != nil {
		e.Name = *name
	}
	if description != nil {
		e.Description = *description
	}
	if discountType != nil {
		e.DiscountType = *discountType
	}
	if discountValue != nil {
		e.DiscountValue = *discountValue
	}
	if currency != nil {
		e.Currency = *currency
	}
	if minOrderAmount != nil {
		e.MinOrderAmount = *minOrderAmount
	}
	if maxDiscountAmount != nil {
		e.MaxDiscountAmount = *maxDiscountAmount
	}
	if usageLimit != nil {
		e.UsageLimit = *usageLimit
	}
	if usedCount != nil {
		e.UsedCount = *usedCount
	}
	if perUserLimit != nil {
		e.PerUserLimit = *perUserLimit
	}
	if startsAt != nil {
		e.StartsAt = startsAt
	}
	if endsAt != nil {
		e.EndsAt = endsAt
	}
	if status != nil {
		e.Status = *status
	}
	if metadata != nil {
		e.Metadata = *metadata
	}
	e.AddDomainEvent(NewPromotionUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Promotion) GetID() ddd.ID {
	return e.ID
}
