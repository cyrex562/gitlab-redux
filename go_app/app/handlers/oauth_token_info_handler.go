package handlers

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/app/services"
	"github.com/gin-gonic/gin"
)

// OAuthTokenInfoHandler handles OAuth token info requests
type OAuthTokenInfoHandler struct {
	*BaseHandler
	oauthService *services.OAuthService
}

// NewOAuthTokenInfoHandler creates a new OAuthTokenInfoHandler
func NewOAuthTokenInfoHandler(baseHandler *BaseHandler, oauthService *services.OAuthService) *OAuthTokenInfoHandler {
	return &OAuthTokenInfoHandler{
		BaseHandler:  baseHandler,
		oauthService: oauthService,
	}
}

// Show returns information about the current OAuth token
func (h *OAuthTokenInfoHandler) Show(c *gin.Context) {
	// Get token from Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_token",
			"error_description": "The access token is invalid",
		})
		return
	}

	// Validate token and get token info
	tokenInfo, err := h.oauthService.GetTokenInfo(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_token",
			"error_description": "The access token is invalid",
		})
		return
	}

	// Return token info with backwards compatibility fields
	c.JSON(http.StatusOK, gin.H{
		"resource_owner_id": tokenInfo.UserID,
		"scopes":           tokenInfo.Scopes,
		"expires_in_seconds": tokenInfo.ExpiresIn,
		"application": gin.H{
			"uid": tokenInfo.ApplicationID,
		},
		"created_at": tokenInfo.CreatedAt.Unix(),
	})
} 