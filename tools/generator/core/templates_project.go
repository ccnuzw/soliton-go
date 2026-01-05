package core

// ============================================================================
// PROJECT TEMPLATES
// ============================================================================

const GoModTemplate = `module {{.ModuleName}}

go {{.GoVersion}}

require (
	github.com/soliton-go/framework {{.FrameworkVersion}}
	github.com/gin-gonic/gin v1.11.0
	github.com/google/uuid v1.6.0
	go.uber.org/fx {{.FxVersion}}
	gorm.io/gorm v1.31.1
)
{{ if .FrameworkReplace }}

replace github.com/soliton-go/framework => {{.FrameworkReplace}}
{{ end }}
`

const MainTemplate = `package main

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

		// soliton-gen:modules

		// soliton-gen:handlers

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
`

const ConfigTemplate = `server:
  host: 0.0.0.0
  port: 8080

database:
  driver: sqlite
  dsn: data.db

log:
  level: info
`

const ConfigExampleTemplate = `# Server Configuration
server:
  host: 0.0.0.0
  port: 8080

# Database Configuration
database:
  # Options: sqlite, postgres
  driver: sqlite
  dsn: data.db

  # PostgreSQL example:
  # driver: postgres
  # dsn: host=localhost user=postgres password=secret dbname=myapp port=5432 sslmode=disable

  # MySQL is not enabled by default. Extend framework/orm/db.go if needed.

# Logging
log:
  level: info  # debug, info, warn, error
`

const ResponseTemplate = `package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误码定义
const (
	CodeSuccess      = 0     // 成功
	CodeBadRequest   = 400   // 请求错误
	CodeUnauthorized = 401   // 未授权
	CodeForbidden    = 403   // 禁止访问
	CodeNotFound     = 404   // 资源不存在
	CodeInternal     = 500   // 服务器内部错误

	// 业务错误码 (1000+)
	CodeValidation   = 1001  // 校验失败
	CodeDuplicate    = 1002  // 重复数据
	CodeConflict     = 1003  // 业务冲突
)

// Response 是标准的 API 响应结构体。
type Response struct {
	Code    int         ` + "`json:\"code\"`" + `
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// BadRequest 返回 400 错误响应。
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeBadRequest,
		Message: message,
	})
}

// NotFound 返回 404 错误响应。
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// InternalError 返回 500 错误响应。
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeInternal,
		Message: message,
	})
}

// ValidationError 返回校验错误响应。
func ValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeValidation,
		Message: message,
	})
}
`

const GitignoreTemplate = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output
*.out

# Dependency directories
vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo

# Build
/bin/
/dist/

# Database
*.db
*.sqlite

# Config (keep example)
configs/config.yaml

# Logs
*.log
logs/

# OS
.DS_Store
Thumbs.db
`

const ReadmeTemplate = `# {{.ProjectName}}

A Go project built with [Soliton-Go](https://github.com/soliton-go/framework) framework.

## Quick Start

` + "```bash" + `
# Install dependencies
GOWORK=off go mod tidy

# Generate domain modules (--wire auto-injects into main.go)
soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

# Enable soft delete (optional)
soliton-gen domain User --fields "username,email" --soft-delete --wire

# Run the server
GOWORK=off go run ./cmd/main.go
` + "```" + `

## Project Structure

` + "```" + `
{{.ProjectName}}/
├── cmd/main.go              # Entry point
├── configs/                 # Configuration
├── internal/
│   ├── domain/              # Domain layer (entities, repos, events)
│   ├── application/         # Application layer (commands, queries)
│   ├── infrastructure/      # Infrastructure layer (repo implementations)
│   └── interfaces/          # Interface layer (HTTP handlers)
└── go.mod
` + "```" + `

## API Endpoints

After generating domains, the following endpoints are available:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |

Domain specific endpoints will be available after generating domains using ` + "`soliton-gen domain`" + `.

### Pagination

List endpoints support pagination:

` + "```bash" + `
curl "http://localhost:8080/api/users?page=1&page_size=20"
` + "```" + `

Response:
` + "```json" + `
{
  "items": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
` + "```" + `

> **Note**: If running in a monorepo with go.work, use ` + "`GOWORK=off`" + ` prefix for go commands.
`

const MakefileTemplate = `.PHONY: run build test clean gen tidy

# Disable go.work by default for monorepo compatibility (override with GOWORK=on).
GOWORK ?= off

# Run the application
run:
	GOWORK=$(GOWORK) go run ./cmd/main.go

# Build the application
build:
	GOWORK=$(GOWORK) go build -o bin/app ./cmd/main.go

# Run tests
test:
	GOWORK=$(GOWORK) go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Generate a domain module
# Usage: make gen NAME=User FIELDS="username,email,status:enum(active|inactive)"
gen:
	soliton-gen domain $(NAME) --fields "$(FIELDS)"

# Tidy dependencies
tidy:
	GOWORK=$(GOWORK) go mod tidy
`
