package explore

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/analytics"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// CatalogController handles catalog-related actions
type CatalogController struct {
	*ApplicationController
	catalogService service.CatalogService
	analyticsService service.AnalyticsService
}

// NewCatalogController creates a new CatalogController
func NewCatalogController(
	appController *ApplicationController,
	catalogService service.CatalogService,
	analyticsService service.AnalyticsService,
) *CatalogController {
	return &CatalogController{
		ApplicationController: appController,
		catalogService:       catalogService,
		analyticsService:     analyticsService,
	}
}

// RegisterRoutes registers the routes for the CatalogController
func (c *CatalogController) RegisterRoutes(router *gin.Engine) {
	explore := router.Group("/explore")
	{
		catalog := explore.Group("/catalog")
		{
			catalog.GET("", c.Index)
			catalog.GET("/:full_path", c.Show)
		}
	}
}

// Show handles the show action
func (c *CatalogController) Show(ctx *gin.Context) {
	// Check resource access
	if !c.checkResourceAccess(ctx) {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Render the show template
	ctx.HTML(http.StatusOK, "explore/catalog/show", gin.H{
		"layout": c.GetLayout(),
	})
}

// Index handles the index action
func (c *CatalogController) Index(ctx *gin.Context) {
	// Track internal event
	analytics.TrackInternalEvent(ctx, "unique_users_visiting_ci_catalog", nil)

	// Render the show template
	ctx.HTML(http.StatusOK, "explore/catalog/show", gin.H{
		"layout": c.GetLayout(),
	})
}

// checkResourceAccess checks if the resource exists
func (c *CatalogController) checkResourceAccess(ctx *gin.Context) bool {
	fullPath := ctx.Param("full_path")
	resource := c.catalogResource(ctx, fullPath)
	return resource != nil
}

// catalogResource gets the catalog resource
func (c *CatalogController) catalogResource(ctx *gin.Context, fullPath string) *model.CatalogResource {
	user := ctx.MustGet("current_user").(*model.User)
	return c.catalogService.FindResource(ctx, user, fullPath)
}

// TrackingNamespaceSource returns the namespace source for tracking
func (c *CatalogController) TrackingNamespaceSource(ctx *gin.Context) *model.Namespace {
	user := ctx.MustGet("current_user").(*model.User)
	return user.Namespace
}

// TrackingProjectSource returns the project source for tracking
func (c *CatalogController) TrackingProjectSource(ctx *gin.Context) *model.Project {
	return nil
}
