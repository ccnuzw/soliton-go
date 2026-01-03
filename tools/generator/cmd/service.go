package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var serviceMethods string

var serviceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Generate an application service",
	Long: `Generate an application service for cross-domain business logic.
Services orchestrate multiple domains to implement complex use cases.

Examples:
  soliton-gen service OrderService
  soliton-gen service OrderService --methods "CreateOrder,CancelOrder,GetUserOrders"
  soliton-gen service PaymentService --methods "ProcessPayment,RefundPayment"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("🚀 Generating service: %s\n\n", name)

		generateService(name, serviceMethods)
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVar(&serviceMethods, "methods", "", "Comma-separated list of service methods (e.g., 'CreateOrder,CancelOrder')")
	serviceCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
}

// ServiceMethod represents a service method
type ServiceMethod struct {
	Name      string // Method name (e.g., "CreateOrder")
	CamelName string // Camel case name (e.g., "createOrder")
}

// ServiceData holds template data for service generation
type ServiceData struct {
	ServiceName string
	PackageName string
	Methods     []ServiceMethod
	ModulePath  string
}

func generateService(name string, methodsStr string) {
	// Normalize service name
	serviceName := name
	if !strings.HasSuffix(serviceName, "Service") {
		serviceName = serviceName + "Service"
	}

	// Parse methods
	var methods []ServiceMethod
	if methodsStr != "" {
		for _, m := range strings.Split(methodsStr, ",") {
			m = strings.TrimSpace(m)
			if m != "" {
				// Keep original case if already PascalCase, otherwise convert
				name := m
				if len(name) > 0 && name[0] >= 'a' && name[0] <= 'z' {
					name = strings.ToUpper(string(name[0])) + name[1:]
				}
				camelName := strings.ToLower(string(name[0])) + name[1:]
				methods = append(methods, ServiceMethod{
					Name:      name,
					CamelName: camelName,
				})
			}
		}
	} else {
		// Default methods
		baseName := strings.TrimSuffix(serviceName, "Service")
		methods = []ServiceMethod{
			{Name: "Create" + baseName, CamelName: "create" + baseName},
			{Name: "Get" + baseName, CamelName: "get" + baseName},
			{Name: "List" + baseName + "s", CamelName: "list" + baseName + "s"},
		}
	}

	packageName := strings.ToLower(strings.TrimSuffix(serviceName, "Service"))

	layout, err := ResolveProjectLayout()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	data := ServiceData{
		ServiceName: serviceName,
		PackageName: packageName,
		Methods:     methods,
		ModulePath:  layout.ModulePath,
	}

	// Determine paths
	serviceDir := filepath.Join(layout.AppDir, "services")
	_ = os.MkdirAll(serviceDir, 0755)

	// Generate files
	fmt.Println("📦 Application Service")
	generateServiceFile(filepath.Join(serviceDir, packageName+"_service.go"), serviceTemplate, data)
	generateServiceFile(filepath.Join(serviceDir, packageName+"_dto.go"), serviceDTOTemplate, data)
	generateServiceFile(filepath.Join(serviceDir, "module.go"), serviceModuleTemplate, data)

	// Summary
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("✅ Service generation complete!")
	fmt.Printf("   📁 Service: %s\n", serviceDir)
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n📋 Generated methods:")
	for _, m := range methods {
		fmt.Printf("   • %s(ctx, req) (*Response, error)\n", m.Name)
	}

	fmt.Println("\n📝 Next steps:")
	fmt.Println("   1. Inject required repositories in the service struct")
	fmt.Println("   2. Implement business logic in each method")
	fmt.Println("   3. Register service in main.go:")
	fmt.Printf("      fx.Provide(services.New%s),\n", serviceName)
}

func generateServiceFile(path string, tmpl string, data ServiceData) {
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

	funcMap := template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
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
// SERVICE TEMPLATES
// ============================================================================

const serviceTemplate = `package services

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

const serviceDTOTemplate = `package services

{{range .Methods}}
// {{.Name}}Request is the request for {{.Name}}.
type {{.Name}}Request struct {
	// TODO: Replace Example with real request fields.
	Example string ` + "`json:\"example\"`" + `
}

// {{.Name}}Response is the response for {{.Name}}.
type {{.Name}}Response struct {
	// TODO: Replace Example with real response fields.
	Example string ` + "`json:\"example\"`" + `
}
{{end}}
`

const serviceModuleTemplate = `package services

import "go.uber.org/fx"

// Module provides application service dependencies for Fx.
var Module = fx.Options(
	// fx.Provide(New{{.ServiceName}}),
)
`
