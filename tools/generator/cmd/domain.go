package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var fieldsFlag string
var forceFlag bool

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
  - Database migration support

Examples:
  soliton-gen domain User
  soliton-gen domain User --fields "username,email,password,status:enum(active|inactive)"
  soliton-gen domain User --fields "..." --force  # Overwrite existing files
  soliton-gen domain Product --fields "name,price:int,stock:int,status:enum(draft|active)"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if forceFlag {
			fmt.Printf("🚀 Generating domain: %s (force mode)\n\n", name)
		} else {
			fmt.Printf("🚀 Generating domain: %s\n\n", name)
		}

		generateDomain(name, fieldsFlag)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.Flags().StringVarP(&fieldsFlag, "fields", "f", "", "Comma-separated list of fields (e.g., 'name,email,status:enum(active|inactive)')")
	domainCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
}

// Field represents a parsed field definition
type Field struct {
	Name       string   // Field name (e.g., "Username")
	SnakeName  string   // Snake case name (e.g., "username")
	CamelName  string   // Camel case name (e.g., "username")
	GoType     string   // Go type in domain package (e.g., "UserRole")
	AppGoType  string   // Go type in app layer with package prefix (e.g., "user.UserRole")
	GormTag    string   // GORM tag
	JsonTag    string   // JSON tag
	IsEnum     bool     // Is this an enum type
	EnumValues []string // Enum values if IsEnum is true
	EnumType   string   // Enum type name (e.g., "UserStatus")
}

// parseFields parses the --fields flag value into Field structs
func parseFields(fieldsStr string, entityName string) []Field {
	if fieldsStr == "" {
		// Default field if none specified
		return []Field{
			{
				Name:      "Name",
				SnakeName: "name",
				CamelName: "name",
				GoType:    "string",
				GormTag:   "`gorm:\"size:255\"`",
				JsonTag:   "`json:\"name\"`",
			},
		}
	}

	var fields []Field
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		field := parseFieldDefinition(part, entityName)
		fields = append(fields, field)
	}

	return fields
}

// parseFieldDefinition parses a single field definition
func parseFieldDefinition(def string, entityName string) Field {
	// Check for type specification (field:type or field:enum(...))
	colonIdx := strings.Index(def, ":")

	var fieldName, fieldType string
	if colonIdx == -1 {
		fieldName = def
		fieldType = "string"
	} else {
		fieldName = def[:colonIdx]
		fieldType = def[colonIdx+1:]
	}

	// Normalize field name
	snakeName := toSnakeCase(fieldName)
	pascalName := toPascalCase(fieldName)
	camelName := toCamelCase(fieldName)

	field := Field{
		Name:      pascalName,
		SnakeName: snakeName,
		CamelName: camelName,
	}

	// Parse type
	if strings.HasPrefix(fieldType, "enum(") && strings.HasSuffix(fieldType, ")") {
		// Enum type: enum(value1,value2,...)
		enumContent := fieldType[5 : len(fieldType)-1]
		field.IsEnum = true
		field.EnumValues = strings.Split(enumContent, "|")
		field.EnumType = entityName + pascalName
		field.GoType = field.EnumType
		field.AppGoType = strings.ToLower(entityName) + "." + field.EnumType // e.g., "user.UserRole"
		field.GormTag = fmt.Sprintf("`gorm:\"size:50;default:'%s'\"`", field.EnumValues[0])
		field.JsonTag = fmt.Sprintf("`json:\"%s\"`", snakeName)
	} else {
		// Regular type
		field.GoType, field.GormTag = mapFieldType(fieldType, snakeName)
		field.AppGoType = field.GoType // Same for non-enum types
		field.JsonTag = fmt.Sprintf("`json:\"%s\"`", snakeName)
	}

	return field
}

