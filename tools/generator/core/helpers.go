package core

import (
	"fmt"
	"strings"
)

// ============================================================================
// String Conversion Helpers
// ============================================================================

// ToSnakeCase converts PascalCase or camelCase to snake_case.
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ToPascalCase converts snake_case to PascalCase.
func ToPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(string(p[0])) + strings.ToLower(p[1:])
		}
	}
	return strings.Join(parts, "")
}

// ToCamelCase converts snake_case to camelCase.
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(string(pascal[0])) + pascal[1:]
}

// Pluralize returns the plural form of a word.
func Pluralize(name string) string {
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") ||
		strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return lower + "es"
	}
	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		prev := lower[len(lower)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return lower[:len(lower)-1] + "ies"
		}
	}
	return lower + "s"
}

// EnumConst converts an enum value to a valid Go constant name.
func EnumConst(value string) string {
	var cleaned strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			cleaned.WriteRune(r)
		} else {
			cleaned.WriteByte('_')
		}
	}

	normalized := strings.Trim(cleaned.String(), "_")
	if normalized == "" {
		return "Value"
	}

	normalized = strings.ToLower(normalized)
	result := ToPascalCase(normalized)
	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		return "Value" + result
	}
	return result
}

// ============================================================================
// Field Parsing Helpers
// ============================================================================

// ParseFields parses a field string into FieldConfig slice.
func ParseFields(fieldsStr string) []FieldConfig {
	if fieldsStr == "" {
		return nil
	}

	var fields []FieldConfig
	parts := strings.Split(fieldsStr, ",")
	seen := make(map[string]struct{})

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		field := parseFieldDefinition(part)
		snakeName := ToSnakeCase(field.Name)
		if isReservedField(snakeName) {
			continue
		}
		if _, exists := seen[snakeName]; exists {
			continue
		}
		seen[snakeName] = struct{}{}
		fields = append(fields, field)
	}

	return fields
}

// parseFieldDefinition parses a single field definition.
func parseFieldDefinition(def string) FieldConfig {
	colonIdx := strings.Index(def, ":")

	var fieldName, fieldType string
	if colonIdx == -1 {
		fieldName = def
		fieldType = "string"
	} else {
		fieldName = def[:colonIdx]
		fieldType = def[colonIdx+1:]
	}

	field := FieldConfig{
		Name: ToPascalCase(fieldName),
		Type: fieldType,
	}

	// Parse enum values if present
	if strings.HasPrefix(fieldType, "enum(") && strings.HasSuffix(fieldType, ")") {
		enumContent := fieldType[5 : len(fieldType)-1]
		field.Type = "enum"
		field.EnumValues = parseEnumValues(enumContent)
	}

	return field
}

// parseEnumValues parses enum values from a pipe-separated string.
func parseEnumValues(raw string) []string {
	parts := strings.Split(raw, "|")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}

// isReservedField checks if a field name is reserved.
func isReservedField(snakeName string) bool {
	switch snakeName {
	case "id", "created_at", "updated_at", "deleted_at":
		return true
	default:
		return false
	}
}

// ConvertToFields converts FieldConfig slice to Field slice for templates.
func ConvertToFields(configs []FieldConfig, entityName, packageName string) []Field {
	if len(configs) == 0 {
		// Default field if none specified
		return []Field{
			{
				Name:      "Name",
				SnakeName: "name",
				CamelName: "name",
				GoType:    "string",
				AppGoType: "string",
				GormTag:   "`gorm:\"size:255\"`",
				JsonTag:   "`json:\"name\"`",
			},
		}
	}

	// Built-in fields that should be skipped (already defined in template)
	builtinFields := map[string]bool{
		"id":        true,
		"createdat": true,
		"updatedat": true,
		"deletedat": true,
	}

	fields := make([]Field, 0, len(configs))
	for _, cfg := range configs {
		// Skip built-in fields to avoid duplication
		fieldNameLower := strings.ToLower(cfg.Name)
		if builtinFields[fieldNameLower] {
			continue
		}
		field := convertFieldConfig(cfg, entityName, packageName)
		fields = append(fields, field)
	}
	return fields
}

