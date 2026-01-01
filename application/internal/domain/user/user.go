package user

import "github.com/soliton-go/framework/ddd"

// UserID is a strong typed ID.
type UserID string

func (id UserID) String() string {
	return string(id)
}

// User is the aggregate root.
type User struct {
	ddd.BaseAggregateRoot
	ID   UserID `gorm:"primaryKey"`
	Name string            `gorm:"size:255"`
	// TODO: Add more fields here
}

// TableName returns the table name for GORM.
func (User) TableName() string {
	return "users"
}

// NewUser creates a new User.
func NewUser(id, name string) *User {
	e := &User{
		ID:   UserID(id),
		Name: name,
	}
	e.AddDomainEvent(NewUserCreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *User) Update(name string) {
	e.Name = name
	e.AddDomainEvent(NewUserUpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *User) GetID() ddd.ID {
	return e.ID
}
