package handlers

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/app/services"
	"github.com/gin-gonic/gin"
)

// OAuthTokenRevocationHandler handles OAuth token revocation requests
type OAuthTokenRevocationHandler struct {
	*BaseHandler
	oauthService *services.OAuthService
}

// NewOAuthTokenRevocationHandler creates a new OAuthTokenRevocationHandler
func NewOAuthTokenRevocationHandler(baseHandler *BaseHandler, oauthService *services.OAuthService) *OAuthTokenRevocationHandler {
	return &OAuthTokenRevocationHandler{
		BaseHandler:  baseHandler,
		oauthService: oauthService,
	}
}

// Revoke handles token revocation requests
func (h *OAuthTokenRevocationHandler) Revoke(c *gin.Context) {
	// Get token from request
	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
			"error_description": "Token parameter is required",
		})
		return
	}

	// Get token type hint (optional)
	tokenTypeHint := c.PostForm("token_type_hint")

	// Validate client credentials
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client",
			"error_description": "Client authentication failed",
		})
		return
	}

	// Revoke the token
	err := h.oauthService.RevokeToken(token, tokenTypeHint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
			"error_description": err.Error(),
		})
		return
	}

	// Return success response
	c.Status(http.StatusOK)
} 