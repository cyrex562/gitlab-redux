package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/feature_flags"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/work_items"
)

// WorkItemsController handles requests for group work items
type WorkItemsController struct {
	featureFlagService *feature_flags.Service
	workItemsFinder    *work_items.WorkItemsFinder
}

// NewWorkItemsController creates a new work items controller
func NewWorkItemsController(
	featureFlagService *feature_flags.Service,
	workItemsFinder *work_items.WorkItemsFinder,
) *WorkItemsController {
	return &WorkItemsController{
		featureFlagService: featureFlagService,
		workItemsFinder:    workItemsFinder,
	}
}

// RegisterRoutes registers the routes for the work items controller
func (c *WorkItemsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/work_items", c.pushFeatureFlags(), c.Index)
	router.GET("/work_items/:iid", c.pushFeatureFlags(), c.handleNewWorkItemPath(), c.Show)
}

// Index handles GET requests for group work items
func (c *WorkItemsController) Index(ctx *gin.Context) {
	if !c.namespaceWorkItemsEnabled(ctx) {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Implementation will be added when needed
	ctx.Status(http.StatusOK)
}

// Show handles GET requests for a specific work item
func (c *WorkItemsController) Show(ctx *gin.Context) {
	if !c.namespaceWorkItemsEnabled(ctx) {
		ctx.Status(http.StatusNotFound)
		return
	}

	group := ctx.MustGet("group").(*models.Group)
	user := ctx.MustGet("current_user").(*models.User)
	iid := ctx.Param("iid")

	// Find the work item
	workItem, err := c.workItemsFinder.Execute(ctx, user, work_items.WorkItemsFinderParams{
		GroupID: group.ID,
	}).WithWorkItemType().FindByIID(iid)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Render the view
	ctx.HTML(http.StatusOK, "groups/work_items/show", gin.H{
		"work_item": workItem,
		"group":     group,
	})
}

// Helper methods for middleware and feature flags

func (c *WorkItemsController) pushFeatureFlags() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)

		// Push feature flags
		c.featureFlagService.PushFeatureFlag(ctx, "notifications_todos_buttons", user)
		c.featureFlagService.PushForceFeatureFlag(ctx, "work_items", group.WorkItemsFeatureFlagEnabled())
		c.featureFlagService.PushForceFeatureFlag(ctx, "work_items_beta", group.WorkItemsBetaFeatureFlagEnabled())
		c.featureFlagService.PushForceFeatureFlag(ctx, "work_items_alpha", group.WorkItemsAlphaFeatureFlagEnabled())
		c.featureFlagService.PushForceFeatureFlag(ctx, "namespace_level_work_items", c.namespaceWorkItemsEnabled(ctx))
		c.featureFlagService.PushForceFeatureFlag(ctx, "create_group_level_work_items", group.CreateGroupLevelWorkItemsFeatureFlagEnabled())
		c.featureFlagService.PushForceFeatureFlag(ctx, "glql_integration", group.GlqlIntegrationFeatureFlagEnabled())
		c.featureFlagService.PushForceFeatureFlag(ctx, "glql_load_on_click", group.GlqlLoadOnClickFeatureFlagEnabled())
		c.featureFlagService.PushForceFeatureFlag(ctx, "continue_indented_text", group.ContinueIndentedTextFeatureFlagEnabled())
		c.featureFlagService.PushFeatureFlag(ctx, "issues_list_drawer", group)

		ctx.Next()
	}
}

func (c *WorkItemsController) handleNewWorkItemPath() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		iid := ctx.Param("iid")
		if iid != "new" {
			ctx.Next()
			return
		}

		if c.namespaceWorkItemsEnabled(ctx) {
			// Render the show view for the new work item
			group := ctx.MustGet("group").(*models.Group)
			ctx.HTML(http.StatusOK, "groups/work_items/show", gin.H{
				"group": group,
			})
		} else {
			ctx.Status(http.StatusNotFound)
		}
	}
}

// Helper methods for feature flag checks

func (c *WorkItemsController) namespaceWorkItemsEnabled(ctx *gin.Context) bool {
	group := ctx.MustGet("group").(*models.Group)
	return group.NamespaceWorkItemsEnabled()
} 