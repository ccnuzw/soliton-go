package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateService generates a service with the given configuration.
func GenerateService(cfg ServiceConfig) (*GenerationResult, error) {
	return generateServiceInternal(cfg, false)
}

// PreviewService previews what files would be created without actually creating them.
func PreviewService(cfg ServiceConfig) (*GenerationResult, error) {
	return generateServiceInternal(cfg, true)
}

// ServiceDetectionResult represents the result of service type detection.
type ServiceDetectionResult struct {
	ServiceName  string
	DomainName   string
	DomainExists bool
	ServiceType  string // "domain_service" | "cross_domain_service"
	TargetDir    string
	Message      string
}

// DetectServiceType detects the service type based on domain existence.
func DetectServiceType(serviceName string) (*ServiceDetectionResult, error) {
	// Normalize service name
	if !strings.HasSuffix(serviceName, "Service") {
		serviceName = serviceName + "Service"
	}
	domainName := strings.ToLower(strings.TrimSuffix(serviceName, "Service"))

	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	// Check if domain exists
	domainPath := filepath.Join(layout.DomainDir, domainName)
	domainExists := IsDir(domainPath)

	result := &ServiceDetectionResult{
		ServiceName:  serviceName,
		DomainName:   domainName,
		DomainExists: domainExists,
		TargetDir:    filepath.Join("internal/application", domainName),
	}

	if domainExists {
		result.ServiceType = "domain_service"
		result.Message = fmt.Sprintf("检测到已有 %s 领域，将生成领域服务并生成独立的服务层 DTO (service_dto.go)", domainName)
	} else {
		result.ServiceType = "cross_domain_service"
		result.Message = fmt.Sprintf("未检测到 %s 领域，将生成跨领域服务", domainName)
	}

	return result, nil
}

func generateServiceInternal(cfg ServiceConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{
		Success: true,
		Files:   []GeneratedFile{},
	}

	// Normalize service name
	serviceName := cfg.Name
	if !strings.HasSuffix(serviceName, "Service") {
		serviceName = serviceName + "Service"
	}

	packageName := strings.ToLower(strings.TrimSuffix(serviceName, "Service"))

	// Resolve project layout
	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	// Parse methods
	methods := ParseServiceMethods(cfg.Methods, serviceName)

	data := ServiceData{
		ServiceName:   serviceName,
		ServiceRemark: strings.TrimSpace(cfg.Remark),
		PackageName:   packageName,
		Methods:       methods,
		ModulePath:    layout.ModulePath,
	}

	// Auto-detect: check if domain exists
	domainPath := filepath.Join(layout.DomainDir, packageName)
	domainExists := IsDir(domainPath)

	// Determine paths based on detection
	var serviceDir string
	var serviceFileName string
	var dtoFileName string
	var shouldGenerateDTO bool

	if domainExists {
		// Domain exists -> generate in application/{domain}/
		serviceDir = filepath.Join(layout.InternalDir, "application", packageName)
		serviceFileName = "service.go"
		dtoFileName = "service_dto.go" // Use separate DTO file to avoid collisions

		// Check if service DTO already exists
		dtoPath := filepath.Join(serviceDir, dtoFileName)
		shouldGenerateDTO = !IsFile(dtoPath) || cfg.Force

		// Update package name to match domain application package
		data.PackageName = packageName + "app"
	} else {
		// Domain doesn't exist -> generate in application/{service}/
		serviceDir = filepath.Join(layout.InternalDir, "application", packageName)
		serviceFileName = "service.go"
		dtoFileName = "service_dto.go"
		shouldGenerateDTO = true

		// Use service name as package
		data.PackageName = packageName + "app"

		// Check if DTO exists (for force logic)
		dtoPath := filepath.Join(serviceDir, dtoFileName)
		shouldGenerateDTO = !IsFile(dtoPath) || cfg.Force
	}

	if !previewOnly {
		_ = os.MkdirAll(serviceDir, 0755)
	}

	// Generate service file
	serviceFile := generateServiceFile(
		filepath.Join(serviceDir, serviceFileName),
		ServiceTemplate,
		data,
		cfg.Force,
		previewOnly,
	)
	result.Files = append(result.Files, serviceFile)

	// Generate DTO file (only if needed)
	if shouldGenerateDTO {
		dtoFile := generateServiceFile(
			filepath.Join(serviceDir, dtoFileName),
			ServiceDTOTemplate,
			data,
			cfg.Force,
			previewOnly,
		)
		result.Files = append(result.Files, dtoFile)
	}

	// Generate or update module.go
	if !previewOnly {
		generateOrUpdateServiceModuleGo(filepath.Join(serviceDir, "module.go"), serviceName, data.PackageName)
	}

	if shouldGenerateDTO {
		result.Message = fmt.Sprintf("Service %s 及 DTO (%s) 生成成功", serviceName, dtoFileName)
	} else {
		result.Message = fmt.Sprintf("Service %s 生成成功", serviceName)
	}

	// Check for errors
	var skippedCount int
	for _, f := range result.Files {
		if f.Status == FileStatusError {
			result.Success = false
			result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", f.Path))
		} else if f.Status == FileStatusSkip {
			skippedCount++
		}
	}

	if skippedCount > 0 {
		var skippedFiles []string
		for _, f := range result.Files {
			if f.Status == FileStatusSkip {
				skippedFiles = append(skippedFiles, filepath.Base(f.Path))
			}
		}
		result.Message += fmt.Sprintf("\n⚠️  跳过更新: %s", strings.Join(skippedFiles, ", "))
		result.Message += "\n   (方法变更不会自动同步，使用 --force 强制覆盖)"
	}

	return result, nil
}

