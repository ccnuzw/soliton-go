package persistence

import (
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

// Migrate creates the table if it doesn't exist.
func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&user.User{})
}
