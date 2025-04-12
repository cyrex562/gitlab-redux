package dashboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/analytics"
	"github.com/jmadden/gitlab-redux/internal/handler/sorting"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// ProjectsController handles project-related actions in the dashboard
type ProjectsController struct {
	ApplicationController
	projectService service.ProjectService
	eventService   service.EventService
	renderService  service.RenderService
}

// NewProjectsController creates a new ProjectsController
func NewProjectsController(
	appController ApplicationController,
	projectService service.ProjectService,
	eventService service.EventService,
	renderService service.RenderService,
) *ProjectsController {
	return &ProjectsController{
		ApplicationController: appController,
		projectService:       projectService,
		eventService:         eventService,
		renderService:        renderService,
	}
}

// RegisterRoutes registers the routes for the ProjectsController
func (c *ProjectsController) RegisterRoutes(router *gin.Engine) {
	dashboard := router.Group("/dashboard")
	{
		dashboard.GET("/projects", c.Index)
		dashboard.GET("/projects/starred", c.Starred)
	}
}

// Index handles the index action for projects
func (c *ProjectsController) Index(ctx *gin.Context) {
	// Authenticate sessionless user for RSS
	if ctx.GetHeader("Accept") == "application/atom+xml" {
		if err := c.AuthenticateSessionlessUser(ctx, "rss"); err != nil {
			c.HandleError(ctx, err)
			return
		}
	}

	// Handle redirects
	if ctx.Query("personal") == "true" {
		ctx.Redirect(http.StatusFound, "/dashboard/projects/personal")
		return
	}
	if ctx.Query("archived") == "only" {
		ctx.Redirect(http.StatusFound, "/dashboard/projects/inactive")
		return
	}

	// Set non-archived parameter
	ctx.Set("non_archived", true)

	// Set sorting
	c.setSorting(ctx)

	// Handle different formats
	switch ctx.GetHeader("Accept") {
	case "application/atom+xml":
		c.loadEvents(ctx)
		ctx.HTML(http.StatusOK, "xml", gin.H{
			"events": ctx.MustGet("events"),
		})
	default:
		ctx.HTML(http.StatusOK, "dashboard/projects/index", nil)
	}
}

// Starred handles the starred action for projects
func (c *ProjectsController) Starred(ctx *gin.Context) {
	// Set non-archived parameter
	ctx.Set("non_archived", true)

	// Set sorting
	c.setSorting(ctx)

	ctx.HTML(http.StatusOK, "dashboard/projects/starred", nil)
}

// loadEvents loads events for the projects
func (c *ProjectsController) loadEvents(ctx *gin.Context) {
	// Get current user
	currentUser, err := c.GetCurrentUser(ctx)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Create project finder params
	finderParams := map[string]interface{}{
		"non_public":            true,
		"not_aimed_for_deletion": true,
	}
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 0 {
			finderParams[k] = v[0]
		}
	}

	// Find projects
	projects, err := c.projectService.FindProjects(ctx, finderParams, currentUser)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get offset
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	// Get event filter
	eventFilter := analytics.GetEventFilter(ctx)

	// Get events
	events, err := c.eventService.GetEventsForProjects(ctx, projects, offset, eventFilter)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Render events
	renderedEvents, err := c.renderService.RenderEvents(ctx, events, currentUser, ctx.GetHeader("Accept") == "application/atom+xml")
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.Set("events", renderedEvents)
}

// setSorting sets the sorting for the projects
func (c *ProjectsController) setSorting(ctx *gin.Context) {
	sortOrder := c.getDefaultSortOrder(ctx)
	ctx.Set("sort", sortOrder)
	ctx.Set("sort_field", model.ProjectSortingPreferenceField)
}

// getDefaultSortOrder gets the default sort order
func (c *ProjectsController) getDefaultSortOrder(ctx *gin.Context) string {
	return sorting.GetSortValueName(ctx)
}
