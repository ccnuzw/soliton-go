package core

// ============================================================================
// DOMAIN TEMPLATES / 领域模板
// ============================================================================

const EntityTemplate = `package {{.PackageName}}

import (
	"time"

	"github.com/soliton-go/framework/ddd"
{{- if .SoftDelete}}
	"gorm.io/gorm"
{{- end}}
)

// {{.EntityName}}ID 是强类型的实体标识符。
type {{.EntityName}}ID string

func (id {{.EntityName}}ID) String() string {
	return string(id)
}

{{- range .Fields}}
{{- if .IsEnum}}

// {{.EnumType}} 表示 {{.Name}} 字段的枚举类型。
type {{.EnumType}} string

const (
{{- $enumType := .EnumType}}
{{- range $i, $v := .EnumValues}}
	{{$enumType}}{{$v | enumConst}} {{$enumType}} = "{{$v}}"
{{- end}}
)
{{- end}}
{{- end}}

// {{.EntityName}} 是聚合根实体。
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

// TableName 返回 GORM 映射的数据库表名。
func ({{.EntityName}}) TableName() string {
	return "{{.TableName}}"
}

// New{{.EntityName}} 创建一个新的 {{.EntityName}} 实体。
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

// Update 更新实体字段。
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

// GetID 返回实体 ID。
func (e *{{.EntityName}}) GetID() ddd.ID {
	return e.ID
}
`

const DomainServiceTemplate = `package {{.PackageName}}

import "context"

// {{.EntityName}}DomainService 提供领域内的复杂业务逻辑封装。
type {{.EntityName}}DomainService struct {
	repo {{.EntityName}}Repository
}

// New{{.EntityName}}DomainService 创建 {{.EntityName}}DomainService 实例。
func New{{.EntityName}}DomainService(repo {{.EntityName}}Repository) *{{.EntityName}}DomainService {
	return &{{.EntityName}}DomainService{repo: repo}
}

// TODO: 在此添加领域服务方法，例如复杂校验、跨实体规则等。
func (s *{{.EntityName}}DomainService) Validate(ctx context.Context) error {
	_ = ctx
	return nil
}
`

const RepoTemplate = `package {{.PackageName}}

import (
	"context"

	"github.com/soliton-go/framework/orm"
)

// {{.EntityName}}Repository 定义 {{.EntityName}} 的持久化接口。
type {{.EntityName}}Repository interface {
	orm.Repository[*{{.EntityName}}, {{.EntityName}}ID]
	// FindPaginated 返回分页数据和总数。
	FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*{{.EntityName}}, int64, error)
}
`

const EventsTemplate = `package {{.PackageName}}

import (
	"time"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

// {{.EntityName}}CreatedEvent 在创建 {{.EntityName}} 时发布。
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

// {{.EntityName}}UpdatedEvent 在更新 {{.EntityName}} 时发布。
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

// {{.EntityName}}DeletedEvent 在删除 {{.EntityName}} 时发布。
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

// init 将事件注册到全局注册表。
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
	"fmt"

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

// FindPaginated 返回分页数据和总数。
func (r *{{.EntityName}}RepoImpl) FindPaginated(ctx context.Context, page, pageSize int, sortBy, sortOrder string) ([]*{{.PackageName}}.{{.EntityName}}, int64, error) {
	var entities []*{{.PackageName}}.{{.EntityName}}
	var total int64

	// 查询总数
	baseQuery := r.db.WithContext(ctx).Model(&{{.PackageName}}.{{.EntityName}}{})
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := r.db.WithContext(ctx).Offset(offset).Limit(pageSize)
	if sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Migrate{{.EntityName}} 创建数据库表（如不存在）。
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

// Create{{.EntityName}}Command 是创建 {{.EntityName}} 的命令。
type Create{{.EntityName}}Command struct {
	ID string
{{- range .Fields}}
	{{.Name}} {{.AppGoType}}
{{- end}}
}

// Create{{.EntityName}}Handler 处理 Create{{.EntityName}}Command。
type Create{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
	service *{{.PackageName}}.{{.EntityName}}DomainService
	// 可选：添加事件总线用于发布领域事件
	// eventBus event.EventBus
}

func NewCreate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository, service *{{.PackageName}}.{{.EntityName}}DomainService) *Create{{.EntityName}}Handler {
	return &Create{{.EntityName}}Handler{repo: repo, service: service}
}

func (h *Create{{.EntityName}}Handler) Handle(ctx context.Context, cmd Create{{.EntityName}}Command) (*{{.PackageName}}.{{.EntityName}}, error) {
	entity := {{.PackageName}}.New{{.EntityName}}(cmd.ID{{range .Fields}}, cmd.{{.Name}}{{end}})
	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	// 可选：发布领域事件
	// 取消注释以启用事件发布：
	// events := entity.PullDomainEvents()
	// if len(events) > 0 {
	//     if err := h.eventBus.Publish(ctx, events...); err != nil {
	//         return nil, err
	//     }
	// }

	return entity, nil
}

// Update{{.EntityName}}Command 是更新 {{.EntityName}} 的命令。
type Update{{.EntityName}}Command struct {
	ID string
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}*{{.AppGoType}}{{else if .IsPointer}}{{.AppGoType}}{{else}}*{{.AppGoType}}{{end}}
{{- end}}
}

// Update{{.EntityName}}Handler 处理 Update{{.EntityName}}Command。
type Update{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
	service *{{.PackageName}}.{{.EntityName}}DomainService
}

func NewUpdate{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository, service *{{.PackageName}}.{{.EntityName}}DomainService) *Update{{.EntityName}}Handler {
	return &Update{{.EntityName}}Handler{repo: repo, service: service}
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

// Delete{{.EntityName}}Command 是删除 {{.EntityName}} 的命令。
type Delete{{.EntityName}}Command struct {
	ID string
}

// Delete{{.EntityName}}Handler 处理 Delete{{.EntityName}}Command。
type Delete{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
	service *{{.PackageName}}.{{.EntityName}}DomainService
}

func NewDelete{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository, service *{{.PackageName}}.{{.EntityName}}DomainService) *Delete{{.EntityName}}Handler {
	return &Delete{{.EntityName}}Handler{repo: repo, service: service}
}

func (h *Delete{{.EntityName}}Handler) Handle(ctx context.Context, cmd Delete{{.EntityName}}Command) error {
	return h.repo.Delete(ctx, {{.PackageName}}.{{.EntityName}}ID(cmd.ID))
}
`

