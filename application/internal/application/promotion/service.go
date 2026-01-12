package promotionapp

import (
	"context"
	"fmt"
	"time"

	"github.com/soliton-go/application/internal/domain/promotion"
)

// PromotionService 处理跨领域的业务逻辑编排。
type PromotionService struct {
	repo promotion.PromotionRepository
}

// NewPromotionService 创建 PromotionService 实例。
func NewPromotionService(
	repo promotion.PromotionRepository,
) *PromotionService {
	return &PromotionService{
		repo: repo,
	}
}

// ApplyPromotion 实现 ApplyPromotion 用例。
func (s *PromotionService) ApplyPromotion(ctx context.Context, req ApplyPromotionServiceRequest) (*ApplyPromotionServiceResponse, error) {
	promo, discount, valid, msg, err := s.evaluatePromotion(ctx, req.Code, req.OrderAmount, req.Currency)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("promotion not valid: %s", msg)
	}

	usedCount := promo.UsedCount + 1
	var status *promotion.PromotionStatus
	if promo.UsageLimit > 0 && usedCount >= promo.UsageLimit {
		expired := promotion.PromotionStatusExpired
		status = &expired
	}
	promo.Update(nil, nil, nil, nil, nil, nil, nil, nil, nil, &usedCount, nil, nil, nil, status, nil)
	if err := s.repo.Save(ctx, promo); err != nil {
		return nil, err
	}

	finalAmount := req.OrderAmount - discount
	if finalAmount < 0 {
		finalAmount = 0
	}

	return &ApplyPromotionServiceResponse{
		Success:        true,
		Message:        "applied",
		Code:           promo.Code,
		DiscountAmount: discount,
		FinalAmount:    finalAmount,
	}, nil
}

// ValidatePromotion 实现 ValidatePromotion 用例。
func (s *PromotionService) ValidatePromotion(ctx context.Context, req ValidatePromotionServiceRequest) (*ValidatePromotionServiceResponse, error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if req.OrderAmount <= 0 {
		return nil, fmt.Errorf("order_amount must be greater than 0")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	promo, discount, valid, msg, err := s.evaluatePromotion(ctx, req.Code, req.OrderAmount, req.Currency)
	if err != nil {
		return nil, err
	}

	return &ValidatePromotionServiceResponse{
		Success:        true,
		Message:        msg,
		Code:           promo.Code,
		DiscountAmount: discount,
		Valid:          valid,
	}, nil
}

// RevokePromotion 实现 RevokePromotion 用例。
func (s *PromotionService) RevokePromotion(ctx context.Context, req RevokePromotionServiceRequest) (*RevokePromotionServiceResponse, error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	promo, err := s.findByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	usedCount := promo.UsedCount
	if usedCount > 0 {
		usedCount--
		promo.Update(nil, nil, nil, nil, nil, nil, nil, nil, nil, &usedCount, nil, nil, nil, nil, nil)
		if err := s.repo.Save(ctx, promo); err != nil {
			return nil, err
		}
	}

	return &RevokePromotionServiceResponse{
		Success:   true,
		Message:   "revoked",
		Code:      promo.Code,
		RevokedAt: req.RevokedAt,
	}, nil
}

func (s *PromotionService) evaluatePromotion(ctx context.Context, code string, orderAmount int64, currency string) (*promotion.Promotion, int64, bool, string, error) {
	promo, err := s.findByCode(ctx, code)
	if err != nil {
		return nil, 0, false, "", err
	}

	now := time.Now()
	if promo.Status != promotion.PromotionStatusActive {
		return promo, 0, false, "promotion is not active", nil
	}
	if promo.StartsAt != nil && now.Before(*promo.StartsAt) {
		return promo, 0, false, "promotion not started", nil
	}
	if promo.EndsAt != nil && now.After(*promo.EndsAt) {
		return promo, 0, false, "promotion expired", nil
	}
	if promo.UsageLimit > 0 && promo.UsedCount >= promo.UsageLimit {
		return promo, 0, false, "promotion usage limit reached", nil
	}
	if promo.MinOrderAmount > 0 && orderAmount < promo.MinOrderAmount {
		return promo, 0, false, "order amount not eligible", nil
	}
	if promo.Currency != "" && promo.Currency != currency {
		return promo, 0, false, "currency mismatch", nil
	}

	discount := calculateDiscount(promo, orderAmount)
	return promo, discount, true, "ok", nil
}

func (s *PromotionService) findByCode(ctx context.Context, code string) (*promotion.Promotion, error) {
	items, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.Code == code {
			return item, nil
		}
	}
	return nil, fmt.Errorf("promotion not found")
}

func calculateDiscount(promo *promotion.Promotion, orderAmount int64) int64 {
	discount := int64(0)
	switch promo.DiscountType {
	case promotion.PromotionDiscountTypePercentage:
		discount = orderAmount * promo.DiscountValue / 100
	case promotion.PromotionDiscountTypeFixed:
		discount = promo.DiscountValue
	case promotion.PromotionDiscountTypeFreeShipping:
		discount = 0
	}
	if promo.MaxDiscountAmount > 0 && discount > promo.MaxDiscountAmount {
		discount = promo.MaxDiscountAmount
	}
	if discount > orderAmount {
		discount = orderAmount
	}
	return discount
}