func generateServiceFile(path string, tmpl string, data ServiceData, force bool, previewOnly bool) GeneratedFile {
	genFile := GeneratedFile{
		Path: path,
	}

	// Check if file exists
	exists := false
	if _, err := os.Stat(path); err == nil {
		exists = true
		if !force && !previewOnly {
			genFile.Status = FileStatusSkip
			return genFile
		}
	}

	// Template functions
	funcMap := template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
	}

	// Render template
	t := template.Must(template.New("file").Funcs(funcMap).Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		genFile.Status = FileStatusError
		return genFile
	}

	content := buf.String()
	genFile.Content = content

	if previewOnly {
		if exists {
			genFile.Status = FileStatusOverwrite
		} else {
			genFile.Status = FileStatusNew
		}
		return genFile
	}

	// Create file
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		genFile.Status = FileStatusError
		return genFile
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		genFile.Status = FileStatusError
		return genFile
	}

	if exists {
		genFile.Status = FileStatusOverwrite
	} else {
		genFile.Status = FileStatusNew
	}
	genFile.Content = "" // Don't include content in non-preview mode

	return genFile
}

func generateOrUpdateServiceModuleGo(path string, serviceName string, packageName string) {
	provideCode := fmt.Sprintf("fx.Provide(New%s),", serviceName)

	// Check if file exists
	content, err := os.ReadFile(path)
	if err == nil {
		// File exists, check if service is already registered
		if strings.Contains(string(content), provideCode) {
			return
		}

		// Append new service to existing module.go
		result := string(content)
		marker := "// soliton-gen:services"

		if strings.Contains(result, marker) {
			// Has marker, insert before it
			result = strings.Replace(result,
				marker,
				provideCode+"\n\t"+marker,
				1)
		} else {
			// No marker, try to find var Module = fx.Options(
			// We look for where Module is defined
			moduleDef := "var Module = fx.Options("
			idx := strings.Index(result, moduleDef)
			if idx >= 0 {
				// Find the matching closing parenthesis for this block
				// This is a naive implementation; for robust parsing we might need AST,
				// but for generated code, counting braces from the start of definition works reasonably well.
				startSearch := idx + len(moduleDef)
				balance := 1
				insertPos := -1

				for i := startSearch; i < len(result); i++ {
					if result[i] == '(' {
						balance++
					} else if result[i] == ')' {
						balance--
					}

					if balance == 0 {
						insertPos = i
						break
					}
				}

				if insertPos > 0 {
					// We found the closing bracket of fx.Options(...)
					// Insert before it
					// Check if previous line has comma, if not add it (though usually it ends with comma or newline)
					insertion := "\t" + provideCode + "\n\t// soliton-gen:services\n"
					result = result[:insertPos] + insertion + result[insertPos:]
				}
			}
		}

		_ = os.WriteFile(path, []byte(result), 0644)
		return
	}

	// File doesn't exist, create new one with marker
	newContent := fmt.Sprintf(`package %s

import "go.uber.org/fx"

// Module provides application service dependencies for Fx.
var Module = fx.Options(
	%s
	// soliton-gen:services
)
`, packageName, provideCode)

	_ = os.WriteFile(path, []byte(newContent), 0644)
}

// ValidateServiceConfig validates the service configuration.
func ValidateServiceConfig(cfg ServiceConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if strings.ContainsAny(cfg.Name, " \t\n/\\") {
		return fmt.Errorf("service name contains invalid characters")
	}
	return nil
}
