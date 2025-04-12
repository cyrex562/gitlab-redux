package explore

import (
	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/auth"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// ApplicationController is the base controller for all explore controllers
type ApplicationController struct {
	authService service.AuthService
}

// NewApplicationController creates a new ApplicationController
func NewApplicationController(authService service.AuthService) *ApplicationController {
	return &ApplicationController{
		authService: authService,
	}
}

// SkipAuthenticationUnlessPublicVisibilityRestricted skips authentication unless public visibility is restricted
func (c *ApplicationController) SkipAuthenticationUnlessPublicVisibilityRestricted(ctx *gin.Context) {
	// Check if public visibility is restricted
	if !c.authService.IsPublicVisibilityRestricted(ctx) {
		// Skip authentication
		auth.SkipAuthentication(ctx)
	}
}

// GetLayout returns the layout for the explore controllers
func (c *ApplicationController) GetLayout() string {
	return "explore"
}
