package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateDomain generates a domain module with the given configuration.
func GenerateDomain(cfg DomainConfig) (*GenerationResult, error) {
	return generateDomainInternal(cfg, false)
}

// PreviewDomain previews what files would be created without actually creating them.
func PreviewDomain(cfg DomainConfig) (*GenerationResult, error) {
	return generateDomainInternal(cfg, true)
}

func generateDomainInternal(cfg DomainConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{
		Success: true,
		Files:   []GeneratedFile{},
	}

	// Normalize name
	entityName := ToPascalCase(cfg.Name)
	packageName := strings.ToLower(cfg.Name)

	// Resolve project layout
	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	// Convert field configs to template fields
	fields := ConvertToFields(cfg.Fields, entityName, packageName)

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

	tableName := cfg.TableName
	if tableName == "" {
		tableName = Pluralize(packageName)
	}

	routeBase := cfg.RouteBase
	if routeBase == "" {
		routeBase = Pluralize(packageName)
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
		SoftDelete:  cfg.SoftDelete,
	}

	// Domain Layer
	domainDir := filepath.Join(layout.DomainDir, packageName)
	if !previewOnly {
		_ = os.MkdirAll(domainDir, 0755)
	}

	domainFiles := []struct {
		path     string
		template string
	}{
		{filepath.Join(domainDir, packageName+".go"), EntityTemplate},
		{filepath.Join(domainDir, "repository.go"), RepoTemplate},
		{filepath.Join(domainDir, "events.go"), EventsTemplate},
	}

	for _, f := range domainFiles {
		genFile := generateDomainFile(f.path, f.template, data, cfg.Force, previewOnly)
		result.Files = append(result.Files, genFile)
	}

	// Infrastructure Layer
	if !previewOnly {
		_ = os.MkdirAll(layout.InfraDir, 0755)
	}
	repoImplFile := generateDomainFile(
		filepath.Join(layout.InfraDir, packageName+"_repo.go"),
		RepoImplTemplate,
		data,
		cfg.Force,
		previewOnly,
	)
	result.Files = append(result.Files, repoImplFile)

	// Application Layer
	appModuleDir := filepath.Join(layout.AppDir, packageName)
	if !previewOnly {
		_ = os.MkdirAll(appModuleDir, 0755)
	}

	appFiles := []struct {
		path     string
		template string
	}{
		{filepath.Join(appModuleDir, "commands.go"), CommandsTemplate},
		{filepath.Join(appModuleDir, "queries.go"), QueriesTemplate},
		{filepath.Join(appModuleDir, "dto.go"), DTOTemplate},
		{filepath.Join(appModuleDir, "module.go"), FxModuleTemplate},
	}

	for _, f := range appFiles {
		genFile := generateDomainFile(f.path, f.template, data, cfg.Force, previewOnly)
		result.Files = append(result.Files, genFile)
	}

	// Interfaces Layer (HTTP)
	if !previewOnly {
		_ = os.MkdirAll(layout.InterfacesDir, 0755)
	}

	// Generate helpers.go if it doesn't exist
	helpersPath := filepath.Join(layout.InterfacesDir, "helpers.go")
	if _, err := os.Stat(helpersPath); os.IsNotExist(err) || previewOnly {
		helpersFile := generateDomainFile(helpersPath, HelpersHTTPTemplate, data, false, previewOnly)
		result.Files = append(result.Files, helpersFile)
	}

	handlerFile := generateDomainFile(
		filepath.Join(layout.InterfacesDir, packageName+"_handler.go"),
		HandlerTemplate,
		data,
		cfg.Force,
		previewOnly,
	)
	result.Files = append(result.Files, handlerFile)

	// Wire main.go if requested
	if cfg.Wire && !previewOnly {
		mainGoPath := filepath.Join(filepath.Dir(layout.InternalDir), "cmd", "main.go")
		if WireMainGo(mainGoPath, entityName, packageName, layout.ModulePath) {
			result.Message = fmt.Sprintf("Domain %s 生成成功，已自动注入到 main.go", entityName)
		} else {
			result.Message = fmt.Sprintf("Domain %s 生成成功（需手动注入到 main.go）", entityName)
		}
	} else {
		result.Message = fmt.Sprintf("Domain %s 生成成功", entityName)
	}

	// Check for errors
	for _, f := range result.Files {
		if f.Status == FileStatusError {
			result.Success = false
			result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", f.Path))
		}
	}

	return result, nil
}

