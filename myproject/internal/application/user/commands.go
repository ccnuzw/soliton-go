package userapp

import (
	"context"

	"com/soliton-go/myproject/internal/domain/user"
)

// CreateUserCommand is the command for creating a User.
type CreateUserCommand struct {
	ID string
	Username string
	Email string
}

// CreateUserHandler handles CreateUserCommand.
type CreateUserHandler struct {
	repo user.UserRepository
	// Optional: Add event bus for domain event publishing
	// eventBus event.EventBus
}

func NewCreateUserHandler(repo user.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
	entity := user.NewUser(cmd.ID, cmd.Username, cmd.Email)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// Optional: Publish domain events
	// Uncomment to enable event publishing:
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// UpdateUserCommand is the command for updating a User.
type UpdateUserCommand struct {
	ID string
	Username *string
	Email *string
}

// UpdateUserHandler handles UpdateUserCommand.
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
	entity.Update(cmd.Username, cmd.Email)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeleteUserCommand is the command for deleting a User.
type DeleteUserCommand struct {
	ID string
}

// DeleteUserHandler handles DeleteUserCommand.
type DeleteUserHandler struct {
	repo user.UserRepository
}

func NewDeleteUserHandler(repo user.UserRepository) *DeleteUserHandler {
	return &DeleteUserHandler{repo: repo}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	return h.repo.Delete(ctx, user.UserID(cmd.ID))
}
