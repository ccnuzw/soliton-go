package core

// ============================================================================
// DDD TEMPLATES / 领域增强模板
// ============================================================================

const ValueObjectTemplate = `package {{.PackageName}}

import (
	"errors"
	"reflect"
{{- if .HasTime}}
	"time"
{{- end}}
)

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

// {{.ValueObjectName}} 是领域值对象。
type {{.ValueObjectName}} struct {
{{- range .Fields}}
	{{.Name}} {{.GoType}} {{.JsonTag}}{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}

// New{{.ValueObjectName}} 创建一个新的 {{.ValueObjectName}}。
func New{{.ValueObjectName}}({{range $i, $f := .Fields}}{{if $i}}, {{end}}{{$f.CamelName}} {{$f.GoType}}{{end}}) ({{.ValueObjectName}}, error) {
	vo := {{.ValueObjectName}}{
{{- range .Fields}}
		{{.Name}}: {{.CamelName}},
{{- end}}
	}
	if err := vo.Validate(); err != nil {
		return {{.ValueObjectName}}{}, err
	}
	return vo, nil
}

// Validate 执行值对象的领域校验规则。
func (v {{.ValueObjectName}}) Validate() error {
	// TODO: 在此添加校验逻辑
	return errors.New("not implemented")
}

// Equals 比较两个值对象是否相等。
func (v {{.ValueObjectName}}) Equals(other {{.ValueObjectName}}) bool {
	return reflect.DeepEqual(v, other)
}
`

const SpecificationTemplate = `package {{.PackageName}}

// {{.SpecificationName}} 表示一个领域规格（Specification）。
type {{.SpecificationName}} struct{}

// IsSatisfiedBy 判断目标对象是否满足规格。
func (s {{.SpecificationName}}) IsSatisfiedBy(target {{if .TargetIsAny}}any{{else}}*{{.TargetType}}{{end}}) bool {
	// TODO: 实现规格校验逻辑
	return true
}
`

const PolicyTemplate = `package {{.PackageName}}

import "errors"

// {{.PolicyName}} 表示一个领域策略（Policy）。
type {{.PolicyName}} struct{}

// Validate 执行策略校验，返回错误表示不满足策略。
func (p {{.PolicyName}}) Validate(target {{if .TargetIsAny}}any{{else}}*{{.TargetType}}{{end}}) error {
	// TODO: 实现策略校验逻辑
	return errors.New("not implemented")
}
`

const EventTemplate = `package {{.PackageName}}

import (
{{- if .HasTime}}
	"time"
{{- end}}
	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
)

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

// {{.EventStructName}} 是领域事件。
type {{.EventStructName}} struct {
	ddd.BaseDomainEvent
{{- range .Fields}}
	{{.Name}} {{.GoType}} {{.JsonTag}}{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}

// EventName 返回事件名称（主题）。
func (e {{.EventStructName}}) EventName() string {
	return "{{.EventTopic}}"
}

// New{{.EventStructName}} 创建一个新的事件实例。
func New{{.EventStructName}}({{range $i, $f := .Fields}}{{if $i}}, {{end}}{{$f.CamelName}} {{$f.GoType}}{{end}}) {{.EventStructName}} {
	return {{.EventStructName}}{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
{{- range .Fields}}
		{{.Name}}: {{.CamelName}},
{{- end}}
	}
}

// init 将事件注册到全局注册表。
func init() {
	event.RegisterEvent("{{.EventTopic}}", func() ddd.DomainEvent {
		return &{{.EventStructName}}{}
	})
}
`

const EventHandlerTemplate = `package {{.PackageName}}

import (
	"context"
	"fmt"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/event"
	"go.uber.org/fx"

	"{{.ModulePath}}/internal/domain/{{.DomainPackage}}"
)

// {{.HandlerName}} 处理 {{.EventStructName}} 事件。
type {{.HandlerName}} struct{}

// New{{.HandlerName}} 创建 {{.HandlerName}} 实例。
func New{{.HandlerName}}() *{{.HandlerName}} {
	return &{{.HandlerName}}{}
}

// Handle 处理领域事件。
func (h *{{.HandlerName}}) Handle(ctx context.Context, evt ddd.DomainEvent) error {
	e, ok := evt.(*{{.DomainPackage}}.{{.EventStructName}})
	if !ok {
		return fmt.Errorf("unexpected event type: %T", evt)
	}
	_ = e
	// TODO: 在此实现事件处理逻辑
	return nil
}

// Register{{.HandlerName}} 注册事件处理器。
func Register{{.HandlerName}}(lc fx.Lifecycle, bus event.EventBus, handler *{{.HandlerName}}) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return bus.Subscribe(ctx, "{{.EventTopic}}", handler.Handle)
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
`
