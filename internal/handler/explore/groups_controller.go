package explore

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// GroupsController handles group-related actions
type GroupsController struct {
	*ApplicationController
	groupService service.GroupService
}

// NewGroupsController creates a new GroupsController
func NewGroupsController(
	appController *ApplicationController,
	groupService service.GroupService,
) *GroupsController {
	return &GroupsController{
		ApplicationController: appController,
		groupService:         groupService,
	}
}

// RegisterRoutes registers the routes for the GroupsController
func (c *GroupsController) RegisterRoutes(router *gin.Engine) {
	explore := router.Group("/explore")
	{
		explore.GET("/groups", c.Index)
	}
}

// Index handles the index action
func (c *GroupsController) Index(ctx *gin.Context) {
	// Get current user
	user := ctx.MustGet("current_user").(*model.User)

	// Find groups
	groups, err := c.groupService.FindGroups(ctx, user, 10000) // MAX_QUERY_SIZE = 10_000
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render group tree
	groups.RenderGroupTree(ctx, groups)
}
