package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain [name]",
	Short: "Generate a complete domain module ready for production",
	Long: `Generate a complete domain entity including:
  - Entity with ID and aggregate root
  - Repository interface and GORM implementation
  - Domain events (Created, Updated, Deleted)
  - Application layer (Commands, Queries, DTOs)
  - HTTP Handler with CRUD endpoints
  - Fx dependency injection module
  - Database migration support`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("🚀 Generating domain: %s\n\n", name)

		generateDomain(name)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
}

func generateDomain(name string) {
	// Normalize name: first letter uppercase
	entityName := strings.Title(name)
	packageName := strings.ToLower(name)

	// Determine base paths
	baseDir := findPath("application/internal/domain")
	infraDir := findPath("application/internal/infrastructure/persistence")
	appDir := findPath("application/internal/application")
	interfacesDir := findPath("application/internal/interfaces/http")

	data := map[string]string{
		"PackageName": packageName,
		"EntityName":  entityName,
	}

	// === Domain Layer ===
	fmt.Println("📦 Domain Layer")
	domainDir := filepath.Join(baseDir, packageName)
	_ = os.MkdirAll(domainDir, 0755)

	generateFile(filepath.Join(domainDir, packageName+".go"), entityTemplate, data)
	generateFile(filepath.Join(domainDir, "repository.go"), repoTemplate, data)
	generateFile(filepath.Join(domainDir, "events.go"), eventsTemplate, data)

	// === Infrastructure Layer ===
	fmt.Println("\n🔧 Infrastructure Layer")
	_ = os.MkdirAll(infraDir, 0755)
	generateFile(filepath.Join(infraDir, packageName+"_repo.go"), repoImplTemplate, data)

	// === Application Layer ===
	fmt.Println("\n⚙️ Application Layer")
	appModuleDir := filepath.Join(appDir, packageName)
	_ = os.MkdirAll(appModuleDir, 0755)

	generateFile(filepath.Join(appModuleDir, "commands.go"), commandsTemplate, data)
	generateFile(filepath.Join(appModuleDir, "queries.go"), queriesTemplate, data)
	generateFile(filepath.Join(appModuleDir, "dto.go"), dtoTemplate, data)

	// === Interfaces Layer (HTTP) ===
	fmt.Println("\n🌐 Interfaces Layer")
	_ = os.MkdirAll(interfacesDir, 0755)
	generateFile(filepath.Join(interfacesDir, packageName+"_handler.go"), handlerTemplate, data)

	// === Fx Module ===
	fmt.Println("\n📌 Fx Module")
	generateFile(filepath.Join(appModuleDir, "module.go"), fxModuleTemplate, data)

	// === Summary ===
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("✅ Domain generation complete!")
	fmt.Printf("   📁 Domain:      %s\n", domainDir)
	fmt.Printf("   📁 Application: %s\n", appModuleDir)
	fmt.Printf("   📁 Persistence: %s\n", infraDir)
	fmt.Printf("   📁 HTTP:        %s\n", interfacesDir)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("\n📝 Next steps:")
	fmt.Println("   1. Add entity fields in " + packageName + ".go")
	fmt.Println("   2. Update DTO fields in dto.go")
	fmt.Println("   3. Import module in main.go:")
	fmt.Printf("      %sapp.Module,\n", packageName)
}

func findPath(relative string) string {
	// Try from tools/generator
	path := "../../" + relative
	if _, err := os.Stat(path); err == nil {
		return path
	}
	// Try from project root
	return relative
}

// generateFile creates a file from template if it doesn't exist (Lock mechanism)
func generateFile(path string, tmpl string, data interface{}) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("   [SKIP] %s (exists)\n", filepath.Base(path))
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
		fmt.Printf("   [NEW] %s\n", filepath.Base(path))
	}
}

// ============================================================================
// TEMPLATES
// ============================================================================

const entityTemplate = `package {{.PackageName}}

import "github.com/soliton-go/framework/ddd"

// {{.EntityName}}ID is a strong typed ID.
type {{.EntityName}}ID string

func (id {{.EntityName}}ID) String() string {
	return string(id)
}

// {{.EntityName}} is the aggregate root.
type {{.EntityName}} struct {
	ddd.BaseAggregateRoot
	ID   {{.EntityName}}ID ` + "`gorm:\"primaryKey\"`" + `
	Name string            ` + "`gorm:\"size:255\"`" + `
	// TODO: Add more fields here
}

// TableName returns the table name for GORM.
func ({{.EntityName}}) TableName() string {
	return "{{.PackageName}}s"
}

// New{{.EntityName}} creates a new {{.EntityName}}.
func New{{.EntityName}}(id, name string) *{{.EntityName}} {
	e := &{{.EntityName}}{
		ID:   {{.EntityName}}ID(id),
		Name: name,
	}
	e.AddDomainEvent(New{{.EntityName}}CreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *{{.EntityName}}) Update(name string) {
	e.Name = name
	e.AddDomainEvent(New{{.EntityName}}UpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *{{.EntityName}}) GetID() ddd.ID {
	return e.ID
}
`

