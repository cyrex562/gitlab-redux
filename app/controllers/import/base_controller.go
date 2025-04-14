package import_controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/rate_limiter"
)

// BaseController handles base import functionality
type BaseController struct {
	rateLimiter *rate_limiter.Service
}

// NewBaseController creates a new base controller
func NewBaseController(rateLimiter *rate_limiter.Service) *BaseController {
	return &BaseController{
		rateLimiter: rateLimiter,
	}
}

// RegisterRoutes registers the routes for the base controller
func (c *BaseController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/status", c.Status)
	router.GET("/realtime_changes", c.RealtimeChanges)
}

// Status handles GET requests for import status
func (c *BaseController) Status(ctx *gin.Context) {
	format := ctx.DefaultQuery("format", "json")
	if format == "json" {
		ctx.JSON(http.StatusOK, gin.H{
			"imported_projects":    c.serializedImportedProjects(ctx),
			"provider_repos":       c.serializedProviderRepos(ctx),
			"incompatible_repos":   c.serializedIncompatibleRepos(ctx),
		})
		return
	}

	// HTML format
	namespaceID := ctx.Query("namespace_id")
	if namespaceID != "" {
		namespace, err := c.findNamespace(ctx, namespaceID)
		if err != nil {
			ctx.HTML(http.StatusNotFound, "errors/404", nil)
			return
		}

		user := ctx.MustGet("current_user").(*models.User)
		if !user.CanImportProjects(namespace) {
			ctx.HTML(http.StatusNotFound, "errors/404", nil)
			return
		}
	}

	ctx.HTML(http.StatusOK, "import/status", nil)
}

// RealtimeChanges handles GET requests for realtime changes
func (c *BaseController) RealtimeChanges(ctx *gin.Context) {
	// Set polling interval header (3 seconds)
	ctx.Header("X-Poll-Interval", "3000")

	projects := c.alreadyAddedProjects(ctx)
	ctx.JSON(http.StatusOK, projects)
}

// Protected methods that should be implemented by child controllers

// ImportableRepos returns the list of importable repositories
func (c *BaseController) ImportableRepos(ctx *gin.Context) ([]interface{}, error) {
	panic("not implemented")
}

// IncompatibleRepos returns the list of incompatible repositories
func (c *BaseController) IncompatibleRepos(ctx *gin.Context) ([]interface{}, error) {
	panic("not implemented")
}

// ProviderName returns the name of the provider
func (c *BaseController) ProviderName(ctx *gin.Context) string {
	panic("not implemented")
}

// ProviderURL returns the URL of the provider
func (c *BaseController) ProviderURL(ctx *gin.Context) string {
	panic("not implemented")
}

// ExtraRepresentationOpts returns additional options for representation
func (c *BaseController) ExtraRepresentationOpts(ctx *gin.Context) map[string]interface{} {
	return make(map[string]interface{})
}

// Private helper methods

func (c *BaseController) sanitizedFilterParam(ctx *gin.Context) string {
	filter := ctx.Query("filter")
	if filter == "" {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(filter))
}

func (c *BaseController) filtered(collection []interface{}, filter string) []interface{} {
	if filter == "" {
		return collection
	}

	var filtered []interface{}
	for _, item := range collection {
		if repo, ok := item.(map[string]interface{}); ok {
			if name, ok := repo["name"].(string); ok {
				if strings.Contains(strings.ToLower(name), filter) {
					filtered = append(filtered, item)
				}
			}
		}
	}
	return filtered
}

func (c *BaseController) serializedProviderRepos(ctx *gin.Context) []interface{} {
	repos, _ := c.ImportableRepos(ctx)
	return c.filtered(repos, c.sanitizedFilterParam(ctx))
}

func (c *BaseController) serializedIncompatibleRepos(ctx *gin.Context) []interface{} {
	repos, _ := c.IncompatibleRepos(ctx)
	return c.filtered(repos, c.sanitizedFilterParam(ctx))
}

func (c *BaseController) serializedImportedProjects(ctx *gin.Context) []interface{} {
	projects := c.alreadyAddedProjects(ctx)
	// TODO: Implement project serialization
	return projects
}

func (c *BaseController) alreadyAddedProjects(ctx *gin.Context) []interface{} {
	user := ctx.MustGet("current_user").(*models.User)
	return c.findAlreadyAddedProjects(ctx, user, c.ProviderName(ctx))
}

func (c *BaseController) findAlreadyAddedProjects(ctx *gin.Context, user *models.User, importType string) []interface{} {
	// TODO: Implement finding already added projects
	return []interface{}{}
}

func (c *BaseController) findNamespace(ctx *gin.Context, namespaceID string) (*models.Namespace, error) {
	// TODO: Implement namespace finding
	return nil, nil
} 