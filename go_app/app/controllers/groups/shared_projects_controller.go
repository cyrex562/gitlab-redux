package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/serializers"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/projects"
)

// SharedProjectsController handles requests for group shared projects
type SharedProjectsController struct {
	groupProjectsFinder *projects.GroupProjectsFinder
	groupChildSerializer *serializers.GroupChildSerializer
}

// NewSharedProjectsController creates a new shared projects controller
func NewSharedProjectsController(
	groupProjectsFinder *projects.GroupProjectsFinder,
	groupChildSerializer *serializers.GroupChildSerializer,
) *SharedProjectsController {
	return &SharedProjectsController{
		groupProjectsFinder: groupProjectsFinder,
		groupChildSerializer: groupChildSerializer,
	}
}

// RegisterRoutes registers the routes for the shared projects controller
func (c *SharedProjectsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/shared_projects", c.Index)
}

// Index handles GET requests for group shared projects
func (c *SharedProjectsController) Index(ctx *gin.Context) {
	group := ctx.MustGet("group").(*models.Group)
	user := ctx.MustGet("current_user").(*models.User)

	// Get finder parameters
	finderParams := c.getFinderParams(ctx)

	// Find shared projects
	sharedProjects, err := c.groupProjectsFinder.Execute(ctx, group, user, finderParams, projects.GroupProjectsFinderOptions{
		OnlyShared: true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serialize the projects
	serializer := c.groupChildSerializer.WithPagination(ctx.Request, ctx.Writer)
	result := serializer.Represent(sharedProjects)

	ctx.JSON(http.StatusOK, result)
}

// getFinderParams extracts and processes finder parameters from the request
func (c *SharedProjectsController) getFinderParams(ctx *gin.Context) projects.GroupProjectsFinderParams {
	params := projects.GroupProjectsFinderParams{}

	// Make the `search` param consistent for the frontend,
	// which will be using `filter`.
	if filter := ctx.Query("filter"); filter != "" {
		params.Search = filter
	} else {
		params.Search = ctx.Query("search")
	}

	// Don't show archived projects
	params.NonArchived = true

	// Add other parameters
	params.Sort = ctx.Query("sort")

	return params
} 