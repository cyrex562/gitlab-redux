package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OauthCallbacksHandler struct {
	// No dependencies needed for this simple handler
}

func NewOauthCallbacksHandler() *OauthCallbacksHandler {
	return &OauthCallbacksHandler{}
}

// Index handles GET /api/v4/jira/connect/oauth/callbacks
// This serves as a landing page after users install and authenticate
// the GitLab.com for Jira Cloud app
func (h *OauthCallbacksHandler) Index(c *gin.Context) {
	// Skip authentication as per the Ruby controller
	// This is handled by the router middleware

	// Render the landing page
	// For now, we'll just return a simple HTML response
	// In a real implementation, this would render a template
	c.HTML(http.StatusOK, "oauth_callbacks/index.html", gin.H{
		"title": "GitLab for Jira - Installation Complete",
	})
} 