package integrations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SlackControllerSettings provides functionality for handling Slack integration settings
type SlackControllerSettings struct {
	integrationService *service.IntegrationService
	slackService      *service.SlackService
}

// NewSlackControllerSettings creates a new instance of SlackControllerSettings
func NewSlackControllerSettings(integrationService *service.IntegrationService, slackService *service.SlackService) *SlackControllerSettings {
	return &SlackControllerSettings{
		integrationService: integrationService,
		slackService:      slackService,
	}
}

// RegisterRoutes registers the routes for Slack integration settings
func (s *SlackControllerSettings) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/slack/auth", s.handleOAuthError, s.checkOAuthState, s.slackAuth)
	router.DELETE("/slack", s.destroy)
}

// slackAuth handles the GET /slack/auth endpoint for Slack OAuth
func (s *SlackControllerSettings) slackAuth(ctx *gin.Context) {
	// Get the installation service
	installationService := s.getInstallationService(ctx)

	// Execute the installation service
	result, err := installationService.Execute(ctx)
	if err != nil {
		ctx.SetCookie("alert", "Failed to install Slack integration", 3600, "/", "", false, true)
		s.redirectToIntegrationPage(ctx)
		return
	}

	// Set flash message if there's an error
	if result.Error {
		ctx.SetCookie("alert", result.Message, 3600, "/", "", false, true)
	}

	// Set session variable
	ctx.SetCookie("slack_install_success", result.Success, 3600, "/", "", false, true)

	// Redirect to integration page
	s.redirectToIntegrationPage(ctx)
}

// destroy handles the DELETE /slack endpoint for destroying a Slack integration
func (s *SlackControllerSettings) destroy(ctx *gin.Context) {
	// Get the Slack integration
	slackIntegration := s.getSlackIntegration(ctx)

	// Destroy the Slack integration
	err := s.slackService.DestroySlackIntegration(ctx, slackIntegration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to destroy Slack integration",
			"errors":  err.Error(),
		})
		return
	}

	// Propagate integration changes if not project level
	integration := s.getIntegration(ctx)
	if !integration.IsProjectLevel() {
		go s.integrationService.PropagateIntegration(ctx, integration.ID)
	}

	// Redirect to integration page
	s.redirectToIntegrationPage(ctx)
}

// handleOAuthError middleware handles OAuth errors
func (s *SlackControllerSettings) handleOAuthError(ctx *gin.Context) {
	if ctx.Query("error") == "access_denied" {
		ctx.SetCookie("alert", "Access request canceled", 3600, "/", "", false, true)
		s.redirectToIntegrationPage(ctx)
		ctx.Abort()
		return
	}
	ctx.Next()
}

// checkOAuthState middleware checks the OAuth state
func (s *SlackControllerSettings) checkOAuthState(ctx *gin.Context) {
	// Get the state from the query
	state := ctx.Query("state")

	// Get the state from the session
	sessionState, exists := ctx.Get("slack_oauth_state")
	if !exists || sessionState.(string) != state {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Next()
}

// getSlackIntegration gets the Slack integration
func (s *SlackControllerSettings) getSlackIntegration(ctx *gin.Context) *model.SlackIntegration {
	integration := s.getIntegration(ctx)
	slackIntegration, err := s.slackService.GetSlackIntegration(ctx, integration)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}
	return slackIntegration
}

// getIntegration gets the integration
func (s *SlackControllerSettings) getIntegration(ctx *gin.Context) *model.Integration {
	integration, exists := ctx.Get("integration")
	if !exists {
		ctx.AbortWithStatus(http.StatusNotFound)
		return nil
	}
	return integration.(*model.Integration)
}

// getInstallationService gets the installation service
func (s *SlackControllerSettings) getInstallationService(ctx *gin.Context) *service.SlackInstallationService {
	// This should be implemented by the controller that uses this handler
	// For now, we'll return a placeholder
	return s.slackService.NewInstallationService(ctx)
}

// redirectToIntegrationPage redirects to the integration page
func (s *SlackControllerSettings) redirectToIntegrationPage(ctx *gin.Context) {
	// This should be implemented by the controller that uses this handler
	// For now, we'll redirect to a placeholder URL
	ctx.Redirect(http.StatusSeeOther, "/integrations/slack")
}
