package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SpamLogsController handles spam log management for GitLab instance
type SpamLogsController struct {
	spamLogService *service.SpamLogService
}

// NewSpamLogsController creates a new instance of SpamLogsController
func NewSpamLogsController(spamLogService *service.SpamLogService) *SpamLogsController {
	return &SpamLogsController{
		spamLogService: spamLogService,
	}
}

// RegisterRoutes registers the routes for the SpamLogsController
func (c *SpamLogsController) RegisterRoutes(r *gin.RouterGroup) {
	spamLogs := r.Group("/admin/spam_logs")
	{
		spamLogs.Use(c.requireAdmin)
		spamLogs.GET("/", c.index)
		spamLogs.DELETE("/:id", c.destroy)
		spamLogs.POST("/:id/mark_as_ham", c.markAsHam)
	}
}

// requireAdmin middleware ensures that only admin users can access these endpoints
func (c *SpamLogsController) requireAdmin(ctx *gin.Context) {
	user := ctx.MustGet("user")
	if user == nil || !user.IsAdmin() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// index handles the GET /admin/spam_logs endpoint
func (c *SpamLogsController) index(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	spamLogs, err := c.spamLogService.List(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spam logs"})
		return
	}

	// TODO: Implement HTML rendering
	ctx.JSON(http.StatusOK, spamLogs)
}

// destroy handles the DELETE /admin/spam_logs/:id endpoint
func (c *SpamLogsController) destroy(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spam log ID"})
		return
	}

	removeUser := ctx.Query("remove_user") == "true"
	user := ctx.MustGet("user").(*model.User)

	if err := c.spamLogService.Destroy(ctx, id, removeUser, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to destroy spam log"})
		return
	}

	if removeUser {
		ctx.Redirect(http.StatusFound, "/admin/spam_logs")
	} else {
		ctx.Status(http.StatusOK)
	}
}

// markAsHam handles the POST /admin/spam_logs/:id/mark_as_ham endpoint
func (c *SpamLogsController) markAsHam(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spam log ID"})
		return
	}

	if err := c.spamLogService.MarkAsHam(ctx, id); err != nil {
		ctx.Redirect(http.StatusFound, "/admin/spam_logs?alert=Error with Akismet. Please check the logs for more info.")
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/spam_logs?notice=Spam log successfully submitted as ham.")
}
