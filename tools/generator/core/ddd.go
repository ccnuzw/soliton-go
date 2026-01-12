package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateValueObject generates a value object for a domain.
func GenerateValueObject(cfg ValueObjectConfig) (*GenerationResult, error) {
	return generateValueObjectInternal(cfg, false)
}

// PreviewValueObject previews value object generation.
func PreviewValueObject(cfg ValueObjectConfig) (*GenerationResult, error) {
	return generateValueObjectInternal(cfg, true)
}

func generateValueObjectInternal(cfg ValueObjectConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{Success: true}

	if err := ValidateValueObjectConfig(cfg); err != nil {
		return nil, err
	}

	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	domainName := strings.ToLower(cfg.Domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !IsDir(domainDir) {
		return nil, fmt.Errorf("domain %s not found in %s", domainName, layout.DomainDir)
	}

	valueObjectName := ToPascalCase(cfg.Name)
	fieldsCfg := cfg.Fields
	if len(fieldsCfg) == 0 {
		fieldsCfg = []FieldConfig{{Name: "Value", Type: "string"}}
	}
	fields := ConvertToFields(fieldsCfg, valueObjectName, domainName)
	hasTime, hasEnums := detectFieldFlags(fields)

	data := ValueObjectData{
		PackageName:     domainName,
		ValueObjectName: valueObjectName,
		Fields:          fields,
		HasTime:         hasTime,
		HasEnums:        hasEnums,
	}

	fileName := fmt.Sprintf("value_object_%s.go", ToSnakeCase(valueObjectName))
	genFile := generateTemplateFile(filepath.Join(domainDir, fileName), ValueObjectTemplate, data, cfg.Force, previewOnly)
	result.Files = append(result.Files, genFile)
	result.Message = fmt.Sprintf("ValueObject %s 生成成功", valueObjectName)

	if genFile.Status == FileStatusError {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", genFile.Path))
	}

	return result, nil
}

// ValidateValueObjectConfig validates value object configuration.
func ValidateValueObjectConfig(cfg ValueObjectConfig) error {
	if cfg.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cfg.Name == "" {
		return fmt.Errorf("value object name is required")
	}
	if strings.ContainsAny(cfg.Name, " \t\n/\\") {
		return fmt.Errorf("value object name contains invalid characters")
	}
	return nil
}

// GenerateSpecification generates a specification for a domain.
func GenerateSpecification(cfg SpecificationConfig) (*GenerationResult, error) {
	return generateSpecificationInternal(cfg, false)
}

// PreviewSpecification previews specification generation.
func PreviewSpecification(cfg SpecificationConfig) (*GenerationResult, error) {
	return generateSpecificationInternal(cfg, true)
}

func generateSpecificationInternal(cfg SpecificationConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{Success: true}

	if err := ValidateSpecificationConfig(cfg); err != nil {
		return nil, err
	}

	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	domainName := strings.ToLower(cfg.Domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !IsDir(domainDir) {
		return nil, fmt.Errorf("domain %s not found in %s", domainName, layout.DomainDir)
	}

	specName := ToPascalCase(cfg.Name)
	targetType, targetIsAny := normalizeTargetType(cfg.Target)

	data := SpecificationData{
		PackageName:       domainName,
		SpecificationName: specName,
		TargetType:        targetType,
		TargetIsAny:       targetIsAny,
	}

	fileName := fmt.Sprintf("spec_%s.go", ToSnakeCase(specName))
	genFile := generateTemplateFile(filepath.Join(domainDir, fileName), SpecificationTemplate, data, cfg.Force, previewOnly)
	result.Files = append(result.Files, genFile)
	result.Message = fmt.Sprintf("Specification %s 生成成功", specName)

	if genFile.Status == FileStatusError {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", genFile.Path))
	}

	return result, nil
}

// ValidateSpecificationConfig validates specification configuration.
func ValidateSpecificationConfig(cfg SpecificationConfig) error {
	if cfg.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cfg.Name == "" {
		return fmt.Errorf("specification name is required")
	}
	if strings.ContainsAny(cfg.Name, " \t\n/\\") {
		return fmt.Errorf("specification name contains invalid characters")
	}
	return nil
}

// GeneratePolicy generates a policy for a domain.
func GeneratePolicy(cfg PolicyConfig) (*GenerationResult, error) {
	return generatePolicyInternal(cfg, false)
}

// PreviewPolicy previews policy generation.
func PreviewPolicy(cfg PolicyConfig) (*GenerationResult, error) {
	return generatePolicyInternal(cfg, true)
}

func generatePolicyInternal(cfg PolicyConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{Success: true}

	if err := ValidatePolicyConfig(cfg); err != nil {
		return nil, err
	}

	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	domainName := strings.ToLower(cfg.Domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !IsDir(domainDir) {
		return nil, fmt.Errorf("domain %s not found in %s", domainName, layout.DomainDir)
	}

	policyName := ToPascalCase(cfg.Name)
	targetType, targetIsAny := normalizeTargetType(cfg.Target)

	data := PolicyData{
		PackageName: domainName,
		PolicyName:  policyName,
		TargetType:  targetType,
		TargetIsAny: targetIsAny,
	}

	fileName := fmt.Sprintf("policy_%s.go", ToSnakeCase(policyName))
	genFile := generateTemplateFile(filepath.Join(domainDir, fileName), PolicyTemplate, data, cfg.Force, previewOnly)
	result.Files = append(result.Files, genFile)
	result.Message = fmt.Sprintf("Policy %s 生成成功", policyName)

	if genFile.Status == FileStatusError {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", genFile.Path))
	}

	return result, nil
}

// ValidatePolicyConfig validates policy configuration.
func ValidatePolicyConfig(cfg PolicyConfig) error {
	if cfg.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cfg.Name == "" {
		return fmt.Errorf("policy name is required")
	}
	if strings.ContainsAny(cfg.Name, " \t\n/\\") {
		return fmt.Errorf("policy name contains invalid characters")
	}
	return nil
}

// GenerateEvent generates a custom domain event for a domain.
func GenerateEvent(cfg EventConfig) (*GenerationResult, error) {
	return generateEventInternal(cfg, false)
}

// PreviewEvent previews domain event generation.
func PreviewEvent(cfg EventConfig) (*GenerationResult, error) {
	return generateEventInternal(cfg, true)
}

func generateEventInternal(cfg EventConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{Success: true}

	if err := ValidateEventConfig(cfg); err != nil {
		return nil, err
	}

	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	domainName := strings.ToLower(cfg.Domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !IsDir(domainDir) {
		return nil, fmt.Errorf("domain %s not found in %s", domainName, layout.DomainDir)
	}

	eventStructName := normalizeEventStructName(cfg.Name)
	eventTopic := normalizeEventTopic(cfg.Topic, domainName, eventStructName)
	fields := ConvertToFieldsAllowReserved(cfg.Fields, eventStructName, domainName)
	hasTime, hasEnums := detectFieldFlags(fields)

	data := EventData{
		PackageName:     domainName,
		EventStructName: eventStructName,
		EventTopic:      eventTopic,
		Fields:          fields,
		HasTime:         hasTime,
		HasEnums:        hasEnums,
	}

	fileName := fmt.Sprintf("event_%s.go", ToSnakeCase(strings.TrimSuffix(eventStructName, "Event")))
	genFile := generateTemplateFile(filepath.Join(domainDir, fileName), EventTemplate, data, cfg.Force, previewOnly)
	result.Files = append(result.Files, genFile)
	result.Message = fmt.Sprintf("Domain Event %s 生成成功", eventStructName)

	if genFile.Status == FileStatusError {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", genFile.Path))
	}

	return result, nil
}

// ValidateEventConfig validates event configuration.
func ValidateEventConfig(cfg EventConfig) error {
	if cfg.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cfg.Name == "" {
		return fmt.Errorf("event name is required")
	}
	if strings.ContainsAny(cfg.Name, " \t\n/\\") {
		return fmt.Errorf("event name contains invalid characters")
	}
	return nil
}

// GenerateEventHandler generates a domain event handler.
func GenerateEventHandler(cfg EventHandlerConfig) (*GenerationResult, error) {
	return generateEventHandlerInternal(cfg, false)
}

// PreviewEventHandler previews event handler generation.
func PreviewEventHandler(cfg EventHandlerConfig) (*GenerationResult, error) {
	return generateEventHandlerInternal(cfg, true)
}

func generateEventHandlerInternal(cfg EventHandlerConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{Success: true}

	if err := ValidateEventHandlerConfig(cfg); err != nil {
		return nil, err
	}

	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, fmt.Errorf("could not resolve project layout: %w", err)
	}

	domainName := strings.ToLower(cfg.Domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !IsDir(domainDir) {
		return nil, fmt.Errorf("domain %s not found in %s", domainName, layout.DomainDir)
	}

	appDir := filepath.Join(layout.AppDir, domainName)
	if !previewOnly {
		_ = os.MkdirAll(appDir, 0755)
	}

	eventStructName := normalizeEventStructName(cfg.EventName)
	handlerName := strings.TrimSuffix(eventStructName, "Event") + "Handler"
	eventTopic := normalizeEventTopic(cfg.Topic, domainName, eventStructName)

	data := EventHandlerData{
		PackageName:     domainName + "app",
		DomainPackage:   domainName,
		EventStructName: eventStructName,
		HandlerName:     handlerName,
		EventTopic:      eventTopic,
		ModulePath:      layout.ModulePath,
	}

	fileName := fmt.Sprintf("event_handler_%s.go", ToSnakeCase(strings.TrimSuffix(eventStructName, "Event")))
	genFile := generateTemplateFile(filepath.Join(appDir, fileName), EventHandlerTemplate, data, cfg.Force, previewOnly)
	result.Files = append(result.Files, genFile)

	modulePath := filepath.Join(appDir, "module.go")
	mainGoPath := filepath.Join(filepath.Dir(layout.InternalDir), "cmd", "main.go")
	if previewOnly {
		if modulePreview := previewModuleForEventHandler(modulePath, handlerName); modulePreview != nil {
			result.Files = append(result.Files, *modulePreview)
		}
		if mainPreview := previewEventBusProvider(mainGoPath); mainPreview != nil {
			result.Files = append(result.Files, *mainPreview)
		}
	} else {
		_ = updateModuleForEventHandler(modulePath, handlerName)
		_ = ensureEventBusProvider(mainGoPath)
	}

	result.Message = fmt.Sprintf("Event Handler %s 生成成功", handlerName)

	if genFile.Status == FileStatusError {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", genFile.Path))
	}

	return result, nil
}

// ValidateEventHandlerConfig validates event handler configuration.
func ValidateEventHandlerConfig(cfg EventHandlerConfig) error {
	if cfg.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if cfg.EventName == "" {
		return fmt.Errorf("event name is required")
	}
	if strings.ContainsAny(cfg.EventName, " \t\n/\\") {
		return fmt.Errorf("event name contains invalid characters")
	}
	return nil
}

func normalizeTargetType(target string) (string, bool) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "any", true
	}
	target = strings.TrimPrefix(target, "*")
	return target, false
}

func normalizeEventStructName(name string) string {
	eventName := ToPascalCase(name)
	if !strings.HasSuffix(strings.ToLower(eventName), "event") {
		return eventName + "Event"
	}
	return eventName[:len(eventName)-5] + "Event"
}

func normalizeEventTopic(topic, domainName, eventStructName string) string {
	if strings.TrimSpace(topic) != "" {
		return topic
	}

	baseName := strings.TrimSuffix(eventStructName, "Event")
	domainPascal := ToPascalCase(domainName)
	if strings.HasPrefix(baseName, domainPascal) {
		baseName = strings.TrimPrefix(baseName, domainPascal)
	}
	baseName = strings.TrimSpace(baseName)
	if baseName == "" {
		baseName = "event"
	}

	return fmt.Sprintf("%s.%s", domainName, ToSnakeCase(baseName))
}

func detectFieldFlags(fields []Field) (bool, bool) {
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
	return hasTime, hasEnums
}

func updateModuleContentForEventHandler(content string, handlerName string) (string, bool) {
	provideCode := fmt.Sprintf("fx.Provide(New%s),", handlerName)
	invokeCode := fmt.Sprintf("fx.Invoke(Register%s),", handlerName)

	result := content
	modified := false

	marker := "// soliton-gen:event-handlers"
	if strings.Contains(result, marker) {
		insert := ""
		if !strings.Contains(result, provideCode) {
			insert += "\t" + provideCode + "\n"
		}
		if !strings.Contains(result, invokeCode) {
			insert += "\t" + invokeCode + "\n"
		}
		if insert != "" {
			result = strings.Replace(result, "\t"+marker, insert+"\t"+marker, 1)
			modified = true
		}
	} else {
		// No marker, insert before closing )
		insertPoint := strings.LastIndex(result, ")")
		if insertPoint > 0 {
			insert := ""
			if !strings.Contains(result, provideCode) {
				insert += "\t" + provideCode + "\n"
			}
			if !strings.Contains(result, invokeCode) {
				insert += "\t" + invokeCode + "\n"
			}
			if insert != "" {
				result = result[:insertPoint] + insert + result[insertPoint:]
				modified = true
			}
		}
	}

	return result, modified
}

func ensureEventBusProviderContent(content string) (string, bool) {
	result := content
	modified := false

	importPath := "github.com/soliton-go/framework/event"
	if !strings.Contains(result, importPath) {
		if strings.Contains(result, "// soliton-gen:imports") {
			result = strings.Replace(result,
				"\t// soliton-gen:imports",
				"\t\""+importPath+"\"\n\t// soliton-gen:imports",
				1)
			modified = true
		} else if strings.Contains(result, "import (") {
			result = strings.Replace(result, "import (\n", "import (\n\t\""+importPath+"\"\n", 1)
			modified = true
		}
	}

	provider := "func() event.EventBus { return event.NewLocalEventBus() },"
	legacyProvider := "event.NewLocalEventBus,"
	if strings.Contains(result, legacyProvider) && !strings.Contains(result, provider) {
		result = strings.Replace(result, legacyProvider, provider, 1)
		modified = true
	} else if !strings.Contains(result, provider) {
		if strings.Contains(result, "// soliton-gen:providers") {
			result = strings.Replace(result,
				"\t\t// soliton-gen:providers",
				"\t\t"+provider+"\n\t\t// soliton-gen:providers",
				1)
			modified = true
		} else if strings.Contains(result, "\t\tNewRouter,") {
			result = strings.Replace(result,
				"\t\tNewRouter,",
				"\t\t"+provider+"\n\t\tNewRouter,",
				1)
			modified = true
		}
	}

	return result, modified
}

func previewModuleForEventHandler(path, handlerName string) *GeneratedFile {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	updated, modified := updateModuleContentForEventHandler(string(content), handlerName)
	if !modified {
		return nil
	}
	return &GeneratedFile{
		Path:    path,
		Status:  FileStatusOverwrite,
		Content: updated,
	}
}

func previewEventBusProvider(mainGoPath string) *GeneratedFile {
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return nil
	}
	updated, modified := ensureEventBusProviderContent(string(content))
	if !modified {
		return nil
	}
	return &GeneratedFile{
		Path:    mainGoPath,
		Status:  FileStatusOverwrite,
		Content: updated,
	}
}

func updateModuleForEventHandler(path, handlerName string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	result, modified := updateModuleContentForEventHandler(string(content), handlerName)

	if !modified {
		return true
	}

	return os.WriteFile(path, []byte(result), 0644) == nil
}

func ensureEventBusProvider(mainGoPath string) bool {
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return false
	}
	result, modified := ensureEventBusProviderContent(string(content))

	if !modified {
		return true
	}

	return os.WriteFile(mainGoPath, []byte(result), 0644) == nil
}
