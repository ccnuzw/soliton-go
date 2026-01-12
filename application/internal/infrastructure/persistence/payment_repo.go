package persistence

import (
	"context"
	"fmt"

	"github.com/soliton-go/application/internal/domain/payment"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type PaymentRepoImpl struct {
	*orm.GormRepository[*payment.Payment, payment.PaymentID]
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) payment.PaymentRepository {
	return &PaymentRepoImpl{
		GormRepository: orm.NewGormRepository[*payment.Payment, payment.PaymentID](db),
		db:             db,
	}
}

// FindPaginated 返回分页数据和总数。
func (r *PaymentRepoImpl) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*payment.Payment, int64, error) {
	var entities []*payment.Payment
	var total int64

	// 查询总数
	baseQuery := r.db.WithContext(ctx).Model(&payment.Payment{})
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

// MigratePayment 创建数据库表（如不存在）。
func MigratePayment(db *gorm.DB) error {
	return db.AutoMigrate(&payment.Payment{})
}