const QueriesTemplate = `package {{.PackageName}}app

import (
	"context"
	"strings"

	"{{.ModulePath}}/internal/domain/{{.PackageName}}"
)

// Get{{.EntityName}}Query 是获取单个 {{.EntityName}} 的查询。
type Get{{.EntityName}}Query struct {
	ID string
}

// Get{{.EntityName}}Handler 处理 Get{{.EntityName}}Query。
type Get{{.EntityName}}Handler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewGet{{.EntityName}}Handler(repo {{.PackageName}}.{{.EntityName}}Repository) *Get{{.EntityName}}Handler {
	return &Get{{.EntityName}}Handler{repo: repo}
}

func (h *Get{{.EntityName}}Handler) Handle(ctx context.Context, query Get{{.EntityName}}Query) (*{{.PackageName}}.{{.EntityName}}, error) {
	return h.repo.Find(ctx, {{.PackageName}}.{{.EntityName}}ID(query.ID))
}

// List{{.EntityName}}sQuery 是分页列表查询。
type List{{.EntityName}}sQuery struct {
	Page     int // 页码（从 1 开始）
	PageSize int // 每页数量（默认: 20, 最大: 100）
	SortBy   string // 排序字段（默认: id）
	SortOrder string // 排序方式（asc/desc）
}

// List{{.EntityName}}sResult 是分页查询结果。
type List{{.EntityName}}sResult struct {
	Items      []*{{.PackageName}}.{{.EntityName}}
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// List{{.EntityName}}sHandler 处理 List{{.EntityName}}sQuery。
type List{{.EntityName}}sHandler struct {
	repo {{.PackageName}}.{{.EntityName}}Repository
}

func NewList{{.EntityName}}sHandler(repo {{.PackageName}}.{{.EntityName}}Repository) *List{{.EntityName}}sHandler {
	return &List{{.EntityName}}sHandler{repo: repo}
}

func (h *List{{.EntityName}}sHandler) Handle(ctx context.Context, query List{{.EntityName}}sQuery) (*List{{.EntityName}}sResult, error) {
	// 规范化分页参数
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

	// 排序字段白名单
	sortBy := strings.ToLower(strings.TrimSpace(query.SortBy))
	sortOrder := strings.ToLower(strings.TrimSpace(query.SortOrder))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	allowedSorts := map[string]struct{}{
		"id":         {},
		"created_at": {},
		"updated_at": {},
{{- range .Fields}}
		"{{.SnakeName}}": {},
{{- end}}
	}
	if _, ok := allowedSorts[sortBy]; !ok {
		sortBy = "id"
	}

	// 获取总数和分页数据
	items, total, err := h.repo.FindPaginated(ctx, page, pageSize, sortBy, sortOrder)
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

// Create{{.EntityName}}Request 是创建 {{.EntityName}} 的请求体。
type Create{{.EntityName}}Request struct {
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}string{{else}}{{.AppGoType}}{{end}} ` + "`json:\"{{.SnakeName}}\"{{createBindingTag .}}`" + `
{{- end}}
}

// Update{{.EntityName}}Request 是更新 {{.EntityName}} 的请求体。
type Update{{.EntityName}}Request struct {
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}*string{{else if .IsPointer}}{{.AppGoType}}{{else}}*{{.AppGoType}}{{end}} ` + "`json:\"{{.SnakeName}},omitempty\"{{updateBindingTag .}}`" + `
{{- end}}
}

// {{.EntityName}}Response 是 {{.EntityName}} 的响应体。
type {{.EntityName}}Response struct {
	ID        string    ` + "`json:\"id\"`" + `
{{- range .Fields}}
	{{.Name}} {{if .IsEnum}}string{{else}}{{.AppGoType}}{{end}} {{.JsonTag}}
{{- end}}
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// To{{.EntityName}}Response 将实体转换为响应体。
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

// To{{.EntityName}}ResponseList 将实体列表转换为响应体列表。
func To{{.EntityName}}ResponseList(entities []*{{.PackageName}}.{{.EntityName}}) []{{.EntityName}}Response {
	result := make([]{{.EntityName}}Response, len(entities))
	for i, e := range entities {
		result[i] = To{{.EntityName}}Response(e)
	}
	return result
}
`

