package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

// ServiceRequest is the request body for service generation.
type ServiceRequest struct {
	Name    string                     `json:"name" binding:"required"`
	Remark  string                     `json:"remark"`
	Methods []core.ServiceMethodConfig `json:"methods"`
	Force   bool                       `json:"force"`
}

// GenerateService handles POST /api/services
func GenerateService(c *gin.Context) {
	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.ServiceConfig{
		Name:    req.Name,
		Remark:  req.Remark,
		Methods: req.Methods,
		Force:   req.Force,
	}

	if err := core.ValidateServiceConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.GenerateService(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PreviewService handles POST /api/services/preview
func PreviewService(c *gin.Context) {
	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.ServiceConfig{
		Name:    req.Name,
		Remark:  req.Remark,
		Methods: req.Methods,
		Force:   req.Force,
	}

	if err := core.ValidateServiceConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewService(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProjectLayout handles GET /api/layout
func GetProjectLayout(c *gin.Context) {
	layout, err := core.ResolveProjectLayout()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"found":   false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"found":          true,
		"module_path":    layout.ModulePath,
		"module_dir":     layout.ModuleDir,
		"internal_dir":   layout.InternalDir,
		"domain_dir":     layout.DomainDir,
		"app_dir":        layout.AppDir,
		"infra_dir":      layout.InfraDir,
		"interfaces_dir": layout.InterfacesDir,
	})
}
