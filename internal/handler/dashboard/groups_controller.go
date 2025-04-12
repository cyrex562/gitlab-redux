package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitlab-org/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/handler/groups"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// GroupsController handles the groups section of the dashboard
type GroupsController struct {
	*ApplicationController
	groupService *service.GroupService
	groupTreeHandler *groups.GroupTreeHandler
}

// NewGroupsController creates a new dashboard groups controller
func NewGroupsController(
	appController *ApplicationController,
	groupService *service.GroupService,
	groupTreeHandler *groups.GroupTreeHandler,
) *GroupsController {
	return &GroupsController{
		ApplicationController: appController,
		groupService: groupService,
		groupTreeHandler: groupTreeHandler,
	}
}

// RegisterRoutes registers the groups routes
func (c *GroupsController) RegisterRoutes(router *gin.RouterGroup) {
	groups := router.Group("/groups")
	{
		groups.GET("", c.index)
	}
}

// index handles the groups index page
func (c *GroupsController) index(ctx *gin.Context) {
	// Get the current user
	user, err := c.authService.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the groups
	groups, err := c.groupService.FindGroups(ctx, user, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render the group tree
	c.renderGroupTree(ctx, groups)
}

// renderGroupTree renders the group tree
func (c *GroupsController) renderGroupTree(ctx *gin.Context, groups []*model.Group) {
	// Use the group tree handler to render the tree
	c.groupTreeHandler.RenderGroupTree(ctx, groups)
}
