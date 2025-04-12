package award

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// ToggleAwardEmojiHandler handles HTTP requests for toggling award emojis
type ToggleAwardEmojiHandler struct {
	awardEmojiService *service.AwardEmojiService
}

// NewToggleAwardEmojiHandler creates a new handler instance
func NewToggleAwardEmojiHandler(awardEmojiService *service.AwardEmojiService) *ToggleAwardEmojiHandler {
	return &ToggleAwardEmojiHandler{
		awardEmojiService: awardEmojiService,
	}
}

// RegisterRoutes registers the handler routes
func (h *ToggleAwardEmojiHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		award := api.Group("/award-emoji")
		{
			award.POST("/toggle", h.toggleAwardEmoji)
		}
	}
}

// toggleAwardEmoji handles POST /api/award-emoji/toggle
func (h *ToggleAwardEmojiHandler) toggleAwardEmoji(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "emoji name is required"})
		return
	}

	// Get the awardable from the context (set by middleware)
	awardable, exists := c.Get("awardable")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "awardable not found"})
		return
	}

	// Toggle the award emoji
	success, err := h.awardEmojiService.Toggle(c.Request.Context(), awardable, name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to toggle award emoji"})
		return
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"ok": false})
	}
}
