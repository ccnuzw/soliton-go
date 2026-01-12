package reviewapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/review"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 Review 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) review.ReviewRepository {
		return persistence.NewReviewRepository(db)
	}),

	// Domain Services
	fx.Provide(review.NewReviewDomainService),

	// Command Handlers
	fx.Provide(NewCreateReviewHandler),
	fx.Provide(NewUpdateReviewHandler),
	fx.Provide(NewDeleteReviewHandler),

	// Query Handlers
	fx.Provide(NewGetReviewHandler),
	fx.Provide(NewListReviewsHandler),
	
	fx.Provide(NewReviewService),
	// soliton-gen:services
	// soliton-gen:event-handlers

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreateReviewHandler,
	//     updateHandler *UpdateReviewHandler,
	//     deleteHandler *DeleteReviewHandler,
	//     getHandler *GetReviewHandler,
	//     listHandler *ListReviewsHandler) {
	//     cmdBus.Register(CreateReviewCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdateReviewCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeleteReviewCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetReviewQuery{}, getHandler.Handle)
	//     queryBus.Register(ListReviewsQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration 注册 Review 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigrateReview(db)
}
