package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct {
	healthCheckService *service.HealthCheckService
}

// NewHealthCheckHandler creates a new HealthCheckHandler instance
func NewHealthCheckHandler(healthCheckService *service.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{
		healthCheckService: healthCheckService,
	}
}

// Show handles the GET request to display health check results
func (h *HealthCheckHandler) Show(c *gin.Context) {
	// Check authorization
	if !h.hasAdminHealthCheckPermission(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	errors, err := h.healthCheckService.ProcessChecks(c, h.getChecks())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process health checks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errors": errors,
	})
}

// hasAdminHealthCheckPermission checks if the user has permission to read admin health checks
func (h *HealthCheckHandler) hasAdminHealthCheckPermission(c *gin.Context) bool {
	// TODO: Implement proper authorization check
	// This should integrate with your authorization system
	return true
}

// getChecks returns the list of health checks to perform
func (h *HealthCheckHandler) getChecks() []string {
	return []string{"standard"}
}
