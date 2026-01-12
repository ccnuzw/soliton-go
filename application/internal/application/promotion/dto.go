package promotionapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/promotion"
	"gorm.io/datatypes"
)

// CreatePromotionRequest 是创建 Promotion 的请求体。
type CreatePromotionRequest struct {
	Code              string         `json:"code" binding:"required"`
	Name              string         `json:"name" binding:"required"`
	Description       string         `json:"description" binding:"required"`
	DiscountType      string         `json:"discount_type" binding:"required,oneof=percentage fixed free_shipping"`
	DiscountValue     int64          `json:"discount_value"`
	Currency          string         `json:"currency" binding:"required"`
	MinOrderAmount    int64          `json:"min_order_amount"`
	MaxDiscountAmount int64          `json:"max_discount_amount"`
	UsageLimit        int            `json:"usage_limit"`
	UsedCount         int            `json:"used_count"`
	PerUserLimit      int            `json:"per_user_limit"`
	StartsAt          *time.Time     `json:"starts_at"`
	EndsAt            *time.Time     `json:"ends_at"`
	Status            string         `json:"status" binding:"required,oneof=draft active expired disabled"`
	Metadata          datatypes.JSON `json:"metadata"`
}

// UpdatePromotionRequest 是更新 Promotion 的请求体。
type UpdatePromotionRequest struct {
	Code              *string         `json:"code,omitempty"`
	Name              *string         `json:"name,omitempty"`
	Description       *string         `json:"description,omitempty"`
	DiscountType      *string         `json:"discount_type,omitempty" binding:"omitempty,oneof=percentage fixed free_shipping"`
	DiscountValue     *int64          `json:"discount_value,omitempty"`
	Currency          *string         `json:"currency,omitempty"`
	MinOrderAmount    *int64          `json:"min_order_amount,omitempty"`
	MaxDiscountAmount *int64          `json:"max_discount_amount,omitempty"`
	UsageLimit        *int            `json:"usage_limit,omitempty"`
	UsedCount         *int            `json:"used_count,omitempty"`
	PerUserLimit      *int            `json:"per_user_limit,omitempty"`
	StartsAt          *time.Time      `json:"starts_at,omitempty"`
	EndsAt            *time.Time      `json:"ends_at,omitempty"`
	Status            *string         `json:"status,omitempty" binding:"omitempty,oneof=draft active expired disabled"`
	Metadata          *datatypes.JSON `json:"metadata,omitempty"`
}

// PromotionResponse 是 Promotion 的响应体。
type PromotionResponse struct {
	ID                string         `json:"id"`
	Code              string         `json:"code"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	DiscountType      string         `json:"discount_type"`
	DiscountValue     int64          `json:"discount_value"`
	Currency          string         `json:"currency"`
	MinOrderAmount    int64          `json:"min_order_amount"`
	MaxDiscountAmount int64          `json:"max_discount_amount"`
	UsageLimit        int            `json:"usage_limit"`
	UsedCount         int            `json:"used_count"`
	PerUserLimit      int            `json:"per_user_limit"`
	StartsAt          *time.Time     `json:"starts_at"`
	EndsAt            *time.Time     `json:"ends_at"`
	Status            string         `json:"status"`
	Metadata          datatypes.JSON `json:"metadata"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// ToPromotionResponse 将实体转换为响应体。
func ToPromotionResponse(e *promotion.Promotion) PromotionResponse {
	return PromotionResponse{
		ID:                string(e.ID),
		Code:              e.Code,
		Name:              e.Name,
		Description:       e.Description,
		DiscountType:      string(e.DiscountType),
		DiscountValue:     e.DiscountValue,
		Currency:          e.Currency,
		MinOrderAmount:    e.MinOrderAmount,
		MaxDiscountAmount: e.MaxDiscountAmount,
		UsageLimit:        e.UsageLimit,
		UsedCount:         e.UsedCount,
		PerUserLimit:      e.PerUserLimit,
		StartsAt:          e.StartsAt,
		EndsAt:            e.EndsAt,
		Status:            string(e.Status),
		Metadata:          e.Metadata,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

// ToPromotionResponseList 将实体列表转换为响应体列表。
func ToPromotionResponseList(entities []*promotion.Promotion) []PromotionResponse {
	result := make([]PromotionResponse, len(entities))
	for i, e := range entities {
		result[i] = ToPromotionResponse(e)
	}
	return result
}
