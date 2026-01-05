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

// FindPaginated returns a page of entities with total count.
func (r *UserRepoImpl) FindPaginated(ctx context.Context, page, pageSize int) ([]*user.User, int64, error) {
	var entities []*user.User
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get page
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// MigrateUser creates the table if it doesn't exist.
func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&user.User{})
}
