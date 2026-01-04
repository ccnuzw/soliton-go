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
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
}

// ServiceMethodDetail represents detailed method information
type ServiceMethodDetail struct {
	Name      string `json:"name"`
	CamelName string `json:"camel_name"`
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

	// Recursively scan for *_service.go files
	err = filepath.Walk(appDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Look for *_service.go files
		fileName := info.Name()
		if !strings.HasSuffix(fileName, "_service.go") {
			return nil
		}

		// Extract service name (e.g., "order_service.go" -> "OrderService")
		// Keep the full name including "Service" suffix
		baseName := strings.TrimSuffix(fileName, ".go")
		serviceName := toPascalCase(baseName)

		methods := parseServiceMethods(path)

		services = append(services, ServiceInfo{
			Name:    serviceName,
			Methods: methods,
		})

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan application directory"})
		return
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

	// Convert service name to file name (e.g., "OrderService" -> "order_service.go")
	fileName := toSnakeCase(serviceName) + ".go"

	fmt.Printf("[DEBUG] Looking for service: %s, fileName: %s\n", serviceName, fileName)

	// Recursively search for the service file
	var serviceFile string
	err = filepath.Walk(layout.AppDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == fileName {
			fmt.Printf("[DEBUG] Found file: %s\n", path)
			serviceFile = path
			return filepath.SkipAll // Found it, stop walking
		}
		return nil
	})

	if err != nil || serviceFile == "" {
		fmt.Printf("[DEBUG] File not found. Searched for: %s in %s\n", fileName, layout.AppDir)
		c.JSON(http.StatusNotFound, gin.H{"error": "service file not found"})
		return
	}

	// Parse methods with details
	methods := parseServiceMethodsDetailed(serviceFile)

	c.JSON(http.StatusOK, gin.H{
		"name":    serviceName,
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

// parseServiceMethodsDetailed parses service file and returns detailed method information
func parseServiceMethodsDetailed(filePath string) []ServiceMethodDetail {
	methodNames := parseServiceMethods(filePath)
	var methods []ServiceMethodDetail

	for _, name := range methodNames {
		methods = append(methods, ServiceMethodDetail{
			Name:      name,
			CamelName: toCamelCase(name),
		})
	}

	return methods
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
