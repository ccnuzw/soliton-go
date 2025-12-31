package user

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// UserRepository is the interface for User persistence.
// It embeds the generic Repository interface but can add specific methods.
type UserRepository interface {
	orm.Repository[*User, UserID]
	FindByEmail(ctx context.Context, email string) (*User, error)
}
