package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

type ValueObjectRequest struct {
	Domain string             `json:"domain" binding:"required"`
	Name   string             `json:"name" binding:"required"`
	Fields []core.FieldConfig `json:"fields"`
	Force  bool               `json:"force"`
}

type SpecificationRequest struct {
	Domain string `json:"domain" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Target string `json:"target"`
	Force  bool   `json:"force"`
}

type PolicyRequest struct {
	Domain string `json:"domain" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Target string `json:"target"`
	Force  bool   `json:"force"`
}

type EventRequest struct {
	Domain string             `json:"domain" binding:"required"`
	Name   string             `json:"name" binding:"required"`
	Fields []core.FieldConfig `json:"fields"`
	Topic  string             `json:"topic"`
	Force  bool               `json:"force"`
}

type EventHandlerRequest struct {
	Domain    string `json:"domain" binding:"required"`
	EventName string `json:"event_name" binding:"required"`
	Topic     string `json:"topic"`
	Force     bool   `json:"force"`
}

func GenerateValueObject(c *gin.Context) {
	var req ValueObjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.ValueObjectConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Fields: req.Fields,
		Force:  req.Force,
	}

	if err := core.ValidateValueObjectConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.GenerateValueObject(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func PreviewValueObject(c *gin.Context) {
	var req ValueObjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.ValueObjectConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Fields: req.Fields,
		Force:  req.Force,
	}

	if err := core.ValidateValueObjectConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewValueObject(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GenerateSpecification(c *gin.Context) {
	var req SpecificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.SpecificationConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Target: req.Target,
		Force:  req.Force,
	}

	if err := core.ValidateSpecificationConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.GenerateSpecification(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func PreviewSpecification(c *gin.Context) {
	var req SpecificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.SpecificationConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Target: req.Target,
		Force:  req.Force,
	}

	if err := core.ValidateSpecificationConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewSpecification(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GeneratePolicy(c *gin.Context) {
	var req PolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.PolicyConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Target: req.Target,
		Force:  req.Force,
	}

	if err := core.ValidatePolicyConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.GeneratePolicy(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func PreviewPolicy(c *gin.Context) {
	var req PolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.PolicyConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Target: req.Target,
		Force:  req.Force,
	}

	if err := core.ValidatePolicyConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewPolicy(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GenerateEvent(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.EventConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Fields: req.Fields,
		Topic:  req.Topic,
		Force:  req.Force,
	}

	if err := core.ValidateEventConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.GenerateEvent(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func PreviewEvent(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.EventConfig{
		Domain: req.Domain,
		Name:   req.Name,
		Fields: req.Fields,
		Topic:  req.Topic,
		Force:  req.Force,
	}

	if err := core.ValidateEventConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewEvent(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GenerateEventHandler(c *gin.Context) {
	var req EventHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.EventHandlerConfig{
		Domain:    req.Domain,
		EventName: req.EventName,
		Topic:     req.Topic,
		Force:     req.Force,
	}

	if err := core.ValidateEventHandlerConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.GenerateEventHandler(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func PreviewEventHandler(c *gin.Context) {
	var req EventHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := core.EventHandlerConfig{
		Domain:    req.Domain,
		EventName: req.EventName,
		Topic:     req.Topic,
		Force:     req.Force,
	}

	if err := core.ValidateEventHandlerConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := core.PreviewEventHandler(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
