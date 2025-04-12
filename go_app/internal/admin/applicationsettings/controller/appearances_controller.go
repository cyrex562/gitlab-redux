package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/routing"
)

// AppearancesController handles appearance settings in the admin panel
type AppearancesController struct {
	*routing.BaseController
}

// NewAppearancesController creates a new instance of AppearancesController
func NewAppearancesController() *AppearancesController {
	return &AppearancesController{
		BaseController: routing.NewBaseController(),
	}
}

// SetupRoutes configures the routes for the appearances controller
func (c *AppearancesController) SetupRoutes(router *gin.Engine) {
	admin := router.Group("/admin/application_settings/appearances")
	{
		admin.GET("", c.Show)
		admin.GET("/preview_sign_in", c.PreviewSignIn)
		admin.POST("", c.Create)
		admin.PUT("", c.Update)
		admin.DELETE("/logo", c.RemoveLogo)
		admin.DELETE("/header_logos", c.RemoveHeaderLogo)
		admin.DELETE("/pwa_icon", c.RemovePWAIcon)
		admin.DELETE("/favicon", c.RemoveFavicon)
	}
}

// Show displays the appearance settings page
func (c *AppearancesController) Show(ctx *gin.Context) {
	appearance := models.GetCurrentAppearance()
	if appearance == nil {
		appearance = &models.Appearance{}
	}
	ctx.HTML(http.StatusOK, "admin/appearances/show", gin.H{
		"appearance": appearance,
	})
}

// PreviewSignIn shows the sign-in page preview
func (c *AppearancesController) PreviewSignIn(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "devise/sessions/new", nil)
}

// Create creates a new appearance setting
func (c *AppearancesController) Create(ctx *gin.Context) {
	var appearance models.Appearance
	if err := ctx.ShouldBind(&appearance); err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, "admin/appearances/show", gin.H{
			"appearance": &appearance,
			"error":     err.Error(),
		})
		return
	}

	if err := appearance.Save(); err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, "admin/appearances/show", gin.H{
			"appearance": &appearance,
			"error":     err.Error(),
		})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/application_settings/appearances")
	ctx.Set("notice", "Appearance was successfully created.")
}

// Update updates the appearance settings
func (c *AppearancesController) Update(ctx *gin.Context) {
	appearance := models.GetCurrentAppearance()
	if appearance == nil {
		appearance = &models.Appearance{}
	}

	if err := ctx.ShouldBind(appearance); err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, "admin/appearances/show", gin.H{
			"appearance": appearance,
			"error":     err.Error(),
		})
		return
	}

	if err := appearance.Save(); err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, "admin/appearances/show", gin.H{
			"appearance": appearance,
			"error":     err.Error(),
		})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/application_settings/appearances")
	ctx.Set("notice", "Appearance was successfully updated.")
}

// RemoveLogo removes the logo from appearance settings
func (c *AppearancesController) RemoveLogo(ctx *gin.Context) {
	appearance := models.GetCurrentAppearance()
	if appearance != nil {
		appearance.RemoveLogo()
		appearance.Save()
	}
	ctx.Redirect(http.StatusFound, "/admin/application_settings/appearances")
	ctx.Set("notice", "Logo was successfully removed.")
}

// RemoveHeaderLogo removes the header logo from appearance settings
func (c *AppearancesController) RemoveHeaderLogo(ctx *gin.Context) {
	appearance := models.GetCurrentAppearance()
	if appearance != nil {
		appearance.RemoveHeaderLogo()
		appearance.Save()
	}
	ctx.Redirect(http.StatusFound, "/admin/application_settings/appearances")
	ctx.Set("notice", "Header logo was successfully removed.")
}

// RemovePWAIcon removes the PWA icon from appearance settings
func (c *AppearancesController) RemovePWAIcon(ctx *gin.Context) {
	appearance := models.GetCurrentAppearance()
	if appearance != nil {
		appearance.RemovePWAIcon()
		appearance.Save()
	}
	ctx.Redirect(http.StatusFound, "/admin/application_settings/appearances")
	ctx.Set("notice", "Progressive Web App (PWA) icon was successfully removed.")
}

// RemoveFavicon removes the favicon from appearance settings
func (c *AppearancesController) RemoveFavicon(ctx *gin.Context) {
	appearance := models.GetCurrentAppearance()
	if appearance != nil {
		appearance.RemoveFavicon()
		appearance.Save()
	}
	ctx.Redirect(http.StatusFound, "/admin/application_settings/appearances")
	ctx.Set("notice", "Favicon was successfully removed.")
}
