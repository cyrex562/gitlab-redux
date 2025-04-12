package events

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// FiltersEvents handles event filtering functionality
type FiltersEvents struct {
	eventFilterService *service.EventFilterService
}

// NewFiltersEvents creates a new instance of FiltersEvents
func NewFiltersEvents(eventFilterService *service.EventFilterService) *FiltersEvents {
	return &FiltersEvents{
		eventFilterService: eventFilterService,
	}
}

// GetEventFilter gets the current event filter from the context or creates a new one
func (f *FiltersEvents) GetEventFilter(ctx *gin.Context) *model.EventFilter {
	// Check if the event filter is already in the context
	if eventFilter, exists := ctx.Get("event_filter"); exists {
		return eventFilter.(*model.EventFilter)
	}

	// Create a new event filter
	eventFilter := f.newEventFilter(ctx)

	// Store the event filter in the context
	ctx.Set("event_filter", eventFilter)

	// Set the event filter in the cookie
	ctx.SetCookie("event_filter", eventFilter.Filter, 0, "/", "", false, true)

	return eventFilter
}

// newEventFilter creates a new event filter based on the request parameters or cookie
func (f *FiltersEvents) newEventFilter(ctx *gin.Context) *model.EventFilter {
	// Get the active filter from the request parameters or cookie
	activeFilter := ctx.Query("event_filter")
	if activeFilter == "" {
		// Try to get the filter from the cookie
		if cookie, err := ctx.Cookie("event_filter"); err == nil {
			activeFilter = cookie
		}
	}

	// Create a new event filter
	return f.eventFilterService.NewEventFilter(activeFilter)
}
