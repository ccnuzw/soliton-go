package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm" // required for fx.Invoke with *gorm.DB

	"github.com/soliton-go/framework/core/config"
	"github.com/soliton-go/framework/core/logger"
	"github.com/soliton-go/framework/orm"

		userapp "github.com/soliton-go/test-project/internal/application/user"
		interfaceshttp "github.com/soliton-go/test-project/internal/interfaces/http"
		productapp "github.com/soliton-go/test-project/internal/application/product"
	// soliton-gen:imports
)

func main() {
	fx.New(
		fx.Provide(
			config.NewConfig,
			logger.NewLogger,
			orm.NewGormDB,
			NewRouter,
		),

				userapp.Module,
				productapp.Module,
		// soliton-gen:modules

				fx.Provide(interfaceshttp.NewUserHandler),
				fx.Provide(interfaceshttp.NewProductHandler),
		// soliton-gen:handlers

				fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.UserHandler) {
			userapp.RegisterMigration(db)
			h.RegisterRoutes(r)
		}),
				fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.ProductHandler) {
			productapp.RegisterMigration(db)
			h.RegisterRoutes(r)
		}),
		// soliton-gen:routes

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
