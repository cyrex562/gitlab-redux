package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/analytics"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// UsageTrendsController handles usage trends and analytics tracking
type UsageTrendsController struct {
	analyticsService *service.AnalyticsService
}

// NewUsageTrendsController creates a new instance of UsageTrendsController
func NewUsageTrendsController(analyticsService *service.AnalyticsService) *UsageTrendsController {
	return &UsageTrendsController{
		analyticsService: analyticsService,
	}
}

// RegisterRoutes registers the routes for the UsageTrendsController
func (c *UsageTrendsController) RegisterRoutes(r *gin.RouterGroup) {
	usageTrends := r.Group("/admin/usage_trends")
	{
		usageTrends.Use(c.requireAdmin)
		usageTrends.GET("/", c.index)
	}
}

// requireAdmin middleware ensures that only admin users can access these endpoints
func (c *UsageTrendsController) requireAdmin(ctx *gin.Context) {
	user := ctx.MustGet("user")
	if user == nil || !user.IsAdmin() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// index handles the GET /admin/usage_trends endpoint
func (c *UsageTrendsController) index(ctx *gin.Context) {
	// Track analytics event
	event := analytics.Event{
		Name:        "i_analytics_instance_statistics",
		Action:      "perform_analytics_usage_action",
		Label:       "redis_hll_counters.analytics.analytics_total_unique_counts_monthly",
		Destinations: []string{"redis_hll", "snowplow"},
	}

	if err := c.analyticsService.TrackEvent(ctx, event); err != nil {
		// Log error but don't fail the request
		// TODO: Implement proper logging
	}

	// TODO: Implement HTML rendering for usage trends view
	ctx.JSON(http.StatusOK, gin.H{"message": "Usage trends view"})
}
