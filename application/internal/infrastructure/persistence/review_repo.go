package persistence

import (
	"context"
	"fmt"

	"github.com/soliton-go/application/internal/domain/review"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type ReviewRepoImpl struct {
	*orm.GormRepository[*review.Review, review.ReviewID]
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) review.ReviewRepository {
	return &ReviewRepoImpl{
		GormRepository: orm.NewGormRepository[*review.Review, review.ReviewID](db),
		db:             db,
	}
}

// FindPaginated 返回分页数据和总数。
func (r *ReviewRepoImpl) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*review.Review, int64, error) {
	var entities []*review.Review
	var total int64

	// 查询总数
	baseQuery := r.db.WithContext(ctx).Model(&review.Review{})
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := r.db.WithContext(ctx).Offset(offset).Limit(pageSize)
	if sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateReview 创建数据库表（如不存在）。
func MigrateReview(db *gorm.DB) error {
	return db.AutoMigrate(&review.Review{})
}
