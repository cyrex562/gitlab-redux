package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AchievementsController handles group achievements
type AchievementsController struct {
	// Add any dependencies here
}

// NewAchievementsController creates a new AchievementsController
func NewAchievementsController() *AchievementsController {
	return &AchievementsController{}
}

// RegisterRoutes registers the routes for the AchievementsController
func (c *AchievementsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/new", c.New)
}

// RegisterMiddleware registers the middleware for the AchievementsController
func (c *AchievementsController) RegisterMiddleware(router *gin.RouterGroup) {
	router.Use(c.AuthorizeReadAchievement())
}

// New handles the new action
func (c *AchievementsController) New(ctx *gin.Context) {
	// Render the index template
	ctx.HTML(http.StatusOK, "groups/achievements/index", gin.H{})
}

// AuthorizeReadAchievement middleware checks if the user has permission to read achievements
func (c *AchievementsController) AuthorizeReadAchievement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			ctx.Abort()
			return
		}

		// Get the group from the context
		group, exists := ctx.Get("group")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			ctx.Abort()
			return
		}

		// Check if the user has permission to read achievements
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the user has the read_achievement permission
		canReadAchievement := true // Replace with actual check

		if !canReadAchievement {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "You don't have permission to read achievements"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
