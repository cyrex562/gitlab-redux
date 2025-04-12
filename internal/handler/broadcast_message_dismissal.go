package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// BroadcastMessageDismissalHandler handles HTTP requests for broadcast message dismissals
type BroadcastMessageDismissalHandler struct {
	dismissalService *service.BroadcastMessageDismissalService
}

// NewBroadcastMessageDismissalHandler creates a new handler instance
func NewBroadcastMessageDismissalHandler(dismissalService *service.BroadcastMessageDismissalService) *BroadcastMessageDismissalHandler {
	return &BroadcastMessageDismissalHandler{
		dismissalService: dismissalService,
	}
}

// RegisterRoutes registers the handler routes
func (h *BroadcastMessageDismissalHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		broadcast := api.Group("/broadcast-messages")
		{
			broadcast.POST("/:messageID/dismiss", h.createDismissal)
			broadcast.GET("/:messageID/dismissed", h.isDismissed)
		}
	}
}

// createDismissal handles POST /api/broadcast-messages/:messageID/dismiss
func (h *BroadcastMessageDismissalHandler) createDismissal(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	messageID, err := strconv.ParseInt(c.Param("messageID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
		return
	}

	dismissal, err := h.dismissalService.CreateDismissal(c.Request.Context(), userID, messageID)
	if err != nil {
		if err == service.ErrInvalidMessageID {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dismissal"})
		return
	}

	// Set cookie with 30-day expiration
	c.SetCookie(
		dismissal.CookieKey(),
		"true",
		int(time.Hour * 24 * 30),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusCreated, dismissal)
}

// isDismissed handles GET /api/broadcast-messages/:messageID/dismissed
func (h *BroadcastMessageDismissalHandler) isDismissed(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	messageID, err := strconv.ParseInt(c.Param("messageID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
		return
	}

	isDismissed, err := h.dismissalService.IsDismissed(c.Request.Context(), userID, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check dismissal status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dismissed": isDismissed})
}
