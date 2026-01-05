package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// DeleteResult holds the result of a delete operation.
type DeleteResult struct {
	Success      bool     `json:"success"`
	DeletedItems []string `json:"deleted_items"`
	Errors       []string `json:"errors,omitempty"`
	Message      string   `json:"message"`
}

// DeleteDomain deletes a domain and all related files.
func DeleteDomain(domainName string) DeleteResult {
	layout, err := ResolveProjectLayout()
	if err != nil {
		return DeleteResult{
			Success: false,
			Errors:  []string{"no project detected"},
			Message: "未检测到项目",
		}
	}

	domainDir := filepath.Join(layout.DomainDir, domainName)
	if !IsDir(domainDir) {
		return DeleteResult{
			Success: false,
			Errors:  []string{"domain not found"},
			Message: fmt.Sprintf("领域 %s 不存在", domainName),
		}
	}

	var deletedItems []string
	var errors []string

	// 1. Delete domain layer directory (internal/domain/<name>/)
	if err := os.RemoveAll(domainDir); err != nil {
		errors = append(errors, fmt.Sprintf("domain dir: %v", err))
	} else {
		deletedItems = append(deletedItems, "domain/"+domainName)
	}

	// 2. Delete application layer directory (internal/application/<name>/)
	appDir := filepath.Join(layout.AppDir, domainName)
	if IsDir(appDir) {
		if err := os.RemoveAll(appDir); err != nil {
			errors = append(errors, fmt.Sprintf("application dir: %v", err))
		} else {
			deletedItems = append(deletedItems, "application/"+domainName)
		}
	}

	// 3. Delete infrastructure persistence file (internal/infrastructure/persistence/<name>_repo.go)
	repoFile := filepath.Join(layout.InfraDir, domainName+"_repo.go")
	if IsFile(repoFile) {
		if err := os.Remove(repoFile); err != nil {
			errors = append(errors, fmt.Sprintf("repo file: %v", err))
		} else {
			deletedItems = append(deletedItems, "persistence/"+domainName+"_repo.go")
		}
	}

	// 4. Delete interfaces HTTP handler file (internal/interfaces/http/<name>_handler.go)
	handlerFile := filepath.Join(layout.InterfacesDir, domainName+"_handler.go")
	if IsFile(handlerFile) {
		if err := os.Remove(handlerFile); err != nil {
			errors = append(errors, fmt.Sprintf("handler file: %v", err))
		} else {
			deletedItems = append(deletedItems, "http/"+domainName+"_handler.go")
		}
	}

	// 5. Remove injection from main.go
	mainGoPath := filepath.Join(filepath.Dir(layout.InternalDir), "cmd", "main.go")
	if IsFile(mainGoPath) {
		if UnwireMainGo(mainGoPath, domainName) {
			deletedItems = append(deletedItems, "main.go injection")
		}
	}

	if len(errors) > 0 {
		return DeleteResult{
			Success:      false,
			DeletedItems: deletedItems,
			Errors:       errors,
			Message:      fmt.Sprintf("领域 %s 删除部分失败", domainName),
		}
	}

	return DeleteResult{
		Success:      true,
		DeletedItems: deletedItems,
		Message:      fmt.Sprintf("领域 %s 删除成功", domainName),
	}
}

// UnwireMainGo removes domain module injection from main.go
func UnwireMainGo(mainGoPath, domainName string) bool {
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return false
	}

	original := string(content)
	modified := original

	// Remove import line for the domain module
	// Pattern: "module-path/internal/application/<domainName>"
	importPattern := fmt.Sprintf(`\n\t"[^"]+/internal/application/%s"`, domainName)
	re := regexp.MustCompile(importPattern)
	modified = re.ReplaceAllString(modified, "")

	// Remove module.Module from fx.Options
	// Pattern: <domainName>.Module,
	modulePattern := fmt.Sprintf(`\n?\s*%s\.Module,?`, domainName)
	re = regexp.MustCompile(modulePattern)
	modified = re.ReplaceAllString(modified, "")

	// Clean up empty lines and trailing commas
	modified = regexp.MustCompile(`\n\n\n+`).ReplaceAllString(modified, "\n\n")

	if modified == original {
		return false
	}

	if err := os.WriteFile(mainGoPath, []byte(modified), 0644); err != nil {
		return false
	}

	return true
}

