package import_controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/bitbucket_server"
)

// BitbucketServerController handles Bitbucket Server import functionality
type BitbucketServerController struct {
	*BaseController
	bitbucketServerService *bitbucket_server.Service
}

// NewBitbucketServerController creates a new Bitbucket Server controller
func NewBitbucketServerController(
	baseController *BaseController,
	bitbucketServerService *bitbucket_server.Service,
) *BitbucketServerController {
	return &BitbucketServerController{
		BaseController:         baseController,
		bitbucketServerService: bitbucketServerService,
	}
}

// RegisterRoutes registers the routes for the Bitbucket Server controller
func (c *BitbucketServerController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/new", c.New)
	router.POST("/create", c.Create)
	router.POST("/configure", c.Configure)
	router.GET("/status", c.Status)
}

// New handles the new import page
func (c *BitbucketServerController) New(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "import/bitbucket_server/new", nil)
}

// Create handles project creation from Bitbucket Server
func (c *BitbucketServerController) Create(ctx *gin.Context) {
	if !c.bitbucketServerImportEnabled() {
		ctx.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	if err := c.normalizeImportParams(ctx); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	if err := c.validateImportParams(ctx); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	projectKey := ctx.PostForm("bitbucket_server_project")
	repoSlug := ctx.PostForm("bitbucket_server_repo")

	repo, err := c.bitbucketServerService.GetRepo(ctx, projectKey, repoSlug)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": c.safeFormat("Project %s/%s could not be found", projectKey, repoSlug),
		})
		return
	}

	result, err := c.bitbucketServerService.ImportProject(ctx, repo, c.getCredentials(ctx))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result.Project)
}

// Configure handles Bitbucket Server configuration
func (c *BitbucketServerController) Configure(ctx *gin.Context) {
	ctx.Set("bitbucket_server_personal_access_token", ctx.PostForm("personal_access_token"))
	ctx.Set("bitbucket_server_username", ctx.PostForm("bitbucket_server_username"))
	ctx.Set("bitbucket_server_url", ctx.PostForm("bitbucket_server_url"))

	ctx.Redirect(http.StatusFound, c.getStatusURL(ctx))
}

// ImportableRepos returns the list of importable Bitbucket Server repositories
func (c *BitbucketServerController) ImportableRepos(ctx *gin.Context) ([]interface{}, error) {
	repos, err := c.bitbucketServerService.GetRepos(
		ctx,
		c.getPageOffset(ctx),
		c.getLimitPerPage(),
		c.sanitizedFilterParam(ctx),
	)
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

// IncompatibleRepos returns the list of incompatible Bitbucket Server repositories
func (c *BitbucketServerController) IncompatibleRepos(ctx *gin.Context) ([]interface{}, error) {
	repos, err := c.bitbucketServerService.GetRepos(
		ctx,
		c.getPageOffset(ctx),
		c.getLimitPerPage(),
		c.sanitizedFilterParam(ctx),
	)
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

// ProviderName returns "bitbucket_server" as the provider name
func (c *BitbucketServerController) ProviderName(ctx *gin.Context) string {
	return "bitbucket_server"
}

// ProviderURL returns the Bitbucket Server URL from session
func (c *BitbucketServerController) ProviderURL(ctx *gin.Context) string {
	return ctx.GetString("bitbucket_server_url")
}

// Private helper methods

func (c *BitbucketServerController) bitbucketServerImportEnabled() bool {
	// TODO: Implement check for Bitbucket Server import feature flag
	return true
}

func (c *BitbucketServerController) bitbucketAuth(ctx *gin.Context) {
	if ctx.GetString("bitbucket_server_url") == "" ||
		ctx.GetString("bitbucket_server_username") == "" ||
		ctx.GetString("bitbucket_server_personal_access_token") == "" {
		ctx.Redirect(http.StatusFound, c.getNewURL(ctx))
	}
}

func (c *BitbucketServerController) normalizeImportParams(ctx *gin.Context) error {
	repoID := ctx.PostForm("repo_id")
	parts := strings.Split(repoID, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repo_id format")
	}

	ctx.Set("bitbucket_server_project", parts[0])
	ctx.Set("bitbucket_server_repo", parts[1])
	return nil
}

func (c *BitbucketServerController) validateImportParams(ctx *gin.Context) error {
	projectKey := ctx.GetString("bitbucket_server_project")
	repoSlug := ctx.GetString("bitbucket_server_repo")

	if projectKey == "" || repoSlug == "" {
		return fmt.Errorf("missing project key or repository slug")
	}

	validProjectChars := regexp.MustCompile(`^~?[\w\-\.\s]+$`)
	validRepoChars := regexp.MustCompile(`^[\w\-\.\s]+$`)

	if !validProjectChars.MatchString(projectKey) {
		return fmt.Errorf("invalid project key")
	}

	if !validRepoChars.MatchString(repoSlug) {
		return fmt.Errorf("invalid repository slug")
	}

	return nil
}

func (c *BitbucketServerController) clearSessionData(ctx *gin.Context) {
	ctx.Set("bitbucket_server_url", "")
	ctx.Set("bitbucket_server_username", "")
	ctx.Set("bitbucket_server_personal_access_token", "")
}

func (c *BitbucketServerController) getCredentials(ctx *gin.Context) *bitbucket_server.Credentials {
	return &bitbucket_server.Credentials{
		BaseURI:  ctx.GetString("bitbucket_server_url"),
		Username: ctx.GetString("bitbucket_server_username"),
		Password: ctx.GetString("bitbucket_server_personal_access_token"),
	}
}

func (c *BitbucketServerController) getPageOffset(ctx *gin.Context) int {
	page := ctx.DefaultQuery("page", "0")
	offset, _ := strconv.Atoi(page)
	if offset < 0 {
		return 0
	}
	return offset
}

func (c *BitbucketServerController) getLimitPerPage() int {
	return 25 // TODO: Make this configurable
}

func (c *BitbucketServerController) getNewURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.bitbucketServerService.GetNewURL(ctx, "bitbucket_server", namespaceID)
}

func (c *BitbucketServerController) getStatusURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.bitbucketServerService.GetStatusURL(ctx, "bitbucket_server", namespaceID)
}

func (c *BitbucketServerController) safeFormat(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (c *BitbucketServerController) handleConnectionError(ctx *gin.Context, err error) {
	errorMsg := c.safeFormat("Unable to connect to server: %v", err)
	c.clearSessionData(ctx)

	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"error": gin.H{
			"message":  errorMsg,
			"redirect": c.getNewURL(ctx),
		},
	})
} 