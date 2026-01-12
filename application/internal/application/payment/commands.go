package paymentapp

import (
	"context"
	"time"

	"github.com/soliton-go/application/internal/domain/payment"
)

// CreatePaymentCommand 是创建 Payment 的命令。
type CreatePaymentCommand struct {
	ID string
	OrderId string
	UserId string
	Amount float64
	Currency string
	Method payment.PaymentMethod
	Status payment.PaymentStatus
	Provider string
	ProviderTxnId string
	PaidAt *time.Time
	RefundedAt *time.Time
	FailureReason string
	Metadata datatypes.JSON
}

// CreatePaymentHandler 处理 CreatePaymentCommand。
type CreatePaymentHandler struct {
	repo payment.PaymentRepository
	service *payment.PaymentDomainService
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreatePaymentHandler(repo payment.PaymentRepository, service *payment.PaymentDomainService) *CreatePaymentHandler {
	return &CreatePaymentHandler{repo: repo, service: service}
}

func (h *CreatePaymentHandler) Handle(ctx context.Context, cmd CreatePaymentCommand) (*payment.Payment, error) {
	entity := payment.NewPayment(cmd.ID, cmd.OrderId, cmd.UserId, cmd.Amount, cmd.Currency, cmd.Method, cmd.Status, cmd.Provider, cmd.ProviderTxnId, cmd.PaidAt, cmd.RefundedAt, cmd.FailureReason, cmd.Metadata)
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

// UpdatePaymentCommand 是更新 Payment 的命令。
type UpdatePaymentCommand struct {
	ID string
	OrderId *string
	UserId *string
	Amount *float64
	Currency *string
	Method *payment.PaymentMethod
	Status *payment.PaymentStatus
	Provider *string
	ProviderTxnId *string
	PaidAt *time.Time
	RefundedAt *time.Time
	FailureReason *string
	Metadata *datatypes.JSON
}

// UpdatePaymentHandler 处理 UpdatePaymentCommand。
type UpdatePaymentHandler struct {
	repo payment.PaymentRepository
	service *payment.PaymentDomainService
}

func NewUpdatePaymentHandler(repo payment.PaymentRepository, service *payment.PaymentDomainService) *UpdatePaymentHandler {
	return &UpdatePaymentHandler{repo: repo, service: service}
}

func (h *UpdatePaymentHandler) Handle(ctx context.Context, cmd UpdatePaymentCommand) (*payment.Payment, error) {
	entity, err := h.repo.Find(ctx, payment.PaymentID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.OrderId, cmd.UserId, cmd.Amount, cmd.Currency, cmd.Method, cmd.Status, cmd.Provider, cmd.ProviderTxnId, cmd.PaidAt, cmd.RefundedAt, cmd.FailureReason, cmd.Metadata)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// DeletePaymentCommand 是删除 Payment 的命令。
type DeletePaymentCommand struct {
	ID string
}

// DeletePaymentHandler 处理 DeletePaymentCommand。
type DeletePaymentHandler struct {
	repo payment.PaymentRepository
	service *payment.PaymentDomainService
}

func NewDeletePaymentHandler(repo payment.PaymentRepository, service *payment.PaymentDomainService) *DeletePaymentHandler {
	return &DeletePaymentHandler{repo: repo, service: service}
}

func (h *DeletePaymentHandler) Handle(ctx context.Context, cmd DeletePaymentCommand) error {
	return h.repo.Delete(ctx, payment.PaymentID(cmd.ID))
}
