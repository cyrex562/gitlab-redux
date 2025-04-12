package event_forward

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/feature_flags"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// EventForwardController handles forwarding events for product usage tracking
type EventForwardController struct {
	analyticsService   service.AnalyticsService
	featureFlagService service.FeatureFlagService
	logger             *Logger
}

// NewEventForwardController creates a new EventForwardController
func NewEventForwardController(
	analyticsService service.AnalyticsService,
	featureFlagService service.FeatureFlagService,
	baseLogger service.Logger,
) *EventForwardController {
	return &EventForwardController{
		analyticsService:   analyticsService,
		featureFlagService: featureFlagService,
		logger:             NewLogger(baseLogger),
	}
}

// RegisterRoutes registers the routes for the EventForwardController
func (c *EventForwardController) RegisterRoutes(router *gin.Engine) {
	eventForward := router.Group("/event_forward")
	{
		eventForward.POST("/forward", c.Forward)
	}
}

// Forward handles the forward action for event forwarding
func (c *EventForwardController) Forward(ctx *gin.Context) {
	// Check if feature flag is enabled
	featureEnabled := feature_flags.PushFrontendFeatureFlag(ctx, "collect_product_usage_events", nil)
	if !featureEnabled {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Process events
	err := c.processEvents(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// processEvents processes the events from the request payload
func (c *EventForwardController) processEvents(ctx *gin.Context) error {
	// Read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	// Parse payload
	var payload struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}

	// Process each event
	for _, event := range payload.Data {
		err := c.analyticsService.TrackInternalEvent(ctx, event["name"].(string), nil)
		if err != nil {
			return err
		}
	}

	// Log the number of events processed
	c.logger.Info("Enqueued events for forwarding", "count", len(payload.Data))

	return nil
}