const repoTemplate = `package {{.PackageName}}

import (
	"github.com/soliton-go/framework/orm"
)

// {{.EntityName}}Repository is the interface for {{.EntityName}} persistence.
type {{.EntityName}}Repository interface {
	orm.Repository[*{{.EntityName}}, {{.EntityName}}ID]
	// TODO: Add custom query methods here
}
`

const eventsTemplate = `package {{.PackageName}}

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// {{.EntityName}}CreatedEvent is published when a new {{.EntityName}} is created.
type {{.EntityName}}CreatedEvent struct {
	ddd.BaseDomainEvent
	{{.EntityName}}ID string ` + "`json:\"{{.PackageName}}_id\"`" + `
}

func (e {{.EntityName}}CreatedEvent) EventName() string {
	return "{{.PackageName}}.created"
}

func New{{.EntityName}}CreatedEvent(id string) {{.EntityName}}CreatedEvent {
	return {{.EntityName}}CreatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		{{.EntityName}}ID: id,
	}
}

// {{.EntityName}}UpdatedEvent is published when a {{.EntityName}} is updated.
type {{.EntityName}}UpdatedEvent struct {
	ddd.BaseDomainEvent
	{{.EntityName}}ID string ` + "`json:\"{{.PackageName}}_id\"`" + `
}

func (e {{.EntityName}}UpdatedEvent) EventName() string {
	return "{{.PackageName}}.updated"
}

func New{{.EntityName}}UpdatedEvent(id string) {{.EntityName}}UpdatedEvent {
	return {{.EntityName}}UpdatedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		{{.EntityName}}ID: id,
	}
}

// {{.EntityName}}DeletedEvent is published when a {{.EntityName}} is deleted.
type {{.EntityName}}DeletedEvent struct {
	ddd.BaseDomainEvent
	{{.EntityName}}ID string    ` + "`json:\"{{.PackageName}}_id\"`" + `
	DeletedAt         time.Time ` + "`json:\"deleted_at\"`" + `
}

func (e {{.EntityName}}DeletedEvent) EventName() string {
	return "{{.PackageName}}.deleted"
}

func New{{.EntityName}}DeletedEvent(id string) {{.EntityName}}DeletedEvent {
	return {{.EntityName}}DeletedEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		{{.EntityName}}ID: id,
		DeletedAt: time.Now(),
	}
}

// init registers events with the global registry.
func init() {
	event.RegisterEvent("{{.PackageName}}.created", func() ddd.DomainEvent {
		return &{{.EntityName}}CreatedEvent{}
	})
	event.RegisterEvent("{{.PackageName}}.updated", func() ddd.DomainEvent {
		return &{{.EntityName}}UpdatedEvent{}
	})
	event.RegisterEvent("{{.PackageName}}.deleted", func() ddd.DomainEvent {
		return &{{.EntityName}}DeletedEvent{}
	})
}
`

const repoImplTemplate = `package persistence

import (
	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type {{.EntityName}}RepoImpl struct {
	*orm.GormRepository[*{{.PackageName}}.{{.EntityName}}, {{.PackageName}}.{{.EntityName}}ID]
	db *gorm.DB
}

func New{{.EntityName}}Repository(db *gorm.DB) {{.PackageName}}.{{.EntityName}}Repository {
	return &{{.EntityName}}RepoImpl{
		GormRepository: orm.NewGormRepository[*{{.PackageName}}.{{.EntityName}}, {{.PackageName}}.{{.EntityName}}ID](db),
		db:             db,
	}
}

// Migrate creates the table if it doesn't exist.
func Migrate{{.EntityName}}(db *gorm.DB) error {
	return db.AutoMigrate(&{{.PackageName}}.{{.EntityName}}{})
}
`

