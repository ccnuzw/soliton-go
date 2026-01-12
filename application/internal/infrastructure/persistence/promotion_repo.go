package persistence

import (
	"context"
	"fmt"

	"github.com/soliton-go/application/internal/domain/promotion"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type PromotionRepoImpl struct {
	*orm.GormRepository[*promotion.Promotion, promotion.PromotionID]
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) promotion.PromotionRepository {
	return &PromotionRepoImpl{
		GormRepository: orm.NewGormRepository[*promotion.Promotion, promotion.PromotionID](db),
		db:             db,
	}
}

// FindPaginated 返回分页数据和总数。
func (r *PromotionRepoImpl) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*promotion.Promotion, int64, error) {
	var entities []*promotion.Promotion
	var total int64

	// 查询总数
	baseQuery := r.db.WithContext(ctx).Model(&promotion.Promotion{})
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

// MigratePromotion 创建数据库表（如不存在）。
func MigratePromotion(db *gorm.DB) error {
	return db.AutoMigrate(&promotion.Promotion{})
}
