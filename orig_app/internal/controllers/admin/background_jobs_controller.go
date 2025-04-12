package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitlab-org/gitlab-redux/internal/controllers"
	"github.com/gitlab-org/gitlab-redux/internal/services"
)

// BackgroundJobsController handles background job management in the admin interface
type BackgroundJobsController struct {
	controllers.BaseController
	sidekiqService *services.SidekiqService
}

// NewBackgroundJobsController creates a new instance of BackgroundJobsController
func NewBackgroundJobsController(sidekiqService *services.SidekiqService) *BackgroundJobsController {
	return &BackgroundJobsController{
		sidekiqService: sidekiqService,
	}
}

// RegisterRoutes registers the routes for this controller
func (c *BackgroundJobsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/background_jobs", c.RequireAdmin, c.Show)
}

// Show displays the background jobs dashboard
func (c *BackgroundJobsController) Show(ctx *gin.Context) {
	// Get Sidekiq statistics
	stats, err := c.sidekiqService.GetStats()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get queue information
	queues, err := c.sidekiqService.GetQueues()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get worker information
	workers, err := c.sidekiqService.GetWorkers()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get scheduled jobs
	scheduled, err := c.sidekiqService.GetScheduledJobs()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get retries
	retries, err := c.sidekiqService.GetRetries()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get dead jobs
	dead, err := c.sidekiqService.GetDeadJobs()
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Track the page load event
	c.TrackPageLoad(ctx, "view_admin_background_jobs_pageload")

	ctx.JSON(http.StatusOK, gin.H{
		"stats":     stats,
		"queues":    queues,
		"workers":   workers,
		"scheduled": scheduled,
		"retries":   retries,
		"dead":      dead,
	})
}

// RequireAdmin is a middleware that ensures the user has admin privileges
func (c *BackgroundJobsController) RequireAdmin(ctx *gin.Context) {
	// TODO: Implement proper admin authorization check
	// This should check if the user has the :read_admin_background_jobs permission
	ctx.Next()
}

// TrackPageLoad tracks a page load event
func (c *BackgroundJobsController) TrackPageLoad(ctx *gin.Context, eventName string) {
	// TODO: Implement proper event tracking
	// This should track the page load event for analytics
}
