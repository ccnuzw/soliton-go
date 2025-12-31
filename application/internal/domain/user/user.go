package user

import (
	"github.com/soliton-go/framework/ddd"
)

// UserID is a strong typed ID.
type UserID string

func (id UserID) String() string {
	return string(id)
}

// User is the aggregate root.
type User struct {
	ddd.BaseAggregateRoot

	ID    UserID `gorm:"primaryKey"`
	Name  string
	Email string
}

func NewUser(id string, name string, email string) *User {
	return &User{
		ID:    UserID(id),
		Name:  name,
		Email: email,
	}
}

func (u *User) GetID() ddd.ID {
	return u.ID
}
