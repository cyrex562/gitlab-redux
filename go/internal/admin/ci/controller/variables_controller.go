package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/routing"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services"
)

// VariablesController handles CI variables in the admin panel
type VariablesController struct {
	*routing.BaseController
}

// NewVariablesController creates a new instance of VariablesController
func NewVariablesController() *VariablesController {
	return &VariablesController{
		BaseController: routing.NewBaseController(),
	}
}

// SetupRoutes configures the routes for the variables controller
func (c *VariablesController) SetupRoutes(router *gin.Engine) {
	admin := router.Group("/admin/ci/variables")
	{
		admin.GET("", c.Show)
		admin.PUT("", c.Update)
	}
}

// Show displays the CI variables
func (c *VariablesController) Show(ctx *gin.Context) {
	variables := models.GetAllInstanceVariables()
	ctx.JSON(http.StatusOK, gin.H{
		"variables": variables,
	})
}

// Update updates the CI variables
func (c *VariablesController) Update(ctx *gin.Context) {
	var params struct {
		VariablesAttributes []models.InstanceVariable `json:"variables_attributes"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := services.NewUpdateInstanceVariablesService(params.VariablesAttributes, c.GetCurrentUser(ctx))
	if err := service.Execute(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variables := models.GetAllInstanceVariables()
	ctx.JSON(http.StatusOK, gin.H{
		"variables": variables,
	})
}