func generateDomainFile(path string, tmpl string, data TemplateData, force bool, previewOnly bool) GeneratedFile {
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
		"title":            strings.Title,
		"lower":            strings.ToLower,
		"upper":            strings.ToUpper,
		"enumConst":        EnumConst,
		"createBindingTag": CreateBindingTag,
		"updateBindingTag": UpdateBindingTag,
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

// ValidateDomainConfig validates the domain configuration.
func ValidateDomainConfig(cfg DomainConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("domain name is required")
	}
	if strings.ContainsAny(cfg.Name, " \t\n/\\") {
		return fmt.Errorf("domain name contains invalid characters")
	}
	return nil
}

// WireMainGo attempts to inject module into main.go using marker comments.
func WireMainGo(mainGoPath, entityName, packageName, modulePath string) bool {
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

func wireMainGoNew(mainGoPath, entityName, packageName, modulePath, original string) bool {
	result := original
	modified := false

	// 0. Replace blank gorm import with normal import
	if strings.Contains(result, "_ \"gorm.io/gorm\"") && !strings.Contains(result, "\t\"gorm.io/gorm\"") {
		result = strings.Replace(result, "_ \"gorm.io/gorm\"", "\"gorm.io/gorm\"", 1)
		modified = true
	}

	// 1. Add app import
	appImport := fmt.Sprintf("%sapp \"%s/internal/application/%s\"", packageName, modulePath, packageName)
	if !strings.Contains(result, appImport) {
		result = strings.Replace(result,
			"\t// soliton-gen:imports",
			"\t"+appImport+"\n\t// soliton-gen:imports",
			1)
		modified = true
	}

	// 2. Add interfaceshttp import
	httpImport := fmt.Sprintf("interfaceshttp \"%s/internal/interfaces/http\"", modulePath)
	if !strings.Contains(result, httpImport) {
		result = strings.Replace(result,
			"\t// soliton-gen:imports",
			fmt.Sprintf("\tinterfaceshttp \"%s/internal/interfaces/http\"\n\t// soliton-gen:imports", modulePath),
			1)
		modified = true
	}

	// 3. Add module
	moduleCode := fmt.Sprintf("%sapp.Module,", packageName)
	if !strings.Contains(result, moduleCode) {
		result = strings.Replace(result,
			"\t\t// soliton-gen:modules",
			"\t\t"+moduleCode+"\n\t\t// soliton-gen:modules",
			1)
		modified = true
	}

	// 4. Add handler
	handlerCode := fmt.Sprintf("fx.Provide(interfaceshttp.New%sHandler),", entityName)
	if !strings.Contains(result, handlerCode) {
		result = strings.Replace(result,
			"\t\t// soliton-gen:handlers",
			"\t\t"+handlerCode+"\n\t\t// soliton-gen:handlers",
			1)
		modified = true
	}

	// 5. Add route registration
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

func wireMainGoLegacy(mainGoPath, entityName, packageName, modulePath, original string) bool {
	result := original
	modified := false

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

	// Uncomment invoke block
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

// GetAvailableFieldTypes returns the list of available field types.
func GetAvailableFieldTypes() []map[string]string {
	return []map[string]string{
		{"type": "string", "description": "String (varchar 255)"},
		{"type": "text", "description": "Text (long text)"},
		{"type": "int", "description": "Integer (32-bit)"},
		{"type": "int64", "description": "Integer (64-bit)"},
		{"type": "float64", "description": "Float (64-bit)"},
		{"type": "bool", "description": "Boolean"},
		{"type": "time", "description": "Timestamp"},
		{"type": "time?", "description": "Optional Timestamp"},
		{"type": "uuid", "description": "UUID (indexed)"},
		{"type": "enum", "description": "Enum (requires enum_values)"},
	}
}
