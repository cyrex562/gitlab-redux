package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SystemInfoController handles system information display for GitLab instance
type SystemInfoController struct {
	systemInfoService *service.SystemInfoService
}

// NewSystemInfoController creates a new instance of SystemInfoController
func NewSystemInfoController(systemInfoService *service.SystemInfoService) *SystemInfoController {
	return &SystemInfoController{
		systemInfoService: systemInfoService,
	}
}

// RegisterRoutes registers the routes for the SystemInfoController
func (c *SystemInfoController) RegisterRoutes(r *gin.RouterGroup) {
	systemInfo := r.Group("/admin/system_info")
	{
		systemInfo.Use(c.requireAdmin)
		systemInfo.GET("/", c.show)
	}
}

// requireAdmin middleware ensures that only admin users can access these endpoints
func (c *SystemInfoController) requireAdmin(ctx *gin.Context) {
	user := ctx.MustGet("user")
	if user == nil || !user.IsAdmin() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// show handles the GET /admin/system_info endpoint
func (c *SystemInfoController) show(ctx *gin.Context) {
	info, err := c.systemInfoService.GetSystemInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch system information"})
		return
	}

	// TODO: Implement HTML rendering
	ctx.JSON(http.StatusOK, info)
}
