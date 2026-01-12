package reviewapp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/soliton-go/application/internal/domain/review"
)

// ReviewService 处理跨领域的业务逻辑编排。
type ReviewService struct {
	repo review.ReviewRepository
}

// NewReviewService 创建 ReviewService 实例。
func NewReviewService(
	repo review.ReviewRepository,
) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

// CreateReview 实现 CreateReview 用例。
func (s *ReviewService) CreateReview(ctx context.Context, req CreateReviewServiceRequest) (*CreateReviewServiceResponse, error) {
	if req.ProductId == "" || req.UserId == "" || req.OrderId == "" {
		return nil, fmt.Errorf("product_id, user_id, order_id are required")
	}
	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}
	if req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	entity := review.NewReview(
		uuid.New().String(),
		req.ProductId,
		req.UserId,
		req.OrderId,
		req.Rating,
		req.Title,
		req.Content,
		review.ReviewStatusPending,
		req.IsAnonymous,
		0,
		"",
		req.Images,
	)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &CreateReviewServiceResponse{
		Success:  true,
		Message:  "created",
		ReviewId: string(entity.ID),
		Status:   string(entity.Status),
	}, nil
}

// ModerateReview 实现 ModerateReview 用例。
func (s *ReviewService) ModerateReview(ctx context.Context, req ModerateReviewServiceRequest) (*ModerateReviewServiceResponse, error) {
	if req.ReviewId == "" {
		return nil, fmt.Errorf("review_id is required")
	}
	if req.Status == "" {
		return nil, fmt.Errorf("status is required")
	}

	entity, err := s.repo.Find(ctx, review.ReviewID(req.ReviewId))
	if err != nil {
		return nil, err
	}

	status := review.ReviewStatus(req.Status)
	entity.Update(nil, nil, nil, nil, nil, nil, &status, nil, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &ModerateReviewServiceResponse{
		Success:  true,
		Message:  "moderated",
		ReviewId: string(entity.ID),
		Status:   string(entity.Status),
	}, nil
}

// ReplyReview 实现 ReplyReview 用例。
func (s *ReviewService) ReplyReview(ctx context.Context, req ReplyReviewServiceRequest) (*ReplyReviewServiceResponse, error) {
	if req.ReviewId == "" {
		return nil, fmt.Errorf("review_id is required")
	}
	if req.Reply == "" {
		return nil, fmt.Errorf("reply is required")
	}

	entity, err := s.repo.Find(ctx, review.ReviewID(req.ReviewId))
	if err != nil {
		return nil, err
	}

	entity.Update(nil, nil, nil, nil, nil, nil, nil, nil, nil, &req.Reply, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	repliedAt := req.RepliedAt
	if repliedAt == nil {
		now := time.Now()
		repliedAt = &now
	}

	return &ReplyReviewServiceResponse{
		Success:   true,
		Message:   "replied",
		ReviewId:  string(entity.ID),
		Reply:     entity.Reply,
		RepliedAt: repliedAt,
	}, nil
}
