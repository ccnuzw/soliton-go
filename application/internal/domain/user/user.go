package user

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)

// UserID 是强类型的实体标识符。
type UserID string

func (id UserID) String() string {
	return string(id)
}

// UserGender 表示 Gender 字段的枚举类型。
type UserGender string

const (
	UserGenderMale UserGender = "male"
	UserGenderFemale UserGender = "female"
	UserGenderOther UserGender = "other"
)

// UserRole 表示 Role 字段的枚举类型。
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleManager UserRole = "manager"
	UserRoleUser UserRole = "user"
	UserRoleGuest UserRole = "guest"
)

// UserStatus 表示 Status 字段的枚举类型。
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned UserStatus = "banned"
)

// User 是聚合根实体。
type User struct {
	ddd.BaseAggregateRoot
	ID UserID `gorm:"primaryKey"`
	Username string `gorm:"size:255"` // 用户名
	Email string `gorm:"size:255"` // 邮箱
	Password string `gorm:"size:255"` // 密码
	FullName string `gorm:"size:255"` // 全名
	Phone string `gorm:"size:255"` // 电话
	Avatar string `gorm:"size:255"` // 头像URL
	Bio string `gorm:"size:255"` // 个人简介
	BirthDate time.Time `gorm:"type:timestamp"` // 生日
	Gender UserGender `gorm:"size:50;default:'male'"` // 性别
	Role UserRole `gorm:"size:50;default:'admin'"` // 角色
	Status UserStatus `gorm:"size:50;default:'active'"` // 状态
	EmailVerified bool `gorm:"default:false"` // 邮箱已验证
	PhoneVerified bool `gorm:"default:false"` // 电话已验证
	LastLoginAt time.Time `gorm:"type:timestamp"` // 最后登录时间
	LoginCount int `gorm:"not null;default:0"` // 登录次数
	FailedLoginCount int `gorm:"not null;default:0"` // 失败登录次数
	Balance int64 `gorm:"not null;default:0"` // 账户余额
	Points int `gorm:"not null;default:0"` // 积分
	VipLevel int `gorm:"not null;default:0"` // VIP等级
	Preferences string `gorm:"size:255"` // 用户偏好设置
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (User) TableName() string {
	return "users"
}

// NewUser 创建一个新的 User 实体。
func NewUser(id string, username string, email string, password string, fullName string, phone string, avatar string, bio string, birthDate time.Time, gender UserGender, role UserRole, status UserStatus, emailVerified bool, phoneVerified bool, lastLoginAt time.Time, loginCount int, failedLoginCount int, balance int64, points int, vipLevel int, preferences string) *User {
	e := &User{
		ID: UserID(id),
		Username: username,
		Email: email,
		Password: password,
		FullName: fullName,
		Phone: phone,
		Avatar: avatar,
		Bio: bio,
		BirthDate: birthDate,
		Gender: gender,
		Role: role,
		Status: status,
		EmailVerified: emailVerified,
		PhoneVerified: phoneVerified,
		LastLoginAt: lastLoginAt,
		LoginCount: loginCount,
		FailedLoginCount: failedLoginCount,
		Balance: balance,
		Points: points,
		VipLevel: vipLevel,
		Preferences: preferences,
	}
	e.AddDomainEvent(NewUserCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *User) Update(username *string, email *string, password *string, fullName *string, phone *string, avatar *string, bio *string, birthDate *time.Time, gender *UserGender, role *UserRole, status *UserStatus, emailVerified *bool, phoneVerified *bool, lastLoginAt *time.Time, loginCount *int, failedLoginCount *int, balance *int64, points *int, vipLevel *int, preferences *string) {
	if username != nil {
		e.Username = *username
	}
	if email != nil {
		e.Email = *email
	}
	if password != nil {
		e.Password = *password
	}
	if fullName != nil {
		e.FullName = *fullName
	}
	if phone != nil {
		e.Phone = *phone
	}
	if avatar != nil {
		e.Avatar = *avatar
	}
	if bio != nil {
		e.Bio = *bio
	}
	if birthDate != nil {
		e.BirthDate = *birthDate
	}
	if gender != nil {
		e.Gender = *gender
	}
	if role != nil {
		e.Role = *role
	}
	if status != nil {
		e.Status = *status
	}
	if emailVerified != nil {
		e.EmailVerified = *emailVerified
	}
	if phoneVerified != nil {
		e.PhoneVerified = *phoneVerified
	}
	if lastLoginAt != nil {
		e.LastLoginAt = *lastLoginAt
	}
	if loginCount != nil {
		e.LoginCount = *loginCount
	}
	if failedLoginCount != nil {
		e.FailedLoginCount = *failedLoginCount
	}
	if balance != nil {
		e.Balance = *balance
	}
	if points != nil {
		e.Points = *points
	}
	if vipLevel != nil {
		e.VipLevel = *vipLevel
	}
	if preferences != nil {
		e.Preferences = *preferences
	}
	e.AddDomainEvent(NewUserUpdatedEvent(string(e.ID)))
}

// GetID 返回实体 ID。
func (e *User) GetID() ddd.ID {
	return e.ID
}
