package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// FeatureCategory represents the feature category for impersonations
const FeatureCategory = "user_management"

// ImpersonationsHandler handles impersonation session requests
type ImpersonationsHandler struct {
	impersonationService *service.ImpersonationService
}

// NewImpersonationsHandler creates a new ImpersonationsHandler instance
func NewImpersonationsHandler(impersonationService *service.ImpersonationService) *ImpersonationsHandler {
	return &ImpersonationsHandler{
		impersonationService: impersonationService,
	}
}

// Destroy handles the DELETE request to stop an impersonation session
func (h *ImpersonationsHandler) Destroy(c *gin.Context) {
	// Set feature category
	c.Set("feature_category", FeatureCategory)

	// Authenticate impersonator
	if !h.authenticateImpersonator(c) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	// Stop impersonation and get original user
	originalUser, err := h.impersonationService.StopImpersonation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop impersonation"})
		return
	}

	// Redirect to original user's admin page
	c.Redirect(http.StatusFound, "/admin/users/"+originalUser.Username)
}

// authenticateImpersonator checks if the current user is a valid impersonator
func (h *ImpersonationsHandler) authenticateImpersonator(c *gin.Context) bool {
	// Get current user from context
	currentUser, exists := c.Get("current_user")
	if !exists {
		return false
	}

	user, ok := currentUser.(*model.User)
	if !ok {
		return false
	}

	// Check if user is admin and not blocked
	return user.IsAdmin && !user.Blocked
}
