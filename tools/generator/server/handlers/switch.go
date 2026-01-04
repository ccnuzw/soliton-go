package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

// SwitchProjectRequest is the request body for switching projects.
type SwitchProjectRequest struct {
	Path string `json:"path" binding:"required"`
}

// SwitchProject handles POST /api/projects/switch
func SwitchProject(c *gin.Context) {
	var req SwitchProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate path exists
	absPath, err := filepath.Abs(req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}

	// Check if path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path does not exist or is not a directory"})
		return
	}

	// Change working directory
	if err := os.Chdir(absPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to switch directory"})
		return
	}

	// Get new layout
	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Switched to directory, but no project detected",
			"path":    absPath,
			"layout": gin.H{
				"found": false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully switched to project",
		"path":    absPath,
		"layout": gin.H{
			"found":          true,
			"module_path":    layout.ModulePath,
			"module_dir":     layout.ModuleDir,
			"internal_dir":   layout.InternalDir,
			"domain_dir":     layout.DomainDir,
			"app_dir":        layout.AppDir,
			"infra_dir":      layout.InfraDir,
			"interfaces_dir": layout.InterfacesDir,
		},
	})
}

// ListProjects handles GET /api/projects/list
func ListProjects(c *gin.Context) {
	cwd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get current directory"})
		return
	}

	var projects []gin.H
	currentProjectDir := ""

	// Check current directory
	if layout, err := core.ResolveProjectLayout(); err == nil {
		currentProjectDir = layout.ModuleDir
		projects = append(projects, gin.H{
			"name":        filepath.Base(layout.ModuleDir),
			"path":        layout.ModuleDir,
			"module_path": layout.ModulePath,
			"is_current":  true,
		})
	}

	// Determine search directory
	// If we're in tools/generator, search in ../../ (soliton-go root)
	// Otherwise, search in parent directory
	searchDir := filepath.Dir(cwd)
	if filepath.Base(cwd) == "generator" && filepath.Base(filepath.Dir(cwd)) == "tools" {
		// We're in tools/generator, go up two levels
		searchDir = filepath.Dir(filepath.Dir(cwd))
	}

	// Check sibling directories
	entries, err := os.ReadDir(searchDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			dirPath := filepath.Join(searchDir, entry.Name())

			// Skip if it's the current project or tools directory
			if dirPath == currentProjectDir || entry.Name() == "tools" {
				continue
			}

			// Check if it's a valid project
			if core.IsDir(filepath.Join(dirPath, "internal")) && core.IsFile(filepath.Join(dirPath, "go.mod")) {
				modulePath, _ := core.ReadModulePath(filepath.Join(dirPath, "go.mod"))
				projects = append(projects, gin.H{
					"name":        entry.Name(),
					"path":        dirPath,
					"module_path": modulePath,
					"is_current":  false,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}
