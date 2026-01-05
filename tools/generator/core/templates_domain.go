package core

// ============================================================================
// DOMAIN TEMPLATES
// ============================================================================

const EntityTemplate = `package {{.PackageName}}

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
	{{.Name}} {{.GoType}} {{.GormTag}}{{if .Comment}} // {{.Comment}}{{end}}
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

const RepoTemplate = `package {{.PackageName}}

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

const EventsTemplate = `package {{.PackageName}}

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

const RepoImplTemplate = `package persistence

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

// Migrate{{.EntityName}} creates the table if it doesn't exist.
func Migrate{{.EntityName}}(db *gorm.DB) error {
	return db.AutoMigrate(&{{.PackageName}}.{{.EntityName}}{})
}
`

const CommandsTemplate = `package {{.PackageName}}app

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

const QueriesTemplate = `package {{.PackageName}}app

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

const DTOTemplate = `package {{.PackageName}}app

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

const HelpersHTTPTemplate = `package http

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

const HandlerTemplate = `package http

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

const FxModuleTemplate = `package {{.PackageName}}app

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
