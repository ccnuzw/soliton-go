package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var moduleName string
var frameworkVersion string
var frameworkReplace string

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Soliton-Go project",
	Long: `Initialize a new Soliton-Go project with complete directory structure:
  - cmd/main.go (Entry point with Fx setup)
  - configs/config.yaml (Configuration)
  - internal/ (DDD layer structure)
  - go.mod (Dependencies)

Examples:
  soliton-gen init my-project
  soliton-gen init my-project --module github.com/myorg/my-project`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		if moduleName == "" {
			moduleName = "github.com/soliton-go/" + projectName
		}
		fmt.Printf("🚀 Initializing project: %s\n", projectName)
		fmt.Printf("   Module: %s\n\n", moduleName)

		initProject(projectName, moduleName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Go module name (default: github.com/soliton-go/<project-name>)")
	initCmd.Flags().StringVar(&frameworkVersion, "framework-version", "", "Framework version (default: auto)")
	initCmd.Flags().StringVar(&frameworkReplace, "framework-replace", "", "Replace github.com/soliton-go/framework with a local path")
}

func initProject(projectName, module string) {
	// Create project root directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		fmt.Printf("❌ Failed to create project directory: %v\n", err)
		return
	}

	frameworkVersionValue := frameworkVersion
	frameworkReplaceValue := frameworkReplace

	if frameworkReplaceValue == "" {
		if info, err := os.Stat("framework"); err == nil && info.IsDir() {
			frameworkReplaceValue = filepath.ToSlash(filepath.Join("..", "framework"))
		}
	}

	if frameworkVersionValue == "" {
		if frameworkReplaceValue != "" {
			frameworkVersionValue = "v0.0.0-00010101000000-000000000000"
		} else {
			frameworkVersionValue = "v0.1.0"
		}
	}

	data := map[string]string{
		"ProjectName":      projectName,
		"ModuleName":       module,
		"FrameworkVersion": frameworkVersionValue,
		"FrameworkReplace": frameworkReplaceValue,
	}

	// === Create directory structure ===
	dirs := []string{
		"cmd",
		"configs",
		"internal/domain",
		"internal/application",
		"internal/infrastructure/persistence",
		"internal/interfaces/http",
	}

	fmt.Println("📁 Creating directories")
	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("   [ERROR] %s: %v\n", dir, err)
		} else {
			fmt.Printf("   [DIR] %s\n", dir)
		}
	}

	// === Create files ===
	fmt.Println("\n📄 Creating files")

	files := []struct {
		path     string
		template string
	}{
		{"go.mod", goModTemplate},
		{"cmd/main.go", mainTemplate},
		{"configs/config.yaml", configTemplate},
		{"configs/config.example.yaml", configExampleTemplate},
		{"internal/interfaces/http/response.go", responseTemplate},
		{".gitignore", gitignoreTemplate},
		{"README.md", readmeTemplate},
		{"Makefile", makefileTemplate},
	}

	for _, f := range files {
		path := filepath.Join(projectName, f.path)
		generateInitFile(path, f.template, data)
	}

	// === Summary ===
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("✅ Project initialized successfully!")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n📝 Next steps:")
	fmt.Printf("   cd %s\n", projectName)
	fmt.Println("   GOWORK=off go mod tidy   # skip go.work if in monorepo")
	fmt.Println("   soliton-gen domain User --fields \"username,email,status:enum(active|inactive)\"")
	fmt.Println("   GOWORK=off go run ./cmd/main.go")
}

func generateInitFile(path string, tmpl string, data map[string]string) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("   [SKIP] %s (exists)\n", filepath.Base(path))
		return
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("   [ERROR] %s: %v\n", filepath.Base(path), err)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("   [ERROR] %s: %v\n", filepath.Base(path), err)
		return
	}
	defer f.Close()

	t := template.Must(template.New("file").Parse(tmpl))
	if err := t.Execute(f, data); err != nil {
		fmt.Printf("   [ERROR] %s: %v\n", filepath.Base(path), err)
	} else {
		fmt.Printf("   [NEW] %s\n", path)
	}
}

// ============================================================================
// INIT TEMPLATES
// ============================================================================

const goModTemplate = `module {{.ModuleName}}

go 1.22

require (
	github.com/soliton-go/framework {{.FrameworkVersion}}
	github.com/gin-gonic/gin v1.11.0
	github.com/google/uuid v1.6.0
	go.uber.org/fx v1.24.0
	gorm.io/gorm v1.31.1
)
{{ if .FrameworkReplace }}

replace github.com/soliton-go/framework => {{.FrameworkReplace}}
{{ end }}
`

const mainTemplate = `package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	_ "gorm.io/gorm" // required for fx.Invoke with *gorm.DB

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
`

const configTemplate = `server:
  host: 0.0.0.0
  port: 8080

database:
  driver: sqlite
  dsn: data.db

log:
  level: info
`

const configExampleTemplate = `# Server Configuration
server:
  host: 0.0.0.0
  port: 8080

# Database Configuration
database:
  # Options: sqlite, postgres, mysql
  driver: sqlite
  dsn: data.db

  # PostgreSQL example:
  # driver: postgres
  # dsn: host=localhost user=postgres password=secret dbname=myapp port=5432 sslmode=disable

  # MySQL example:
  # driver: mysql
  # dsn: user:password@tcp(127.0.0.1:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local

# Logging
log:
  level: info  # debug, info, warn, error
`

const responseTemplate = `package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response.
type Response struct {
	Code    int         ` + "`json:\"code\"`" + `
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
}

// Success returns a successful response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// BadRequest returns a 400 error response.
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
	})
}

// NotFound returns a 404 error response.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
	})
}

// InternalError returns a 500 error response.
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
	})
}
`

const gitignoreTemplate = `# Binaries
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

const readmeTemplate = `# {{.ProjectName}}

A Go project built with [Soliton-Go](https://github.com/soliton-go/framework) framework.

## Quick Start

` + "```bash" + `
# Install dependencies
GOWORK=off go mod tidy

# Generate domain modules (--wire auto-injects into main.go)
soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

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
| POST | /api/users | Create user |
| GET | /api/users | List users |
| GET | /api/users/:id | Get user |
| PUT | /api/users/:id | Update user |
| PATCH | /api/users/:id | Partial update user |
| DELETE | /api/users/:id | Delete user |

> **Note**: If running in a monorepo with go.work, use ` + "`GOWORK=off`" + ` prefix for go commands.
`

const makefileTemplate = `.PHONY: run build test clean gen

# Run the application
run:
	go run ./cmd/main.go

# Build the application
build:
	go build -o bin/app ./cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Generate a domain module
# Usage: make gen NAME=User FIELDS="username,email,status:enum(active|inactive)"
gen:
	soliton-gen domain $(NAME) --fields "$(FIELDS)"

# Tidy dependencies
tidy:
	go mod tidy
`
