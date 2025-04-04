package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// VersionCheckController handles version checking for GitLab instance
type VersionCheckController struct {
	versionService *service.VersionService
}

// NewVersionCheckController creates a new instance of VersionCheckController
func NewVersionCheckController(versionService *service.VersionService) *VersionCheckController {
	return &VersionCheckController{
		versionService: versionService,
	}
}

// RegisterRoutes registers the routes for the VersionCheckController
func (c *VersionCheckController) RegisterRoutes(r *gin.RouterGroup) {
	version := r.Group("/admin/version_check")
	{
		version.Use(c.requireAdmin)
		version.GET("/", c.versionCheck)
	}
}

// versionCheck handles the GET /admin/version_check endpoint
func (c *VersionCheckController) versionCheck(ctx *gin.Context) {
	response, err := c.versionService.CheckVersion(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check version"})
		return
	}

	// Set cache expiration to 1 minute if response exists
	if response != nil {
		ctx.Header("Cache-Control", "public, max-age=60")
		ctx.Header("Expires", time.Now().Add(time.Minute).Format(time.RFC1123))
	}

	ctx.JSON(http.StatusOK, response)
}
