package userapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/user"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module provides all User dependencies for Fx.
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) user.UserRepository {
		return persistence.NewUserRepository(db)
	}),

	// Command Handlers
	fx.Provide(NewCreateUserHandler),
	fx.Provide(NewUpdateUserHandler),
	fx.Provide(NewDeleteUserHandler),

	// Query Handlers
	fx.Provide(NewGetUserHandler),
	fx.Provide(NewListUsersHandler),

	// Optional: Register with CQRS bus
	// Uncomment to enable CQRS pattern:
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreateUserHandler,
	//     updateHandler *UpdateUserHandler,
	//     deleteHandler *DeleteUserHandler,
	//     getHandler *GetUserHandler,
	//     listHandler *ListUsersHandler) {
	//     cmdBus.Register(CreateUserCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdateUserCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeleteUserCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetUserQuery{}, getHandler.Handle)
	//     queryBus.Register(ListUsersQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration registers the User table migration.
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateUser(db)
}
