package core

// ============================================================================
// Project Configuration Types
// ============================================================================

// ProjectConfig holds configuration for project initialization.
type ProjectConfig struct {
	Name             string `json:"name"`
	ModuleName       string `json:"module_name"`
	FrameworkVersion string `json:"framework_version,omitempty"`
	FrameworkReplace string `json:"framework_replace,omitempty"`
}

// ============================================================================
// Domain Configuration Types
// ============================================================================

// DomainConfig holds configuration for domain generation.
type DomainConfig struct {
	Name       string        `json:"name"`
	Fields     []FieldConfig `json:"fields"`
	TableName  string        `json:"table_name,omitempty"`
	RouteBase  string        `json:"route_base,omitempty"`
	SoftDelete bool          `json:"soft_delete"`
	Wire       bool          `json:"wire"`
	Force      bool          `json:"force"`
}

// FieldConfig holds configuration for a single field.
type FieldConfig struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`                  // string, int, int64, text, uuid, time, time?, enum
	Comment    string   `json:"comment,omitempty"`     // Field comment/description
	EnumValues []string `json:"enum_values,omitempty"` // For enum type only
}

// ============================================================================
// Service Configuration Types
// ============================================================================

// ServiceConfig holds configuration for service generation.
type ServiceConfig struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods,omitempty"`
	Force   bool     `json:"force"`
}

// ValueObjectConfig holds configuration for value object generation.
type ValueObjectConfig struct {
	Domain string        `json:"domain"`
	Name   string        `json:"name"`
	Fields []FieldConfig `json:"fields,omitempty"`
	Force  bool          `json:"force"`
}

// SpecificationConfig holds configuration for specification generation.
type SpecificationConfig struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Target string `json:"target,omitempty"`
	Force  bool   `json:"force"`
}

// PolicyConfig holds configuration for policy generation.
type PolicyConfig struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Target string `json:"target,omitempty"`
	Force  bool   `json:"force"`
}

// EventConfig holds configuration for domain event generation.
type EventConfig struct {
	Domain string        `json:"domain"`
	Name   string        `json:"name"`
	Topic  string        `json:"topic,omitempty"`
	Fields []FieldConfig `json:"fields,omitempty"`
	Force  bool          `json:"force"`
}

// EventHandlerConfig holds configuration for event handler generation.
type EventHandlerConfig struct {
	Domain    string `json:"domain"`
	EventName string `json:"event_name"`
	Topic     string `json:"topic,omitempty"`
	Force     bool   `json:"force"`
}

// ============================================================================
// Generation Result Types
// ============================================================================

// GenerationResult holds the result of a generation operation.
type GenerationResult struct {
	Success bool            `json:"success"`
	Files   []GeneratedFile `json:"files"`
	Errors  []string        `json:"errors,omitempty"`
	Message string          `json:"message,omitempty"`
}

// GeneratedFile represents information about a generated file.
type GeneratedFile struct {
	Path    string `json:"path"`
	Status  string `json:"status"`            // new, overwrite, skip, error
	Content string `json:"content,omitempty"` // Populated only for preview
}

// FileStatus constants
const (
	FileStatusNew       = "new"
	FileStatusOverwrite = "overwrite"
	FileStatusSkip      = "skip"
	FileStatusError     = "error"
)

// ============================================================================
// Internal Types (for template rendering)
// ============================================================================

// Field represents a parsed field definition for template use.
type Field struct {
	Name       string   // Field name (e.g., "Username")
	SnakeName  string   // Snake case name (e.g., "username")
	CamelName  string   // Camel case name (e.g., "username")
	GoType     string   // Go type in domain package (e.g., "UserRole")
	AppGoType  string   // Go type in app layer with package prefix (e.g., "user.UserRole")
	GormTag    string   // GORM tag
	JsonTag    string   // JSON tag
	Comment    string   // Field comment/description
	IsEnum     bool     // Is this an enum type
	EnumValues []string // Enum values if IsEnum is true
	EnumType   string   // Enum type name (e.g., "UserStatus")
	IsPointer  bool     // True if GoType is a pointer type
}

// TemplateData holds all data for domain template generation.
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

// ServiceMethod represents a service method for template use.
type ServiceMethod struct {
	Name      string // Method name (e.g., "CreateOrder")
	CamelName string // Camel case name (e.g., "createOrder")
}

// ServiceData holds template data for service generation.
type ServiceData struct {
	ServiceName string
	PackageName string
	Methods     []ServiceMethod
	ModulePath  string
}

// ProjectData holds template data for project initialization.
type ProjectData struct {
	ProjectName      string
	ModuleName       string
	FrameworkVersion string
	FrameworkReplace string
	GoVersion        string
	FxVersion        string
}

// ValueObjectData holds template data for value object generation.
type ValueObjectData struct {
	PackageName     string
	ValueObjectName string
	Fields          []Field
	HasTime         bool
	HasEnums        bool
}

// SpecificationData holds template data for specification generation.
type SpecificationData struct {
	PackageName       string
	SpecificationName string
	TargetType        string
	TargetIsAny       bool
}

// PolicyData holds template data for policy generation.
type PolicyData struct {
	PackageName string
	PolicyName  string
	TargetType  string
	TargetIsAny bool
}

// EventData holds template data for domain event generation.
type EventData struct {
	PackageName     string
	EventStructName string
	EventTopic      string
	Fields          []Field
	HasTime         bool
	HasEnums        bool
}

// EventHandlerData holds template data for event handler generation.
type EventHandlerData struct {
	PackageName     string
	DomainPackage   string
	EventStructName string
	HandlerName     string
	EventTopic      string
	ModulePath      string
}
