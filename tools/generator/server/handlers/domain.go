package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

// DomainRequest is the request body for domain generation.
type DomainRequest struct {
	Name       string             `json:"name" binding:"required"`
	Remark     string             `json:"remark"`
	Fields     []core.FieldConfig `json:"fields"`
	TableName  string             `json:"table_name"`
	RouteBase  string             `json:"route_base"`
	SoftDelete bool               `json:"soft_delete"`
	Wire       bool               `json:"wire"`
	Force      bool               `json:"force"`
}

// GenerateDomain handles POST /api/domains
func GenerateDomain(c *gin.Context) {
	var req DomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.DomainConfig{
		Name:       req.Name,
		Remark:     req.Remark,
		Fields:     req.Fields,
		TableName:  req.TableName,
		RouteBase:  req.RouteBase,
		SoftDelete: req.SoftDelete,
		Wire:       req.Wire,
		Force:      req.Force,
	}

	if err := core.ValidateDomainConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.GenerateDomain(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PreviewDomain handles POST /api/domains/preview
func PreviewDomain(c *gin.Context) {
	var req DomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.DomainConfig{
		Name:       req.Name,
		Remark:     req.Remark,
		Fields:     req.Fields,
		TableName:  req.TableName,
		RouteBase:  req.RouteBase,
		SoftDelete: req.SoftDelete,
		Wire:       req.Wire,
		Force:      req.Force,
	}

	if err := core.ValidateDomainConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewDomain(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetFieldTypes handles GET /api/field-types
func GetFieldTypes(c *gin.Context) {
	types := core.GetAvailableFieldTypes()
	c.JSON(http.StatusOK, gin.H{"types": types})
}
