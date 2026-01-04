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
var tableNameFlag string
var routeBaseFlag string
var wireFlag bool
var softDeleteFlag bool

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
	domainCmd.Flags().StringVar(&tableNameFlag, "table", "", "Override database table name")
	domainCmd.Flags().StringVar(&routeBaseFlag, "route", "", "Override route base path (e.g., users)")
	domainCmd.Flags().BoolVar(&wireFlag, "wire", false, "Auto-wire module into main.go (requires init template structure)")
	domainCmd.Flags().BoolVar(&softDeleteFlag, "soft-delete", false, "Enable soft delete (adds deleted_at field)")
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
	IsPointer  bool     // True if GoType is a pointer type
}

// parseFields parses the --fields flag value into Field structs
func parseFields(fieldsStr string, entityName string, packageName string) []Field {
	if fieldsStr == "" {
		// Default field if none specified
		return []Field{
			{
				Name:      "Name",
				SnakeName: "name",
				CamelName: "name",
				GoType:    "string",
				AppGoType: "string",
				GormTag:   "`gorm:\"size:255\"`",
				JsonTag:   "`json:\"name\"`",
			},
		}
	}

	var fields []Field
	parts := strings.Split(fieldsStr, ",")
	seen := make(map[string]struct{})

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		field := parseFieldDefinition(part, entityName, packageName)
		if isReservedField(field) {
			fmt.Printf("   [WARN] field %q conflicts with built-in fields, skipping\n", field.SnakeName)
			continue
		}
		if _, exists := seen[field.SnakeName]; exists {
			fmt.Printf("   [WARN] field %q is duplicated, skipping\n", field.SnakeName)
			continue
		}
		seen[field.SnakeName] = struct{}{}
		fields = append(fields, field)
	}

	return fields
}

// parseFieldDefinition parses a single field definition
func parseFieldDefinition(def string, entityName string, packageName string) Field {
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
		field.EnumValues = parseEnumValues(enumContent)
		if len(field.EnumValues) == 0 {
			fmt.Printf("   [WARN] enum field %q has no values, defaulting to \"default\"\n", fieldName)
			field.EnumValues = []string{"default"}
		}
		field.EnumType = entityName + pascalName
		field.GoType = field.EnumType
		field.AppGoType = packageName + "." + field.EnumType // e.g., "user.UserRole"
		field.GormTag = fmt.Sprintf("`gorm:\"size:50;default:'%s'\"`", field.EnumValues[0])
		field.JsonTag = fmt.Sprintf("`json:\"%s\"`", snakeName)
	} else {
		// Regular type
		field.GoType, field.GormTag = mapFieldType(fieldType, snakeName)
		field.AppGoType = field.GoType // Same for non-enum types
		field.IsPointer = strings.HasPrefix(field.GoType, "*")
		field.JsonTag = fmt.Sprintf("`json:\"%s\"`", snakeName)
	}

	return field
}

func parseEnumValues(raw string) []string {
	parts := strings.Split(raw, "|")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}

func isReservedField(field Field) bool {
	switch field.SnakeName {
	case "id", "created_at", "updated_at":
		return true
	default:
		return false
	}
}

func createBindingTag(field Field) string {
	if field.IsPointer {
		return ""
	}
	if field.IsEnum {
		return oneOfTag("required", field.EnumValues)
	}
	if isStringType(field.GoType) {
		return " binding:\"required\""
	}
	return ""
}

func updateBindingTag(field Field) string {
	if !field.IsEnum {
		return ""
	}
	return oneOfTag("omitempty", field.EnumValues)
}

func oneOfTag(prefix string, values []string) string {
	if len(values) == 0 {
		if prefix == "required" {
			return " binding:\"required\""
		}
		return ""
	}
	return fmt.Sprintf(" binding:\"%s,oneof=%s\"", prefix, strings.Join(values, " "))
}

