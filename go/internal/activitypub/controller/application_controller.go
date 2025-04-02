package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/feature"
	"gitlab.com/gitlab-org/gitlab-redux/internal/routing"
)

// ApplicationController handles ActivityPub related requests
type ApplicationController struct {
	*routing.BaseController
}

// NewApplicationController creates a new instance of ApplicationController
func NewApplicationController() *ApplicationController {
	return &ApplicationController{
		BaseController: routing.NewBaseController(),
	}
}

// SetupRoutes configures the routes for the ActivityPub controller
func (c *ApplicationController) SetupRoutes(router *gin.Engine) {
	// ActivityPub routes will be added here
	group := router.Group("/activity-pub")
	{
		group.Use(c.ensureFeatureFlag)
		group.Use(c.setContentType)
		// Add specific routes here
	}
}

// ensureFeatureFlag middleware checks if the ActivityPub feature is enabled
func (c *ApplicationController) ensureFeatureFlag(ctx *gin.Context) {
	if !feature.IsEnabled("activity_pub") {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.Next()
}

// setContentType middleware sets the content type for ActivityPub responses
func (c *ApplicationController) setContentType(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/activity+json")
	ctx.Next()
}

// Can checks if an object has permission to perform an action
func (c *ApplicationController) Can(object interface{}, action string, subject interface{}) bool {
	// TODO: Implement ability checking logic
	return false
}

// RouteNotFound handles 404 responses
func (c *ApplicationController) RouteNotFound(ctx *gin.Context) {
	ctx.Status(http.StatusNotFound)
}
