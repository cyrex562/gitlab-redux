package harbor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// Tag provides Harbor tag listing and querying functionality
type Tag struct {
	harborService *service.HarborService
}

// NewTag creates a new instance of Tag
func NewTag(harborService *service.HarborService) *Tag {
	return &Tag{
		harborService: harborService,
	}
}

// RegisterRoutes registers the routes for Harbor tags
func (t *Tag) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", t.index)
}

// QueryParams represents the query parameters for tag listing
type QueryParams struct {
	RepositoryID *string `form:"repository_id"`
	ArtifactID   *string `form:"artifact_id"`
	Sort         *string `form:"sort"`
	Page         *int    `form:"page"`
	Limit        *int    `form:"limit"`
}

// index handles the GET / endpoint for listing tags
func (t *Tag) index(ctx *gin.Context) {
	// Parse query parameters
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters",
			"errors":  err.Error(),
		})
		return
	}

	// Get the container from context (should be set by previous middleware)
	container, exists := ctx.Get("container")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Container not found"})
		return
	}

	// Create query
	query, err := t.harborService.NewQuery(ctx, container, &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create query",
			"errors":  err.Error(),
		})
		return
	}

	// Validate query
	if !query.IsValid() {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid parameters",
			"errors":  query.GetErrors(),
		})
		return
	}

	// Get tags
	tags, err := query.GetTags()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get tags",
			"errors":  err.Error(),
		})
		return
	}

	// Get Harbor integration details
	integration, err := t.harborService.GetIntegration(ctx, container)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get Harbor integration",
			"errors":  err.Error(),
		})
		return
	}

	// Serialize response with pagination
	serializer := model.NewHarborTagSerializer()
	response, err := serializer.SerializeWithPagination(ctx, tags, &model.HarborSerializerOptions{
		URL:         integration.URL,
		ProjectName: integration.ProjectName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to serialize tags",
			"errors":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
