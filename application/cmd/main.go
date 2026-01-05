package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/soliton-go/framework/core/config"
	"github.com/soliton-go/framework/core/logger"
	"github.com/soliton-go/framework/orm"

	userapp "github.com/soliton-go/application/internal/application/user"
	interfaceshttp "github.com/soliton-go/application/internal/interfaces/http"
	orderapp "github.com/soliton-go/application/internal/application/order"
	productapp "github.com/soliton-go/application/internal/application/product"
	"github.com/soliton-go/framework/event"
	// soliton-gen:imports
)

func main() {
	fx.New(
		fx.Provide(
			config.NewConfig,
			logger.NewLogger,
			orm.NewGormDB,
			func() event.EventBus { return event.NewLocalEventBus() },
		// soliton-gen:providers
			NewRouter,
		),

		userapp.Module,
		orderapp.Module,
		productapp.Module,
		// soliton-gen:modules

		fx.Provide(interfaceshttp.NewUserHandler),
		fx.Provide(interfaceshttp.NewOrderHandler),
		fx.Provide(interfaceshttp.NewProductHandler),
		// soliton-gen:handlers

		fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.UserHandler) error {
			if err := userapp.RegisterMigration(db); err != nil {
				return err
			}
			h.RegisterRoutes(r)
			return nil
		}),
		fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.OrderHandler) error {
			if err := orderapp.RegisterMigration(db); err != nil {
				return err
			}
			h.RegisterRoutes(r)
			return nil
		}),
		fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.ProductHandler) error {
			if err := productapp.RegisterMigration(db); err != nil {
				return err
			}
			h.RegisterRoutes(r)
			return nil
		}),
		// soliton-gen:routes

		// 启动服务器
		fx.Invoke(StartServer),
	).Run()
}

// NewRouter 创建 Gin 引擎并注册基础路由。
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

// StartServer 启动 HTTP 服务器（带 Fx 生命周期管理）。
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
