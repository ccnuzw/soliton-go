package userapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/user"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 User 的所有 Fx 依赖。
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
	
	fx.Provide(NewUserService),
	// soliton-gen:services

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
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

// RegisterMigration 注册 User 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateUser(db)
}
