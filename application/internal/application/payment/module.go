package paymentapp

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/payment"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module 提供 Payment 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) payment.PaymentRepository {
		return persistence.NewPaymentRepository(db)
	}),

	// Domain Services
	fx.Provide(payment.NewPaymentDomainService),

	// Command Handlers
	fx.Provide(NewCreatePaymentHandler),
	fx.Provide(NewUpdatePaymentHandler),
	fx.Provide(NewDeletePaymentHandler),

	// Query Handlers
	fx.Provide(NewGetPaymentHandler),
	fx.Provide(NewListPaymentsHandler),
	
	fx.Provide(NewPaymentService),
	// soliton-gen:services
	// soliton-gen:event-handlers

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *CreatePaymentHandler,
	//     updateHandler *UpdatePaymentHandler,
	//     deleteHandler *DeletePaymentHandler,
	//     getHandler *GetPaymentHandler,
	//     listHandler *ListPaymentsHandler) {
	//     cmdBus.Register(CreatePaymentCommand{}, createHandler.Handle)
	//     cmdBus.Register(UpdatePaymentCommand{}, updateHandler.Handle)
	//     cmdBus.Register(DeletePaymentCommand{}, deleteHandler.Handle)
	//     queryBus.Register(GetPaymentQuery{}, getHandler.Handle)
	//     queryBus.Register(ListPaymentsQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration 注册 Payment 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.MigratePayment(db)
}
