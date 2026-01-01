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

// UserRole represents the Role enum.
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleSeller UserRole = "seller"
	UserRoleCustomer UserRole = "customer"
)

// UserStatus represents the Status enum.
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned UserStatus = "banned"
)

// User is the aggregate root.
type User struct {
	ddd.BaseAggregateRoot
	ID UserID `gorm:"primaryKey"`
	Username string `gorm:"size:255"`
	Email string `gorm:"size:255"`
	PasswordHash string `gorm:"size:255"`
	Phone string `gorm:"size:255"`
	Avatar string `gorm:"size:255"`
	Nickname string `gorm:"size:255"`
	Role UserRole `gorm:"size:50;default:'admin'"`
	Status UserStatus `gorm:"size:50;default:'active'"`
	LastLoginAt *time.Time 
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (User) TableName() string {
	return "users"
}

// NewUser creates a new User.
func NewUser(id string, username string, email string, passwordHash string, phone string, avatar string, nickname string, role UserRole, status UserStatus, lastLoginAt *time.Time) *User {
	e := &User{
		ID: UserID(id),
		Username: username,
		Email: email,
		PasswordHash: passwordHash,
		Phone: phone,
		Avatar: avatar,
		Nickname: nickname,
		Role: role,
		Status: status,
		LastLoginAt: lastLoginAt,
	}
	e.AddDomainEvent(NewUserCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *User) Update(username string, email string, passwordHash string, phone string, avatar string, nickname string, role UserRole, status UserStatus, lastLoginAt *time.Time) {
	e.Username = username
	e.Email = email
	e.PasswordHash = passwordHash
	e.Phone = phone
	e.Avatar = avatar
	e.Nickname = nickname
	e.Role = role
	e.Status = status
	e.LastLoginAt = lastLoginAt
	e.AddDomainEvent(NewUserUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *User) GetID() ddd.ID {
	return e.ID
}
