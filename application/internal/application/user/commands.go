package userapp

import (
	"context"

	"github.com/soliton-go/application/internal/domain/user"
)

type CreateUserCommand struct {
	ID    string
	Name  string
	Email string
}

type CreateUserHandler struct {
	repo user.UserRepository
}

func NewCreateUserHandler(repo user.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
	u := user.NewUser(cmd.ID, cmd.Name, cmd.Email)
	return h.repo.Save(ctx, u)
}
