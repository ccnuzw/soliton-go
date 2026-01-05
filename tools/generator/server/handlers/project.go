package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

// InitProjectRequest is the request body for project initialization.
type InitProjectRequest struct {
	Name             string `json:"name" binding:"required"`
	ModuleName       string `json:"module_name"`
	FrameworkVersion string `json:"framework_version"`
	FrameworkReplace string `json:"framework_replace"`
}

// InitProject handles POST /api/projects/init
func InitProject(c *gin.Context) {
	var req InitProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.ProjectConfig{
		Name:             req.Name,
		ModuleName:       req.ModuleName,
		FrameworkVersion: req.FrameworkVersion,
		FrameworkReplace: req.FrameworkReplace,
	}

	if err := core.ValidateProjectConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.InitProject(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PreviewInitProject handles POST /api/projects/init/preview
func PreviewInitProject(c *gin.Context) {
	var req InitProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.ProjectConfig{
		Name:             req.Name,
		ModuleName:       req.ModuleName,
		FrameworkVersion: req.FrameworkVersion,
		FrameworkReplace: req.FrameworkReplace,
	}

	if err := core.ValidateProjectConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewInitProject(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
