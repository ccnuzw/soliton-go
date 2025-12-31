package main

import (
	"context"

	"github.com/soliton-go/framework/core/config"
	"github.com/soliton-go/framework/core/logger"
	"github.com/soliton-go/framework/orm"
	"github.com/soliton-go/framework/web"
	"go.uber.org/fx"

	userapp "github.com/soliton-go/application/internal/application/user"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
)

func main() {
	fx.New(
		fx.Provide(
			// Framework modules
			config.NewConfig,
			logger.NewLogger,
			orm.NewGormDB,
			web.NewServer,

			// Application modules
			persistence.NewUserRepository,
			userapp.NewCreateUserHandler,
		),
		fx.Invoke(func(lc fx.Lifecycle, server *web.Server) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go server.Run(":8080")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	).Run()
}