// convertFieldConfig converts a single FieldConfig to Field.
func convertFieldConfig(cfg FieldConfig, entityName, packageName string) Field {
	snakeName := ToSnakeCase(cfg.Name)
	pascalName := ToPascalCase(cfg.Name)
	camelName := ToCamelCase(cfg.Name)

	field := Field{
		Name:      pascalName,
		SnakeName: snakeName,
		CamelName: camelName,
		Comment:   cfg.Comment,
	}

	if cfg.Type == "enum" || (len(cfg.EnumValues) > 0) {
		field.IsEnum = true
		field.EnumValues = cfg.EnumValues
		if len(field.EnumValues) == 0 {
			field.EnumValues = []string{"default"}
		}
		field.EnumType = entityName + pascalName
		field.GoType = field.EnumType
		field.AppGoType = packageName + "." + field.EnumType
		field.GormTag = fmt.Sprintf("`gorm:\"size:50;default:'%s'\"`", field.EnumValues[0])
		field.JsonTag = fmt.Sprintf("`json:\"%s\"`", snakeName)
	} else {
		field.GoType, field.GormTag = mapFieldType(cfg.Type, snakeName)
		field.AppGoType = field.GoType
		field.IsPointer = strings.HasPrefix(field.GoType, "*")
		field.JsonTag = fmt.Sprintf("`json:\"%s\"`", snakeName)
	}

	return field
}

// mapFieldType maps type shorthand to Go type and GORM tag.
func mapFieldType(typeName string, snakeName string) (goType, gormTag string) {
	switch strings.ToLower(typeName) {
	case "string", "str", "":
		return "string", "`gorm:\"size:255\"`"
	case "text":
		return "string", "`gorm:\"type:text\"`"
	case "int", "integer":
		return "int", "`gorm:\"not null;default:0\"`"
	case "int64", "long":
		return "int64", "`gorm:\"not null;default:0\"`"
	case "float", "float64", "double":
		return "float64", "`gorm:\"default:0\"`"
	case "decimal":
		return "float64", "`gorm:\"type:decimal(10,2);default:0\"`"
	case "bool", "boolean":
		return "bool", "`gorm:\"default:false\"`"
	case "time", "datetime", "timestamp":
		return "time.Time", "`gorm:\"type:timestamp\"`"
	case "time?", "datetime?":
		return "*time.Time", ""
	case "date":
		return "time.Time", "`gorm:\"type:date\"`"
	case "date?":
		return "*time.Time", "`gorm:\"type:date\"`"
	case "uuid", "id":
		return "string", "`gorm:\"size:36;index\"`"
	case "json":
		return "datatypes.JSON", ""
	case "jsonb":
		return "datatypes.JSON", "`gorm:\"type:jsonb\"`"
	case "bytes", "binary", "blob":
		return "[]byte", "`gorm:\"type:bytes\"`"
	default:
		return "string", "`gorm:\"size:255\"`"
	}
}

// ============================================================================
// Template Helpers (for template.FuncMap)
// ============================================================================

// CreateBindingTag creates a binding tag for create requests.
func CreateBindingTag(field Field) string {
	if field.IsPointer {
		return ""
	}
	if field.IsEnum {
		return OneOfTag("required", field.EnumValues)
	}
	if isStringType(field.GoType) {
		return " binding:\"required\""
	}
	return ""
}

// UpdateBindingTag creates a binding tag for update requests.
func UpdateBindingTag(field Field) string {
	if !field.IsEnum {
		return ""
	}
	return OneOfTag("omitempty", field.EnumValues)
}

// OneOfTag creates a oneof binding tag.
func OneOfTag(prefix string, values []string) string {
	if len(values) == 0 {
		if prefix == "required" {
			return " binding:\"required\""
		}
		return ""
	}
	return fmt.Sprintf(" binding:\"%s,oneof=%s\"", prefix, strings.Join(values, " "))
}

func isStringType(goType string) bool {
	return strings.TrimPrefix(goType, "*") == "string"
}

// ============================================================================
// Service Helpers
// ============================================================================

// ParseServiceMethods parses service methods from a comma-separated string.
func ParseServiceMethods(methodsStr, serviceName string) []ServiceMethod {
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

	return methods
}