// mapFieldType maps type shorthand to Go type and GORM tag
func mapFieldType(typeName string, snakeName string) (goType, gormTag string) {
	switch strings.ToLower(typeName) {
	case "string", "str", "":
		return "string", "`gorm:\"size:255\"`"
	case "text":
		return "string", "`gorm:\"type:text\"`"
	case "int", "integer":
		return "int", "`gorm:\"not null;default:0\"`"
	case "int64", "long":
		return "int64", "`gorm:\"not null;default:0\"`"
	case "float", "float64", "double":
		return "float64", "`gorm:\"default:0\"`"
	case "bool", "boolean":
		return "bool", "`gorm:\"default:false\"`"
	case "time", "datetime", "timestamp":
		return "time.Time", "`gorm:\"autoCreateTime\"`"
	case "time?", "datetime?":
		return "*time.Time", ""
	case "uuid", "id":
		return "string", "`gorm:\"size:36;index\"`"
	default:
		return "string", "`gorm:\"size:255\"`"
	}
}

// toSnakeCase converts PascalCase or camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// toPascalCase converts snake_case to PascalCase
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(string(p[0])) + strings.ToLower(p[1:])
		}
	}
	return strings.Join(parts, "")
}

// toCamelCase converts snake_case to camelCase
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(string(pascal[0])) + pascal[1:]
}

// TemplateData holds all data for template generation
type TemplateData struct {
	PackageName string
	EntityName  string
	Fields      []Field
	HasTime     bool
	HasEnums    bool
}

func generateDomain(name string, fieldsStr string) {
	// Normalize name: first letter uppercase
	entityName := toPascalCase(name)
	packageName := strings.ToLower(name)

	// Parse fields
	fields := parseFields(fieldsStr, entityName)

	// Check for special types
	hasTime := false
	hasEnums := false
	for _, f := range fields {
		if strings.Contains(f.GoType, "time.Time") {
			hasTime = true
		}
		if f.IsEnum {
			hasEnums = true
		}
	}

	// Determine base paths
	baseDir := findPath("application/internal/domain")
	infraDir := findPath("application/internal/infrastructure/persistence")
	appDir := findPath("application/internal/application")
	interfacesDir := findPath("application/internal/interfaces/http")

	data := TemplateData{
		PackageName: packageName,
		EntityName:  entityName,
		Fields:      fields,
		HasTime:     hasTime,
		HasEnums:    hasEnums,
	}

	// === Domain Layer ===
	fmt.Println("📦 Domain Layer")
	domainDir := filepath.Join(baseDir, packageName)
	_ = os.MkdirAll(domainDir, 0755)

	generateFileWithData(filepath.Join(domainDir, packageName+".go"), entityTemplateV2, data)
	generateFileWithData(filepath.Join(domainDir, "repository.go"), repoTemplateV2, data)
	generateFileWithData(filepath.Join(domainDir, "events.go"), eventsTemplateV2, data)

	// === Infrastructure Layer ===
	fmt.Println("\n🔧 Infrastructure Layer")
	_ = os.MkdirAll(infraDir, 0755)
	generateFileWithData(filepath.Join(infraDir, packageName+"_repo.go"), repoImplTemplateV2, data)

	// === Application Layer ===
	fmt.Println("\n⚙️ Application Layer")
	appModuleDir := filepath.Join(appDir, packageName)
	_ = os.MkdirAll(appModuleDir, 0755)

	generateFileWithData(filepath.Join(appModuleDir, "commands.go"), commandsTemplateV2, data)
	generateFileWithData(filepath.Join(appModuleDir, "queries.go"), queriesTemplateV2, data)
	generateFileWithData(filepath.Join(appModuleDir, "dto.go"), dtoTemplateV2, data)

	// === Interfaces Layer (HTTP) ===
	fmt.Println("\n🌐 Interfaces Layer")
	_ = os.MkdirAll(interfacesDir, 0755)
	generateFileWithData(filepath.Join(interfacesDir, packageName+"_handler.go"), handlerTemplateV2, data)

	// === Fx Module ===
	fmt.Println("\n📌 Fx Module")
	generateFileWithData(filepath.Join(appModuleDir, "module.go"), fxModuleTemplateV2, data)

	// === Summary ===
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("✅ Domain generation complete!")
	fmt.Printf("   📁 Domain:      %s\n", domainDir)
	fmt.Printf("   📁 Application: %s\n", appModuleDir)
	fmt.Printf("   📁 Persistence: %s\n", infraDir)
	fmt.Printf("   📁 HTTP:        %s\n", interfacesDir)
	fmt.Println(strings.Repeat("─", 50))

	// Print fields info
	if len(fields) > 0 {
		fmt.Println("\n📋 Generated fields:")
		for _, f := range fields {
			if f.IsEnum {
				fmt.Printf("   • %s: %s (%s)\n", f.Name, f.GoType, strings.Join(f.EnumValues, ", "))
			} else {
				fmt.Printf("   • %s: %s\n", f.Name, f.GoType)
			}
		}
	}

	fmt.Println("\n📝 Next steps:")
	fmt.Println("   1. Review generated code")
	fmt.Println("   2. Import module in main.go:")
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

// generateFileWithData creates a file from template with TemplateData
func generateFileWithData(path string, tmpl string, data TemplateData) {
	exists := false
	if _, err := os.Stat(path); err == nil {
		if !forceFlag {
			fmt.Printf("   [SKIP] %s (exists, use --force to overwrite)\n", filepath.Base(path))
			return
		}
		exists = true
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("   [ERROR] %s: %v\n", filepath.Base(path), err)
		return
	}
	defer f.Close()

	// Create template with helper functions
	funcMap := template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}

	t := template.Must(template.New("file").Funcs(funcMap).Parse(tmpl))
	if err := t.Execute(f, data); err != nil {
		fmt.Printf("   [ERROR] %s: %v\n", filepath.Base(path), err)
	} else {
		if exists {
			fmt.Printf("   [OVERWRITE] %s\n", filepath.Base(path))
		} else {
			fmt.Printf("   [NEW] %s\n", filepath.Base(path))
		}
	}
}

