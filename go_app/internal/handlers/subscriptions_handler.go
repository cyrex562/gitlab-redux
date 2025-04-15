package handlers

import (
	"net/http"
	"strings"

	"github.com/cyrex562/gitlab-redux/internal/services"
	"github.com/gin-gonic/gin"
)

// ALLOWED_IFRAME_ANCESTORS defines the allowed iframe ancestors for CSP
var ALLOWED_IFRAME_ANCESTORS = []string{"'self'", "https://*.atlassian.net", "https://*.jira.com"}

type SubscriptionsHandler struct {
	jiraService *services.JiraService
}

func NewSubscriptionsHandler(jiraService *services.JiraService) *SubscriptionsHandler {
	return &SubscriptionsHandler{
		jiraService: jiraService,
	}
}

// Index handles GET /api/v4/jira/connect/subscriptions
// Lists all subscriptions for the current Jira installation
func (h *SubscriptionsHandler) Index(c *gin.Context) {
	// Allow rendering in iframe
	c.Header("X-Frame-Options", "")

	// Set Content Security Policy
	h.setContentSecurityPolicy(c)

	// Get the current Jira installation
	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get subscriptions for the installation
	subscriptions, err := h.jiraService.GetSubscriptions(installation.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle different response formats
	format := c.DefaultQuery("format", "html")
	if format == "json" {
		c.JSON(http.StatusOK, gin.H{
			"subscriptions": subscriptions,
		})
	} else {
		// Render HTML template
		c.HTML(http.StatusOK, "subscriptions/index.html", gin.H{
			"subscriptions": subscriptions,
		})
	}
}

// Create handles POST /api/v4/jira/connect/subscriptions
// Creates a new subscription
func (h *SubscriptionsHandler) Create(c *gin.Context) {
	// Get the current user
	user, err := h.jiraService.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get the current Jira installation
	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get request parameters
	var params struct {
		NamespacePath string `json:"namespace_path"`
		JiraUser      string `json:"jira_user"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the subscription
	_, err = h.jiraService.CreateSubscription(installation.ID, user.ID, params.NamespacePath, params.JiraUser)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Destroy handles DELETE /api/v4/jira/connect/subscriptions/:id
// Deletes a subscription
func (h *SubscriptionsHandler) Destroy(c *gin.Context) {
	// Get the subscription ID
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	// Get the current Jira installation
	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the subscription
	err = h.jiraService.DeleteSubscription(installation.ID, subscriptionID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// setContentSecurityPolicy sets the Content Security Policy headers
func (h *SubscriptionsHandler) setContentSecurityPolicy(c *gin.Context) {
	// Get the current Jira installation
	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		return
	}

	// Base CSP directives
	c.Header("Content-Security-Policy", "default-src 'self'")

	// Add frame-ancestors directive
	frameAncestors := ALLOWED_IFRAME_ANCESTORS
	if installation.InstanceURL != "" {
		// Add additional iframe ancestors from configuration
		additionalAncestors := h.jiraService.GetAdditionalIframeAncestors()
		frameAncestors = append(frameAncestors, additionalAncestors...)
	}
	c.Header("Content-Security-Policy", c.GetHeader("Content-Security-Policy")+"; frame-ancestors "+strings.Join(frameAncestors, " "))

	// Add script-src directive
	c.Header("Content-Security-Policy", c.GetHeader("Content-Security-Policy")+"; script-src 'self' https://connect-cdn.atl-paas.net")

	// Add style-src directive
	c.Header("Content-Security-Policy", c.GetHeader("Content-Security-Policy")+"; style-src 'self' 'unsafe-inline'")

	// Add connect-src directive for instance URL if available
	if installation.InstanceURL != "" {
		connectSrc := []string{
			installation.InstanceURL + "/-/jira_connect/",
			installation.InstanceURL + "/api/",
			installation.InstanceURL + "/oauth/token",
		}
		c.Header("Content-Security-Policy", c.GetHeader("Content-Security-Policy")+"; connect-src 'self' "+strings.Join(connectSrc, " "))
	}
} 