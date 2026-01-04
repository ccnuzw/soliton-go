package core

// ============================================================================
// SERVICE TEMPLATES
// ============================================================================

const ServiceTemplate = `package services

import (
	"context"
	"errors"

	// Import your domain repositories here:
	// "{{.ModulePath}}/internal/domain/user"
	// "{{.ModulePath}}/internal/domain/order"
)

// {{.ServiceName}} handles cross-domain business logic.
type {{.ServiceName}} struct {
	// Add your repositories here:
	// userRepo  user.UserRepository
	// orderRepo order.OrderRepository
}

// New{{.ServiceName}} creates a new {{.ServiceName}}.
func New{{.ServiceName}}(
	// Add your repository parameters here:
	// userRepo user.UserRepository,
	// orderRepo order.OrderRepository,
) *{{.ServiceName}} {
	return &{{.ServiceName}}{
		// userRepo:  userRepo,
		// orderRepo: orderRepo,
	}
}

{{range .Methods}}
// {{.Name}} implements the {{.Name}} use case.
func (s *{{$.ServiceName}}) {{.Name}}(ctx context.Context, req {{.Name}}Request) (*{{.Name}}Response, error) {
	// TODO: Implement business logic
	// Example:
	// 1. Validate request
	// 2. Load entities from repositories
	// 3. Execute domain logic
	// 4. Save changes
	// 5. Publish domain events
	// 6. Return response

	return nil, errors.New("not implemented")
}
{{end}}
`

const ServiceDTOTemplate = `package services

{{range .Methods}}
// {{.Name}}Request is the request for {{.Name}}.
type {{.Name}}Request struct {
	// Add your request fields here.
	// Common patterns:
	ID string ` + "`json:\"id,omitempty\"`" + ` // Entity ID (for Get/Update/Delete operations)
	// Data   any    ` + "`json:\"data,omitempty\"`" + ` // Payload for Create/Update operations
}

// {{.Name}}Response is the response for {{.Name}}.
type {{.Name}}Response struct {
	Success bool   ` + "`json:\"success\"`" + `           // Operation success flag
	Message string ` + "`json:\"message,omitempty\"`" + ` // Human-readable message
	Data    any    ` + "`json:\"data,omitempty\"`" + `    // Response payload
}
{{end}}
`

const ServiceModuleTemplate = `package services

import "go.uber.org/fx"

// Module provides application service dependencies for Fx.
var Module = fx.Options(
	// fx.Provide(New{{.ServiceName}}),
)
`
