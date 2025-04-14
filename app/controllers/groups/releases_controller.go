package groups

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/releases"
)

// ReleasesController handles requests for group releases
type ReleasesController struct {
	groupReleasesFinder *releases.GroupReleasesFinder
}

// NewReleasesController creates a new releases controller
func NewReleasesController(groupReleasesFinder *releases.GroupReleasesFinder) *ReleasesController {
	return &ReleasesController{
		groupReleasesFinder: groupReleasesFinder,
	}
}

// RegisterRoutes registers the routes for the releases controller
func (c *ReleasesController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/releases", c.Index)
}

// Index handles GET requests for group releases
func (c *ReleasesController) Index(ctx *gin.Context) {
	group := ctx.MustGet("group").(*models.Group)
	user := ctx.MustGet("current_user").(*models.User)

	// Get pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage := 30 // Default per page as in Ruby version

	// Find releases
	releases, err := c.groupReleasesFinder.Execute(ctx, group, user, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, releases)
} 