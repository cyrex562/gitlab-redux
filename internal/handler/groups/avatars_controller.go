package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AvatarsController handles group avatars
type AvatarsController struct {
	// Embed the ApplicationController to inherit its functionality
	*ApplicationController

	// Add any additional dependencies here
	groupService *GroupService
}

// NewAvatarsController creates a new AvatarsController
func NewAvatarsController(
	applicationController *ApplicationController,
	groupService *GroupService,
) *AvatarsController {
	return &AvatarsController{
		ApplicationController: applicationController,
		groupService:          groupService,
	}
}

// RegisterRoutes registers the routes for the AvatarsController
func (c *AvatarsController) RegisterRoutes(router *gin.RouterGroup) {
	router.DELETE("/", c.Destroy)
}

// RegisterMiddleware registers the middleware for the AvatarsController
func (c *AvatarsController) RegisterMiddleware(router *gin.RouterGroup) {
	// Register the middleware from the ApplicationController
	c.ApplicationController.RegisterMiddleware(router)

	// Add authorization middleware
	router.Use(c.AuthorizeAdminGroup())

	// Skip cross project access check for destroy action
	// This is handled in the route registration
}

// Destroy handles the destroy action
func (c *AvatarsController) Destroy(ctx *gin.Context) {
	// Get the group from the context
	group := c.GetGroup(ctx)

	// Remove the avatar
	err := c.groupService.RemoveAvatar(group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove avatar"})
		return
	}

	// Save the group
	err = c.groupService.Save(group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save group"})
		return
	}

	// Redirect to the edit group path
	ctx.Redirect(http.StatusFound, "/groups/"+group.ID+"/edit")
}

// AuthorizeAdminGroup middleware checks if the user has permission to admin the group
func (c *AvatarsController) AuthorizeAdminGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the user has permission to admin the group
		if !c.AuthorizeAdminGroup(ctx) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "You don't have permission to admin the group"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// GroupService handles group operations
type GroupService struct {
	// Add any dependencies here
}

// RemoveAvatar removes the avatar from a group
func (s *GroupService) RemoveAvatar(group interface{}) error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would remove the avatar from the group
	return nil
}

// Save saves a group
func (s *GroupService) Save(group interface{}) error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would save the group
	return nil
}
