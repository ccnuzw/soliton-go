package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/soliton-go/framework/core/config"
	"github.com/soliton-go/framework/core/logger"
	"github.com/soliton-go/framework/orm"

	// Import your modules here:
	// userapp "github.com/soliton-go/test-project/internal/application/user"
	// interfaceshttp "github.com/soliton-go/test-project/internal/interfaces/http"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewConfig,
			logger.NewLogger,
			orm.NewGormDB,
			NewRouter,
		),

		// Modules - uncomment after generating domains:
		// userapp.Module,

		// HTTP Handlers - uncomment after generating domains:
		// fx.Provide(interfaceshttp.NewUserHandler),

		// Register routes and migrations - uncomment after generating domains:
		// fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.UserHandler) {
		// 	userapp.RegisterMigration(db)
		// 	h.RegisterRoutes(r)
		// }),

		// Start server
		fx.Invoke(StartServer),
	).Run()
}

// NewRouter creates the Gin engine and registers base routes.
func NewRouter() *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

// StartServer starts the HTTP server with Fx lifecycle.
func StartServer(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger, r *gin.Engine) {
	addr := fmt.Sprintf("%s:%d", cfg.GetString("server.host"), cfg.GetInt("server.port"))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("server starting", zap.String("addr", addr))
			go func() {
				if err := r.Run(addr); err != nil {
					logger.Fatal("server stopped", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
