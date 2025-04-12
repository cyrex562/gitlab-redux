package googleapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// AuthorizationsController handles Google API authorization callbacks
type AuthorizationsController struct {
	cloudPlatformClient *CloudPlatformClient
}

// NewAuthorizationsController creates a new AuthorizationsController
func NewAuthorizationsController(cloudPlatformClient *CloudPlatformClient) *AuthorizationsController {
	return &AuthorizationsController{
		cloudPlatformClient: cloudPlatformClient,
	}
}

// Callback handles the response from Google after the user
// goes through authentication and authorization process
func (c *AuthorizationsController) Callback(ctx *gin.Context) {
	// Get redirect URI from session
	redirectURI, err := c.getRedirectURIFromSession(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if there's an error parameter (user declined authorization)
	if errorParam := ctx.Query("error"); errorParam != "" {
		ctx.SetFlash("alert", "Google Cloud authorizations required")
		// Get error URI from session
		errorURI, _ := ctx.Get("error_uri")
		if errorURI != nil {
			redirectURI = errorURI.(string)
		}
		ctx.Redirect(http.StatusFound, redirectURI)
		return
	}

	// Check if there's a code parameter (successful authorization)
	if code := ctx.Query("code"); code != "" {
		// Get callback URL for Google API
		callbackURL := fmt.Sprintf("%s/google_api/auth/callback", ctx.Request.Host)

		// Get token from Google API
		token, expiresAt, err := c.cloudPlatformClient.GetToken(code, callbackURL)
		if err != nil {
			// Handle timeout or connection errors
			if errors.Is(err, ErrTimeout) || errors.Is(err, ErrConnectionFailed) {
				ctx.SetFlash("alert", "Timeout connecting to the Google API. Please try again.")
				ctx.Redirect(http.StatusFound, redirectURI)
				return
			}
			// Handle other errors
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Store token and expiration in session
		state := ctx.Query("state")
		if state != "" {
			ctx.Set(CloudPlatformClient.SessionKeyForToken(state), token)
			ctx.Set(CloudPlatformClient.SessionKeyForExpiresAt(state), expiresAt.Format(time.RFC3339))
		}
	}

	// Redirect to the appropriate URI
	ctx.Redirect(http.StatusFound, redirectURI)
}

// getRedirectURIFromSession retrieves the redirect URI from the session
func (c *AuthorizationsController) getRedirectURIFromSession(ctx *gin.Context) (string, error) {
	state := ctx.Query("state")
	if state == "" {
		return "", fmt.Errorf("state parameter is required")
	}

	sessionKey := CloudPlatformClient.SessionKeyForRedirectURI(state)
	redirectURI, exists := ctx.Get(sessionKey)
	if !exists {
		return "", fmt.Errorf("redirect URI not found in session")
	}

	return redirectURI.(string), nil
}

// CloudPlatformClient represents a client for Google Cloud Platform API
type CloudPlatformClient struct {
	// Add fields as needed
}

// NewCloudPlatformClient creates a new CloudPlatformClient
func NewCloudPlatformClient() *CloudPlatformClient {
	return &CloudPlatformClient{}
}

// GetToken retrieves a token from Google API
func (c *CloudPlatformClient) GetToken(code, redirectURI string) (string, time.Time, error) {
	// Implementation would go here
	// This is a placeholder for the actual implementation
	return "token", time.Now().Add(1 * time.Hour), nil
}

// SessionKeyForToken returns the session key for storing the token
func (c *CloudPlatformClient) SessionKeyForToken(state string) string {
	return fmt.Sprintf("google_api_token_%s", state)
}

// SessionKeyForExpiresAt returns the session key for storing the token expiration
func (c *CloudPlatformClient) SessionKeyForExpiresAt(state string) string {
	return fmt.Sprintf("google_api_expires_at_%s", state)
}

// SessionKeyForRedirectURI returns the session key for storing the redirect URI
func (c *CloudPlatformClient) SessionKeyForRedirectURI(state string) string {
	return fmt.Sprintf("google_api_redirect_uri_%s", state)
}

// Error types
var (
	ErrTimeout           = errors.New("timeout connecting to Google API")
	ErrConnectionFailed  = errors.New("connection to Google API failed")
)
