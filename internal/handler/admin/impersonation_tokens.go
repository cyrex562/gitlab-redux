package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// FeatureCategory represents the feature category for impersonation tokens
const FeatureCategory = "user_management"

// TokenParams represents the parameters for creating a token
type TokenParams struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
	Scopes      []string  `json:"scopes" binding:"required,min=1"`
}

// ImpersonationTokensHandler handles impersonation token requests
type ImpersonationTokensHandler struct {
	impersonationService *service.ImpersonationService
}

// NewImpersonationTokensHandler creates a new ImpersonationTokensHandler instance
func NewImpersonationTokensHandler(impersonationService *service.ImpersonationService) *ImpersonationTokensHandler {
	return &ImpersonationTokensHandler{
		impersonationService: impersonationService,
	}
}

// Index handles the GET request to list impersonation tokens
func (h *ImpersonationTokensHandler) Index(c *gin.Context) {
	// Set feature category
	c.Set("feature_category", FeatureCategory)

	// Verify impersonation is enabled
	if !h.verifyImpersonationEnabled(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Impersonation is not enabled"})
		return
	}

	username := c.Param("user_id")
	user, err := h.impersonationService.GetUserByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get available scopes
	scopes, err := h.impersonationService.GetAvailableScopes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get available scopes"})
		return
	}

	// Get active tokens
	tokens, err := h.impersonationService.GetActiveTokens(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active tokens"})
		return
	}

	// Check if user can be impersonated
	canImpersonate := h.canImpersonateUser(c, user)
	impersonationError := ""
	if canImpersonate {
		impersonationError = h.getImpersonationError(c, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"scopes":                scopes,
		"active_tokens":         tokens,
		"can_impersonate":       canImpersonate,
		"impersonation_error":   impersonationError,
	})
}

// Create handles the POST request to create a new impersonation token
func (h *ImpersonationTokensHandler) Create(c *gin.Context) {
	// Set feature category
	c.Set("feature_category", FeatureCategory)

	username := c.Param("user_id")
	var params TokenParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token data"})
		return
	}

	user, err := h.impersonationService.GetUserByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Create token from params
	token := &model.ImpersonationToken{
		Name:        params.Name,
		Description: params.Description,
		ExpiresAt:   params.ExpiresAt,
		Scopes:      params.Scopes,
		OrganizationID: c.GetUint("current_organization_id"),
	}

	createdToken, err := h.impersonationService.CreateToken(c, user.ID, token)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Get updated active tokens
	activeTokens, err := h.impersonationService.GetActiveTokens(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"new_token":         createdToken.Token,
		"active_tokens":     activeTokens,
		"total":            len(activeTokens),
	})
}

// Revoke handles the DELETE request to revoke an impersonation token
func (h *ImpersonationTokensHandler) Revoke(c *gin.Context) {
	// Set feature category
	c.Set("feature_category", FeatureCategory)

	username := c.Param("user_id")
	tokenID := c.Param("id")

	user, err := h.impersonationService.GetUserByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token, err := h.impersonationService.RevokeToken(c, user.ID, tokenID)
	if err != nil {
		c.SetFlash("alert", "Could not revoke impersonation token "+token.Name)
	} else {
		c.SetFlash("notice", "Revoked impersonation token "+token.Name+"!")
	}

	c.Redirect(http.StatusFound, "/admin/users/"+username+"/impersonation_tokens")
}

// Rotate handles the POST request to rotate an impersonation token
func (h *ImpersonationTokensHandler) Rotate(c *gin.Context) {
	// Set feature category
	c.Set("feature_category", FeatureCategory)

	username := c.Param("user_id")
	tokenID := c.Param("id")

	user, err := h.impersonationService.GetUserByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	result, err := h.impersonationService.RotateToken(c, user.ID, tokenID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Get updated active tokens
	activeTokens, err := h.impersonationService.GetActiveTokens(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"new_token":     result.Token,
		"active_tokens": activeTokens,
		"total":        len(activeTokens),
	})
}

// verifyImpersonationEnabled checks if impersonation is enabled
func (h *ImpersonationTokensHandler) verifyImpersonationEnabled(c *gin.Context) bool {
	// TODO: Implement proper impersonation check
	// This should check if impersonation is enabled in the system configuration
	return true
}

// canImpersonateUser checks if the current user can impersonate the target user
func (h *ImpersonationTokensHandler) canImpersonateUser(c *gin.Context, user *model.User) bool {
	// TODO: Implement proper impersonation check
	// This should check:
	// 1. If the current user has admin privileges
	// 2. If the target user is not already being impersonated
	// 3. If the target user is not an admin
	return true
}

// getImpersonationError returns any error message related to user impersonation
func (h *ImpersonationTokensHandler) getImpersonationError(c *gin.Context, user *model.User) string {
	// TODO: Implement proper impersonation error check
	// This should return appropriate error messages for:
	// 1. User already being impersonated
	// 2. User is an admin
	// 3. Other impersonation restrictions
	return ""
}
