package userapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/user"
)

// CreateUserCommand 是创建 User 的命令。
type CreateUserCommand struct {
	ID string
	Username string
	Email string
	Password string
	FullName string
	Phone string
	Avatar string
	Bio string
	BirthDate time.Time
	Gender user.UserGender
	Role user.UserRole
	Status user.UserStatus
	EmailVerified bool
	PhoneVerified bool
	LastLoginAt time.Time
	LoginCount int
	FailedLoginCount int
	Balance int64
	Points int
	VipLevel int
	Preferences string
}

// CreateUserHandler 处理 CreateUserCommand。
type CreateUserHandler struct {
	repo user.UserRepository
	service *user.UserDomainService
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreateUserHandler(repo user.UserRepository, service *user.UserDomainService) *CreateUserHandler {
	return &CreateUserHandler{repo: repo, service: service}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
	entity := user.NewUser(cmd.ID, cmd.Username, cmd.Email, cmd.Password, cmd.FullName, cmd.Phone, cmd.Avatar, cmd.Bio, cmd.BirthDate, cmd.Gender, cmd.Role, cmd.Status, cmd.EmailVerified, cmd.PhoneVerified, cmd.LastLoginAt, cmd.LoginCount, cmd.FailedLoginCount, cmd.Balance, cmd.Points, cmd.VipLevel, cmd.Preferences)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// 可选：发布领域事件
	// 取消注释以启用事件发布：
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateUserCommand 是更新 User 的命令。
type UpdateUserCommand struct {
	ID string
	Username *string
	Email *string
	Password *string
	FullName *string
	Phone *string
	Avatar *string
	Bio *string
	BirthDate *time.Time
	Gender *user.UserGender
	Role *user.UserRole
	Status *user.UserStatus
	EmailVerified *bool
	PhoneVerified *bool
	LastLoginAt *time.Time
	LoginCount *int
	FailedLoginCount *int
	Balance *int64
	Points *int
	VipLevel *int
	Preferences *string
}

// UpdateUserHandler 处理 UpdateUserCommand。
type UpdateUserHandler struct {
	repo user.UserRepository
	service *user.UserDomainService
}

func NewUpdateUserHandler(repo user.UserRepository, service *user.UserDomainService) *UpdateUserHandler {
	return &UpdateUserHandler{repo: repo, service: service}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*user.User, error) {
	entity, err := h.repo.Find(ctx, user.UserID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.Username, cmd.Email, cmd.Password, cmd.FullName, cmd.Phone, cmd.Avatar, cmd.Bio, cmd.BirthDate, cmd.Gender, cmd.Role, cmd.Status, cmd.EmailVerified, cmd.PhoneVerified, cmd.LastLoginAt, cmd.LoginCount, cmd.FailedLoginCount, cmd.Balance, cmd.Points, cmd.VipLevel, cmd.Preferences)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteUserCommand 是删除 User 的命令。
type DeleteUserCommand struct {
	ID string
}

// DeleteUserHandler 处理 DeleteUserCommand。
type DeleteUserHandler struct {
	repo user.UserRepository
	service *user.UserDomainService
}

func NewDeleteUserHandler(repo user.UserRepository, service *user.UserDomainService) *DeleteUserHandler {
	return &DeleteUserHandler{repo: repo, service: service}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	return h.repo.Delete(ctx, user.UserID(cmd.ID))
}
