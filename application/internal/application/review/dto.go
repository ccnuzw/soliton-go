package reviewapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/review"
)

// CreateReviewRequest 是创建 Review 的请求体。
type CreateReviewRequest struct {
	ProductId string `json:"product_id" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
	OrderId string `json:"order_id" binding:"required"`
	Rating int `json:"rating"`
	Title string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status string `json:"status" binding:"required,oneof=pending approved rejected hidden"`
	IsAnonymous bool `json:"is_anonymous"`
	HelpfulCount int `json:"helpful_count"`
	Reply string `json:"reply" binding:"required"`
	Images datatypes.JSON `json:"images"`
}

// UpdateReviewRequest 是更新 Review 的请求体。
type UpdateReviewRequest struct {
	ProductId *string `json:"product_id,omitempty"`
	UserId *string `json:"user_id,omitempty"`
	OrderId *string `json:"order_id,omitempty"`
	Rating *int `json:"rating,omitempty"`
	Title *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=pending approved rejected hidden"`
	IsAnonymous *bool `json:"is_anonymous,omitempty"`
	HelpfulCount *int `json:"helpful_count,omitempty"`
	Reply *string `json:"reply,omitempty"`
	Images *datatypes.JSON `json:"images,omitempty"`
}

// ReviewResponse 是 Review 的响应体。
type ReviewResponse struct {
	ID        string    `json:"id"`
	ProductId string `json:"product_id"`
	UserId string `json:"user_id"`
	OrderId string `json:"order_id"`
	Rating int `json:"rating"`
	Title string `json:"title"`
	Content string `json:"content"`
	Status string `json:"status"`
	IsAnonymous bool `json:"is_anonymous"`
	HelpfulCount int `json:"helpful_count"`
	Reply string `json:"reply"`
	Images datatypes.JSON `json:"images"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToReviewResponse 将实体转换为响应体。
func ToReviewResponse(e *review.Review) ReviewResponse {
	return ReviewResponse{
		ID:        string(e.ID),
		ProductId: e.ProductId,
		UserId: e.UserId,
		OrderId: e.OrderId,
		Rating: e.Rating,
		Title: e.Title,
		Content: e.Content,
		Status: string(e.Status),
		IsAnonymous: e.IsAnonymous,
		HelpfulCount: e.HelpfulCount,
		Reply: e.Reply,
		Images: e.Images,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToReviewResponseList 将实体列表转换为响应体列表。
func ToReviewResponseList(entities []*review.Review) []ReviewResponse {
	result := make([]ReviewResponse, len(entities))
	for i, e := range entities {
		result[i] = ToReviewResponse(e)
	}
	return result
}
