package import_controller

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/bitbucket"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/oauth"
)

// BitbucketController handles Bitbucket import functionality
type BitbucketController struct {
	*BaseController
	bitbucketService *bitbucket.Service
	oauthService     *oauth.Service
}

// NewBitbucketController creates a new Bitbucket controller
func NewBitbucketController(
	baseController *BaseController,
	bitbucketService *bitbucket.Service,
	oauthService *oauth.Service,
) *BitbucketController {
	return &BitbucketController{
		BaseController:   baseController,
		bitbucketService: bitbucketService,
		oauthService:     oauthService,
	}
}

// RegisterRoutes registers the routes for the Bitbucket controller
func (c *BitbucketController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/callback", c.Callback)
	router.GET("/status", c.Status)
	router.POST("/create", c.Create)
}

// Callback handles the OAuth callback from Bitbucket
func (c *BitbucketController) Callback(ctx *gin.Context) {
	authState := ctx.GetString("bitbucket_auth_state")
	ctx.Set("bitbucket_auth_state", "")

	if authState == "" || !secureCompare(authState, ctx.Query("state")) {
		c.goToBitbucketForPermissions(ctx)
		return
	}

	token, err := c.oauthService.ExchangeCode(
		ctx,
		ctx.Query("code"),
		c.getCallbackURL(ctx),
	)
	if err != nil {
		c.bitbucketUnauthorized(ctx, err)
		return
	}

	// Store token information in session
	ctx.Set("bitbucket_token", token.AccessToken)
	ctx.Set("bitbucket_expires_at", token.ExpiresAt)
	ctx.Set("bitbucket_expires_in", token.ExpiresIn)
	ctx.Set("bitbucket_refresh_token", token.RefreshToken)

	ctx.Redirect(http.StatusFound, c.getStatusURL(ctx))
}

// Create handles project creation from Bitbucket
func (c *BitbucketController) Create(ctx *gin.Context) {
	if !c.bitbucketImportEnabled() {
		ctx.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	repoID := ctx.PostForm("repo_id")
	name := strings.ReplaceAll(repoID, "___", "/")
	repo, err := c.bitbucketService.GetRepo(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projectName := ctx.PostForm("new_name")
	if projectName == "" {
		projectName = repo.Name
	}

	repoOwner := repo.Owner
	if repoOwner == c.bitbucketService.GetCurrentUser(ctx).Username {
		repoOwner = ctx.MustGet("current_user").(*models.User).Username
	}

	namespacePath := ctx.PostForm("new_namespace")
	if namespacePath == "" {
		namespacePath = repoOwner
	}

	targetNamespace, err := c.findOrCreateNamespace(ctx, namespacePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := ctx.MustGet("current_user").(*models.User)
	if !user.CanImportProjects(targetNamespace) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": "You are not allowed to import projects in this namespace.",
		})
		return
	}

	// Update token in session
	ctx.Set("bitbucket_token", c.bitbucketService.GetConnection(ctx).Token)

	project, err := c.bitbucketService.CreateProject(ctx, repo, projectName, targetNamespace, user)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, project)
}

// ImportableRepos returns the list of importable Bitbucket repositories
func (c *BitbucketController) ImportableRepos(ctx *gin.Context) ([]interface{}, error) {
	repos, err := c.bitbucketService.GetRepos(ctx, c.sanitizedFilterParam(ctx))
	if err != nil {
		return nil, err
	}

	var validRepos []interface{}
	for _, repo := range repos {
		if repo.IsValid() {
			validRepos = append(validRepos, repo)
		}
	}
	return validRepos, nil
}

// IncompatibleRepos returns the list of incompatible Bitbucket repositories
func (c *BitbucketController) IncompatibleRepos(ctx *gin.Context) ([]interface{}, error) {
	repos, err := c.bitbucketService.GetRepos(ctx, c.sanitizedFilterParam(ctx))
	if err != nil {
		return nil, err
	}

	var invalidRepos []interface{}
	for _, repo := range repos {
		if !repo.IsValid() {
			invalidRepos = append(invalidRepos, repo)
		}
	}
	return invalidRepos, nil
}

// ProviderName returns "bitbucket" as the provider name
func (c *BitbucketController) ProviderName(ctx *gin.Context) string {
	return "bitbucket"
}

// ProviderURL returns nil as there's no specific provider URL
func (c *BitbucketController) ProviderURL(ctx *gin.Context) string {
	return ""
}

// Private helper methods

func (c *BitbucketController) bitbucketImportEnabled() bool {
	// TODO: Implement check for Bitbucket import feature flag
	return true
}

func (c *BitbucketController) bitbucketAuth(ctx *gin.Context) {
	if ctx.GetString("bitbucket_token") == "" {
		c.goToBitbucketForPermissions(ctx)
	}
}

func (c *BitbucketController) goToBitbucketForPermissions(ctx *gin.Context) {
	state := generateRandomState()
	ctx.Set("bitbucket_auth_state", state)

	authURL := c.oauthService.GetAuthorizationURL(
		ctx,
		c.getCallbackURL(ctx),
		state,
	)
	ctx.Redirect(http.StatusFound, authURL)
}

func (c *BitbucketController) bitbucketUnauthorized(ctx *gin.Context, err error) {
	// TODO: Log the exception
	c.goToBitbucketForPermissions(ctx)
}

func (c *BitbucketController) getCredentials(ctx *gin.Context) *bitbucket.Credentials {
	return &bitbucket.Credentials{
		Token:        ctx.GetString("bitbucket_token"),
		ExpiresAt:    ctx.GetTime("bitbucket_expires_at"),
		ExpiresIn:    ctx.GetInt("bitbucket_expires_in"),
		RefreshToken: ctx.GetString("bitbucket_refresh_token"),
	}
}

func (c *BitbucketController) getCallbackURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.oauthService.GetCallbackURL(ctx, "bitbucket", namespaceID)
}

func (c *BitbucketController) getStatusURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.oauthService.GetStatusURL(ctx, "bitbucket", namespaceID)
}

func generateRandomState() string {
	b := make([]byte, 64)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func secureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
} 