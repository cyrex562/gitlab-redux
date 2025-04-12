package internal_events

import (
	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// TrackInternalEvent tracks an internal event
func TrackInternalEvent(ctx *gin.Context, eventName string, user *model.User) {
	// Get analytics service from context
	analyticsService, ok := ctx.MustGet("analytics_service").(service.AnalyticsService)
	if !ok {
		return
	}

	// Track the event
	analyticsService.TrackInternalEvent(ctx, eventName, user)
}

// TrackUserEvent tracks a user event
func TrackUserEvent(ctx *gin.Context, eventName string, user *model.User) {
	// Get analytics service from context
	analyticsService, ok := ctx.MustGet("analytics_service").(service.AnalyticsService)
	if !ok {
		return
	}

	// Track the event
	analyticsService.TrackUserEvent(ctx, eventName, user)
}

// TrackGroupEvent tracks a group event
func TrackGroupEvent(ctx *gin.Context, eventName string, group *model.Group) {
	// Get analytics service from context
	analyticsService, ok := ctx.MustGet("analytics_service").(service.AnalyticsService)
	if !ok {
		return
	}

	// Track the event
	analyticsService.TrackGroupEvent(ctx, eventName, group)
}

// TrackProjectEvent tracks a project event
func TrackProjectEvent(ctx *gin.Context, eventName string, project *model.Project) {
	// Get analytics service from context
	analyticsService, ok := ctx.MustGet("analytics_service").(service.AnalyticsService)
	if !ok {
		return
	}

	// Track the event
	analyticsService.TrackProjectEvent(ctx, eventName, project)
}