// ============================================================================
// TEMPLATES V2 - With dynamic fields support
// ============================================================================

const entityTemplateV2 = `package {{.PackageName}}

import (
	"time"

	"github.com/soliton-go/framework/ddd"
)

// {{.EntityName}}ID is a strong typed ID.
type {{.EntityName}}ID string

func (id {{.EntityName}}ID) String() string {
	return string(id)
}

{{- range .Fields}}
{{- if .IsEnum}}

// {{.EnumType}} represents the {{.Name}} enum.
type {{.EnumType}} string

const (
{{- $enumType := .EnumType}}
{{- range $i, $v := .EnumValues}}
	{{$enumType}}{{$v | title}} {{$enumType}} = "{{$v}}"
{{- end}}
)
{{- end}}
{{- end}}

// {{.EntityName}} is the aggregate root.
type {{.EntityName}} struct {
	ddd.BaseAggregateRoot
	ID {{.EntityName}}ID ` + "`gorm:\"primaryKey\"`" + `
{{- range .Fields}}
	{{.Name}} {{.GoType}} {{.GormTag}}
{{- end}}
	CreatedAt time.Time ` + "`gorm:\"autoCreateTime\"`" + `
	UpdatedAt time.Time ` + "`gorm:\"autoUpdateTime\"`" + `
}

// TableName returns the table name for GORM.
func ({{.EntityName}}) TableName() string {
	return "{{.PackageName}}s"
}

// New{{.EntityName}} creates a new {{.EntityName}}.
func New{{.EntityName}}(id string{{range .Fields}}, {{.CamelName}} {{.GoType}}{{end}}) *{{.EntityName}} {
	e := &{{.EntityName}}{
		ID: {{.EntityName}}ID(id),
{{- range .Fields}}
		{{.Name}}: {{.CamelName}},
{{- end}}
	}
	e.AddDomainEvent(New{{.EntityName}}CreatedEvent(id))
	return e
}

// Update updates the entity fields.
func (e *{{.EntityName}}) Update({{range $i, $f := .Fields}}{{if $i}}, {{end}}{{$f.CamelName}} {{$f.GoType}}{{end}}) {
{{- range .Fields}}
	e.{{.Name}} = {{.CamelName}}
{{- end}}
	e.AddDomainEvent(New{{.EntityName}}UpdatedEvent(string(e.ID)))
}

// GetID returns the entity ID.
func (e *{{.EntityName}}) GetID() ddd.ID {
	return e.ID
}
`

