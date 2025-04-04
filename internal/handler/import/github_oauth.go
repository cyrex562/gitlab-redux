package import

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// OAuthConfigMissingError represents an error when OAuth configuration is missing
type OAuthConfigMissingError struct{}

func (e *OAuthConfigMissingError) Error() string {
	return "Missing OAuth configuration for GitHub"
}

// GithubOauth provides GitHub OAuth functionality for imports
type GithubOauth struct {
	config *oauth2.Config
}

// NewGithubOauth creates a new instance of GithubOauth
func NewGithubOauth(config *oauth2.Config) *GithubOauth {
	return &GithubOauth{
		config: config,
	}
}

// RegisterRoutes registers the routes for GitHub OAuth
func (g *GithubOauth) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/auth", g.providerAuth)
	router.GET("/callback", g.callback)
}

// providerAuth handles the OAuth authorization flow
func (g *GithubOauth) providerAuth(ctx *gin.Context) {
	// Check if access token exists in session
	if token, exists := ctx.Get("access_token"); exists && token != nil {
		return
	}

	// Check if CI/CD only mode is enabled
	if g.isCICDOnly(ctx) {
		return
	}

	// Generate state
	state := g.generateState()
	ctx.SetCookie("auth_state", state, 3600, "/", "", false, true)

	// Store failure path
	ctx.SetCookie("auth_on_failure_path", "/new#import_project", 3600, "/", "", false, true)

	// Redirect to GitHub
	url := g.config.AuthCodeURL(state, oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("scope", "repo read:org"))
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// callback handles the OAuth callback
func (g *GithubOauth) callback(ctx *gin.Context) {
	// Verify state
	state, err := ctx.Cookie("auth_state")
	if err != nil {
		g.handleError(ctx, "Invalid state")
		return
	}

	if state != ctx.Query("state") {
		g.handleError(ctx, "State mismatch")
		return
	}

	// Exchange code for token
	code := ctx.Query("code")
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		g.handleError(ctx, "Failed to exchange code for token")
		return
	}

	// Store token in session
	ctx.Set("access_token", token.AccessToken)

	// Redirect back to import page
	failurePath, _ := ctx.Cookie("auth_on_failure_path")
	if failurePath == "" {
		failurePath = "/new#import_project"
	}
	ctx.Redirect(http.StatusTemporaryRedirect, failurePath)
}

// isCICDOnly checks if CI/CD only mode is enabled
func (g *GithubOauth) isCICDOnly(ctx *gin.Context) bool {
	cicdOnly := ctx.Query("ci_cd_only")
	return cicdOnly == "1" || cicdOnly == "true"
}

// generateState generates a random state string
func (g *GithubOauth) generateState() string {
	b := make([]byte, 64)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// handleError handles OAuth errors
func (g *GithubOauth) handleError(ctx *gin.Context, message string) {
	// Clear access token from session
	ctx.Set("access_token", nil)

	// Handle different response formats
	switch ctx.GetHeader("Accept") {
	case "application/json":
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"errors": message,
		})
	default:
		ctx.Redirect(http.StatusTemporaryRedirect, "/new/import?alert="+message)
	}
}

// GetToken retrieves the access token from the session
func (g *GithubOauth) GetToken(ctx *gin.Context) (string, error) {
	token, exists := ctx.Get("access_token")
	if !exists || token == nil {
		return "", &OAuthConfigMissingError{}
	}
	return token.(string), nil
}
