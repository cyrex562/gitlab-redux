package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/app/models"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// OAuthApplicationsHandler handles OAuth application requests
type OAuthApplicationsHandler struct {
	*BaseHandler
	oauthService *services.OAuthService
}

// NewOAuthApplicationsHandler creates a new OAuthApplicationsHandler
func NewOAuthApplicationsHandler(baseHandler *BaseHandler, oauthService *services.OAuthService) *OAuthApplicationsHandler {
	return &OAuthApplicationsHandler{
		BaseHandler:  baseHandler,
		oauthService: oauthService,
	}
}

// Index handles the index page for OAuth applications
func (h *OAuthApplicationsHandler) Index(c *gin.Context) {
	// TODO: Implement pagination with cursor
	currentUser := h.GetCurrentUser(c.Request)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	applications, err := h.oauthService.GetUserApplications(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authorizedTokens, err := h.oauthService.GetUserAuthorizedTokens(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := gin.H{
		"applications":      applications,
		"authorized_tokens": authorizedTokens,
		"application":       &models.OAuthApplication{}, // New application for form
	}

	c.HTML(http.StatusOK, "oauth/applications/index", data)
}

// Show handles displaying a single OAuth application
func (h *OAuthApplicationsHandler) Show(c *gin.Context) {
	// TODO: Implement show page
	c.HTML(http.StatusOK, "oauth/applications/show", nil)
}

// Create handles creating a new OAuth application
func (h *OAuthApplicationsHandler) Create(c *gin.Context) {
	currentUser := h.GetCurrentUser(c.Request)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	application := &models.OAuthApplication{
		Name:         c.PostForm("name"),
		RedirectURI:  c.PostForm("redirect_uri"),
		Scopes:       c.PostForm("scopes"),
		OwnerID:      currentUser.ID,
		// TODO: Add other fields as needed
	}

	createdApp, err := h.oauthService.CreateApplication(application)
	if err != nil {
		// TODO: Handle validation errors
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Set flash message
	// TODO: Implement flash messages

	data := gin.H{
		"application": createdApp,
		"created":     true,
	}

	c.HTML(http.StatusOK, "oauth/applications/show", data)
}

// RenewSecret handles renewing an OAuth application secret
func (h *OAuthApplicationsHandler) RenewSecret(c *gin.Context) {
	currentUser := h.GetCurrentUser(c.Request)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	appID, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	// Verify the application belongs to the current user
	application, err := h.oauthService.GetUserApplication(currentUser.ID, uint(appID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	newSecret, err := h.oauthService.RenewApplicationSecret(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": newSecret})
} 