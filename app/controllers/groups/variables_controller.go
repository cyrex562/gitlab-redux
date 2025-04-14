package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/serializers"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/ci"
)

// VariablesController handles requests for group CI/CD variables
type VariablesController struct {
	changeVariablesService *ci.ChangeVariablesService
	groupVariableSerializer *serializers.GroupVariableSerializer
}

// NewVariablesController creates a new variables controller
func NewVariablesController(
	changeVariablesService *ci.ChangeVariablesService,
	groupVariableSerializer *serializers.GroupVariableSerializer,
) *VariablesController {
	return &VariablesController{
		changeVariablesService: changeVariablesService,
		groupVariableSerializer: groupVariableSerializer,
	}
}

// RegisterRoutes registers the routes for the variables controller
func (c *VariablesController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/variables", c.authorizeAdminGroup(), c.Show)
	router.PUT("/variables", c.authorizeAdminCicdVariables(), c.Update)
}

// Show handles GET requests for group variables
func (c *VariablesController) Show(ctx *gin.Context) {
	group := ctx.MustGet("group").(*models.Group)

	// Serialize the variables
	variables := c.groupVariableSerializer.Represent(group.Variables)

	ctx.JSON(http.StatusOK, gin.H{
		"variables": variables,
	})
}

// Update handles PUT requests for updating group variables
func (c *VariablesController) Update(ctx *gin.Context) {
	group := ctx.MustGet("group").(*models.Group)
	user := ctx.MustGet("current_user").(*models.User)

	// Get variables parameters
	var params ci.VariablesParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update variables
	updateResult, err := c.changeVariablesService.Execute(ctx, group, user, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateResult {
		// Render group variables
		c.renderGroupVariables(ctx, group)
	} else {
		// Render error
		c.renderError(ctx, group)
	}
}

// Helper methods for middleware and authorization

func (c *VariablesController) authorizeAdminGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)

		if !user.CanAdminGroup(group) {
			ctx.Status(http.StatusForbidden)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c *VariablesController) authorizeAdminCicdVariables() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("current_user").(*models.User)
		group := ctx.MustGet("group").(*models.Group)

		if !user.CanAdminCicdVariables(group) {
			ctx.Status(http.StatusForbidden)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// Helper methods for rendering

func (c *VariablesController) renderGroupVariables(ctx *gin.Context, group *models.Group) {
	// Serialize the variables
	variables := c.groupVariableSerializer.Represent(group.Variables)

	ctx.JSON(http.StatusOK, gin.H{
		"variables": variables,
	})
}

func (c *VariablesController) renderError(ctx *gin.Context, group *models.Group) {
	ctx.JSON(http.StatusBadRequest, group.Errors)
} 