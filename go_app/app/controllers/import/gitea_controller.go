package import_controller

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/gitea"
)

// GiteaController handles Gitea import functionality
type GiteaController struct {
	*GithubController
	giteaService *gitea.Service
}

// NewGiteaController creates a new Gitea controller
func NewGiteaController(
	githubController *GithubController,
	giteaService *gitea.Service,
) *GiteaController {
	return &GiteaController{
		GithubController: githubController,
		giteaService:    giteaService,
	}
}

// RegisterRoutes registers the routes for the Gitea controller
func (c *GiteaController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/new", c.New)
	router.POST("/personal_access_token", c.PersonalAccessToken)
	router.GET("/status", c.Status)
}

// New handles the new import page
func (c *GiteaController) New(ctx *gin.Context) {
	if c.hasAccessToken(ctx) && c.getProviderURL(ctx) != "" {
		ctx.Redirect(http.StatusFound, c.getStatusURL(ctx))
		return
	}

	ctx.HTML(http.StatusOK, "import/gitea/new", nil)
}

// PersonalAccessToken handles personal access token submission
func (c *GiteaController) PersonalAccessToken(ctx *gin.Context) {
	ctx.Set(c.getHostKey(), ctx.PostForm(c.getHostKey()))
	c.GithubController.PersonalAccessToken(ctx)
}

// Status handles the import status page
func (c *GiteaController) Status(ctx *gin.Context) {
	// Request repos to display error page if provider token is invalid
	// Improving in https://gitlab.com/gitlab-org/gitlab-redux/-/issues/25859
	_, err := c.getClientRepos(ctx)
	if err != nil {
		ctx.HTML(http.StatusOK, "import/gitea/error", gin.H{
			"error": err.Error(),
		})
		return
	}

	if ctx.GetHeader("Accept") == "application/json" {
		importedProjects, err := c.getSerializedImportedProjects(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		providerRepos, err := c.getSerializedProviderRepos(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		incompatibleRepos, err := c.getSerializedIncompatibleRepos(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"imported_projects":  importedProjects,
			"provider_repos":     providerRepos,
			"incompatible_repos": incompatibleRepos,
		})
		return
	}

	namespaceID := ctx.Query("namespace_id")
	if namespaceID != "" {
		namespace, err := c.getNamespace(ctx, namespaceID)
		if err != nil {
			ctx.HTML(http.StatusNotFound, "errors/404", nil)
			return
		}

		user := ctx.MustGet("current_user").(*models.User)
		if !c.canImportProjects(ctx, user, namespace) {
			ctx.HTML(http.StatusNotFound, "errors/404", nil)
			return
		}
	}

	ctx.HTML(http.StatusOK, "import/gitea/status", nil)
}

// ProviderName returns "gitea" as the provider name
func (c *GiteaController) ProviderName(ctx *gin.Context) string {
	return "gitea"
}

// ProviderURL returns the Gitea host URL from session
func (c *GiteaController) ProviderURL(ctx *gin.Context) string {
	return ctx.GetString(c.getHostKey())
}

// LoggedInWithProvider returns false as Gitea is not yet an OAuth provider
func (c *GiteaController) LoggedInWithProvider(ctx *gin.Context) bool {
	return false
}

// ProviderAuth checks if both access token and host URL are present
func (c *GiteaController) ProviderAuth(ctx *gin.Context) {
	if !c.hasAccessToken(ctx) || c.getProviderURL(ctx) == "" {
		ctx.Redirect(http.StatusFound, c.getNewURL(ctx))
		ctx.Set("alert", "You need to specify both an access token and a Host URL.")
		return
	}
}

// GetSerializedImportedProjects returns serialized imported projects
func (c *GiteaController) GetSerializedImportedProjects(ctx *gin.Context, projects []*models.Project) ([]interface{}, error) {
	return c.giteaService.SerializeProjects(ctx, projects, c.getProviderURL(ctx))
}

// GetClientRepos returns filtered client repositories
func (c *GiteaController) GetClientRepos(ctx *gin.Context) ([]interface{}, error) {
	client, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	repos, err := client.GetRepos(ctx)
	if err != nil {
		return nil, err
	}

	return c.filterRepos(ctx, repos)
}

// Private helper methods

func (c *GiteaController) getHostKey() string {
	return fmt.Sprintf("%s_host_url", c.ProviderName(nil))
}

func (c *GiteaController) getClient(ctx *gin.Context) (*gitea.Client, error) {
	options, err := c.getClientOptions(ctx)
	if err != nil {
		return nil, err
	}

	return c.giteaService.NewClient(ctx, c.getAccessToken(ctx), options)
}

func (c *GiteaController) getClientOptions(ctx *gin.Context) (*gitea.ClientOptions, error) {
	verifiedURL, providerHostname, err := c.verifyBlockedURI(ctx)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(verifiedURL)
	if err != nil {
		return nil, err
	}

	return &gitea.ClientOptions{
		Host:       parsedURL.Scheme == "https" ? c.getProviderURL(ctx) : verifiedURL,
		APIVersion: "v1",
		Hostname:   providerHostname,
	}, nil
}

func (c *GiteaController) verifyBlockedURI(ctx *gin.Context) (string, string, error) {
	return c.giteaService.ValidateURL(ctx, c.getProviderURL(ctx), c.allowLocalRequests(ctx))
}

func (c *GiteaController) allowLocalRequests(ctx *gin.Context) bool {
	// TODO: Implement settings check
	return true
}

func (c *GiteaController) hasAccessToken(ctx *gin.Context) bool {
	return ctx.GetString(c.getAccessTokenKey()) != ""
}

func (c *GiteaController) getAccessToken(ctx *gin.Context) string {
	return ctx.GetString(c.getAccessTokenKey())
}

func (c *GiteaController) getAccessTokenKey() string {
	return fmt.Sprintf("%s_access_token", c.ProviderName(nil))
}

func (c *GiteaController) getNewURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.giteaService.GetNewURL(ctx, "gitea", namespaceID)
}

func (c *GiteaController) getStatusURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.giteaService.GetStatusURL(ctx, "gitea", namespaceID)
}

func (c *GiteaController) getNamespace(ctx *gin.Context, id string) (*models.Namespace, error) {
	// TODO: Implement namespace lookup
	return nil, nil
}

func (c *GiteaController) canImportProjects(ctx *gin.Context, user *models.User, namespace *models.Namespace) bool {
	// TODO: Implement permission check
	return true
}

func (c *GiteaController) filterRepos(ctx *gin.Context, repos []interface{}) ([]interface{}, error) {
	// TODO: Implement repository filtering
	return repos, nil
} 