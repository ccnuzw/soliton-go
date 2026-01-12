package reviewapp

import (
	"context"
	"gorm.io/datatypes"

	"github.com/soliton-go/application/internal/domain/review"
)

// CreateReviewCommand 是创建 Review 的命令。
type CreateReviewCommand struct {
	ID string
	ProductId string
	UserId string
	OrderId string
	Rating int
	Title string
	Content string
	Status review.ReviewStatus
	IsAnonymous bool
	HelpfulCount int
	Reply string
	Images datatypes.JSON
}

// CreateReviewHandler 处理 CreateReviewCommand。
type CreateReviewHandler struct {
	repo review.ReviewRepository
	service *review.ReviewDomainService
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreateReviewHandler(repo review.ReviewRepository, service *review.ReviewDomainService) *CreateReviewHandler {
	return &CreateReviewHandler{repo: repo, service: service}
}

func (h *CreateReviewHandler) Handle(ctx context.Context, cmd CreateReviewCommand) (*review.Review, error) {
	entity := review.NewReview(cmd.ID, cmd.ProductId, cmd.UserId, cmd.OrderId, cmd.Rating, cmd.Title, cmd.Content, cmd.Status, cmd.IsAnonymous, cmd.HelpfulCount, cmd.Reply, cmd.Images)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// 可选：发布领域事件
	// 取消注释以启用事件发布：
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateReviewCommand 是更新 Review 的命令。
type UpdateReviewCommand struct {
	ID string
	ProductId *string
	UserId *string
	OrderId *string
	Rating *int
	Title *string
	Content *string
	Status *review.ReviewStatus
	IsAnonymous *bool
	HelpfulCount *int
	Reply *string
	Images *datatypes.JSON
}

// UpdateReviewHandler 处理 UpdateReviewCommand。
type UpdateReviewHandler struct {
	repo review.ReviewRepository
	service *review.ReviewDomainService
}

func NewUpdateReviewHandler(repo review.ReviewRepository, service *review.ReviewDomainService) *UpdateReviewHandler {
	return &UpdateReviewHandler{repo: repo, service: service}
}

func (h *UpdateReviewHandler) Handle(ctx context.Context, cmd UpdateReviewCommand) (*review.Review, error) {
	entity, err := h.repo.Find(ctx, review.ReviewID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.ProductId, cmd.UserId, cmd.OrderId, cmd.Rating, cmd.Title, cmd.Content, cmd.Status, cmd.IsAnonymous, cmd.HelpfulCount, cmd.Reply, cmd.Images)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteReviewCommand 是删除 Review 的命令。
type DeleteReviewCommand struct {
	ID string
}

// DeleteReviewHandler 处理 DeleteReviewCommand。
type DeleteReviewHandler struct {
	repo review.ReviewRepository
	service *review.ReviewDomainService
}

func NewDeleteReviewHandler(repo review.ReviewRepository, service *review.ReviewDomainService) *DeleteReviewHandler {
	return &DeleteReviewHandler{repo: repo, service: service}
}

func (h *DeleteReviewHandler) Handle(ctx context.Context, cmd DeleteReviewCommand) error {
	return h.repo.Delete(ctx, review.ReviewID(cmd.ID))
}
