package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// GitalyServersHandler handles Gitaly server related requests
type GitalyServersHandler struct {
	gitalyService *service.GitalyService
}

// NewGitalyServersHandler creates a new GitalyServersHandler instance
func NewGitalyServersHandler(gitalyService *service.GitalyService) *GitalyServersHandler {
	return &GitalyServersHandler{
		gitalyService: gitalyService,
	}
}

// Index handles the GET request to list all Gitaly servers
func (h *GitalyServersHandler) Index(c *gin.Context) {
	// Check authorization
	if !h.hasAdminGitalyServersPermission(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	servers, err := h.gitalyService.GetAllServers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Gitaly servers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gitaly_servers": servers,
	})
}

// hasAdminGitalyServersPermission checks if the user has permission to read admin Gitaly servers
func (h *GitalyServersHandler) hasAdminGitalyServersPermission(c *gin.Context) bool {
	// TODO: Implement proper authorization check
	// This should integrate with your authorization system
	return true
}
