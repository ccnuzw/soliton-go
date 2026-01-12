package reviewapp

import (
	"time"

	"gorm.io/datatypes"
)

// CreateReviewServiceRequest 是 CreateReview 方法的请求参数。
type CreateReviewServiceRequest struct {
	ProductId   string         `json:"product_id"`       // 商品ID
	UserId      string         `json:"user_id"`          // 用户ID
	OrderId     string         `json:"order_id"`         // 订单ID
	Rating      int            `json:"rating"`           // 评分
	Title       string         `json:"title"`            // 标题
	Content     string         `json:"content"`          // 评价内容
	IsAnonymous bool           `json:"is_anonymous"`     // 是否匿名
	Images      datatypes.JSON `json:"images,omitempty"` // 图片列表
}

// CreateReviewServiceResponse 是 CreateReview 方法的响应结果。
type CreateReviewServiceResponse struct {
	Success  bool   `json:"success"`             // 操作是否成功
	Message  string `json:"message,omitempty"`   // 提示消息
	ReviewId string `json:"review_id,omitempty"` // 评价ID
	Status   string `json:"status,omitempty"`    // 审核状态
}

// ModerateReviewServiceRequest 是 ModerateReview 方法的请求参数。
type ModerateReviewServiceRequest struct {
	ReviewId    string `json:"review_id"`              // 评价ID
	Status      string `json:"status"`                 // 审核状态
	Reason      string `json:"reason,omitempty"`       // 审核原因
	ModeratorId string `json:"moderator_id,omitempty"` // 审核人ID
}

// ModerateReviewServiceResponse 是 ModerateReview 方法的响应结果。
type ModerateReviewServiceResponse struct {
	Success  bool   `json:"success"`             // 操作是否成功
	Message  string `json:"message,omitempty"`   // 提示消息
	ReviewId string `json:"review_id,omitempty"` // 评价ID
	Status   string `json:"status,omitempty"`    // 审核状态
}

// ReplyReviewServiceRequest 是 ReplyReview 方法的请求参数。
type ReplyReviewServiceRequest struct {
	ReviewId  string     `json:"review_id"`            // 评价ID
	Reply     string     `json:"reply"`                // 回复内容
	RepliedBy string     `json:"replied_by,omitempty"` // 回复人
	RepliedAt *time.Time `json:"replied_at,omitempty"` // 回复时间
}

// ReplyReviewServiceResponse 是 ReplyReview 方法的响应结果。
type ReplyReviewServiceResponse struct {
	Success   bool       `json:"success"`              // 操作是否成功
	Message   string     `json:"message,omitempty"`    // 提示消息
	ReviewId  string     `json:"review_id,omitempty"`  // 评价ID
	Reply     string     `json:"reply,omitempty"`      // 回复内容
	RepliedAt *time.Time `json:"replied_at,omitempty"` // 回复时间
}
