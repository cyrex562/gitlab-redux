package dashboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/access"
	"github.com/jmadden/gitlab-redux/internal/handler/user"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// ApplicationController is the base controller for dashboard handlers
type ApplicationController struct {
	authService *service.AuthService
	projectService *service.ProjectService
	userService *service.UserService
	crossProjectAccessHandler *access.CrossProjectAccessHandler
	recordUserLastActivityHandler *user.RecordUserLastActivityHandler
}

// NewApplicationController creates a new dashboard application controller
func NewApplicationController(
	authService *service.AuthService,
	projectService *service.ProjectService,
	userService *service.UserService,
	crossProjectAccessHandler *access.CrossProjectAccessHandler,
	recordUserLastActivityHandler *user.RecordUserLastActivityHandler,
) *ApplicationController {
	return &ApplicationController{
		authService: authService,
		projectService: projectService,
		userService: userService,
		crossProjectAccessHandler: crossProjectAccessHandler,
		recordUserLastActivityHandler: recordUserLastActivityHandler,
	}
}

// RegisterMiddleware registers the middleware for the dashboard application controller
func (c *ApplicationController) RegisterMiddleware(router *gin.Engine) {
	// Register the cross-project access check middleware
	c.crossProjectAccessHandler.RegisterMiddleware(router)

	// Register the record user last activity middleware
	c.recordUserLastActivityHandler.RegisterMiddleware(router)
}

// GetCurrentUser gets the current user from the context
func (c *ApplicationController) GetCurrentUser(ctx *gin.Context) (*model.User, error) {
	return c.authService.GetCurrentUser(ctx)
}

// GetProjectIDs gets the project IDs from the context
func (c *ApplicationController) GetProjectIDs(ctx *gin.Context) ([]int64, error) {
	// Get project IDs from query parameters
	projectIDStrs := ctx.QueryArray("project_ids[]")
	if len(projectIDStrs) == 0 {
		// If no project IDs are provided, get all authorized projects
		projects := c.GetProjects(ctx)
		projectIDs := make([]int64, 0, len(projects))
		for _, project := range projects {
			projectIDs = append(projectIDs, project.ID)
		}
		return projectIDs, nil
	}

	// Convert project ID strings to int64
	projectIDs := make([]int64, 0, len(projectIDStrs))
	for _, projectIDStr := range projectIDStrs {
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		projectIDs = append(projectIDs, projectID)
	}

	return projectIDs, nil
}

// GetProjects gets the authorized projects for the current user
func (c *ApplicationController) GetProjects(ctx *gin.Context) []*model.Project {
	// Get the current user
	user, err := c.authService.GetCurrentUser(ctx)
	if err != nil {
		return []*model.Project{}
	}

	// Get the authorized projects
	projects, err := c.projectService.GetAuthorizedProjects(ctx, user)
	if err != nil {
		return []*model.Project{}
	}

	// Filter out archived projects
	var nonArchivedProjects []*model.Project
	for _, project := range projects {
		if !project.Archived {
			nonArchivedProjects = append(nonArchivedProjects, project)
		}
	}

	// Sort projects by activity
	return c.projectService.SortProjectsByActivity(nonArchivedProjects)
}

// RenderDashboard renders the dashboard layout
func (c *ApplicationController) RenderDashboard(ctx *gin.Context, template string, data gin.H) {
	// Add the layout to the data
	data["layout"] = "dashboard"

	// Render the template
	ctx.HTML(http.StatusOK, template, data)
}

// RenderError renders an error response
func (c *ApplicationController) RenderError(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"error": message,
	})
}
