package promotion

import "context"

// PromotionDomainService 提供领域内的复杂业务逻辑封装。
type PromotionDomainService struct {
	repo PromotionRepository
}

// NewPromotionDomainService 创建 PromotionDomainService 实例。
func NewPromotionDomainService(repo PromotionRepository) *PromotionDomainService {
	return &PromotionDomainService{repo: repo}
}

// TODO: 在此添加领域服务方法，例如复杂校验、跨实体规则等。
func (s *PromotionDomainService) Validate(ctx context.Context) error {
	_ = ctx
	return nil
}
