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
	Username string `gorm:"size:255"`
	Email string `gorm:"size:255"`
	Password string `gorm:"size:255"`
	Fullname string `gorm:"size:255"`
	Phone string `gorm:"size:255"`
	Avatar string `gorm:"size:255"`
	Bio string `gorm:"size:255"`
	Birthdate time.Time `gorm:"type:timestamp"`
	Gender UserGender `gorm:"size:50;default:'male'"`
	Role UserRole `gorm:"size:50;default:'admin'"`
	Status UserStatus `gorm:"size:50;default:'active'"`
	Emailverified bool `gorm:"default:false"`
	Phoneverified bool `gorm:"default:false"`
	Lastloginat time.Time `gorm:"type:timestamp"`
	Logincount int `gorm:"not null;default:0"`
	Failedlogincount int `gorm:"not null;default:0"`
	Balance int64 `gorm:"not null;default:0"`
	Points int `gorm:"not null;default:0"`
	Viplevel int `gorm:"not null;default:0"`
	Preferences string `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回 GORM 映射的数据库表名。
func (User) TableName() string {
	return "users"
}

// NewUser 创建一个新的 User 实体。
func NewUser(id string, username string, email string, password string, fullname string, phone string, avatar string, bio string, birthdate time.Time, gender UserGender, role UserRole, status UserStatus, emailverified bool, phoneverified bool, lastloginat time.Time, logincount int, failedlogincount int, balance int64, points int, viplevel int, preferences string) *User {
	e := &User{
		ID: UserID(id),
		Username: username,
		Email: email,
		Password: password,
		Fullname: fullname,
		Phone: phone,
		Avatar: avatar,
		Bio: bio,
		Birthdate: birthdate,
		Gender: gender,
		Role: role,
		Status: status,
		Emailverified: emailverified,
		Phoneverified: phoneverified,
		Lastloginat: lastloginat,
		Logincount: logincount,
		Failedlogincount: failedlogincount,
		Balance: balance,
		Points: points,
		Viplevel: viplevel,
		Preferences: preferences,
	}
	e.AddDomainEvent(NewUserCreatedEvent(id))
	return e
}

// Update 更新实体字段。
func (e *User) Update(username *string, email *string, password *string, fullname *string, phone *string, avatar *string, bio *string, birthdate *time.Time, gender *UserGender, role *UserRole, status *UserStatus, emailverified *bool, phoneverified *bool, lastloginat *time.Time, logincount *int, failedlogincount *int, balance *int64, points *int, viplevel *int, preferences *string) {
	if username != nil {
		e.Username = *username
	}
	if email != nil {
		e.Email = *email
	}
	if password != nil {
		e.Password = *password
	}
	if fullname != nil {
		e.Fullname = *fullname
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
	if birthdate != nil {
		e.Birthdate = *birthdate
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
	if emailverified != nil {
		e.Emailverified = *emailverified
	}
	if phoneverified != nil {
		e.Phoneverified = *phoneverified
	}
	if lastloginat != nil {
		e.Lastloginat = *lastloginat
	}
	if logincount != nil {
		e.Logincount = *logincount
	}
	if failedlogincount != nil {
		e.Failedlogincount = *failedlogincount
	}
	if balance != nil {
		e.Balance = *balance
	}
	if points != nil {
		e.Points = *points
	}
	if viplevel != nil {
		e.Viplevel = *viplevel
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
