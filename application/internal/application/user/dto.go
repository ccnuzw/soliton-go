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
	Fullname string `json:"fullname" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
	Bio string `json:"bio" binding:"required"`
	Birthdate time.Time `json:"birthdate"`
	Gender string `json:"gender" binding:"required,oneof=male female other"`
	Role string `json:"role" binding:"required,oneof=admin manager user guest"`
	Status string `json:"status" binding:"required,oneof=active inactive suspended banned"`
	Emailverified bool `json:"emailverified"`
	Phoneverified bool `json:"phoneverified"`
	Lastloginat time.Time `json:"lastloginat"`
	Logincount int `json:"logincount"`
	Failedlogincount int `json:"failedlogincount"`
	Balance int64 `json:"balance"`
	Points int `json:"points"`
	Viplevel int `json:"viplevel"`
	Preferences string `json:"preferences" binding:"required"`
}

// UpdateUserRequest 是更新 User 的请求体。
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Fullname *string `json:"fullname,omitempty"`
	Phone *string `json:"phone,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
	Bio *string `json:"bio,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Gender *string `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Role *string `json:"role,omitempty" binding:"omitempty,oneof=admin manager user guest"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive suspended banned"`
	Emailverified *bool `json:"emailverified,omitempty"`
	Phoneverified *bool `json:"phoneverified,omitempty"`
	Lastloginat *time.Time `json:"lastloginat,omitempty"`
	Logincount *int `json:"logincount,omitempty"`
	Failedlogincount *int `json:"failedlogincount,omitempty"`
	Balance *int64 `json:"balance,omitempty"`
	Points *int `json:"points,omitempty"`
	Viplevel *int `json:"viplevel,omitempty"`
	Preferences *string `json:"preferences,omitempty"`
}

// UserResponse 是 User 的响应体。
type UserResponse struct {
	ID        string    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
	Bio string `json:"bio"`
	Birthdate time.Time `json:"birthdate"`
	Gender string `json:"gender"`
	Role string `json:"role"`
	Status string `json:"status"`
	Emailverified bool `json:"emailverified"`
	Phoneverified bool `json:"phoneverified"`
	Lastloginat time.Time `json:"lastloginat"`
	Logincount int `json:"logincount"`
	Failedlogincount int `json:"failedlogincount"`
	Balance int64 `json:"balance"`
	Points int `json:"points"`
	Viplevel int `json:"viplevel"`
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
		Fullname: e.Fullname,
		Phone: e.Phone,
		Avatar: e.Avatar,
		Bio: e.Bio,
		Birthdate: e.Birthdate,
		Gender: string(e.Gender),
		Role: string(e.Role),
		Status: string(e.Status),
		Emailverified: e.Emailverified,
		Phoneverified: e.Phoneverified,
		Lastloginat: e.Lastloginat,
		Logincount: e.Logincount,
		Failedlogincount: e.Failedlogincount,
		Balance: e.Balance,
		Points: e.Points,
		Viplevel: e.Viplevel,
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
