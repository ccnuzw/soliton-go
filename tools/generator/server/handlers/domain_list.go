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

// DomainInfo represents information about a generated domain
type DomainInfo struct {
	Name       string   `json:"name"`
	ModulePath string   `json:"module_path"`
	TableName  string   `json:"table_name"`
	Fields     []string `json:"fields"`
	HasFiles   bool     `json:"has_files"`
}

// ListDomains handles GET /api/domains/list
func ListDomains(c *gin.Context) {
	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"domains": []DomainInfo{},
			"message": "No project detected",
		})
		return
	}

	var domains []DomainInfo

	// Read domain directory
	domainDir := layout.DomainDir
	if !core.IsDir(domainDir) {
		c.JSON(http.StatusOK, gin.H{
			"domains": domains,
		})
		return
	}

	entries, err := os.ReadDir(domainDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read domain directory"})
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		domainName := entry.Name()
		domainPath := filepath.Join(domainDir, domainName)

		// Check if domain files exist
		entityFile := filepath.Join(domainPath, domainName+".go")
		repoFile := filepath.Join(domainPath, "repository.go")

		if !core.IsFile(entityFile) {
			continue
		}

		// Parse entity file to get fields
		fields := parseEntityFields(entityFile)

		domains = append(domains, DomainInfo{
			Name:       domainName,
			ModulePath: layout.ModulePath,
			Fields:     fields,
			HasFiles:   core.IsFile(repoFile),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"domains": domains,
	})
}

// parseEntityFields parses a domain entity file to extract field names
func parseEntityFields(filePath string) []string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}
	}

	var fields []string
	lines := strings.Split(string(content), "\n")
	inStruct := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Find struct definition (e.g., "type User struct {")
		if strings.HasPrefix(line, "type ") && strings.Contains(line, " struct {") {
			inStruct = true
			continue
		}

		// End of struct
		if inStruct && line == "}" {
			break
		}

		// Parse field line
		if inStruct && line != "" && !strings.HasPrefix(line, "//") {
			// Skip embedded structs and base fields
			if strings.Contains(line, "ddd.AggregateRoot") ||
				strings.Contains(line, "ddd.Entity") ||
				strings.Contains(line, "ID") && strings.Contains(line, "uuid.UUID") {
				continue
			}

			// Extract field name (first word)
			parts := strings.Fields(line)
			if len(parts) > 0 {
				fieldName := parts[0]
				// Skip if it's a comment or special character
				if !strings.HasPrefix(fieldName, "//") && fieldName != "{" && fieldName != "}" {
					fields = append(fields, fieldName)
				}
			}
		}
	}

	return fields
}

// GetDomainDetail handles GET /api/domains/:name
func GetDomainDetail(c *gin.Context) {
	domainName := c.Param("name")
	if domainName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain name is required"})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no project detected"})
		return
	}

	domainPath := filepath.Join(layout.DomainDir, domainName)
	if !core.IsDir(domainPath) {
		c.JSON(http.StatusNotFound, gin.H{"error": "domain not found"})
		return
	}

	// Read entity file
	entityFile := filepath.Join(domainPath, domainName+".go")
	if !core.IsFile(entityFile) {
		c.JSON(http.StatusNotFound, gin.H{"error": "entity file not found"})
		return
	}

	// Parse fields with more details
	fields := parseEntityFieldsDetailed(entityFile)

	c.JSON(http.StatusOK, gin.H{
		"name":   domainName,
		"fields": fields,
		"files": gin.H{
			"entity":     core.IsFile(entityFile),
			"repository": core.IsFile(filepath.Join(domainPath, "repository.go")),
			"events":     core.IsFile(filepath.Join(domainPath, "events.go")),
		},
	})
}

// FieldDetail represents detailed field information
type FieldDetail struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsEnum    bool   `json:"is_enum"`
	GormTag   string `json:"gorm_tag"`
	JsonTag   string `json:"json_tag"`
	SnakeName string `json:"snake_name"`
}

// parseEntityFieldsDetailed parses entity file and returns detailed field information
func parseEntityFieldsDetailed(filePath string) []FieldDetail {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []FieldDetail{}
	}

	var fields []FieldDetail
	lines := strings.Split(string(content), "\n")
	inStruct := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Find struct definition
		if strings.HasPrefix(trimmed, "type ") && strings.Contains(trimmed, " struct {") {
			inStruct = true
			continue
		}

		// End of struct
		if inStruct && trimmed == "}" {
			break
		}

		// Parse field line
		if inStruct && trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			// Skip base fields
			if strings.Contains(trimmed, "ddd.AggregateRoot") ||
				strings.Contains(trimmed, "ddd.Entity") ||
				(strings.Contains(trimmed, "ID") && strings.Contains(trimmed, "uuid.UUID")) {
				continue
			}

			// Parse field: Name Type `tags`
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				fieldName := parts[0]
				fieldType := parts[1]

				// Extract tags
				gormTag := ""
				jsonTag := ""
				if len(parts) > 2 {
					tags := strings.Join(parts[2:], " ")
					if strings.Contains(tags, "gorm:") {
						gormTag = extractTag(tags, "gorm")
					}
					if strings.Contains(tags, "json:") {
						jsonTag = extractTag(tags, "json")
					}
				}

				// Determine snake_name from json tag or convert from field name
				snakeName := ""
				if jsonTag != "" {
					snakeName = strings.Trim(jsonTag, `"`)
				} else {
					snakeName = toSnakeCase(fieldName)
				}

				fields = append(fields, FieldDetail{
					Name:      fieldName,
					Type:      fieldType,
					IsEnum:    !strings.Contains(fieldType, "string") && !strings.Contains(fieldType, "int") && !strings.Contains(fieldType, "time"),
					GormTag:   gormTag,
					JsonTag:   jsonTag,
					SnakeName: snakeName,
				})
			}
		}
	}

	return fields
}

// extractTag extracts a specific tag value from a tag string
func extractTag(tags, tagName string) string {
	start := strings.Index(tags, tagName+":")
	if start == -1 {
		return ""
	}
	start += len(tagName) + 1

	// Find the closing quote
	end := strings.Index(tags[start:], `"`)
	if end == -1 {
		return ""
	}

	return tags[start : start+end]
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// DeleteDomain handles DELETE /api/domains/:name
func DeleteDomain(c *gin.Context) {
	domainName := c.Param("name")
	if domainName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain name is required"})
		return
	}

	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no project detected"})
		return
	}

	// Domain directory path
	domainDir := filepath.Join(layout.DomainDir, domainName)

	if !core.IsDir(domainDir) {
		c.JSON(http.StatusNotFound, gin.H{"error": "domain not found"})
		return
	}

	// Delete the entire domain directory
	if err := os.RemoveAll(domainDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Domain %s deleted successfully", domainName),
	})
}
