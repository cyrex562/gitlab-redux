package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitlab-org/gitlab-redux/internal/controllers"
	"github.com/gitlab-org/gitlab-redux/internal/models"
	"github.com/gitlab-org/gitlab-redux/internal/services"
)

// ApplicationsController handles OAuth application management in the admin interface
type ApplicationsController struct {
	controllers.BaseController
	applicationService *services.ApplicationService
}

// NewApplicationsController creates a new instance of ApplicationsController
func NewApplicationsController(applicationService *services.ApplicationService) *ApplicationsController {
	return &ApplicationsController{
		applicationService: applicationService,
	}
}

// Index lists all OAuth applications
func (c *ApplicationsController) Index(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	applications, totalCount, err := c.applicationService.FindAll(cursor)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"applications": applications,
		"total_count": totalCount,
	})
}

// Show displays a single OAuth application
func (c *ApplicationsController) Show(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	application, err := c.applicationService.FindByID(id)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, application)
}

// New renders the form for creating a new OAuth application
func (c *ApplicationsController) New(ctx *gin.Context) {
	scopes := c.applicationService.GetAvailableScopes()
	ctx.JSON(http.StatusOK, gin.H{
		"scopes": scopes,
	})
}

// Create creates a new OAuth application
func (c *ApplicationsController) Create(ctx *gin.Context) {
	var params models.ApplicationParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		c.HandleError(ctx, err)
		return
	}

	application, err := c.applicationService.Create(params, ctx.Request)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"application": application,
		"message":     "Application was successfully created.",
	})
}

// Update updates an existing OAuth application
func (c *ApplicationsController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	var params models.ApplicationParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		c.HandleError(ctx, err)
		return
	}

	application, err := c.applicationService.Update(id, params)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"application": application,
		"message":     "Application was successfully updated.",
	})
}

// RenewSecret generates a new secret for an OAuth application
func (c *ApplicationsController) RenewSecret(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	secret, err := c.applicationService.RenewSecret(id)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"secret": secret,
	})
}

// Destroy deletes an OAuth application
func (c *ApplicationsController) Destroy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	if err := c.applicationService.Delete(id); err != nil {
		c.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusFound, gin.H{
		"message": "Application was successfully destroyed.",
	})
}

// ResetWebIdeOAuthApplicationSettings resets the Web IDE OAuth application settings
func (c *ApplicationsController) ResetWebIdeOAuthApplicationSettings(ctx *gin.Context) {
	success := c.applicationService.ResetWebIdeOAuthApplicationSettings()
	if !success {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
