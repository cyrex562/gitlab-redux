package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/labels"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// LabelsController handles requests related to labels in the dashboard
type LabelsController struct {
	*ApplicationController
	labelService service.LabelService
}

// NewLabelsController creates a new LabelsController
func NewLabelsController(appController *ApplicationController, labelService service.LabelService) *LabelsController {
	return &LabelsController{
		ApplicationController: appController,
		labelService:         labelService,
	}
}

// Index handles GET requests to /dashboard/labels
// It returns a JSON response with distinct labels based on project IDs
func (c *LabelsController) Index(ctx *gin.Context) {
	// Get the current user from the context
	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		c.RenderError(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get project IDs from the context
	projectIDs, err := c.GetProjectIDs(ctx)
	if err != nil {
		c.RenderError(ctx, http.StatusBadRequest, "Invalid project IDs")
		return
	}

	// Find distinct labels based on project IDs
	foundLabels, err := c.labelService.FindDistinctLabelsByProjects(ctx, user, projectIDs)
	if err != nil {
		c.RenderError(ctx, http.StatusInternalServerError, "Failed to retrieve labels")
		return
	}

	// Serialize the labels for appearance
	serializedLabels := labels.SerializeAppearance(foundLabels)

	// Return JSON response
	ctx.JSON(http.StatusOK, serializedLabels)
}

// RegisterRoutes registers the routes for the LabelsController
func (c *LabelsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/labels", c.Index)
}
