package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/feature_flags"
)

// UsageQuotasController handles requests for group usage quotas
type UsageQuotasController struct {
	featureFlagService *feature_flags.Service
}

// NewUsageQuotasController creates a new usage quotas controller
func NewUsageQuotasController(featureFlagService *feature_flags.Service) *UsageQuotasController {
	return &UsageQuotasController{
		featureFlagService: featureFlagService,
	}
}

// RegisterRoutes registers the routes for the usage quotas controller
func (c *UsageQuotasController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/usage_quotas", c.authorizeReadUsageQuotas(), c.verifyUsageQuotasEnabled(), c.pushFeatureFlags(), c.Root)
}

// Root handles GET requests for the usage quotas root page
func (c *UsageQuotasController) Root(ctx *gin.Context) {
	group := ctx.MustGet("group").(*models.Group)
	user := ctx.MustGet("current_user").(*models.User)

	// Get seat count data
	seatCountData := c.getSeatCountData(ctx, group, user)

	// Render the view
	ctx.HTML(http.StatusOK, "groups/usage_quotas/root", gin.H{
		"group":          group,
		"seat_count_data": seatCountData,
	})
}

// Helper methods for middleware and authorization

func (c *UsageQuotasController) authorizeReadUsageQuotas() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)

		if !user.CanReadUsageQuotas(group) {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c *UsageQuotasController) verifyUsageQuotasEnabled() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		group := ctx.MustGet("group").(*models.Group)

		if !group.UsageQuotasEnabled() {
			ctx.Status(http.StatusNotFound)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c *UsageQuotasController) pushFeatureFlags() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)

		// Push the virtual registry maven feature flag
		c.featureFlagService.PushFeatureFlag(ctx, "virtual_registry_maven", user)

		ctx.Next()
	}
}

// getSeatCountData returns the seat count data for the group
// This method can be overridden in enterprise editions
func (c *UsageQuotasController) getSeatCountData(ctx *gin.Context, group *models.Group, user *models.User) interface{} {
	// TODO: Implement the actual seat count data logic
	return nil
} 