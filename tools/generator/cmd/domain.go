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
	Short: "Generate a new domain entity with full DDD structure",
	Long: `Generate a complete domain entity including:
  - Entity with ID and aggregate root
  - Repository interface
  - Domain events (Created, Updated, Deleted)
  - Application layer (Commands and Queries)
  - Infrastructure implementation`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Generating domain entity: %s\n", name)

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

	// Try to find application directory
	baseDir := "../../application/internal/domain"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = "application/internal/domain"
	}

	targetDir := filepath.Join(baseDir, packageName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
		return
	}

	data := map[string]string{
		"PackageName": packageName,
		"EntityName":  entityName,
	}

	// 1. Generate Entity
	entityFile := filepath.Join(targetDir, packageName+".go")
	generateFile(entityFile, entityTemplate, data)

	// 2. Generate Repository
	repoFile := filepath.Join(targetDir, "repository.go")
	generateFile(repoFile, repoTemplate, data)

	// 3. Generate Domain Events
	eventsFile := filepath.Join(targetDir, "events.go")
	generateFile(eventsFile, eventsTemplate, data)

	// 4. Generate Infrastructure Implementation (Persistence)
	infraDir := "../../application/internal/infrastructure/persistence"
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		infraDir = "application/internal/infrastructure/persistence"
	}
	_ = os.MkdirAll(infraDir, 0755)

	repoImplFile := filepath.Join(infraDir, packageName+"_repo.go")
	generateFile(repoImplFile, repoImplTemplate, data)

	// 5. Generate Application Layer (Commands & Queries)
	appBaseDir := "../../application/internal/application"
	if _, err := os.Stat(appBaseDir); os.IsNotExist(err) {
		appBaseDir = "application/internal/application"
	}
	appDir := filepath.Join(appBaseDir, packageName)
	_ = os.MkdirAll(appDir, 0755)

	commandsFile := filepath.Join(appDir, "commands.go")
	generateFile(commandsFile, commandsTemplate, data)

	queriesFile := filepath.Join(appDir, "queries.go")
	generateFile(queriesFile, queriesTemplate, data)

	fmt.Println("\n✅ Domain generation complete!")
	fmt.Printf("   Domain:      %s\n", targetDir)
	fmt.Printf("   Application: %s\n", appDir)
	fmt.Printf("   Persistence: %s\n", infraDir)
}

// generateFile creates a file from template if it doesn't exist (Lock mechanism)
func generateFile(path string, tmpl string, data interface{}) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("[LOCK] Skipping %s: file already exists\n", path)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return
	}
	defer f.Close()

	t := template.Must(template.New("file").Parse(tmpl))
	if err := t.Execute(f, data); err != nil {
		fmt.Printf("Error executing template for %s: %v\n", path, err)
	} else {
		fmt.Printf("[CREATED] %s\n", path)
	}
}

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
	ID {{.EntityName}}ID ` + "`gorm:\"primaryKey\"`" + `
	// Add your fields here
}

// New{{.EntityName}} creates a new {{.EntityName}}.
func New{{.EntityName}}(id string) *{{.EntityName}} {
	e := &{{.EntityName}}{
		ID: {{.EntityName}}ID(id),
	}
	e.AddDomainEvent(New{{.EntityName}}CreatedEvent(id))
	return e
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
	// Add custom query methods here
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
`

const commandsTemplate = `package {{.PackageName}}app

import (
	"context"

	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
)

// Create{{.EntityName}}Command is the command for creating a {{.EntityName}}.
type Create{{.EntityName}}Command struct {
	ID string
	// Add command fields here
}

// Create{{.EntityName}}Handler handles Create{{.EntityName}}Command.
type Create{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

// NewCreate{{.EntityName}}Handler creates a new handler.
func NewCreate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Create{{.EntityName}}Handler {
	return &Create{{.EntityName}}Handler{repo: repo}
}

// Handle processes the command.
func (h *Create{{.EntityName}}Handler) Handle(ctx context.Context, cmd Create{{.EntityName}}Command) error {
	entity := {{.PackageName}}.New{{.EntityName}}(cmd.ID)
	return h.repo.Save(ctx, entity)
}

// Update{{.EntityName}}Command is the command for updating a {{.EntityName}}.
type Update{{.EntityName}}Command struct {
	ID string
	// Add update fields here
}

// Update{{.EntityName}}Handler handles Update{{.EntityName}}Command.
type Update{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

// NewUpdate{{.EntityName}}Handler creates a new handler.
func NewUpdate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Update{{.EntityName}}Handler {
	return &Update{{.EntityName}}Handler{repo: repo}
}

// Handle processes the command.
func (h *Update{{.EntityName}}Handler) Handle(ctx context.Context, cmd Update{{.EntityName}}Command) error {
	entity, err := h.repo.Find(ctx, {{.PackageName}}.{{.EntityName}}ID(cmd.ID))
	if err != nil {
		return err
	}
	// Update entity fields here
	entity.AddDomainEvent({{.PackageName}}.New{{.EntityName}}UpdatedEvent(cmd.ID))
	return h.repo.Save(ctx, entity)
}

// Delete{{.EntityName}}Command is the command for deleting a {{.EntityName}}.
type Delete{{.EntityName}}Command struct {
	ID string
}

// Delete{{.EntityName}}Handler handles Delete{{.EntityName}}Command.
type Delete{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

// NewDelete{{.EntityName}}Handler creates a new handler.
func NewDelete{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Delete{{.EntityName}}Handler {
	return &Delete{{.EntityName}}Handler{repo: repo}
}

// Handle processes the command.
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

// NewGet{{.EntityName}}Handler creates a new handler.
func NewGet{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Get{{.EntityName}}Handler {
	return &Get{{.EntityName}}Handler{repo: repo}
}

// Handle processes the query.
func (h *Get{{.EntityName}}Handler) Handle(ctx context.Context, query Get{{.EntityName}}Query) (*{{.PackageName}}.{{.EntityName}}, error) {
	return h.repo.Find(ctx, {{.PackageName}}.{{.EntityName}}ID(query.ID))
}

// List{{.EntityName}}sQuery is the query for listing all {{.EntityName}}s.
type List{{.EntityName}}sQuery struct{}

// List{{.EntityName}}sHandler handles List{{.EntityName}}sQuery.
type List{{.EntityName}}sHandler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

// NewList{{.EntityName}}sHandler creates a new handler.
func NewList{{.EntityName}}sHandler(repo {{.PackageName}}.{{.EntityName}}Repository) *List{{.EntityName}}sHandler {
	return &List{{.EntityName}}sHandler{repo: repo}
}

// Handle processes the query.
func (h *List{{.EntityName}}sHandler) Handle(ctx context.Context, query List{{.EntityName}}sQuery) ([]*{{.PackageName}}.{{.EntityName}}, error) {
	return h.repo.FindAll(ctx)
}
`
