package persistence

import (
	"context"

	"github.com/soliton-go/application/internal/domain/user"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	*orm.GormRepository[*user.User, user.UserID]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepositoryImpl{
		GormRepository: orm.NewGormRepository[*user.User, user.UserID](db),
		db:             db,
	}
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
