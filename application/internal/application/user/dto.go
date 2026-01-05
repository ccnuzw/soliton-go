package userapp

import (
	"time"

	"github.com/soliton-go/application/internal/domain/user"
)

// CreateUserRequest 是创建 User 的请求体。
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
	Bio string `json:"bio" binding:"required"`
	BirthDate *time.Time `json:"birth_date"`
	Gender string `json:"gender" binding:"required,oneof=male female other"`
	Role string `json:"role" binding:"required,oneof=admin manager user guest"`
	Status string `json:"status" binding:"required,oneof=active inactive suspended banned"`
	EmailVerified bool `json:"email_verified"`
	PhoneVerified bool `json:"phone_verified"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LoginCount int `json:"login_count"`
	FailedLoginCount int `json:"failed_login_count"`
	Balance int64 `json:"balance"`
	Points int `json:"points"`
	VipLevel int `json:"vip_level"`
	Preferences string `json:"preferences" binding:"required"`
}

// UpdateUserRequest 是更新 User 的请求体。
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	Phone *string `json:"phone,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
	Bio *string `json:"bio,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender *string `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Role *string `json:"role,omitempty" binding:"omitempty,oneof=admin manager user guest"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive suspended banned"`
	EmailVerified *bool `json:"email_verified,omitempty"`
	PhoneVerified *bool `json:"phone_verified,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LoginCount *int `json:"login_count,omitempty"`
	FailedLoginCount *int `json:"failed_login_count,omitempty"`
	Balance *int64 `json:"balance,omitempty"`
	Points *int `json:"points,omitempty"`
	VipLevel *int `json:"vip_level,omitempty"`
	Preferences *string `json:"preferences,omitempty"`
}

// UserResponse 是 User 的响应体。
type UserResponse struct {
	ID        string    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
	Bio string `json:"bio"`
	BirthDate *time.Time `json:"birth_date"`
	Gender string `json:"gender"`
	Role string `json:"role"`
	Status string `json:"status"`
	EmailVerified bool `json:"email_verified"`
	PhoneVerified bool `json:"phone_verified"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LoginCount int `json:"login_count"`
	FailedLoginCount int `json:"failed_login_count"`
	Balance int64 `json:"balance"`
	Points int `json:"points"`
	VipLevel int `json:"vip_level"`
	Preferences string `json:"preferences"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse 将实体转换为响应体。
func ToUserResponse(e *user.User) UserResponse {
	return UserResponse{
		ID:        string(e.ID),
		Username: e.Username,
		Email: e.Email,
		Password: e.Password,
		FullName: e.FullName,
		Phone: e.Phone,
		Avatar: e.Avatar,
		Bio: e.Bio,
		BirthDate: e.BirthDate,
		Gender: string(e.Gender),
		Role: string(e.Role),
		Status: string(e.Status),
		EmailVerified: e.EmailVerified,
		PhoneVerified: e.PhoneVerified,
		LastLoginAt: e.LastLoginAt,
		LoginCount: e.LoginCount,
		FailedLoginCount: e.FailedLoginCount,
		Balance: e.Balance,
		Points: e.Points,
		VipLevel: e.VipLevel,
		Preferences: e.Preferences,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToUserResponseList 将实体列表转换为响应体列表。
func ToUserResponseList(entities []*user.User) []UserResponse {
	result := make([]UserResponse, len(entities))
	for i, e := range entities {
		result[i] = ToUserResponse(e)
	}
	return result
}
