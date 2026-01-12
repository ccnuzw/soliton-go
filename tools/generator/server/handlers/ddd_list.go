package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

type DddListItem struct {
	Name      string `json:"name"`
	File      string `json:"file"`
	Target    string `json:"target,omitempty"`
	Topic     string `json:"topic,omitempty"`
	EventName string `json:"event_name,omitempty"`
}

type DddListResponse struct {
	ValueObjects   []DddListItem `json:"value_objects"`
	Specifications []DddListItem `json:"specs"`
	Policies       []DddListItem `json:"policies"`
	Events         []DddListItem `json:"events"`
	EventHandlers  []DddListItem `json:"event_handlers"`
}

type DddDetailResponse struct {
	Name      string             `json:"name,omitempty"`
	Fields    []core.FieldConfig `json:"fields,omitempty"`
	Target    string             `json:"target,omitempty"`
	Topic     string             `json:"topic,omitempty"`
	EventName string             `json:"event_name,omitempty"`
}

type DddDeleteRequest struct {
	Domain string `json:"domain" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type DddRenameRequest struct {
	Domain  string `json:"domain" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Name    string `json:"name" binding:"required"`
	NewName string `json:"new_name" binding:"required"`
	Force   bool   `json:"force"`
}

// ListDDD handles GET /api/ddd/list?domain=xxx
func ListDDD(c *gin.Context) {
	domain := strings.TrimSpace(c.Query("domain"))
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain is required"})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainName := strings.ToLower(domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !core.IsDir(domainDir) {
		c.JSON(http.StatusNotFound, gin.H{"error": "domain not found"})
		return
	}

	valueObjects := listByPrefix(domainDir, "value_object_")
	specs := listByPrefix(domainDir, "spec_")
	policies := listByPrefix(domainDir, "policy_")
	events := listByPrefix(domainDir, "event_")

	response := DddListResponse{
		ValueObjects:   make([]DddListItem, 0, len(valueObjects)),
		Specifications: make([]DddListItem, 0, len(specs)),
		Policies:       make([]DddListItem, 0, len(policies)),
		Events:         make([]DddListItem, 0, len(events)),
		EventHandlers:  []DddListItem{},
	}

	for _, file := range valueObjects {
		name := core.ToPascalCase(strings.TrimSuffix(strings.TrimPrefix(file, "value_object_"), ".go"))
		response.ValueObjects = append(response.ValueObjects, DddListItem{Name: name, File: file})
	}
	for _, file := range specs {
		name := core.ToPascalCase(strings.TrimSuffix(strings.TrimPrefix(file, "spec_"), ".go"))
		response.Specifications = append(response.Specifications, DddListItem{Name: name, File: file})
	}
	for _, file := range policies {
		name := core.ToPascalCase(strings.TrimSuffix(strings.TrimPrefix(file, "policy_"), ".go"))
		response.Policies = append(response.Policies, DddListItem{Name: name, File: file})
	}
	for _, file := range events {
		path := filepath.Join(domainDir, file)
		lines, err := readLines(path)
		if err != nil {
			continue
		}
		eventStructName := parseEventStructName(lines)
		eventName := normalizeEventName(eventStructName)
		topic := parseEventTopic(lines)
		response.Events = append(response.Events, DddListItem{Name: eventName, File: file, Topic: topic})
	}

	handlerFiles := listByPrefix(layout.AppDir, "event_handler_")
	for _, file := range handlerFiles {
		path := filepath.Join(layout.AppDir, file)
		lines, err := readLines(path)
		if err != nil {
			continue
		}
		eventName, topic, domainPackage := parseEventHandlerSummary(lines)
		if eventName == "" || domainPackage == "" {
			continue
		}
		if strings.ToLower(domainPackage) != domainName {
			continue
		}
		response.EventHandlers = append(response.EventHandlers, DddListItem{
			Name:      eventName,
			EventName: eventName,
			File:      file,
			Topic:     topic,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetDDDDetail handles GET /api/ddd/detail?domain=xxx&type=...&name=...
func GetDDDDetail(c *gin.Context) {
	domain := strings.TrimSpace(c.Query("domain"))
	itemType := strings.TrimSpace(c.Query("type"))
	name := strings.TrimSpace(c.Query("name"))
	if domain == "" || itemType == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain, type, and name are required"})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainName := strings.ToLower(domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !core.IsDir(domainDir) {
		c.JSON(http.StatusNotFound, gin.H{"error": "domain not found"})
		return
	}

	switch strings.ToLower(itemType) {
	case "valueobject":
		file := filepath.Join(domainDir, "value_object_"+core.ToSnakeCase(name)+".go")
		detail, err := loadStructDetail(file, core.ToPascalCase(name))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, detail)
	case "spec":
		file := filepath.Join(domainDir, "spec_"+core.ToSnakeCase(name)+".go")
		target, err := loadTargetDetail(file, "IsSatisfiedBy")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, DddDetailResponse{Name: core.ToPascalCase(name), Target: target})
	case "policy":
		file := filepath.Join(domainDir, "policy_"+core.ToSnakeCase(name)+".go")
		target, err := loadTargetDetail(file, "Validate")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, DddDetailResponse{Name: core.ToPascalCase(name), Target: target})
	case "event":
		file := filepath.Join(domainDir, "event_"+core.ToSnakeCase(normalizeEventName(name))+".go")
		lines, err := readLines(file)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		eventStructName := parseEventStructName(lines)
		if eventStructName == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "event struct not found"})
			return
		}
		fields := parseStructFields(lines, eventStructName)
		topic := parseEventTopic(lines)
		c.JSON(http.StatusOK, DddDetailResponse{
			Name:   normalizeEventName(eventStructName),
			Fields: fields,
			Topic:  topic,
		})
	case "event_handler":
		eventName, topic, err := loadEventHandlerDetail(layout.AppDir, normalizeEventName(name), domainName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, DddDetailResponse{
			EventName: eventName,
			Topic:     topic,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported type"})
	}
}

