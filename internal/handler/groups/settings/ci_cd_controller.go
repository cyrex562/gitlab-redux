package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CiCdController handles group CI/CD settings
type CiCdController struct {
	ciCdSettingsService *CiCdSettingsService
	autoDevOpsService   *AutoDevOpsService
}

// NewCiCdController creates a new CiCdController
func NewCiCdController(
	ciCdSettingsService *CiCdSettingsService,
	autoDevOpsService *AutoDevOpsService,
) *CiCdController {
	return &CiCdController{
		ciCdSettingsService: ciCdSettingsService,
		autoDevOpsService:   autoDevOpsService,
	}
}

// RegisterRoutes registers the routes for the CiCdController
func (c *CiCdController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", c.Show)
	router.PUT("/", c.Update)
	router.PUT("/auto_devops", c.UpdateAutoDevOps)
}

// RegisterMiddleware registers the middleware for the CiCdController
func (c *CiCdController) RegisterMiddleware(router *gin.RouterGroup) {
	router.Use(c.AuthorizeAdminGroup())
	router.Use(c.SetVariableLimit())
}

// Show handles the show action
func (c *CiCdController) Show(ctx *gin.Context) {
	// Get the group from the context
	group, _ := ctx.Get("group")

	// Get the entity limit
	entityLimit, _ := ctx.Get("entity_limit")

	// Render the show template
	ctx.HTML(http.StatusOK, "groups/settings/ci_cd/show", gin.H{
		"group": group,
		"entity_limit": entityLimit,
	})
}

// Update handles the update action
func (c *CiCdController) Update(ctx *gin.Context) {
	// Get the group from the context
	group, _ := ctx.Get("group")

	// Get the CI/CD settings parameters
	settingsParams := c.GetCiCdSettingsParams(ctx)

	// Update the CI/CD settings
	err := c.ciCdSettingsService.Execute(group, settingsParams)
	if err != nil {
		ctx.SetFlash("alert", "Failed to update CI/CD settings.")
		ctx.Redirect(http.StatusFound, "/groups/:id/settings/ci_cd")
		return
	}

	// Set flash notice
	ctx.SetFlash("notice", "CI/CD settings were successfully updated.")

	// Redirect to the show page
	ctx.Redirect(http.StatusFound, "/groups/:id/settings/ci_cd")
}

// UpdateAutoDevOps handles the update auto devops action
func (c *CiCdController) UpdateAutoDevOps(ctx *gin.Context) {
	// Get the group from the context
	group, _ := ctx.Get("group")

	// Get the auto devops parameters
	autoDevOpsParams := c.GetAutoDevOpsParams(ctx)

	// Update the auto devops settings
	err := c.autoDevOpsService.Execute(group, autoDevOpsParams)
	if err != nil {
		ctx.SetFlash("alert", "Failed to update Auto DevOps settings.")
		ctx.Redirect(http.StatusFound, "/groups/:id/settings/ci_cd")
		return
	}

	// Set flash notice
	ctx.SetFlash("notice", "Auto DevOps settings were successfully updated.")

	// Redirect to the show page
	ctx.Redirect(http.StatusFound, "/groups/:id/settings/ci_cd")
}

// GetCiCdSettingsParams gets the CI/CD settings parameters
func (c *CiCdController) GetCiCdSettingsParams(ctx *gin.Context) map[string]interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would get the CI/CD settings parameters
	return map[string]interface{}{
		"group_runners_enabled": ctx.PostForm("group_runners_enabled") == "true",
		"shared_runners_enabled": ctx.PostForm("shared_runners_enabled") == "true",
	}
}

// GetAutoDevOpsParams gets the auto devops parameters
func (c *CiCdController) GetAutoDevOpsParams(ctx *gin.Context) map[string]interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would get the auto devops parameters
	return map[string]interface{}{
		"auto_devops_enabled": ctx.PostForm("auto_devops_enabled") == "true",
		"auto_devops_deploy_strategy": ctx.PostForm("auto_devops_deploy_strategy"),
	}
}

// AuthorizeAdminGroup middleware checks if the user has permission to admin the group
func (c *CiCdController) AuthorizeAdminGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			ctx.Abort()
			return
		}

		// Get the group from the context
		group, exists := ctx.Get("group")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			ctx.Abort()
			return
		}

		// Check if the user has permission to admin the group
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the user has the admin_group permission
		canAdminGroup := true // Replace with actual check

		if !canAdminGroup {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "You don't have permission to admin the group"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// SetVariableLimit middleware sets the variable limit
func (c *CiCdController) SetVariableLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the group from the context
		group, _ := ctx.Get("group")

		// Get the variable limit
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would get the variable limit from the group
		variableLimit := 0 // Replace with actual implementation

		// Set the variable limit
		ctx.Set("entity_limit", variableLimit)

		ctx.Next()
	}
}

// CiCdSettingsService handles CI/CD settings
type CiCdSettingsService struct {
	// Add fields as needed
}

// Execute executes the service
func (s *CiCdSettingsService) Execute(group interface{}, params map[string]interface{}) error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would update the CI/CD settings
	return nil
}

// AutoDevOpsService handles Auto DevOps settings
type AutoDevOpsService struct {
	// Add fields as needed
}

// Execute executes the service
func (s *AutoDevOpsService) Execute(group interface{}, params map[string]interface{}) error {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would update the Auto DevOps settings
	return nil
}
