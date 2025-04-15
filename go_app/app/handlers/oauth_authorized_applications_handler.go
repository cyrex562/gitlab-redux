package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// OAuthAuthorizedApplicationsHandler handles OAuth authorized applications requests
type OAuthAuthorizedApplicationsHandler struct {
	*BaseHandler
	oauthService *services.OAuthService
}

// NewOAuthAuthorizedApplicationsHandler creates a new OAuthAuthorizedApplicationsHandler
func NewOAuthAuthorizedApplicationsHandler(baseHandler *BaseHandler, oauthService *services.OAuthService) *OAuthAuthorizedApplicationsHandler {
	return &OAuthAuthorizedApplicationsHandler{
		BaseHandler:  baseHandler,
		oauthService: oauthService,
	}
}

// Index handles the index page for authorized applications
func (h *OAuthAuthorizedApplicationsHandler) Index(c *gin.Context) {
	// In the original Ruby code, this returns a 404 Not Found
	// We'll implement the same behavior
	c.HTML(http.StatusNotFound, "errors/not_found", nil)
}

// Destroy handles revoking an authorized application
func (h *OAuthAuthorizedApplicationsHandler) Destroy(c *gin.Context) {
	currentUser := h.GetCurrentUser(c.Request)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if we're revoking a specific token or an entire application
	tokenID := c.Query("token_id")
	if tokenID != "" {
		// Revoke a specific token
		tokenIDUint, err := strconv.ParseUint(tokenID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
			return
		}

		err = h.oauthService.RevokeToken(currentUser.ID, uint(tokenIDUint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Revoke all tokens for an application
		appID := c.Param("id")
		if appID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
			return
		}

		appIDUint, err := strconv.ParseUint(appID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
			return
		}

		err = h.oauthService.RevokeApplicationTokens(currentUser.ID, uint(appIDUint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Set flash message
	// TODO: Implement flash messages

	// Redirect to applications page
	c.Redirect(http.StatusFound, "/user/settings/applications")
} 