// GetDDDSource handles GET /api/ddd/source?domain=xxx&type=...&name=...
func GetDDDSource(c *gin.Context) {
	domain := strings.TrimSpace(c.Query("domain"))
	itemType := strings.TrimSpace(c.Query("type"))
	name := strings.TrimSpace(c.Query("name"))
	if domain == "" || itemType == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain, type, and name are required"})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path, err := dddFilePath(layout, domain, itemType, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file":    filepath.Base(path),
		"content": string(content),
	})
}

// DeleteDDD handles POST /api/ddd/delete
func DeleteDDD(c *gin.Context) {
	var req DddDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path, err := dddFilePath(layout, req.Domain, req.Type, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.Remove(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if strings.ToLower(req.Type) == "event_handler" {
		handlerName := strings.TrimSuffix(normalizeEventStructName(core.ToPascalCase(req.Name)), "Event") + "Handler"
		modulePath := filepath.Join(layout.AppDir, strings.ToLower(req.Domain), "module.go")
		_ = removeEventHandlerFromModule(modulePath, handlerName)
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// RenameDDD handles POST /api/ddd/rename
func RenameDDD(c *gin.Context) {
	var req DddRenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldPath, err := dddFilePath(layout, req.Domain, req.Type, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newPath, err := dddFilePath(layout, req.Domain, req.Type, req.NewName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if oldPath == newPath {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new name is the same as old name"})
		return
	}
	if core.IsFile(newPath) && !req.Force {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target file exists"})
		return
	}

	content, err := os.ReadFile(oldPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	updated := string(content)
	oldName := core.ToPascalCase(req.Name)
	newName := core.ToPascalCase(req.NewName)

	switch strings.ToLower(req.Type) {
	case "valueobject", "spec", "policy":
		updated = strings.ReplaceAll(updated, oldName, newName)
	case "event":
		oldStruct := normalizeEventStructName(oldName)
		newStruct := normalizeEventStructName(newName)
		updated = strings.ReplaceAll(updated, oldStruct, newStruct)
	case "event_handler":
		oldStruct := normalizeEventStructName(oldName)
		newStruct := normalizeEventStructName(newName)
		oldHandler := strings.TrimSuffix(oldStruct, "Event") + "Handler"
		newHandler := strings.TrimSuffix(newStruct, "Event") + "Handler"
		updated = strings.ReplaceAll(updated, oldHandler, newHandler)
		updated = strings.ReplaceAll(updated, oldStruct, newStruct)
		modulePath := filepath.Join(layout.AppDir, strings.ToLower(req.Domain), "module.go")
		_ = renameEventHandlerInModule(modulePath, oldHandler, newHandler)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported type"})
		return
	}

	if err := os.WriteFile(newPath, []byte(updated), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := os.Remove(oldPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func listByPrefix(dir, prefix string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var result []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasPrefix(name, prefix) {
			result = append(result, name)
		}
	}
	return result
}

func dddFilePath(layout core.ProjectLayout, domain, itemType, name string) (string, error) {
	domainName := strings.ToLower(domain)
	domainDir := filepath.Join(layout.DomainDir, domainName)
	appDir := filepath.Join(layout.AppDir, domainName)
	if !core.IsDir(domainDir) {
		return "", fmt.Errorf("domain not found")
	}

	baseName := core.ToSnakeCase(core.ToPascalCase(name))
	switch strings.ToLower(itemType) {
	case "valueobject":
		return filepath.Join(domainDir, "value_object_"+baseName+".go"), nil
	case "spec":
		return filepath.Join(domainDir, "spec_"+baseName+".go"), nil
	case "policy":
		return filepath.Join(domainDir, "policy_"+baseName+".go"), nil
	case "event":
		eventStruct := normalizeEventStructName(core.ToPascalCase(name))
		snake := core.ToSnakeCase(strings.TrimSuffix(eventStruct, "Event"))
		return filepath.Join(domainDir, "event_"+snake+".go"), nil
	case "event_handler":
		eventStruct := normalizeEventStructName(core.ToPascalCase(name))
		snake := core.ToSnakeCase(strings.TrimSuffix(eventStruct, "Event"))
		return filepath.Join(appDir, "event_handler_"+snake+".go"), nil
	default:
		return "", fmt.Errorf("unsupported type")
	}
}

func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func loadStructDetail(path string, structName string) (DddDetailResponse, error) {
	lines, err := readLines(path)
	if err != nil {
		return DddDetailResponse{}, err
	}
	fields := parseStructFields(lines, structName)
	if len(fields) == 0 {
		return DddDetailResponse{}, fmt.Errorf("no fields found")
	}
	return DddDetailResponse{Name: structName, Fields: fields}, nil
}

func loadTargetDetail(path string, methodName string) (string, error) {
	lines, err := readLines(path)
	if err != nil {
		return "", err
	}
	target := parseTargetType(lines, methodName)
	if target == "" {
		return "", nil
	}
	return target, nil
}

func loadEventHandlerDetail(appDir, eventName, domainName string) (string, string, error) {
	files := listByPrefix(appDir, "event_handler_")
	for _, file := range files {
		lines, err := readLines(filepath.Join(appDir, file))
		if err != nil {
			continue
		}
		parsedEvent, topic, domainPackage := parseEventHandlerSummary(lines)
		if strings.ToLower(domainPackage) != strings.ToLower(domainName) {
			continue
		}
		if normalizeEventName(parsedEvent) == normalizeEventName(eventName) {
			return normalizeEventName(parsedEvent), topic, nil
		}
	}
	return "", "", fmt.Errorf("event handler not found")
}

func parseStructFields(lines []string, structName string) []core.FieldConfig {
	enumValues := parseEnumValuesFromLines(lines)
	inStruct := false
	var fields []core.FieldConfig
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "type "+structName+" struct") {
			inStruct = true
			continue
		}
		if !inStruct {
			continue
		}
		if strings.HasPrefix(trimmed, "}") {
			break
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		comment := ""
		if idx := strings.Index(trimmed, "//"); idx >= 0 {
			comment = strings.TrimSpace(trimmed[idx+2:])
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
		if trimmed == "" {
			continue
		}
		if strings.Contains(trimmed, "`") {
			trimmed = strings.TrimSpace(strings.Split(trimmed, "`")[0])
		}
		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}
		fieldName := parts[0]
		fieldType := parts[1]

		fieldConfig := core.FieldConfig{
			Name:    core.ToSnakeCase(fieldName),
			Type:    mapGoType(fieldType, enumValues),
			Comment: comment,
		}
		if fieldConfig.Type == "enum" {
			fieldConfig.EnumValues = enumValues[fieldType]
		}
		fields = append(fields, fieldConfig)
	}
	return fields
}

func parseEnumValuesFromLines(lines []string) map[string][]string {
	enumValues := map[string][]string{}
	inConst := false
	lastType := ""
	reValue := regexp.MustCompile(`"([^"]+)"`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "const (") {
			inConst = true
			continue
		}
		if inConst && strings.HasPrefix(trimmed, ")") {
			inConst = false
			lastType = ""
			continue
		}
		if !inConst || trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		trimmed = strings.SplitN(trimmed, "//", 2)[0]
		if !strings.Contains(trimmed, "\"") {
			continue
		}
		value := ""
		if match := reValue.FindStringSubmatch(trimmed); len(match) > 1 {
			value = match[1]
		}
		parts := strings.Fields(trimmed)
		if len(parts) >= 3 {
			lastType = parts[1]
		}
		if lastType == "" {
			continue
		}
		enumValues[lastType] = append(enumValues[lastType], value)
	}

	return enumValues
}

func mapGoType(goType string, enumValues map[string][]string) string {
	if _, ok := enumValues[goType]; ok {
		return "enum"
	}
	original := goType
	if strings.HasPrefix(goType, "*") {
		goType = strings.TrimPrefix(goType, "*")
	}
	switch goType {
	case "string":
		return "string"
	case "int":
		return "int"
	case "int64":
		return "int64"
	case "float64":
		return "float64"
	case "bool":
		return "bool"
	case "time.Time":
		if strings.HasPrefix(original, "*") {
			return "time?"
		}
		return "time"
	case "datatypes.JSON":
		return "json"
	case "[]byte":
		return "bytes"
	default:
		return "string"
	}
}

func parseTargetType(lines []string, methodName string) string {
	re := regexp.MustCompile(`target\s+([^\),]+)`)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, methodName+"(") {
			continue
		}
		match := re.FindStringSubmatch(trimmed)
		if len(match) < 2 {
			continue
		}
		target := strings.TrimSpace(match[1])
		target = strings.TrimPrefix(target, "*")
		if target == "any" {
			return ""
		}
		if strings.Contains(target, ".") {
			parts := strings.Split(target, ".")
			target = parts[len(parts)-1]
		}
		return target
	}
	return ""
}

func parseEventStructName(lines []string) string {
	inStruct := false
	structName := ""
	hasBase := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "type ") && strings.Contains(trimmed, " struct") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				structName = parts[1]
				inStruct = true
				hasBase = false
			}
			continue
		}
		if inStruct && strings.HasPrefix(trimmed, "}") {
			if hasBase {
				return structName
			}
			inStruct = false
			continue
		}
		if inStruct && strings.Contains(trimmed, "ddd.BaseDomainEvent") {
			hasBase = true
		}
	}
	return ""
}

func normalizeEventStructName(name string) string {
	if strings.HasSuffix(name, "Event") {
		return name
	}
	return name + "Event"
}

func parseEventTopic(lines []string) string {
	inFunc := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "EventName() string") {
			inFunc = true
			continue
		}
		if !inFunc {
			continue
		}
		if strings.HasPrefix(trimmed, "return ") {
			if match := regexp.MustCompile(`"([^"]+)"`).FindStringSubmatch(trimmed); len(match) > 1 {
				return match[1]
			}
		}
		if strings.HasPrefix(trimmed, "}") {
			inFunc = false
		}
	}
	return ""
}

