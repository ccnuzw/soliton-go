package review

import "context"

// ReviewDomainService 提供领域内的复杂业务逻辑封装。
type ReviewDomainService struct {
	repo ReviewRepository
}

// NewReviewDomainService 创建 ReviewDomainService 实例。
func NewReviewDomainService(repo ReviewRepository) *ReviewDomainService {
	return &ReviewDomainService{repo: repo}
}

// TODO: 在此添加领域服务方法，例如复杂校验、跨实体规则等。
func (s *ReviewDomainService) Validate(ctx context.Context) error {
	_ = ctx
	return nil
}
