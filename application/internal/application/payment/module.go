package paymentapp

import "go.uber.org/fx"

// Module provides application service dependencies for Fx.
var Module = fx.Options(
	fx.Provide(NewPaymentService),
	// soliton-gen:services
)
