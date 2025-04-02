package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services"
)

// ApplicationSettingsController handles application settings in the admin panel
type ApplicationSettingsController struct {
	*ApplicationController
}

// NewApplicationSettingsController creates a new instance of ApplicationSettingsController
func NewApplicationSettingsController() *ApplicationSettingsController {
	return &ApplicationSettingsController{
		ApplicationController: NewApplicationController(),
	}
}

// SetupRoutes configures the routes for the application settings controller
func (c *ApplicationSettingsController) SetupRoutes(router *gin.Engine) {
	admin := router.Group("/admin/application_settings")
	{
		// General settings
		admin.GET("/general", c.General)
		admin.PATCH("/general", c.UpdateGeneral)

		// Repository settings
		admin.GET("/repository", c.Repository)
		admin.PATCH("/repository", c.UpdateRepository)

		// CI/CD settings
		admin.GET("/ci_cd", c.CICD)
		admin.PATCH("/ci_cd", c.UpdateCICD)

		// Reporting settings
		admin.GET("/reporting", c.Reporting)
		admin.PATCH("/reporting", c.UpdateReporting)

		// Metrics and profiling settings
		admin.GET("/metrics_and_profiling", c.MetricsAndProfiling)
		admin.PATCH("/metrics_and_profiling", c.UpdateMetricsAndProfiling)

		// Network settings
		admin.GET("/network", c.Network)
		admin.PATCH("/network", c.UpdateNetwork)

		// Preferences settings
		admin.GET("/preferences", c.Preferences)
		admin.PATCH("/preferences", c.UpdatePreferences)

		// Special actions
		admin.POST("/reset_registration_token", c.ResetRegistrationToken)
		admin.POST("/reset_health_check_token", c.ResetHealthCheckToken)
		admin.POST("/reset_error_tracking_access_token", c.ResetErrorTrackingAccessToken)
		admin.POST("/clear_repository_check_states", c.ClearRepositoryCheckStates)
		admin.GET("/lets_encrypt_terms_of_service", c.LetsEncryptTermsOfService)
		admin.GET("/slack_app_manifest_share", c.SlackAppManifestShare)
		admin.GET("/slack_app_manifest_download", c.SlackAppManifestDownload)
		admin.GET("/usage_data", c.UsageData)
	}
}

// General displays the general settings page
func (c *ApplicationSettingsController) General(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	ctx.HTML(http.StatusOK, "admin/application_settings/general", gin.H{
		"settings": settings,
	})
}

// UpdateGeneral updates the general settings
func (c *ApplicationSettingsController) UpdateGeneral(ctx *gin.Context) {
	c.performUpdate(ctx, "general")
}

// Repository displays the repository settings page
func (c *ApplicationSettingsController) Repository(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	ctx.HTML(http.StatusOK, "admin/application_settings/repository", gin.H{
		"settings": settings,
	})
}

// UpdateRepository updates the repository settings
func (c *ApplicationSettingsController) UpdateRepository(ctx *gin.Context) {
	c.performUpdate(ctx, "repository")
}

// CICD displays the CI/CD settings page
func (c *ApplicationSettingsController) CICD(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	ctx.HTML(http.StatusOK, "admin/application_settings/ci_cd", gin.H{
		"settings": settings,
	})
}

// UpdateCICD updates the CI/CD settings
func (c *ApplicationSettingsController) UpdateCICD(ctx *gin.Context) {
	c.performUpdate(ctx, "ci_cd")
}

// Reporting displays the reporting settings page
func (c *ApplicationSettingsController) Reporting(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	ctx.HTML(http.StatusOK, "admin/application_settings/reporting", gin.H{
		"settings": settings,
	})
}

// UpdateReporting updates the reporting settings
func (c *ApplicationSettingsController) UpdateReporting(ctx *gin.Context) {
	c.performUpdate(ctx, "reporting")
}

// MetricsAndProfiling displays the metrics and profiling settings page
func (c *ApplicationSettingsController) MetricsAndProfiling(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	ctx.HTML(http.StatusOK, "admin/application_settings/metrics_and_profiling", gin.H{
		"settings": settings,
	})
}

// UpdateMetricsAndProfiling updates the metrics and profiling settings
func (c *ApplicationSettingsController) UpdateMetricsAndProfiling(ctx *gin.Context) {
	c.performUpdate(ctx, "metrics_and_profiling")
}

// Network displays the network settings page
func (c *ApplicationSettingsController) Network(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	ctx.HTML(http.StatusOK, "admin/application_settings/network", gin.H{
		"settings": settings,
	})
}

// UpdateNetwork updates the network settings
func (c *ApplicationSettingsController) UpdateNetwork(ctx *gin.Context) {
	c.performUpdate(ctx, "network")
}

