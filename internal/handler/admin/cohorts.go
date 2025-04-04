package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

type CohortsController struct {
	cohortsService *service.CohortsService
	analyticsService *service.AnalyticsService
}

func NewCohortsController(
	cohortsService *service.CohortsService,
	analyticsService *service.AnalyticsService,
) *CohortsController {
	return &CohortsController{
		cohortsService: cohortsService,
		analyticsService: analyticsService,
	}
}

// Index displays the list of cohorts
func (c *CohortsController) Index(ctx *gin.Context) {
	// Track analytics event
	if err := c.analyticsService.TrackEvent(ctx, service.AnalyticsEvent{
		Name: "i_analytics_cohorts",
		Action: "perform_analytics_usage_action",
		Label: "redis_hll_counters.analytics.analytics_total_unique_counts_monthly",
		Destinations: []string{"redis_hll", "snowplow"},
	}); err != nil {
		// Log error but continue with the request
		ctx.Error(err)
	}

	cohorts, err := c.loadCohorts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cohorts": cohorts,
	})
}

// loadCohorts retrieves and caches cohort data
func (c *CohortsController) loadCohorts(ctx *gin.Context) ([]service.Cohort, error) {
	// Try to get from cache first
	cacheKey := "cohorts"
	if cached, err := c.cohortsService.GetFromCache(ctx, cacheKey); err == nil {
		return cached, nil
	}

	// If not in cache, load from service
	cohorts, err := c.cohortsService.Execute(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the results for 24 hours
	if err := c.cohortsService.SetCache(ctx, cacheKey, cohorts, 24*time.Hour); err != nil {
		// Log cache error but don't fail the request
		ctx.Error(err)
	}

	return cohorts, nil
}
