package user

import (
	"time"

	"github.com/soliton-go/framework/ddd"
)

// UserID is a strong typed ID.
type UserID string

func (id UserID) String() string {
	return string(id)
}

// UserStatus represents the Status enum.
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

// User is the aggregate root.
type User struct {
	ddd.BaseAggregateRoot
	ID UserID `gorm:"primaryKey"`
	Username string `gorm:"size:255"`
	Email string `gorm:"size:255"`
	Status UserStatus `gorm:"size:50;default:'active'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (User) TableName() string {
	return "users"
}

// NewUser creates a new User.
func NewUser(id string, username string, email string, status UserStatus) *User {
	e := &User{
		ID: UserID(id),
		Username: username,
		Email: email,
		Status: status,
	}
	e.AddDomainEvent(NewUserCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *User) Update(username *string, email *string, status *UserStatus) {
	if username != nil {
		e.Username = *username
	}
	if email != nil {
		e.Email = *email
	}
	if status != nil {
		e.Status = *status
	}
	e.AddDomainEvent(NewUserUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *User) GetID() ddd.ID {
	return e.ID
}
