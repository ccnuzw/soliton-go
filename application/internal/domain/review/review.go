package review

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)
// DomainRemark: 评价领域

// ReviewID 是强类型的实体标识符。
type ReviewID string

func (id ReviewID) String() string {
	return string(id)
}

// ReviewStatus 表示 Status 字段的枚举类型。
type ReviewStatus string

const (
	ReviewStatusPending ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
	ReviewStatusHidden ReviewStatus = "hidden"
)

// Review 是聚合根实体。
type Review struct {
	ddd.BaseAggregateRoot
	ID ReviewID `gorm:"primaryKey"`
	ProductId string `gorm:"size:255"` // 商品ID
	UserId string `gorm:"size:255"` // 用户ID
	OrderId string `gorm:"size:255"` // 订单ID
	Rating int `gorm:"not null;default:0"` // 评分
	Title string `gorm:"size:255"` // 标题
	Content string `gorm:"size:255"` // 评价内容
	Status ReviewStatus `gorm:"size:50;default:'pending'"` // 审核状态
	IsAnonymous bool `gorm:"default:false"` // 是否匿名
	HelpfulCount int `gorm:"not null;default:0"` // 有用数
	Reply string `gorm:"size:255"` // 官方回复
	Images datatypes.JSON  // 图片列表
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (Review) TableName() string {
	return "reviews"
}

// NewReview 创建一个新的 Review 实体。
func NewReview(id string, productId string, userId string, orderId string, rating int, title string, content string, status ReviewStatus, isAnonymous bool, helpfulCount int, reply string, images datatypes.JSON) *Review {
	e := &Review{
		ID: ReviewID(id),
		ProductId: productId,
		UserId: userId,
		OrderId: orderId,
		Rating: rating,
		Title: title,
		Content: content,
		Status: status,
		IsAnonymous: isAnonymous,
		HelpfulCount: helpfulCount,
		Reply: reply,
		Images: images,
	}
	e.AddDomainEvent(NewReviewCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *Review) Update(productId *string, userId *string, orderId *string, rating *int, title *string, content *string, status *ReviewStatus, isAnonymous *bool, helpfulCount *int, reply *string, images *datatypes.JSON) {
	if productId != nil {
		e.ProductId = *productId
	}
	if userId != nil {
		e.UserId = *userId
	}
	if orderId != nil {
		e.OrderId = *orderId
	}
	if rating != nil {
		e.Rating = *rating
	}
	if title != nil {
		e.Title = *title
	}
	if content != nil {
		e.Content = *content
	}
	if status != nil {
		e.Status = *status
	}
	if isAnonymous != nil {
		e.IsAnonymous = *isAnonymous
	}
	if helpfulCount != nil {
		e.HelpfulCount = *helpfulCount
	}
	if reply != nil {
		e.Reply = *reply
	}
	if images != nil {
		e.Images = *images
	}
	e.AddDomainEvent(NewReviewUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *Review) GetID() ddd.ID {
	return e.ID
}
