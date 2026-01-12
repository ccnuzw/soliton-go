package promotionapp

import "time"

// ApplyPromotionServiceRequest 是 ApplyPromotion 方法的请求参数。
type ApplyPromotionServiceRequest struct {
	Code        string     `json:"code"`                 // 优惠码
	UserId      string     `json:"user_id"`              // 用户ID
	OrderId     string     `json:"order_id"`             // 订单ID
	OrderAmount int64      `json:"order_amount"`         // 订单金额
	Currency    string     `json:"currency"`             // 币种
	AppliedAt   *time.Time `json:"applied_at,omitempty"` // 应用时间
}

// ApplyPromotionServiceResponse 是 ApplyPromotion 方法的响应结果。
type ApplyPromotionServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	Code           string `json:"code,omitempty"`            // 优惠码
	DiscountAmount int64  `json:"discount_amount,omitempty"` // 优惠金额
	FinalAmount    int64  `json:"final_amount,omitempty"`    // 优惠后金额
}

// ValidatePromotionServiceRequest 是 ValidatePromotion 方法的请求参数。
type ValidatePromotionServiceRequest struct {
	Code        string `json:"code"`         // 优惠码
	UserId      string `json:"user_id"`      // 用户ID
	OrderAmount int64  `json:"order_amount"` // 订单金额
	Currency    string `json:"currency"`     // 币种
}

// ValidatePromotionServiceResponse 是 ValidatePromotion 方法的响应结果。
type ValidatePromotionServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	Code           string `json:"code,omitempty"`            // 优惠码
	DiscountAmount int64  `json:"discount_amount,omitempty"` // 优惠金额
	Valid          bool   `json:"valid"`                     // 是否可用
}

// RevokePromotionServiceRequest 是 RevokePromotion 方法的请求参数。
type RevokePromotionServiceRequest struct {
	Code      string     `json:"code"`                 // 优惠码
	OrderId   string     `json:"order_id"`             // 订单ID
	UserId    string     `json:"user_id"`              // 用户ID
	Reason    string     `json:"reason,omitempty"`     // 撤销原因
	RevokedAt *time.Time `json:"revoked_at,omitempty"` // 撤销时间
}

// RevokePromotionServiceResponse 是 RevokePromotion 方法的响应结果。
type RevokePromotionServiceResponse struct {
	Success   bool       `json:"success"`              // 操作是否成功
	Message   string     `json:"message,omitempty"`    // 提示消息
	Code      string     `json:"code,omitempty"`       // 优惠码
	RevokedAt *time.Time `json:"revoked_at,omitempty"` // 撤销时间
}
