package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/core"
)

// ServiceDetectionResult represents the result of service type detection
type ServiceDetectionResult struct {
	ServiceName     string `json:"service_name"`
	DomainName      string `json:"domain_name"`
	DomainExists    bool   `json:"domain_exists"`
	ServiceType     string `json:"service_type"` // "domain_service" | "cross_domain_service"
	TargetDir       string `json:"target_dir"`
	ShouldReuseDTO  bool   `json:"should_reuse_dto"`
	ExistingDTOPath string `json:"existing_dto_path,omitempty"`
	Message         string `json:"message"`
}

// DetectServiceType handles GET /api/services/detect/:name
func DetectServiceType(c *gin.Context) {
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

	// 1. Extract domain name from service name
	// "OrderService" -> "order"
	domainName := extractDomainNameFromService(serviceName)

	// 2. Check if domain exists
	domainPath := filepath.Join(layout.DomainDir, domainName)
	domainExists := core.IsDir(domainPath)

	// 3. Check if DTO exists
	dtoPath := filepath.Join(layout.InternalDir, "application", domainName, "dto.go")
	dtoExists := core.IsFile(dtoPath)

	// 4. Build detection result
	result := ServiceDetectionResult{
		ServiceName:  serviceName,
		DomainName:   domainName,
		DomainExists: domainExists,
	}

	if domainExists {
		// Domain service
		result.ServiceType = "domain_service"
		result.TargetDir = filepath.Join("internal/application", domainName)
		result.ShouldReuseDTO = dtoExists

		if dtoExists {
			result.ExistingDTOPath = dtoPath
			result.Message = fmt.Sprintf("检测到已有 %s 领域，将生成到 %s/service.go 并复用现有 DTO", domainName, result.TargetDir)
		} else {
			result.Message = fmt.Sprintf("检测到已有 %s 领域，将生成到 %s/service.go", domainName, result.TargetDir)
		}
	} else {
		// Cross-domain service
		result.ServiceType = "cross_domain_service"
		result.TargetDir = filepath.Join("internal/application", domainName)
		result.ShouldReuseDTO = false
		result.Message = fmt.Sprintf("未检测到 %s 领域，将生成到 %s/service.go（跨领域服务）", domainName, result.TargetDir)
	}

	c.JSON(http.StatusOK, result)
}

// extractDomainNameFromService extracts domain name from service name
// "OrderService" -> "order"
// "PaymentService" -> "payment"
func extractDomainNameFromService(serviceName string) string {
	name := strings.TrimSuffix(serviceName, "Service")
	return strings.ToLower(name)
}