// Preferences displays the preferences settings page
func (c *ApplicationSettingsController) Preferences(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	ctx.HTML(http.StatusOK, "admin/application_settings/preferences", gin.H{
		"settings": settings,
	})
}

// UpdatePreferences updates the preferences settings
func (c *ApplicationSettingsController) UpdatePreferences(ctx *gin.Context) {
	c.performUpdate(ctx, "preferences")
}

// ResetRegistrationToken resets the registration token
func (c *ApplicationSettingsController) ResetRegistrationToken(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	service := services.NewResetRegistrationTokenService(settings, c.GetCurrentUser(ctx))
	if err := service.Execute(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Set("notice", "New runners registration token has been generated!")
	ctx.Redirect(http.StatusFound, "/admin/runners")
}

// ResetHealthCheckToken resets the health check token
func (c *ApplicationSettingsController) ResetHealthCheckToken(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	if err := settings.ResetHealthCheckAccessToken(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Set("notice", "New health check access token has been generated!")
	ctx.Redirect(http.StatusFound, ctx.Request.Header.Get("Referer"))
}

// ResetErrorTrackingAccessToken resets the error tracking access token
func (c *ApplicationSettingsController) ResetErrorTrackingAccessToken(ctx *gin.Context) {
	settings := models.GetCurrentApplicationSettings()
	if err := settings.ResetErrorTrackingAccessToken(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Set("notice", "New error tracking access token has been generated!")
	ctx.Redirect(http.StatusFound, "/admin/application_settings/general")
}

// ClearRepositoryCheckStates clears all repository check states
func (c *ApplicationSettingsController) ClearRepositoryCheckStates(ctx *gin.Context) {
	// TODO: Implement repository check state clearing
	ctx.Set("notice", "Started asynchronous removal of all repository check states.")
	ctx.Redirect(http.StatusFound, "/admin/application_settings/general")
}

// LetsEncryptTermsOfService redirects to Let's Encrypt terms of service
func (c *ApplicationSettingsController) LetsEncryptTermsOfService(ctx *gin.Context) {
	// TODO: Implement Let's Encrypt terms of service URL
	ctx.Redirect(http.StatusFound, "https://letsencrypt.org/repository/")
}

// SlackAppManifestShare redirects to Slack app manifest share URL
func (c *ApplicationSettingsController) SlackAppManifestShare(ctx *gin.Context) {
	// TODO: Implement Slack app manifest share URL
	ctx.Redirect(http.StatusFound, "https://api.slack.com/apps")
}

// SlackAppManifestDownload downloads the Slack app manifest
func (c *ApplicationSettingsController) SlackAppManifestDownload(ctx *gin.Context) {
	// TODO: Implement Slack app manifest generation
	manifest := `{"name":"GitLab","description":"GitLab integration for Slack","short_description":"GitLab for Slack","guidelines":"","url":"https://gitlab.com","oauth_config":{"scopes":{"bot":["channels:history","channels:read","chat:write","commands","files:read","files:write","groups:history","groups:read","im:history","im:read","im:write","mpim:history","mpim:read","mpim:write","reactions:read","reactions:write","team:read","usergroups:read","users:read","users:read.email"]}},"settings":{"org_deploy_enabled":false,"socket_mode_enabled":false,"is_hosted":false},"features":{"bot_user":{"display_name":"GitLab","always_online":true},"oauth_config":{"scopes":{"bot":["channels:history","channels:read","chat:write","commands","files:read","files:write","groups:history","groups:read","im:history","im:read","im:write","mpim:history","mpim:read","mpim:write","reactions:read","reactions:write","team:read","usergroups:read","users:read","users:read.email"]}}}`
	ctx.Header("Content-Disposition", "attachment; filename=slack_manifest.json")
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, manifest)
}

// UsageData displays the usage data
func (c *ApplicationSettingsController) UsageData(ctx *gin.Context) {
	// TODO: Implement usage data retrieval
	ctx.HTML(http.StatusOK, "admin/application_settings/usage_data", gin.H{
		"usage_data": "{}",
	})
}

// Helper methods

func (c *ApplicationSettingsController) performUpdate(ctx *gin.Context, action string) {
	var params models.ApplicationSettingParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings := models.GetCurrentApplicationSettings()
	service := services.NewUpdateApplicationSettingsService(settings, c.GetCurrentUser(ctx), params)
	success, err := service.Execute()
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if success {
		ctx.Set("notice", "Application settings saved successfully")
		ctx.Redirect(http.StatusFound, "/admin/application_settings/"+action)
	} else {
		ctx.Set("error", "Application settings update failed")
		ctx.HTML(http.StatusUnprocessableEntity, "admin/application_settings/"+action, gin.H{
			"settings": settings,
		})
	}
}
