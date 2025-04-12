package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/middleware"
	"gitlab.com/gitlab-org/gitlab-redux/internal/routing"
)

// ApplicationController provides a base class for Admin controllers to subclass
// Automatically sets the layout and ensures an administrator is logged in
type ApplicationController struct {
	*routing.BaseController
}

// NewApplicationController creates a new instance of ApplicationController
func NewApplicationController() *ApplicationController {
	return &ApplicationController{
		BaseController: routing.NewBaseController(),
	}
}

// SetupRoutes configures the base routes and middleware for admin controllers
func (c *ApplicationController) SetupRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		// Apply admin-specific middleware
		admin.Use(c.ensureAdmin)
		admin.Use(c.setAdminLayout)
	}
}

// ensureAdmin middleware ensures that the current user is an administrator
func (c *ApplicationController) ensureAdmin(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil || !user.IsAdmin() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Administrator privileges required."})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// setAdminLayout middleware sets the admin layout for all admin pages
func (c *ApplicationController) setAdminLayout(ctx *gin.Context) {
	ctx.Set("layout", "admin")
	ctx.Next()
}

// GetCurrentUser returns the current authenticated user
func (c *ApplicationController) GetCurrentUser(ctx *gin.Context) *models.User {
	// TODO: Implement user retrieval from session/token
	return nil
}
