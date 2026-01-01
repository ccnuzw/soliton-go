package user

import (
	"github.com/soliton-go/framework/orm"
)

// UserRepository is the interface for User persistence.
type UserRepository interface {
	orm.Repository[*User, UserID]
	// TODO: Add custom query methods here
}
