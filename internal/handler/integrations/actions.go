package integrations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// Actions provides integration management functionality
type Actions struct {
	integrationService *service.IntegrationService
	paramsHandler     *Params
}

// NewActions creates a new instance of Actions
func NewActions(integrationService *service.IntegrationService, paramsHandler *Params) *Actions {
	return &Actions{
		integrationService: integrationService,
		paramsHandler:     paramsHandler,
	}
}

// RegisterRoutes registers the routes for integration actions
func (a *Actions) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/:id/edit", a.ensureIntegrationEnabled, a.edit)
	router.PUT("/:id", a.ensureIntegrationEnabled, a.update)
	router.POST("/:id/test", a.ensureIntegrationEnabled, a.test)
	router.POST("/:id/reset", a.reset)
}

// edit handles the GET /:id/edit endpoint for editing an integration
func (a *Actions) edit(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "shared/integrations/edit", gin.H{
		"integration": ctx.MustGet("integration"),
	})
}

// update handles the PUT /:id endpoint for updating an integration
func (a *Actions) update(ctx *gin.Context) {
	integration := ctx.MustGet("integration").(*model.Integration)
	params, err := a.paramsHandler.ParseParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters",
			"errors":  err.Error(),
		})
		return
	}

	saved, err := a.integrationService.UpdateIntegration(ctx, integration, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update integration",
			"errors":  err.Error(),
		})
		return
	}

	// Handle different response formats
	switch ctx.GetHeader("Accept") {
	case "application/json":
		status := http.StatusOK
		if !saved {
			status = http.StatusUnprocessableEntity
		}
		ctx.JSON(status, a.serializeAsJSON(integration))
	default:
		if saved {
			// Propagate integration changes
			go a.integrationService.PropagateIntegration(ctx, integration.ID)
			ctx.Redirect(http.StatusSeeOther, a.getScopedEditPath(integration))
			ctx.SetCookie("notice", a.getSuccessMessage(integration), 3600, "/", "", false, true)
		} else {
			ctx.HTML(http.StatusUnprocessableEntity, "shared/integrations/edit", gin.H{
				"integration": integration,
			})
		}
	}
}

// test handles the POST /:id/test endpoint for testing an integration
func (a *Actions) test(ctx *gin.Context) {
	integration := ctx.MustGet("integration").(*model.Integration)
	params, err := a.paramsHandler.ParseParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters",
			"errors":  err.Error(),
		})
		return
	}

	if !integration.IsTestable() {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	// Update integration with test parameters
	integration.UpdateAttributes(params)

	// Test integration based on level
	var result *model.TestResult
	switch {
	case integration.IsProjectLevel():
		result, err = a.integrationService.TestProjectIntegration(ctx, integration, ctx.MustGet("user").(*model.User), ctx.Query("event"))
	case integration.IsGroupLevel():
		result, err = a.integrationService.TestGroupIntegration(ctx, integration, ctx.MustGet("user").(*model.User), ctx.Query("event"))
	default:
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to test integration",
			"errors":  err.Error(),
		})
		return
	}

	if !result.Success {
		ctx.JSON(http.StatusOK, gin.H{
			"error":           true,
			"message":         "Connection failed. Check your integration settings.",
			"service_response": result.Result,
			"test_failed":     true,
		})
		return
	}

	if result.Data != nil {
		ctx.JSON(http.StatusOK, result.Data)
	} else {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

// reset handles the POST /:id/reset endpoint for resetting an integration
func (a *Actions) reset(ctx *gin.Context) {
	integration := ctx.MustGet("integration").(*model.Integration)

	if integration.IsManualActivation() {
		err := a.integrationService.DestroyIntegration(ctx, integration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to reset integration",
				"errors":  err.Error(),
			})
			return
		}

		ctx.SetCookie("notice", "This integration, and inheriting projects were reset.", 3600, "/", "", false, true)
		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Integration cannot be reset.",
		})
	}
}

// ensureIntegrationEnabled middleware ensures the integration exists and is enabled
func (a *Actions) ensureIntegrationEnabled(ctx *gin.Context) {
	integration, err := a.integrationService.FindOrInitializeIntegration(ctx, ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Check if integration is Prometheus and feature flag is enabled
	if integration.Type == "prometheus" && a.integrationService.IsFeatureEnabled(ctx, "remove_monitor_metrics") {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Set("integration", integration)
	ctx.Next()
}

// getSuccessMessage returns the success message for an integration update
func (a *Actions) getSuccessMessage(integration *model.Integration) string {
	if integration.IsActive() {
		return integration.Title + " settings saved and active."
	}
	return integration.Title + " settings saved, but not active."
}

// getScopedEditPath returns the scoped edit path for an integration
func (a *Actions) getScopedEditPath(integration *model.Integration) string {
	if integration.ProjectID != 0 {
		return "/projects/" + integration.ProjectID + "/integrations/" + integration.ID + "/edit"
	}
	if integration.GroupID != 0 {
		return "/groups/" + integration.GroupID + "/integrations/" + integration.ID + "/edit"
	}
	return "/admin/integrations/" + integration.ID + "/edit"
}

// serializeAsJSON serializes an integration as JSON
func (a *Actions) serializeAsJSON(integration *model.Integration) gin.H {
	return gin.H{
		"id":          integration.ID,
		"type":        integration.Type,
		"active":      integration.IsActive(),
		"properties":  integration.Properties,
		"errors":      integration.Errors,
	}
}