const HelpersHTTPTemplate = `package http

// EnumPtr 是一个辅助函数，用于将 *string 转换为枚举类型的 *T。
// 适用于处理更新请求中的可选枚举字段。
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

// {{.EntityName}}Handler 处理 {{.EntityName}} 相关的 HTTP 请求。
type {{.EntityName}}Handler struct {
	createHandler *{{.PackageName}}app.Create{{.EntityName}}Handler
	updateHandler *{{.PackageName}}app.Update{{.EntityName}}Handler
	deleteHandler *{{.PackageName}}app.Delete{{.EntityName}}Handler
	getHandler    *{{.PackageName}}app.Get{{.EntityName}}Handler
	listHandler   *{{.PackageName}}app.List{{.EntityName}}sHandler
}

// New{{.EntityName}}Handler 创建 {{.EntityName}}Handler 实例。
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

// RegisterRoutes 注册 {{.EntityName}} 相关路由。
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

// Create 处理 POST /api/{{.PackageName}}s
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

// Get 处理 GET /api/{{.PackageName}}s/:id
func (h *{{.EntityName}}Handler) Get(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.getHandler.Handle(c.Request.Context(), {{.PackageName}}app.Get{{.EntityName}}Query{ID: id})
	if err != nil {
		NotFound(c, "{{.PackageName}} not found")
		return
	}

	Success(c, {{.PackageName}}app.To{{.EntityName}}Response(entity))
}

// List 处理 GET /api/{{.PackageName}}s?page=1&page_size=20&sort_by=id&sort_order=desc
func (h *{{.EntityName}}Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	result, err := h.listHandler.Handle(c.Request.Context(), {{.PackageName}}app.List{{.EntityName}}sQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortOrder: sortOrder,
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

// Update 处理 PUT /api/{{.PackageName}}s/:id
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

// Delete 处理 DELETE /api/{{.PackageName}}s/:id
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

// Module 提供 {{.EntityName}} 的所有 Fx 依赖。
var Module = fx.Options(
	// Repository
	fx.Provide(func(db *gorm.DB) {{.PackageName}}.{{.EntityName}}Repository {
		return persistence.New{{.EntityName}}Repository(db)
	}),

	// Domain Services
	fx.Provide({{.PackageName}}.New{{.EntityName}}DomainService),

	// Command Handlers
	fx.Provide(NewCreate{{.EntityName}}Handler),
	fx.Provide(NewUpdate{{.EntityName}}Handler),
	fx.Provide(NewDelete{{.EntityName}}Handler),

	// Query Handlers
	fx.Provide(NewGet{{.EntityName}}Handler),
	fx.Provide(NewList{{.EntityName}}sHandler),
	
	// soliton-gen:services
	// soliton-gen:event-handlers

	// 可选：注册到 CQRS 总线
	// 取消注释以启用 CQRS 模式：
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

// RegisterMigration 注册 {{.EntityName}} 表的数据库迁移。
func RegisterMigration(db *gorm.DB) error {
	return persistence.Migrate{{.EntityName}}(db)
}
`