const commandsTemplate = `package {{.PackageName}}app

import (
	"context"

	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
)

// Create{{.EntityName}}Command is the command for creating a {{.EntityName}}.
type Create{{.EntityName}}Command struct {
	ID   string
	Name string
}

// Create{{.EntityName}}Handler handles Create{{.EntityName}}Command.
type Create{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewCreate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Create{{.EntityName}}Handler {
	return &Create{{.EntityName}}Handler{repo: repo}
}

func (h *Create{{.EntityName}}Handler) Handle(ctx context.Context, cmd Create{{.EntityName}}Command) (*{{.PackageName}}.{{.EntityName}}, error) {
	entity := {{.PackageName}}.New{{.EntityName}}(cmd.ID, cmd.Name)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// Update{{.EntityName}}Command is the command for updating a {{.EntityName}}.
type Update{{.EntityName}}Command struct {
	ID   string
	Name string
}

// Update{{.EntityName}}Handler handles Update{{.EntityName}}Command.
type Update{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewUpdate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Update{{.EntityName}}Handler {
	return &Update{{.EntityName}}Handler{repo: repo}
}

func (h *Update{{.EntityName}}Handler) Handle(ctx context.Context, cmd Update{{.EntityName}}Command) (*{{.PackageName}}.{{.EntityName}}, error) {
	entity, err := h.repo.Find(ctx, {{.PackageName}}.{{.EntityName}}ID(cmd.ID))
	if err != nil {
		return nil, err
	}
	entity.Update(cmd.Name)
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// Delete{{.EntityName}}Command is the command for deleting a {{.EntityName}}.
type Delete{{.EntityName}}Command struct {
	ID string
}

// Delete{{.EntityName}}Handler handles Delete{{.EntityName}}Command.
type Delete{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewDelete{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Delete{{.EntityName}}Handler {
	return &Delete{{.EntityName}}Handler{repo: repo}
}

func (h *Delete{{.EntityName}}Handler) Handle(ctx context.Context, cmd Delete{{.EntityName}}Command) error {
	return h.repo.Delete(ctx, {{.PackageName}}.{{.EntityName}}ID(cmd.ID))
}
`

const queriesTemplate = `package {{.PackageName}}app

import (
	"context"

	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
)

// Get{{.EntityName}}Query is the query for getting a single {{.EntityName}}.
type Get{{.EntityName}}Query struct {
	ID string
}

// Get{{.EntityName}}Handler handles Get{{.EntityName}}Query.
type Get{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewGet{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Get{{.EntityName}}Handler {
	return &Get{{.EntityName}}Handler{repo: repo}
}

func (h *Get{{.EntityName}}Handler) Handle(ctx context.Context, query Get{{.EntityName}}Query) (*{{.PackageName}}.{{.EntityName}}, error) {
	return h.repo.Find(ctx, {{.PackageName}}.{{.EntityName}}ID(query.ID))
}

// List{{.EntityName}}sQuery is the query for listing all {{.EntityName}}s.
type List{{.EntityName}}sQuery struct{}

// List{{.EntityName}}sHandler handles List{{.EntityName}}sQuery.
type List{{.EntityName}}sHandler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewList{{.EntityName}}sHandler(repo {{.PackageName}}.{{.EntityName}}Repository) *List{{.EntityName}}sHandler {
	return &List{{.EntityName}}sHandler{repo: repo}
}

func (h *List{{.EntityName}}sHandler) Handle(ctx context.Context, query List{{.EntityName}}sQuery) ([]*{{.PackageName}}.{{.EntityName}}, error) {
	return h.repo.FindAll(ctx)
}
`

const dtoTemplate = `package {{.PackageName}}app

import "github.com/soliton-go/application/internal/domain/{{.PackageName}}"

// Create{{.EntityName}}Request is the request body for creating a {{.EntityName}}.
type Create{{.EntityName}}Request struct {
	Name string ` + "`json:\"name\" binding:\"required\"`" + `
}

// Update{{.EntityName}}Request is the request body for updating a {{.EntityName}}.
type Update{{.EntityName}}Request struct {
	Name string ` + "`json:\"name\" binding:\"required\"`" + `
}

// {{.EntityName}}Response is the response body for {{.EntityName}} data.
type {{.EntityName}}Response struct {
	ID   string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

// To{{.EntityName}}Response converts entity to response.
func To{{.EntityName}}Response(e *{{.PackageName}}.{{.EntityName}}) {{.EntityName}}Response {
	return {{.EntityName}}Response{
		ID:   string(e.ID),
		Name: e.Name,
	}
}

// To{{.EntityName}}ResponseList converts entities to response list.
func To{{.EntityName}}ResponseList(entities []*{{.PackageName}}.{{.EntityName}}) []{{.EntityName}}Response {
	result := make([]{{.EntityName}}Response, len(entities))
	for i, e := range entities {
		result[i] = To{{.EntityName}}Response(e)
	}
	return result
}
`

