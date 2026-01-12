package user

import (
	"time"

	"github.com/soliton-go/framework/ddd"
)
// DomainRemark: 用户领域

// UserID 是强类型的实体标识符。
type UserID string

func (id UserID) String() string {
	return string(id)
}

// User 是聚合根实体。
type User struct {
	ddd.BaseAggregateRoot
	ID UserID `gorm:"primaryKey"`
	Username string `gorm:"size:255"` // 用户名
	Email string `gorm:"size:255"` // 邮箱
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName 返回 GORM 映射的数据库表名。
func (User) TableName() string {
	return "users"
}

// NewUser 创建一个新的 User 实体。
func NewUser(id string, username string, email string) *User {
	e := &User{
		ID: UserID(id),
		Username: username,
		Email: email,
	}
	e.AddDomainEvent(NewUserCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *User) Update(username *string, email *string) {
	if username != nil {
		e.Username = *username
	}
	if email != nil {
		e.Email = *email
	}
	e.AddDomainEvent(NewUserUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *User) GetID() ddd.ID {
	return e.ID
}
