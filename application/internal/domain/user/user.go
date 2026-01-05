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

// UserGender represents the Gender enum.
type UserGender string

const (
	UserGenderMale UserGender = "male"
	UserGenderFemale UserGender = "female"
	UserGenderOther UserGender = "other"
)

// UserRole represents the Role enum.
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleManager UserRole = "manager"
	UserRoleUser UserRole = "user"
	UserRoleGuest UserRole = "guest"
)

// UserStatus represents the Status enum.
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned UserStatus = "banned"
)

// User is the aggregate root.
type User struct {
	ddd.BaseAggregateRoot
	ID UserID `gorm:"primaryKey"`
	Username string `gorm:"size:255"`
	Email string `gorm:"size:255"`
	Password string `gorm:"size:255"`
	FullName string `gorm:"size:255"`
	Phone string `gorm:"size:255"`
	Avatar string `gorm:"size:255"`
	Bio string `gorm:"type:text"`
	BirthDate *time.Time 
	Gender UserGender `gorm:"size:50;default:'male'"`
	Role UserRole `gorm:"size:50;default:'admin'"`
	Status UserStatus `gorm:"size:50;default:'active'"`
	EmailVerified bool `gorm:"default:false"`
	PhoneVerified bool `gorm:"default:false"`
	LastLoginAt *time.Time 
	LoginCount int `gorm:"not null;default:0"`
	FailedLoginCount int `gorm:"not null;default:0"`
	Balance int64 `gorm:"not null;default:0"`
	Points int `gorm:"not null;default:0"`
	VipLevel int `gorm:"not null;default:0"`
	Preferences string `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (User) TableName() string {
	return "users"
}

// NewUser creates a new User.
func NewUser(id string, username string, email string, password string, fullName string, phone string, avatar string, bio string, birthDate *time.Time, gender UserGender, role UserRole, status UserStatus, emailVerified bool, phoneVerified bool, lastLoginAt *time.Time, loginCount int, failedLoginCount int, balance int64, points int, vipLevel int, preferences string) *User {
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

// Update updates the entity fields.
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
		e.BirthDate = birthDate
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
		e.LastLoginAt = lastLoginAt
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

// GetID returns the entity ID.
func (e *User) GetID() ddd.ID {
	return e.ID
}
