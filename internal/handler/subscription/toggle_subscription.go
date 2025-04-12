package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// ToggleSubscriptionHandler handles HTTP requests for toggling subscriptions
type ToggleSubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

// NewToggleSubscriptionHandler creates a new handler instance
func NewToggleSubscriptionHandler(subscriptionService *service.SubscriptionService) *ToggleSubscriptionHandler {
	return &ToggleSubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// RegisterRoutes registers the handler routes
func (h *ToggleSubscriptionHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		subscription := api.Group("/subscription")
		{
			subscription.POST("/toggle", h.toggleSubscription)
		}
	}
}

// toggleSubscription handles POST /api/subscription/toggle
func (h *ToggleSubscriptionHandler) toggleSubscription(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Get the subscribable resource from the context (set by middleware)
	subscribable, exists := c.Get("subscribable")
	if !exists {
		c.Status(http.StatusBadRequest)
		return
	}

	// Get the project from the context (set by middleware)
	project, exists := c.Get("project")
	if !exists {
		c.Status(http.StatusBadRequest)
		return
	}

	// Toggle the subscription
	err := h.subscriptionService.ToggleSubscription(c.Request.Context(), subscribable, project, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
