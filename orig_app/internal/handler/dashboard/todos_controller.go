package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/feature_flags"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// TodosController handles todo-related actions in the dashboard
type TodosController struct {
	ApplicationController
	analyticsService service.AnalyticsService
	featureFlagService service.FeatureFlagService
}

// NewTodosController creates a new TodosController
func NewTodosController(
	appController ApplicationController,
	analyticsService service.AnalyticsService,
	featureFlagService service.FeatureFlagService,
) *TodosController {
	return &TodosController{
		ApplicationController: appController,
		analyticsService:     analyticsService,
		featureFlagService:   featureFlagService,
	}
}

// RegisterRoutes registers the routes for the TodosController
func (c *TodosController) RegisterRoutes(router *gin.Engine) {
	dashboard := router.Group("/dashboard")
	{
		dashboard.GET("/todos", c.Index)
	}
}

// Index handles the index action for todos
func (c *TodosController) Index(ctx *gin.Context) {
	// Get current user
	currentUser, err := c.GetCurrentUser(ctx)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Track internal event
	err = c.analyticsService.TrackInternalEvent(ctx, "view_todo_list", currentUser)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Push frontend feature flag
	featureFlag := feature_flags.PushFrontendFeatureFlag(ctx, "todos_bulk_actions", currentUser)
	ctx.Set("todos_bulk_actions", featureFlag)

	// Render the view
	ctx.HTML(http.StatusOK, "dashboard/todos/index", gin.H{
		"todos_bulk_actions": featureFlag,
	})
}
