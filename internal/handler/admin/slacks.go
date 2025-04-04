package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/integration/slack"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SlacksController handles Slack integration settings for GitLab instance
type SlacksController struct {
	slackService *service.SlackService
}

// NewSlacksController creates a new instance of SlacksController
func NewSlacksController(slackService *service.SlackService) *SlacksController {
	return &SlacksController{
		slackService: slackService,
	}
}

// RegisterRoutes registers the routes for the SlacksController
func (c *SlacksController) RegisterRoutes(r *gin.RouterGroup) {
	slacks := r.Group("/admin/slacks")
	{
		slacks.Use(c.requireAdmin)
		slacks.GET("/", c.index)
		slacks.POST("/install", c.install)
	}
}

// requireAdmin middleware ensures that only admin users can access these endpoints
func (c *SlacksController) requireAdmin(ctx *gin.Context) {
	user := ctx.MustGet("user")
	if user == nil || !user.IsAdmin() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// index handles the GET /admin/slacks endpoint
func (c *SlacksController) index(ctx *gin.Context) {
	c.redirectToIntegrationPage(ctx)
}

// install handles the POST /admin/slacks/install endpoint
func (c *SlacksController) install(ctx *gin.Context) {
	code := ctx.PostForm("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Installation code is required"})
		return
	}

	installationService := c.installationService(ctx).WithCode(code)
	if err := installationService.Execute(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to install Slack integration"})
		return
	}

	c.redirectToIntegrationPage(ctx)
}

// getIntegration retrieves the GitLab Slack application integration
func (c *SlacksController) getIntegration() (*slack.GitlabSlackApplication, error) {
	return slack.ForInstance()
}

// redirectToIntegrationPage redirects to the integration settings page
func (c *SlacksController) redirectToIntegrationPage(ctx *gin.Context) {
	integration, err := c.getIntegration()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get integration"})
		return
	}

	redirectPath := "/admin/application_settings/integrations"
	if integration != nil {
		redirectPath += "/" + integration.ID + "/edit"
	} else {
		redirectPath += "/new"
	}

	ctx.Redirect(http.StatusFound, redirectPath)
}

// installationService creates a new Slack installation service
func (c *SlacksController) installationService(ctx *gin.Context) *service.SlackInstallationService {
	return service.NewSlackInstallationService(ctx)
}
