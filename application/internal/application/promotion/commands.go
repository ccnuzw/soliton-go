package promotionapp

import (
	"context"
	"time"
	"gorm.io/datatypes"

	"github.com/soliton-go/application/internal/domain/promotion"
)

// CreatePromotionCommand 是创建 Promotion 的命令。
type CreatePromotionCommand struct {
	ID string
	Code string
	Name string
	Description string
	DiscountType promotion.PromotionDiscountType
	DiscountValue int64
	Currency string
	MinOrderAmount int64
	MaxDiscountAmount int64
	UsageLimit int
	UsedCount int
	PerUserLimit int
	StartsAt *time.Time
	EndsAt *time.Time
	Status promotion.PromotionStatus
	Metadata datatypes.JSON
}

// CreatePromotionHandler 处理 CreatePromotionCommand。
type CreatePromotionHandler struct {
	repo promotion.PromotionRepository
	service *promotion.PromotionDomainService
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreatePromotionHandler(repo promotion.PromotionRepository, service *promotion.PromotionDomainService) *CreatePromotionHandler {
	return &CreatePromotionHandler{repo: repo, service: service}
}

func (h *CreatePromotionHandler) Handle(ctx context.Context, cmd CreatePromotionCommand) (*promotion.Promotion, error) {
	entity := promotion.NewPromotion(cmd.ID, cmd.Code, cmd.Name, cmd.Description, cmd.DiscountType, cmd.DiscountValue, cmd.Currency, cmd.MinOrderAmount, cmd.MaxDiscountAmount, cmd.UsageLimit, cmd.UsedCount, cmd.PerUserLimit, cmd.StartsAt, cmd.EndsAt, cmd.Status, cmd.Metadata)
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

// UpdatePromotionCommand 是更新 Promotion 的命令。
type UpdatePromotionCommand struct {
	ID string
	Code *string
	Name *string
	Description *string
	DiscountType *promotion.PromotionDiscountType
	DiscountValue *int64
	Currency *string
	MinOrderAmount *int64
	MaxDiscountAmount *int64
	UsageLimit *int
	UsedCount *int
	PerUserLimit *int
	StartsAt *time.Time
	EndsAt *time.Time
	Status *promotion.PromotionStatus
	Metadata *datatypes.JSON
}

// UpdatePromotionHandler 处理 UpdatePromotionCommand。
type UpdatePromotionHandler struct {
	repo promotion.PromotionRepository
	service *promotion.PromotionDomainService
}

func NewUpdatePromotionHandler(repo promotion.PromotionRepository, service *promotion.PromotionDomainService) *UpdatePromotionHandler {
	return &UpdatePromotionHandler{repo: repo, service: service}
}

func (h *UpdatePromotionHandler) Handle(ctx context.Context, cmd UpdatePromotionCommand) (*promotion.Promotion, error) {
	entity, err := h.repo.Find(ctx, promotion.PromotionID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.Code, cmd.Name, cmd.Description, cmd.DiscountType, cmd.DiscountValue, cmd.Currency, cmd.MinOrderAmount, cmd.MaxDiscountAmount, cmd.UsageLimit, cmd.UsedCount, cmd.PerUserLimit, cmd.StartsAt, cmd.EndsAt, cmd.Status, cmd.Metadata)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeletePromotionCommand 是删除 Promotion 的命令。
type DeletePromotionCommand struct {
	ID string
}

// DeletePromotionHandler 处理 DeletePromotionCommand。
type DeletePromotionHandler struct {
	repo promotion.PromotionRepository
	service *promotion.PromotionDomainService
}

func NewDeletePromotionHandler(repo promotion.PromotionRepository, service *promotion.PromotionDomainService) *DeletePromotionHandler {
	return &DeletePromotionHandler{repo: repo, service: service}
}

func (h *DeletePromotionHandler) Handle(ctx context.Context, cmd DeletePromotionCommand) error {
	return h.repo.Delete(ctx, promotion.PromotionID(cmd.ID))
}