func parseEventHandlerSummary(lines []string) (string, string, string) {
	eventName := ""
	topic := ""
	domainPackage := ""
	reTopic := regexp.MustCompile(`"([^"]+)"`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "evt.(*") {
			start := strings.Index(trimmed, "evt.(*")
			if start >= 0 {
				after := trimmed[start+len("evt.(*"):]
				end := strings.Index(after, ")")
				if end > 0 {
					typeName := after[:end]
					if dot := strings.LastIndex(typeName, "."); dot >= 0 {
						domainPackage = typeName[:dot]
						eventName = strings.TrimSuffix(typeName[dot+1:], "Event")
					} else {
						eventName = strings.TrimSuffix(typeName, "Event")
					}
				}
			}
		}
		if strings.Contains(trimmed, "Subscribe(") {
			if match := reTopic.FindStringSubmatch(trimmed); len(match) > 1 {
				topic = match[1]
			}
		}
	}
	return eventName, topic, domainPackage
}

func normalizeEventName(name string) string {
	if strings.HasSuffix(name, "Event") {
		return strings.TrimSuffix(name, "Event")
	}
	return name
}

func removeEventHandlerFromModule(path, handlerName string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	lines := strings.Split(string(content), "\n")
	updated := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, "New"+handlerName) || strings.Contains(line, "Register"+handlerName) {
			continue
		}
		updated = append(updated, line)
	}
	return os.WriteFile(path, []byte(strings.Join(updated, "\n")), 0644) == nil
}

func renameEventHandlerInModule(path, oldHandler, newHandler string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	result := strings.ReplaceAll(string(content), "New"+oldHandler, "New"+newHandler)
	result = strings.ReplaceAll(result, "Register"+oldHandler, "Register"+newHandler)
	return os.WriteFile(path, []byte(result), 0644) == nil
}
