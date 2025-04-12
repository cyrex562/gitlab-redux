package analytics

import (
	"github.com/gin-gonic/gin"
)

// GetEventFilter gets the event filter from the context
func GetEventFilter(ctx *gin.Context) string {
	filter := ctx.DefaultQuery("filter", "")
	if filter == "" {
		filter = ctx.DefaultQuery("event_filter", "")
	}
	return filter
}

// SetEventFilter sets the event filter in the context
func SetEventFilter(ctx *gin.Context, defaultFilter string) string {
	filter := GetEventFilter(ctx)
	if filter == "" {
		filter = defaultFilter
	}
	ctx.Set("event_filter", filter)
	return filter
}

// GetEventFilterOptions gets the event filter options
func GetEventFilterOptions() map[string]string {
	return map[string]string{
		"all":      "All Events",
		"push":     "Push Events",
		"merged":   "Merge Events",
		"issue":    "Issue Events",
		"note":     "Note Events",
		"wiki":     "Wiki Events",
		"pipeline": "Pipeline Events",
		"build":    "Build Events",
	}
}