func isStringType(goType string) bool {
	return strings.TrimPrefix(goType, "*") == "string"
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
		return "time.Time", "`gorm:\"type:timestamp\"`"
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

func pluralize(name string) string {
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") ||
		strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return lower + "es"
	}
	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		prev := lower[len(lower)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return lower[:len(lower)-1] + "ies"
		}
	}
	return lower + "s"
}

func enumConst(value string) string {
	var cleaned strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			cleaned.WriteRune(r)
		} else {
			cleaned.WriteByte('_')
		}
	}

	normalized := strings.Trim(cleaned.String(), "_")
	if normalized == "" {
		return "Value"
	}

	normalized = strings.ToLower(normalized)
	result := toPascalCase(normalized)
	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		return "Value" + result
	}
	return result
}

// TemplateData holds all data for template generation
type TemplateData struct {
	PackageName string
	EntityName  string
	Fields      []Field
	HasTime     bool
	HasEnums    bool
	ModulePath  string
	TableName   string
	RouteBase   string
	SoftDelete  bool
}

func generateDomain(name string, fieldsStr string) {
	// Normalize name: first letter uppercase
	entityName := toPascalCase(name)
	packageName := strings.ToLower(name)

	layout, err := ResolveProjectLayout()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	// Parse fields
	fields := parseFields(fieldsStr, entityName, packageName)

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

	tableName := tableNameFlag
	if tableName == "" {
		tableName = pluralize(packageName)
	}

	routeBase := routeBaseFlag
	if routeBase == "" {
		routeBase = pluralize(packageName)
	}

	data := TemplateData{
		PackageName: packageName,
		EntityName:  entityName,
		Fields:      fields,
		HasTime:     hasTime,
		HasEnums:    hasEnums,
		ModulePath:  layout.ModulePath,
		TableName:   tableName,
		RouteBase:   routeBase,
		SoftDelete:  softDeleteFlag,
	}

	// === Domain Layer ===
	fmt.Println("📦 Domain Layer")
	domainDir := filepath.Join(layout.DomainDir, packageName)
	_ = os.MkdirAll(domainDir, 0755)

	generateFileWithData(filepath.Join(domainDir, packageName+".go"), entityTemplateV2, data)
	generateFileWithData(filepath.Join(domainDir, "repository.go"), repoTemplateV2, data)
	generateFileWithData(filepath.Join(domainDir, "events.go"), eventsTemplateV2, data)

	// === Infrastructure Layer ===
	fmt.Println("\n🔧 Infrastructure Layer")
	_ = os.MkdirAll(layout.InfraDir, 0755)
	generateFileWithData(filepath.Join(layout.InfraDir, packageName+"_repo.go"), repoImplTemplateV2, data)

	// === Application Layer ===
	fmt.Println("\n⚙️ Application Layer")
	appModuleDir := filepath.Join(layout.AppDir, packageName)
	_ = os.MkdirAll(appModuleDir, 0755)

	generateFileWithData(filepath.Join(appModuleDir, "commands.go"), commandsTemplateV2, data)
	generateFileWithData(filepath.Join(appModuleDir, "queries.go"), queriesTemplateV2, data)
	generateFileWithData(filepath.Join(appModuleDir, "dto.go"), dtoTemplateV2, data)

	// === Interfaces Layer (HTTP) ===
	fmt.Println("\n🌐 Interfaces Layer")
	_ = os.MkdirAll(layout.InterfacesDir, 0755)

	// Generate helpers.go if it doesn't exist
	helpersPath := filepath.Join(layout.InterfacesDir, "helpers.go")
	if _, err := os.Stat(helpersPath); os.IsNotExist(err) {
		generateFileWithData(helpersPath, helpersTemplate, data)
	}

	generateFileWithData(filepath.Join(layout.InterfacesDir, packageName+"_handler.go"), handlerTemplateV2, data)

	// === Fx Module ===
	fmt.Println("\n📌 Fx Module")
	generateFileWithData(filepath.Join(appModuleDir, "module.go"), fxModuleTemplateV2, data)

	// === Summary ===
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("✅ Domain generation complete!")
	fmt.Printf("   📁 Domain:      %s\n", domainDir)
	fmt.Printf("   📁 Application: %s\n", appModuleDir)
	fmt.Printf("   📁 Persistence: %s\n", layout.InfraDir)
	fmt.Printf("   📁 HTTP:        %s\n", layout.InterfacesDir)
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
	if wireFlag {
		// Try to wire main.go
		mainGoPath := filepath.Join(filepath.Dir(layout.InternalDir), "cmd", "main.go")
		if wireMainGo(mainGoPath, entityName, packageName, layout.ModulePath) {
			fmt.Println("   ✅ main.go updated automatically!")
			fmt.Println("   Run: GOWORK=off go run ./cmd/main.go")
		} else {
			fmt.Println("   ⚠️  Could not auto-wire main.go (structure not recognized)")
			fmt.Println("   1. Review generated code")
			fmt.Println("   2. Import module in main.go manually:")
			fmt.Printf("      %sapp.Module,\n", packageName)
		}
	} else {
		fmt.Println("   1. Review generated code")
		fmt.Println("   2. Import module in main.go:")
		fmt.Printf("      %sapp.Module,\n", packageName)
		fmt.Println("   Tip: Use --wire flag to auto-inject into main.go")
	}
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
		"title":            strings.Title,
		"lower":            strings.ToLower,
		"upper":            strings.ToUpper,
		"enumConst":        enumConst,
		"createBindingTag": createBindingTag,
		"updateBindingTag": updateBindingTag,
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
{{- if .SoftDelete}}
	"gorm.io/gorm"
{{- end}}
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
	{{$enumType}}{{$v | enumConst}} {{$enumType}} = "{{$v}}"
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
{{- if .SoftDelete}}
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\"`" + `
{{- end}}
}

// TableName returns the table name for GORM.
func ({{.EntityName}}) TableName() string {
	return "{{.TableName}}"
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
func (e *{{.EntityName}}) Update({{range $i, $f := .Fields}}{{if $i}}, {{end}}{{$f.CamelName}} {{if $f.IsEnum}}*{{$f.GoType}}{{else if $f.IsPointer}}{{$f.GoType}}{{else}}*{{$f.GoType}}{{end}}{{end}}) {
{{- range .Fields}}
{{- if .IsPointer}}
	if {{.CamelName}} != nil {
		e.{{.Name}} = {{.CamelName}}
	}
{{- else}}
	if {{.CamelName}} != nil {
		e.{{.Name}} = *{{.CamelName}}
	}
{{- end}}
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
	"context"

	"github.com/soliton-go/framework/orm"
)

// {{.EntityName}}Repository is the interface for {{.EntityName}} persistence.
type {{.EntityName}}Repository interface {
	orm.Repository[*{{.EntityName}}, {{.EntityName}}ID]
	// FindPaginated returns a page of entities with total count.
	FindPaginated(ctx context.Context, page, pageSize int) ([]*{{.EntityName}}, int64, error)
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
	"context"

	"{{.ModulePath}}/internal/domain/{{.PackageName}}"
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

// FindPaginated returns a page of entities with total count.
func (r *{{.EntityName}}RepoImpl) FindPaginated(ctx context.Context, page, pageSize int) ([]*{{.PackageName}}.{{.EntityName}}, int64, error) {
	var entities []*{{.PackageName}}.{{.EntityName}}
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&{{.PackageName}}.{{.EntityName}}{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get page
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Migrate creates the table if it doesn't exist.
func Migrate{{.EntityName}}(db *gorm.DB) error {
	return db.AutoMigrate(&{{.PackageName}}.{{.EntityName}}{})
}
`

const commandsTemplateV2 = `package {{.PackageName}}app

import (
	"context"
{{- if .HasTime}}
	"time"
{{- end}}

	"{{.ModulePath}}/internal/domain/{{.PackageName}}"
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
	// Optional: Add event bus for domain event publishing
	// eventBus event.EventBus
}

func NewCreate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Create{{.EntityName}}Handler {
	return &Create{{.EntityName}}Handler{repo: repo}
}

func (h *Create{{.EntityName}}Handler) Handle(ctx context.Context, cmd Create{{.EntityName}}Command) (*{{.PackageName}}.{{.EntityName}}, error) {
	entity := {{.PackageName}}.New{{.EntityName}}(cmd.ID{{range .Fields}}, cmd.{{.Name}}{{end}})
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// Optional: Publish domain events
	// Uncomment to enable event publishing:
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// Update{{.EntityName}}Command is the command for updating a {{.EntityName}}.
type Update{{.EntityName}}Command struct {
	ID string
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}*{{.AppGoType}}{{else if .IsPointer}}{{.AppGoType}}{{else}}*{{.AppGoType}}{{end}}
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

	"{{.ModulePath}}/internal/domain/{{.PackageName}}"
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

// List{{.EntityName}}sQuery is the query for listing {{.EntityName}}s with pagination.
type List{{.EntityName}}sQuery struct {
	Page     int // Page number (1-based)
	PageSize int // Items per page (default: 20, max: 100)
}

// List{{.EntityName}}sResult is the paginated result for List{{.EntityName}}sQuery.
type List{{.EntityName}}sResult struct {
	Items      []*{{.PackageName}}.{{.EntityName}}
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// List{{.EntityName}}sHandler handles List{{.EntityName}}sQuery.
type List{{.EntityName}}sHandler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewList{{.EntityName}}sHandler(repo {{.PackageName}}.{{.EntityName}}Repository) *List{{.EntityName}}sHandler {
	return &List{{.EntityName}}sHandler{repo: repo}
}

func (h *List{{.EntityName}}sHandler) Handle(ctx context.Context, query List{{.EntityName}}sQuery) (*List{{.EntityName}}sResult, error) {
	// Normalize pagination parameters
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get total count and items
	items, total, err := h.repo.FindPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &List{{.EntityName}}sResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
`

const dtoTemplateV2 = `package {{.PackageName}}app

import (
	"time"

	"{{.ModulePath}}/internal/domain/{{.PackageName}}"
)

// Create{{.EntityName}}Request is the request body for creating a {{.EntityName}}.
type Create{{.EntityName}}Request struct {
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}string{{else}}{{.AppGoType}}{{end}} ` + "`json:\"{{.SnakeName}}\"{{createBindingTag .}}`" + `
{{- end}}
}

// Update{{.EntityName}}Request is the request body for updating a {{.EntityName}}.
type Update{{.EntityName}}Request struct {
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}*string{{else if .IsPointer}}{{.AppGoType}}{{else}}*{{.AppGoType}}{{end}} ` + "`json:\"{{.SnakeName}},omitempty\"{{updateBindingTag .}}`" + `
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

const helpersTemplate = `package http

// EnumPtr is a helper function to convert *string to *T for enum types.
// This is useful for handling optional enum fields in update requests.
func EnumPtr[T any](v *string, parse func(string) T) *T {
	if v == nil {
		return nil
	}
	parsed := parse(*v)
	return &parsed
}
`

const handlerTemplateV2 = `package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	{{.PackageName}}app "{{.ModulePath}}/internal/application/{{.PackageName}}"
{{- if .HasEnums}}
	"{{.ModulePath}}/internal/domain/{{.PackageName}}"
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
	api := r.Group("/api/{{.RouteBase}}")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id", h.Update)
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

// List handles GET /api/{{.PackageName}}s?page=1&page_size=20
func (h *{{.EntityName}}Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.listHandler.Handle(c.Request.Context(), {{.PackageName}}app.List{{.EntityName}}sQuery{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{
		"items":       {{.PackageName}}app.To{{.EntityName}}ResponseList(result.Items),
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
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
	{{- if .IsEnum}}
		{{.Name}}: EnumPtr(req.{{.Name}}, func(v string) {{.AppGoType}} { return {{$.PackageName}}.{{.EnumType}}(v) }),
	{{- else}}
		{{.Name}}: req.{{.Name}},
	{{- end}}
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

	"{{.ModulePath}}/internal/domain/{{.PackageName}}"
	"{{.ModulePath}}/internal/infrastructure/persistence"
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

	// Optional: Register with CQRS bus
	// Uncomment to enable CQRS pattern:
	// fx.Invoke(func(cmdBus *cqrs.InMemoryCommandBus, queryBus *cqrs.InMemoryQueryBus,
	//     createHandler *Create{{.EntityName}}Handler,
	//     updateHandler *Update{{.EntityName}}Handler,
	//     deleteHandler *Delete{{.EntityName}}Handler,
	//     getHandler *Get{{.EntityName}}Handler,
	//     listHandler *List{{.EntityName}}sHandler) {
	//     cmdBus.Register(Create{{.EntityName}}Command{}, createHandler.Handle)
	//     cmdBus.Register(Update{{.EntityName}}Command{}, updateHandler.Handle)
	//     cmdBus.Register(Delete{{.EntityName}}Command{}, deleteHandler.Handle)
	//     queryBus.Register(Get{{.EntityName}}Query{}, getHandler.Handle)
	//     queryBus.Register(List{{.EntityName}}sQuery{}, listHandler.Handle)
	// }),
)

// RegisterMigration registers the {{.EntityName}} table migration.
func RegisterMigration(db *gorm.DB) error {
	return persistence.Migrate{{.EntityName}}(db)
}
`

// wireMainGo attempts to inject module into main.go using marker comments.
// Supports multiple modules by appending to marker lines.
// Returns true if successful, false if main.go structure is not recognized.
func wireMainGo(mainGoPath, entityName, packageName, modulePath string) bool {
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return false
	}

	original := string(content)

	// Check for new template markers
	if strings.Contains(original, "// soliton-gen:imports") {
		return wireMainGoNew(mainGoPath, entityName, packageName, modulePath, original)
	}

	// Try legacy mode for old templates
	if strings.Contains(original, "// Uncomment these imports after generating domains:") {
		return wireMainGoLegacy(mainGoPath, entityName, packageName, modulePath, original)
	}

	return false
}

// wireMainGoNew handles new template format with marker comments
func wireMainGoNew(mainGoPath, entityName, packageName, modulePath, original string) bool {
	result := original
	modified := false

	// 0. Replace blank gorm import with normal import (needed for fx.Invoke)
	if strings.Contains(result, "_ \"gorm.io/gorm\"") && !strings.Contains(result, "\t\"gorm.io/gorm\"") {
		result = strings.Replace(result, "_ \"gorm.io/gorm\"", "\"gorm.io/gorm\"", 1)
		modified = true
	}

	// 1. Add app import after // soliton-gen:imports
	appImport := fmt.Sprintf("%sapp \"%s/internal/application/%s\"", packageName, modulePath, packageName)
	if !strings.Contains(result, appImport) {
		result = strings.Replace(result,
			"\t// soliton-gen:imports",
			"\t"+appImport+"\n\t// soliton-gen:imports",
			1)
		modified = true
	}

	// 2. Add interfaceshttp import if not present
	httpImport := fmt.Sprintf("interfaceshttp \"%s/internal/interfaces/http\"", modulePath)
	if !strings.Contains(result, httpImport) {
		result = strings.Replace(result,
			"\t// soliton-gen:imports",
			fmt.Sprintf("\tinterfaceshttp \"%s/internal/interfaces/http\"\n\t// soliton-gen:imports", modulePath),
			1)
		modified = true
	}

	// 3. Add module after // soliton-gen:modules
	moduleCode := fmt.Sprintf("%sapp.Module,", packageName)
	if !strings.Contains(result, moduleCode) {
		result = strings.Replace(result,
			"\t\t// soliton-gen:modules",
			"\t\t"+moduleCode+"\n\t\t// soliton-gen:modules",
			1)
		modified = true
	}

	// 4. Add handler after // soliton-gen:handlers
	handlerCode := fmt.Sprintf("fx.Provide(interfaceshttp.New%sHandler),", entityName)
	if !strings.Contains(result, handlerCode) {
		result = strings.Replace(result,
			"\t\t// soliton-gen:handlers",
			"\t\t"+handlerCode+"\n\t\t// soliton-gen:handlers",
			1)
		modified = true
	}

	// 5. Add route registration after // soliton-gen:routes
	routeCheck := fmt.Sprintf("h *interfaceshttp.%sHandler", entityName)
	if !strings.Contains(result, routeCheck) {
		routeCode := fmt.Sprintf("fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.%sHandler) error {\n\t\t\tif err := %sapp.RegisterMigration(db); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\th.RegisterRoutes(r)\n\t\t\treturn nil\n\t\t}),", entityName, packageName)
		result = strings.Replace(result,
			"\t\t// soliton-gen:routes",
			"\t\t"+routeCode+"\n\t\t// soliton-gen:routes",
			1)
		modified = true
	}

	if !modified {
		return true // Already wired
	}

	return os.WriteFile(mainGoPath, []byte(result), 0644) == nil
}

// wireMainGoLegacy handles old template format with commented placeholders
func wireMainGoLegacy(mainGoPath, entityName, packageName, modulePath, original string) bool {
	result := original
	modified := false

	// Uncomment gorm, app, http imports, module, handler, invoke
	replacements := []struct{ old, new string }{
		{"\t// \"gorm.io/gorm\"", "\t\"gorm.io/gorm\""},
		{fmt.Sprintf("\t// %sapp \"%s/internal/application/%s\"", packageName, modulePath, packageName),
			fmt.Sprintf("\t%sapp \"%s/internal/application/%s\"", packageName, modulePath, packageName)},
		{fmt.Sprintf("\t// interfaceshttp \"%s/internal/interfaces/http\"", modulePath),
			fmt.Sprintf("\tinterfaceshttp \"%s/internal/interfaces/http\"", modulePath)},
		{fmt.Sprintf("\t\t// %sapp.Module,", packageName), fmt.Sprintf("\t\t%sapp.Module,", packageName)},
		{fmt.Sprintf("\t\t// fx.Provide(interfaceshttp.New%sHandler),", entityName),
			fmt.Sprintf("\t\tfx.Provide(interfaceshttp.New%sHandler),", entityName)},
	}

	for _, r := range replacements {
		if strings.Contains(result, r.old) {
			result = strings.Replace(result, r.old, r.new, 1)
			modified = true
		}
	}

	// Uncomment invoke block (support both legacy and updated commented blocks)
	oldInvoke := fmt.Sprintf("\t\t// fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.%sHandler) {\n\t\t// \t%sapp.RegisterMigration(db)\n\t\t// \th.RegisterRoutes(r)\n\t\t// }),", entityName, packageName)
	legacyInvoke := fmt.Sprintf("\t\t// fx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.%sHandler) error {\n\t\t// \tif err := %sapp.RegisterMigration(db); err != nil {\n\t\t// \t\treturn err\n\t\t// \t}\n\t\t// \th.RegisterRoutes(r)\n\t\t// \treturn nil\n\t\t// }),", entityName, packageName)
	newInvoke := fmt.Sprintf("\t\tfx.Invoke(func(db *gorm.DB, r *gin.Engine, h *interfaceshttp.%sHandler) error {\n\t\t\tif err := %sapp.RegisterMigration(db); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\th.RegisterRoutes(r)\n\t\t\treturn nil\n\t\t}),", entityName, packageName)
	if strings.Contains(result, oldInvoke) {
		result = strings.Replace(result, oldInvoke, newInvoke, 1)
		modified = true
	} else if strings.Contains(result, legacyInvoke) {
		result = strings.Replace(result, legacyInvoke, newInvoke, 1)
		modified = true
	}

	if !modified {
		return false
	}

	return os.WriteFile(mainGoPath, []byte(result), 0644) == nil
}