const repoTemplateV2 = `package {{.PackageName}}

import (
	"github.com/soliton-go/framework/orm"
)

// {{.EntityName}}Repository is the interface for {{.EntityName}} persistence.
type {{.EntityName}}Repository interface {
	orm.Repository[*{{.EntityName}}, {{.EntityName}}ID]
	// TODO: Add custom query methods here
}
`

const eventsTemplateV2 = `package {{.PackageName}}

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

const repoImplTemplateV2 = `package persistence

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

const commandsTemplateV2 = `package {{.PackageName}}app

import (
	"context"

	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
)

// Create{{.EntityName}}Command is the command for creating a {{.EntityName}}.
type Create{{.EntityName}}Command struct {
	ID string
{{- range .Fields}}
	{{.Name}} {{.AppGoType}}
{{- end}}
}

// Create{{.EntityName}}Handler handles Create{{.EntityName}}Command.
type Create{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewCreate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Create{{.EntityName}}Handler {
	return &Create{{.EntityName}}Handler{repo: repo}
}

func (h *Create{{.EntityName}}Handler) Handle(ctx context.Context, cmd Create{{.EntityName}}Command) (*{{.PackageName}}.{{.EntityName}}, error) {
	entity := {{.PackageName}}.New{{.EntityName}}(cmd.ID{{range .Fields}}, cmd.{{.Name}}{{end}})
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// Update{{.EntityName}}Command is the command for updating a {{.EntityName}}.
type Update{{.EntityName}}Command struct {
	ID string
{{- range .Fields}}
	{{.Name}} {{.AppGoType}}
{{- end}}
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
	entity.Update({{range $i, $f := .Fields}}{{if $i}}, {{end}}cmd.{{$f.Name}}{{end}})
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

const queriesTemplateV2 = `package {{.PackageName}}app

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

const dtoTemplateV2 = `package {{.PackageName}}app

import (
	"time"

	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
)

// Create{{.EntityName}}Request is the request body for creating a {{.EntityName}}.
type Create{{.EntityName}}Request struct {
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}string{{else}}{{.AppGoType}}{{end}} {{.JsonTag}}
{{- end}}
}

// Update{{.EntityName}}Request is the request body for updating a {{.EntityName}}.
type Update{{.EntityName}}Request struct {
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}string{{else}}{{.AppGoType}}{{end}} {{.JsonTag}}
{{- end}}
}

// {{.EntityName}}Response is the response body for {{.EntityName}} data.
type {{.EntityName}}Response struct {
	ID        string    ` + "`json:\"id\"`" + `
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}string{{else}}{{.AppGoType}}{{end}} {{.JsonTag}}
{{- end}}
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// To{{.EntityName}}Response converts entity to response.
func To{{.EntityName}}Response(e *{{.PackageName}}.{{.EntityName}}) {{.EntityName}}Response {
	return {{.EntityName}}Response{
		ID:        string(e.ID),
{{- range .Fields}}
		{{.Name}}: {{if .IsEnum}}string(e.{{.Name}}){{else}}e.{{.Name}}{{end}},
{{- end}}
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
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

const handlerTemplateV2 = `package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	{{.PackageName}}app "github.com/soliton-go/application/internal/application/{{.PackageName}}"
{{- if .HasEnums}}
	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
{{- end}}
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
		ID: uuid.New().String(),
{{- range .Fields}}
		{{.Name}}: {{if .IsEnum}}{{$.PackageName}}.{{.EnumType}}(req.{{.Name}}){{else}}req.{{.Name}}{{end}},
{{- end}}
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
		ID: id,
{{- range .Fields}}
		{{.Name}}: {{if .IsEnum}}{{$.PackageName}}.{{.EnumType}}(req.{{.Name}}){{else}}req.{{.Name}}{{end}},
{{- end}}
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

const fxModuleTemplateV2 = `package {{.PackageName}}app

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
