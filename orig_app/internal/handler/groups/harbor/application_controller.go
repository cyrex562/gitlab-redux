package harbor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApplicationController is the base controller for Harbor
type ApplicationController struct {
	// Add any dependencies here
}

// NewApplicationController creates a new ApplicationController
func NewApplicationController() *ApplicationController {
	return &ApplicationController{}
}

// RegisterMiddleware registers the middleware for the ApplicationController
func (c *ApplicationController) RegisterMiddleware(router *gin.RouterGroup) {
	// Register the authorization middleware
	router.Use(c.AuthorizeReadHarborRegistry())
}

// AuthorizeReadHarborRegistry middleware checks if the user has permission to read Harbor registry
func (c *ApplicationController) AuthorizeReadHarborRegistry() gin.HandlerFunc {
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

		// Check if the user has permission to read Harbor registry
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the user has the read_harbor_registry permission
		canReadHarborRegistry := true // Replace with actual check

		if !canReadHarborRegistry {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "You don't have permission to read Harbor registry"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
