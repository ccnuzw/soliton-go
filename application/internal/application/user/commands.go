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
	Fullname string
	Phone string
	Avatar string
	Bio string
	Birthdate time.Time
	Gender user.UserGender
	Role user.UserRole
	Status user.UserStatus
	Emailverified bool
	Phoneverified bool
	Lastloginat time.Time
	Logincount int
	Failedlogincount int
	Balance int64
	Points int
	Viplevel int
	Preferences string
}

// CreateUserHandler 处理 CreateUserCommand。
type CreateUserHandler struct {
	repo user.UserRepository
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreateUserHandler(repo user.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
	entity := user.NewUser(cmd.ID, cmd.Username, cmd.Email, cmd.Password, cmd.Fullname, cmd.Phone, cmd.Avatar, cmd.Bio, cmd.Birthdate, cmd.Gender, cmd.Role, cmd.Status, cmd.Emailverified, cmd.Phoneverified, cmd.Lastloginat, cmd.Logincount, cmd.Failedlogincount, cmd.Balance, cmd.Points, cmd.Viplevel, cmd.Preferences)
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
	Fullname *string
	Phone *string
	Avatar *string
	Bio *string
	Birthdate *time.Time
	Gender *user.UserGender
	Role *user.UserRole
	Status *user.UserStatus
	Emailverified *bool
	Phoneverified *bool
	Lastloginat *time.Time
	Logincount *int
	Failedlogincount *int
	Balance *int64
	Points *int
	Viplevel *int
	Preferences *string
}

// UpdateUserHandler 处理 UpdateUserCommand。
type UpdateUserHandler struct {
	repo user.UserRepository
}

func NewUpdateUserHandler(repo user.UserRepository) *UpdateUserHandler {
	return &UpdateUserHandler{repo: repo}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*user.User, error) {
	entity, err := h.repo.Find(ctx, user.UserID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.Username, cmd.Email, cmd.Password, cmd.Fullname, cmd.Phone, cmd.Avatar, cmd.Bio, cmd.Birthdate, cmd.Gender, cmd.Role, cmd.Status, cmd.Emailverified, cmd.Phoneverified, cmd.Lastloginat, cmd.Logincount, cmd.Failedlogincount, cmd.Balance, cmd.Points, cmd.Viplevel, cmd.Preferences)
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
}

func NewDeleteUserHandler(repo user.UserRepository) *DeleteUserHandler {
	return &DeleteUserHandler{repo: repo}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	return h.repo.Delete(ctx, user.UserID(cmd.ID))
}
