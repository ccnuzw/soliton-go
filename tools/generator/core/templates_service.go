package core

// ============================================================================
// SERVICE TEMPLATES / 服务模板
// ============================================================================

// ServiceTemplate is the template for generating a service file.
const ServiceTemplate = `package {{.PackageName}}

import (
	"context"
	"errors"

	// 在此导入领域层的 Repository：
	// "{{.ModulePath}}/internal/domain/user"
	// "{{.ModulePath}}/internal/domain/order"
)

// {{.ServiceName}} 处理跨领域的业务逻辑编排。
type {{.ServiceName}} struct {
	// 在此添加依赖的 Repository：
	// userRepo  user.UserRepository
	// orderRepo order.OrderRepository
}

// New{{.ServiceName}} 创建 {{.ServiceName}} 实例。
func New{{.ServiceName}}(
	// 在此添加 Repository 参数：
	// userRepo user.UserRepository,
	// orderRepo order.OrderRepository,
) *{{.ServiceName}} {
	return &{{.ServiceName}}{
		// userRepo:  userRepo,
		// orderRepo: orderRepo,
	}
}

{{range .Methods}}
// {{.Name}} 实现 {{.Name}} 用例。
func (s *{{$.ServiceName}}) {{.Name}}(ctx context.Context, req {{.Name}}Request) (*{{.Name}}Response, error) {
	// TODO: 实现业务逻辑
	// 示例步骤：
	// 1. 校验请求参数
	// 2. 从 Repository 加载实体
	// 3. 执行领域逻辑
	// 4. 保存变更
	// 5. 发布领域事件
	// 6. 返回响应

	return nil, errors.New("not implemented")
}
{{end}}
`

const ServiceDTOTemplate = `package {{.PackageName}}

{{range .Methods}}
// {{.Name}}Request 是 {{.Name}} 方法的请求参数。
type {{.Name}}Request struct {
	// 在此添加请求字段：
	ID string ` + "`json:\"id,omitempty\"`" + ` // 实体 ID（用于 Get/Update/Delete 操作）
	// Data   any    ` + "`json:\"data,omitempty\"`" + ` // 请求数据（用于 Create/Update 操作）
}

// {{.Name}}Response 是 {{.Name}} 方法的响应结果。
type {{.Name}}Response struct {
	Success bool   ` + "`json:\"success\"`" + `           // 操作是否成功
	Message string ` + "`json:\"message,omitempty\"`" + ` // 提示消息
	Data    any    ` + "`json:\"data,omitempty\"`" + `    // 响应数据
}
{{end}}
`

const ServiceModuleTemplate = `package services

import "go.uber.org/fx"

// Module 提供应用服务的依赖注入配置。
var Module = fx.Options(
	// fx.Provide(New{{.ServiceName}}),
)
`
