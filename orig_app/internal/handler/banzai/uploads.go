package banzai

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// UploadsController handles file uploads for projects and groups
type UploadsController struct {
	uploadService *service.UploadService
	authService  *service.AuthService
}

// NewUploadsController creates a new instance of UploadsController
func NewUploadsController(uploadService *service.UploadService, authService *service.AuthService) *UploadsController {
	return &UploadsController{
		uploadService: uploadService,
		authService:  authService,
	}
}

// RegisterRoutes registers the routes for the UploadsController
func (c *UploadsController) RegisterRoutes(r *gin.RouterGroup) {
	uploads := r.Group("/uploads")
	{
		uploads.POST("/:model/:model_id", c.handleUpload)
	}
}

// handleUpload handles file uploads for both projects and groups
func (c *UploadsController) handleUpload(ctx *gin.Context) {
	// Verify upload model class
	modelClass, err := c.verifyUploadModelClass(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid upload model"})
		return
	}

	// Find the model (project or group)
	model, err := c.findModel(ctx, modelClass)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	// Check authorization
	if err := c.authorizeAccess(ctx, model); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Handle the file upload
	if err := c.uploadService.HandleUpload(ctx, model); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Upload successful"})
}

// verifyUploadModelClass verifies that the upload model class is valid
func (c *UploadsController) verifyUploadModelClass(ctx *gin.Context) (string, error) {
	modelType := ctx.Param("model")
	switch modelType {
	case "project", "group":
		return modelType, nil
	default:
		return "", fmt.Errorf("invalid model type")
	}
}

// findModel finds the model (project or group) by ID
func (c *UploadsController) findModel(ctx *gin.Context, modelClass string) (interface{}, error) {
	modelID := ctx.Param("model_id")
	switch modelClass {
	case "project":
		return c.uploadService.GetProject(ctx, modelID)
	case "group":
		return c.uploadService.GetGroup(ctx, modelID)
	default:
		return nil, fmt.Errorf("invalid model class")
	}
}

// authorizeAccess checks if the user has permission to access the model
func (c *UploadsController) authorizeAccess(ctx *gin.Context, model interface{}) error {
	// Skip auth checks if configured to bypass
	if c.uploadService.ShouldBypassAuthChecks() {
		return nil
	}

	// Get current user from context
	user, exists := ctx.Get("current_user")
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Check permissions based on model type
	switch m := model.(type) {
	case *model.Project:
		return c.authService.CheckProjectAccess(ctx, user.(*model.User), m)
	case *model.Group:
		return c.authService.CheckGroupAccess(ctx, user.(*model.User), m)
	default:
		return fmt.Errorf("invalid model type")
	}
}
