package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

// ServiceInfo represents information about a generated service
type ServiceInfo struct {
	Name    string                 `json:"name"`
	Remark  string                 `json:"remark,omitempty"`
	Methods []ServiceMethodSummary `json:"methods"`
	Type    string                 `json:"type"` // "domain_service" or "cross_domain_service"
}

// ServiceMethodSummary represents method info for list view.
type ServiceMethodSummary struct {
	Name   string `json:"name"`
	Remark string `json:"remark,omitempty"`
}

// ServiceMethodDetail represents detailed method information
type ServiceMethodDetail struct {
	Name      string `json:"name"`
	CamelName string `json:"camel_name"`
	Remark    string `json:"remark,omitempty"`
}

// ListServices handles GET /api/services/list
func ListServices(c *gin.Context) {
	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"services": []ServiceInfo{},
			"message":  "No project detected",
		})
		return
	}

	var services []ServiceInfo

	// Read application directory
	appDir := layout.AppDir
	if !core.IsDir(appDir) {
		c.JSON(http.StatusOK, gin.H{
			"services": services,
		})
		return
	}

	// Scan application subdirectories for service.go files
	entries, err := os.ReadDir(appDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"services": services,
		})
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		dirPath := filepath.Join(appDir, entry.Name())
		serviceFile := filepath.Join(dirPath, "service.go")

		// Check if this directory has a service.go
		if !core.IsFile(serviceFile) {
			continue
		}

		serviceName := toPascalCase(entry.Name()) + "Service"
		methods := parseServiceMethodsSummary(serviceFile)
		remark := parseServiceRemark(serviceFile)

		// Detect service type: check if corresponding domain exists
		domainDir := filepath.Join(layout.DomainDir, entry.Name())
		serviceType := "cross_domain_service"
		if core.IsDir(domainDir) {
			serviceType = "domain_service"
		}

		services = append(services, ServiceInfo{
			Name:    serviceName,
			Remark:  remark,
			Methods: methods,
			Type:    serviceType,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}

// GetServiceDetail handles GET /api/services/:name
func GetServiceDetail(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service name is required"})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no project detected"})
		return
	}

	// Convert service name to directory name (e.g., "OrderService" -> "order")
	dirName := toSnakeCase(serviceName)
	// Remove "_service" suffix if present
	if len(dirName) > 8 && dirName[len(dirName)-8:] == "_service" {
		dirName = dirName[:len(dirName)-8]
	}

	// Look for service.go in the directory
	serviceFile := filepath.Join(layout.AppDir, dirName, "service.go")

	fmt.Printf("[DEBUG] Looking for service: %s, file: %s\n", serviceName, serviceFile)

	if !core.IsFile(serviceFile) {
		fmt.Printf("[DEBUG] File not found: %s\n", serviceFile)
		c.JSON(http.StatusNotFound, gin.H{"error": "service file not found"})
		return
	}

	// Parse methods with details
	methods := parseServiceMethodsDetailed(serviceFile)
	remark := parseServiceRemark(serviceFile)

	c.JSON(http.StatusOK, gin.H{
		"name":    serviceName,
		"remark":  remark,
		"methods": methods,
	})
}

// parseServiceMethods parses a service file to extract method names
func parseServiceMethods(filePath string) []string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}
	}

	var methods []string
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip comments and empty lines
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		// Look for method definitions: func (s *ServiceName) MethodName(
		// Example: func (s *OrderService) CreateOrder(ctx context.Context, ...
		if strings.HasPrefix(trimmed, "func (") {
			// Find the closing parenthesis of the receiver
			if idx := strings.Index(trimmed, ") "); idx > 0 {
				// Get everything after the receiver
				methodPart := trimmed[idx+2:]
				// Find the method name (before the opening parenthesis)
				if openParen := strings.Index(methodPart, "("); openParen > 0 {
					methodName := strings.TrimSpace(methodPart[:openParen])
					if methodName != "" {
						methods = append(methods, methodName)
					}
				}
			}
		}
	}

	return methods
}

// parseServiceMethodsSummary parses service file and returns method summary information
func parseServiceMethodsSummary(filePath string) []ServiceMethodSummary {
	methodNames := parseServiceMethods(filePath)
	remarks := parseServiceMethodRemarks(filePath)
	var methods []ServiceMethodSummary

	for _, name := range methodNames {
		methods = append(methods, ServiceMethodSummary{
			Name:   name,
			Remark: remarks[name],
		})
	}

	return methods
}

// parseServiceMethodsDetailed parses service file and returns detailed method information
func parseServiceMethodsDetailed(filePath string) []ServiceMethodDetail {
	methodNames := parseServiceMethods(filePath)
	remarks := parseServiceMethodRemarks(filePath)
	var methods []ServiceMethodDetail

	for _, name := range methodNames {
		methods = append(methods, ServiceMethodDetail{
			Name:      name,
			CamelName: toCamelCase(name),
			Remark:    remarks[name],
		})
	}

	return methods
}

func parseServiceMethodRemarks(filePath string) map[string]string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return map[string]string{}
	}

	remarks := make(map[string]string)
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "// MethodRemark:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "// MethodRemark:"))
		name, remark := splitMethodRemark(payload)
		if name != "" {
			remarks[name] = remark
		}
	}

	return remarks
}

func splitMethodRemark(payload string) (string, string) {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return "", ""
	}
	if strings.Contains(payload, "::") {
		parts := strings.SplitN(payload, "::", 2)
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	if idx := strings.Index(payload, " "); idx >= 0 {
		return strings.TrimSpace(payload[:idx]), strings.TrimSpace(payload[idx+1:])
	}
	return strings.TrimSpace(payload), ""
}

func parseServiceRemark(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "// ServiceRemark:") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "// ServiceRemark:"))
		}
	}
	return ""
}

// toPascalCase converts snake_case to PascalCase
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// toCamelCase converts PascalCase to camelCase
func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

// DeleteService handles DELETE /api/services/:name
func DeleteService(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service name is required"})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no project detected"})
		return
	}

	// Convert service name to file name
	fileName := toSnakeCase(serviceName) + ".go"

	// Recursively search for the service file
	var serviceFile string
	err = filepath.Walk(layout.AppDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == fileName {
			serviceFile = path
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil || serviceFile == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "service file not found"})
		return
	}

	// Delete the service file
	if err := os.Remove(serviceFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Service %s deleted successfully", serviceName),
	})
}
