package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/clusters"
)

// ClustersController extends the base ClustersController and adds group-specific functionality
type ClustersController struct {
	// Embed the base ClustersController to inherit its functionality
	*clusters.ClustersController

	// Add any additional dependencies here
	featureFlagService *FeatureFlagService
	clusterablePresenter *ClusterablePresenter
	groupFinder *GroupFinder
}

// NewClustersController creates a new ClustersController
func NewClustersController(
	baseController *clusters.ClustersController,
	featureFlagService *FeatureFlagService,
	clusterablePresenter *ClusterablePresenter,
	groupFinder *GroupFinder,
) *ClustersController {
	return &ClustersController{
		ClustersController: baseController,
		featureFlagService: featureFlagService,
		clusterablePresenter: clusterablePresenter,
		groupFinder: groupFinder,
	}
}

// RegisterMiddleware registers the middleware for the ClustersController
func (c *ClustersController) RegisterMiddleware(router *gin.RouterGroup) {
	// Register the middleware from the base ClustersController
	c.ClustersController.RegisterMiddleware(router)

	// Add cross project access check middleware
	router.Use(c.RequiresCrossProjectAccess())

	// Add feature flag middleware for all actions except index
	router.Use(c.EnsureFeatureEnabledExceptIndex())
}

// EnsureFeatureEnabledExceptIndex middleware ensures the feature is enabled for all actions except index
func (c *ClustersController) EnsureFeatureEnabledExceptIndex() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip the check for the index action
		if ctx.Request.Method == "GET" && ctx.Request.URL.Path == "/" {
			ctx.Next()
			return
		}

		// Check if the feature is enabled
		if !c.featureFlagService.IsEnabled("kubernetes_clusters", c.GetCurrentUser(ctx)) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Feature not found"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// RequiresCrossProjectAccess middleware requires cross project access
func (c *ClustersController) RequiresCrossProjectAccess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set the cross project access flag
		ctx.Set("requires_cross_project_access", true)
		ctx.Next()
	}
}

// GetClusterable gets the clusterable object
func (c *ClustersController) GetClusterable(ctx *gin.Context) interface{} {
	// Get the group from the context
	group := c.GetGroup(ctx)
	if group == nil {
		return nil
	}

	// Get the current user from the context
	currentUser := c.GetCurrentUser(ctx)

	// Fabricate the clusterable presenter
	return c.clusterablePresenter.Fabricate(group, currentUser)
}

// GetGroup gets the group from the request
func (c *ClustersController) GetGroup(ctx *gin.Context) interface{} {
	// Get the group ID from the request
	groupID := ctx.Param("group_id")
	if groupID == "" {
		groupID = ctx.Param("id")
	}

	// Find the group
	group := c.groupFinder.Execute(c.GetCurrentUser(ctx), map[string]interface{}{
		"id": groupID,
	})

	// Set the group in the context
	ctx.Set("group", group)

	return group
}

// GetCurrentUser gets the current user from the context
func (c *ClustersController) GetCurrentUser(ctx *gin.Context) interface{} {
	currentUser, _ := ctx.Get("current_user")
	return currentUser
}

// FeatureFlagService handles feature flags
type FeatureFlagService struct {
	// Add any dependencies here
}

// IsEnabled checks if a feature is enabled
func (s *FeatureFlagService) IsEnabled(feature string, currentUser interface{}) bool {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would check if a feature is enabled
	return true
}

// ClusterablePresenter presents a clusterable object
type ClusterablePresenter struct {
	// Add any dependencies here
}

// Fabricate fabricates a clusterable presenter
func (p *ClusterablePresenter) Fabricate(clusterable interface{}, currentUser interface{}) interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would fabricate a clusterable presenter
	return clusterable
}

// GroupFinder finds groups
type GroupFinder struct {
	// Add any dependencies here
}

// Execute executes the finder
func (f *GroupFinder) Execute(currentUser interface{}, params map[string]interface{}) interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would find a group
	return nil
}
