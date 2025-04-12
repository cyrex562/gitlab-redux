package groups

import (
	"github.com/gin-gonic/gin"
)

// CustomEmojiController handles custom emoji for groups
type CustomEmojiController struct {
	// Embed the ApplicationController to inherit its functionality
	*ApplicationController
}

// NewCustomEmojiController creates a new CustomEmojiController
func NewCustomEmojiController(applicationController *ApplicationController) *CustomEmojiController {
	return &CustomEmojiController{
		ApplicationController: applicationController,
	}
}

// RegisterRoutes registers the routes for the CustomEmojiController
func (c *CustomEmojiController) RegisterRoutes(router *gin.RouterGroup) {
	// Add routes here
}

// RegisterMiddleware registers the middleware for the CustomEmojiController
func (c *CustomEmojiController) RegisterMiddleware(router *gin.RouterGroup) {
	// Register the middleware from the ApplicationController
	c.ApplicationController.RegisterMiddleware(router)
}

// GetFeatureCategory returns the feature category for the controller
func (c *CustomEmojiController) GetFeatureCategory() string {
	return "code_review_workflow"
}

// GetUrgency returns the urgency for the controller
func (c *CustomEmojiController) GetUrgency() string {
	return "low"
}
