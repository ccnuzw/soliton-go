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
}

func initProject(projectName, module string) {
	// Create project root directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		fmt.Printf("❌ Failed to create project directory: %v\n", err)
		return
	}

	data := map[string]string{
		"ProjectName": projectName,
		"ModuleName":  module,
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
	fmt.Println("   go mod tidy")
	fmt.Println("   soliton-gen domain User --fields \"username,email,status:enum(active|inactive)\"")
	fmt.Println("   go run ./cmd/main.go")
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
	github.com/soliton-go/framework v0.1.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.6.0
	go.uber.org/fx v1.22.0
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gorm v1.25.7
)
`

const mainTemplate = `package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// Import your modules here:
	// userapp "{{.ModuleName}}/internal/application/user"
	// "{{.ModuleName}}/internal/interfaces/http"
)

func main() {
	fx.New(
		// Database
		fx.Provide(NewDB),

		// Modules - uncomment after generating domains:
		// userapp.Module,

		// HTTP Handlers - uncomment after generating domains:
		// fx.Provide(http.NewUserHandler),

		// Start server
		fx.Invoke(StartServer),
	).Run()
}

// NewDB creates a new GORM database connection.
func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("✅ Database connected")
	return db
}

// StartServer starts the HTTP server.
func StartServer(db *gorm.DB) {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register routes - uncomment after generating domains:
	// userHandler.RegisterRoutes(r)

	log.Println("🚀 Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
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
go mod tidy

# Generate domain modules
soliton-gen domain User --fields "username,email,status:enum(active|inactive)"

# Run the server
go run ./cmd/main.go
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
| DELETE | /api/users/:id | Delete user |
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
