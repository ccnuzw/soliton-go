package promotionapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/promotion"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 Promotion 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) promotion.PromotionRepository {
		return persistence.NewPromotionRepository(db)
	}),

	// Domain Services
	fx.Provide(promotion.NewPromotionDomainService),

	// Command Handlers
	fx.Provide(NewCreatePromotionHandler),
	fx.Provide(NewUpdatePromotionHandler),
	fx.Provide(NewDeletePromotionHandler),

	// Query Handlers
	fx.Provide(NewGetPromotionHandler),
	fx.Provide(NewListPromotionsHandler),
	
	fx.Provide(NewPromotionService),
	// soliton-gen:services
	// soliton-gen:event-handlers

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreatePromotionHandler,
	//     updateHandler *UpdatePromotionHandler,
	//     deleteHandler *DeletePromotionHandler,
	//     getHandler *GetPromotionHandler,
	//     listHandler *ListPromotionsHandler) {
	//     cmdBus.Register(CreatePromotionCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdatePromotionCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeletePromotionCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetPromotionQuery{}, getHandler.Handle)
	//     queryBus.Register(ListPromotionsQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration 注册 Promotion 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigratePromotion(db)
}
