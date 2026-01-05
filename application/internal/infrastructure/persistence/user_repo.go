package persistence

import (
	"context"

	"github.com/soliton-go/application/internal/domain/user"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type UserRepoImpl struct {
	*orm.GormRepository[*user.User, user.UserID]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepoImpl{
		GormRepository: orm.NewGormRepository[*user.User, user.UserID](db),
		db:             db,
	}
}

// FindPaginated 返回分页数据和总数。
func (r *UserRepoImpl) FindPaginated(ctx context.Context, page, pageSize int) ([]*user.User, int64, error) {
	var entities []*user.User
	var total int64

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateUser 创建数据库表（如不存在）。
func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&user.User{})
}
