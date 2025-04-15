package handlers

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/internal/services"
	"github.com/gin-gonic/gin"
)

type PublicKeysHandler struct {
	jiraService *services.JiraService
}

func NewPublicKeysHandler(jiraService *services.JiraService) *PublicKeysHandler {
	return &PublicKeysHandler{
		jiraService: jiraService,
	}
}

// ShowPublicKey handles GET /api/v4/jira/connect/public_keys/:id
// This endpoint serves public keys for Jira Connect
func (h *PublicKeysHandler) ShowPublicKey(c *gin.Context) {
	// Skip authentication as per the Ruby controller
	// This is handled by the router middleware

	// Check if public key storage is enabled
	if !h.jiraService.IsPublicKeyStorageEnabled() {
		c.Status(http.StatusNotFound)
		return
	}

	// Get the public key ID from the URL
	id := c.Param("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// Find the public key
	key, err := h.jiraService.FindPublicKey(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Render the public key as plain text
	c.String(http.StatusOK, key)
} 