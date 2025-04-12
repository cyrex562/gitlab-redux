package dashboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// MilestonesController handles requests related to milestones in the dashboard
type MilestonesController struct {
	*ApplicationController
	milestoneService service.MilestoneService
	groupService     service.GroupService
}

// NewMilestonesController creates a new MilestonesController
func NewMilestonesController(
	appController *ApplicationController,
	milestoneService service.MilestoneService,
	groupService service.GroupService,
) *MilestonesController {
	return &MilestonesController{
		ApplicationController: appController,
		milestoneService:     milestoneService,
		groupService:         groupService,
	}
}

// Index handles GET requests to /dashboard/milestones
// It returns either HTML or JSON based on the Accept header
func (c *MilestonesController) Index(ctx *gin.Context) {
	// Get the current user
	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		c.RenderError(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get projects
	projects := c.GetProjects(ctx)
	projectIDs := make([]int64, 0, len(projects))
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ID)
	}

	// Get groups
	groups, err := c.getGroups(ctx, user)
	if err != nil {
		c.RenderError(ctx, http.StatusInternalServerError, "Failed to retrieve groups")
		return
	}

	groupIDs := make([]int64, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}

	// Create search parameters
	searchParams := c.createSearchParams(ctx, projectIDs, groupIDs)

	// Handle different response formats
	switch ctx.GetHeader("Accept") {
	case "application/json":
		c.handleJSONResponse(ctx, searchParams)
	default:
		c.handleHTMLResponse(ctx, searchParams, projectIDs, groupIDs)
	}
}

// handleJSONResponse handles JSON response for the index action
func (c *MilestonesController) handleJSONResponse(ctx *gin.Context, params *service.MilestoneSearchParams) {
	// Get milestones in JSON format
	milestones, err := c.milestoneService.GetMilestoneJSON(ctx, params)
	if err != nil {
		c.RenderError(ctx, http.StatusInternalServerError, "Failed to retrieve milestones")
		return
	}

	// Return JSON response
	ctx.JSON(http.StatusOK, milestones)
}

// handleHTMLResponse handles HTML response for the index action
func (c *MilestonesController) handleHTMLResponse(ctx *gin.Context, params *service.MilestoneSearchParams, projectIDs []int64, groupIDs []int64) {
	// Get milestone state count
	stateCount, err := c.milestoneService.GetMilestoneStateCount(ctx, projectIDs, groupIDs)
	if err != nil {
		c.RenderError(ctx, http.StatusInternalServerError, "Failed to retrieve milestone state count")
		return
	}

	// Get milestones
	milestones, err := c.milestoneService.FindMilestones(ctx, params)
	if err != nil {
		c.RenderError(ctx, http.StatusInternalServerError, "Failed to retrieve milestones")
		return
	}

	// Render dashboard with milestones
	c.RenderDashboard(ctx, "dashboard/milestones/index", gin.H{
		"milestone_states": stateCount,
		"milestones":       milestones,
	})
}

// getGroups gets the groups for the current user
func (c *MilestonesController) getGroups(ctx *gin.Context, user *model.User) ([]*model.Group, error) {
	return c.groupService.FindByUser(ctx, user.ID, false)
}

// createSearchParams creates search parameters from the request
func (c *MilestonesController) createSearchParams(ctx *gin.Context, projectIDs []int64, groupIDs []int64) *service.MilestoneSearchParams {
	// Get state from query parameters
	stateStr := ctx.Query("state")
	var state *model.MilestoneState
	if stateStr != "" {
		stateValue := model.MilestoneState(stateStr)
		state = &stateValue
	}

	// Get search title from query parameters
	searchTitle := ctx.Query("search_title")

	// Get pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	return &service.MilestoneSearchParams{
		State:       state,
		SearchTitle: searchTitle,
		GroupIDs:    groupIDs,
		ProjectIDs:  projectIDs,
		Page:        page,
		PerPage:     perPage,
	}
}

// RegisterRoutes registers the routes for the MilestonesController
func (c *MilestonesController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/milestones", c.Index)
}