const handlerTemplate = `package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	{{.PackageName}}app "github.com/soliton-go/application/internal/application/{{.PackageName}}"
)

// {{.EntityName}}Handler handles HTTP requests for {{.EntityName}} operations.
type {{.EntityName}}Handler struct {
	createHandler *{{.PackageName}}app.Create{{.EntityName}}Handler
	updateHandler *{{.PackageName}}app.Update{{.EntityName}}Handler
	deleteHandler *{{.PackageName}}app.Delete{{.EntityName}}Handler
	getHandler    *{{.PackageName}}app.Get{{.EntityName}}Handler
	listHandler   *{{.PackageName}}app.List{{.EntityName}}sHandler
}

// New{{.EntityName}}Handler creates a new {{.EntityName}}Handler.
func New{{.EntityName}}Handler(
	createHandler *{{.PackageName}}app.Create{{.EntityName}}Handler,
	updateHandler *{{.PackageName}}app.Update{{.EntityName}}Handler,
	deleteHandler *{{.PackageName}}app.Delete{{.EntityName}}Handler,
	getHandler *{{.PackageName}}app.Get{{.EntityName}}Handler,
	listHandler *{{.PackageName}}app.List{{.EntityName}}sHandler,
) *{{.EntityName}}Handler {
	return &{{.EntityName}}Handler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// RegisterRoutes registers {{.EntityName}} routes.
func (h *{{.EntityName}}Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/{{.PackageName}}s")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /api/{{.PackageName}}s
func (h *{{.EntityName}}Handler) Create(c *gin.Context) {
	var req {{.PackageName}}app.Create{{.EntityName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := {{.PackageName}}app.Create{{.EntityName}}Command{
		ID:   uuid.New().String(),
		Name: req.Name,
	}

	entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, {{.PackageName}}app.To{{.EntityName}}Response(entity))
}

// Get handles GET /api/{{.PackageName}}s/:id
func (h *{{.EntityName}}Handler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), {{.PackageName}}app.Get{{.EntityName}}Query{ID: id})
	if err != nil {
		NotFound(c, "{{.PackageName}} not found")
		return
	}

	Success(c, {{.PackageName}}app.To{{.EntityName}}Response(entity))
}

// List handles GET /api/{{.PackageName}}s
func (h *{{.EntityName}}Handler) List(c *gin.Context) {
	entities, err := h.listHandler.Handle(c.Request.Context(), {{.PackageName}}app.List{{.EntityName}}sQuery{})
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, {{.PackageName}}app.To{{.EntityName}}ResponseList(entities))
}

// Update handles PUT /api/{{.PackageName}}s/:id
func (h *{{.EntityName}}Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req {{.PackageName}}app.Update{{.EntityName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	cmd := {{.PackageName}}app.Update{{.EntityName}}Command{
		ID:   id,
		Name: req.Name,
	}

	entity, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, {{.PackageName}}app.To{{.EntityName}}Response(entity))
}

// Delete handles DELETE /api/{{.PackageName}}s/:id
func (h *{{.EntityName}}Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	cmd := {{.PackageName}}app.Delete{{.EntityName}}Command{ID: id}
	if err := h.deleteHandler.Handle(c.Request.Context(), cmd); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}
`

const fxModuleTemplate = `package {{.PackageName}}app

import (
	"go.uber.org/fx"

	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
	"github.com/soliton-go/application/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// Module provides all {{.EntityName}} dependencies for Fx.
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) {{.PackageName}}.{{.EntityName}}Repository {
		return persistence.New{{.EntityName}}Repository(db)
	}),

	// Command Handlers
	fx.Provide(NewCreate{{.EntityName}}Handler),
	fx.Provide(NewUpdate{{.EntityName}}Handler),
	fx.Provide(NewDelete{{.EntityName}}Handler),

	// Query Handlers
	fx.Provide(NewGet{{.EntityName}}Handler),
	fx.Provide(NewList{{.EntityName}}sHandler),
)

// RegisterMigration registers the {{.EntityName}} table migration.
func RegisterMigration(db *gorm.DB) error {
	return persistence.Migrate{{.EntityName}}(db)
}
`