// ListDomains returns a list of all domains in the project.
func ListDomains() ([]string, error) {
	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(layout.DomainDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read domain directory: %w", err)
	}

	var domains []string
	for _, entry := range entries {
		if entry.IsDir() && !startsWith(entry.Name(), ".") {
			domains = append(domains, entry.Name())
		}
	}

	return domains, nil
}

// ListServices returns a list of all application services in the project.
func ListServices() ([]string, error) {
	layout, err := ResolveProjectLayout()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(layout.AppDir)
	if err != nil {
		return []string{}, nil
	}

	// Get list of domains to exclude domain-only modules
	domains, _ := ListDomains()
	domainSet := make(map[string]bool)
	for _, d := range domains {
		domainSet[d] = true
	}

	var services []string
	for _, entry := range entries {
		if !entry.IsDir() || startsWith(entry.Name(), ".") {
			continue
		}

		dirPath := filepath.Join(layout.AppDir, entry.Name())

		// Check if this directory contains a service.go file
		serviceFile := filepath.Join(dirPath, "service.go")
		if IsFile(serviceFile) {
			// This is a service module
			serviceName := ToPascalCase(entry.Name()) + "Service"
			services = append(services, serviceName)
		}
	}

	return services, nil
}

// DeleteService deletes a service and its related files.
func DeleteService(serviceName string) DeleteResult {
	layout, err := ResolveProjectLayout()
	if err != nil {
		return DeleteResult{
			Success: false,
			Errors:  []string{"no project detected"},
			Message: "未检测到项目",
		}
	}

	// Convert service name to directory name (e.g., OrderService -> order)
	dirName := ToSnakeCase(serviceName)
	// Remove "_service" suffix if present
	if len(dirName) > 8 && dirName[len(dirName)-8:] == "_service" {
		dirName = dirName[:len(dirName)-8]
	}

	serviceDir := filepath.Join(layout.AppDir, dirName)

	// Check if service.go exists in this directory
	serviceFile := filepath.Join(serviceDir, "service.go")
	if !IsFile(serviceFile) {
		return DeleteResult{
			Success: false,
			Errors:  []string{"service not found"},
			Message: fmt.Sprintf("服务 %s 不存在", serviceName),
		}
	}

	var deletedItems []string
	var errors []string

	// Check if this directory also has domain files (commands.go, queries.go)
	// If so, only delete service.go and not the whole directory
	commandsFile := filepath.Join(serviceDir, "commands.go")
	queriesFile := filepath.Join(serviceDir, "queries.go")

	if IsFile(commandsFile) || IsFile(queriesFile) {
		// This is a domain module with a service, only delete service.go
		if err := os.Remove(serviceFile); err != nil {
			errors = append(errors, fmt.Sprintf("service file: %v", err))
		} else {
			deletedItems = append(deletedItems, "application/"+dirName+"/service.go")
		}
	} else {
		// This is a standalone service, delete the whole directory
		if err := os.RemoveAll(serviceDir); err != nil {
			errors = append(errors, fmt.Sprintf("service directory: %v", err))
		} else {
			deletedItems = append(deletedItems, "application/"+dirName+"/")
		}
	}

	if len(errors) > 0 {
		return DeleteResult{
			Success:      false,
			DeletedItems: deletedItems,
			Errors:       errors,
			Message:      fmt.Sprintf("服务 %s 删除部分失败", serviceName),
		}
	}

	return DeleteResult{
		Success:      true,
		DeletedItems: deletedItems,
		Message:      fmt.Sprintf("服务 %s 删除成功", serviceName),
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
