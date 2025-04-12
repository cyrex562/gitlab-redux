package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

type DashboardController struct {
	dashboardService *service.DashboardService
	kasService      *service.KasService
	redisService    *service.RedisService
}

func NewDashboardController(
	dashboardService *service.DashboardService,
	kasService *service.KasService,
	redisService *service.RedisService,
) *DashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
		kasService:      kasService,
		redisService:    redisService,
	}
}

// Index displays the main admin dashboard
func (c *DashboardController) Index(ctx *gin.Context) {
	// Get approximate counts
	counts, err := c.dashboardService.GetApproximateCounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get recent projects
	projects, err := c.dashboardService.GetRecentProjects(ctx, 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get recent users
	users, err := c.dashboardService.GetRecentUsers(ctx, 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get recent groups
	groups, err := c.dashboardService.GetRecentGroups(ctx, 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get system notices
	notices, err := c.dashboardService.GetSystemNotices(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get KAS server info if enabled
	var kasInfo *service.KasServerInfo
	if enabled, _ := c.kasService.IsEnabled(ctx); enabled {
		kasInfo, err = c.kasService.GetServerInfo(ctx)
		if err != nil {
			ctx.Error(err) // Log error but don't fail the request
		}
	}

	// Get Redis versions
	redisVersions, err := c.redisService.GetVersions(ctx)
	if err != nil {
		ctx.Error(err) // Log error but don't fail the request
	}

	ctx.JSON(http.StatusOK, gin.H{
		"counts":        counts,
		"projects":      projects,
		"users":         users,
		"groups":        groups,
		"notices":       notices,
		"kas_info":      kasInfo,
		"redis_versions": redisVersions,
	})
}

// Stats displays user statistics
func (c *DashboardController) Stats(ctx *gin.Context) {
	stats, err := c.dashboardService.GetUserStatistics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users_statistics": stats,
	})
}
