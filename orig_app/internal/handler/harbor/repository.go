package harbor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// Repository provides Harbor repository listing and querying functionality
type Repository struct {
	harborService *service.HarborService
}

// NewRepository creates a new instance of Repository
func NewRepository(harborService *service.HarborService) *Repository {
	return &Repository{
		harborService: harborService,
	}
}

// RegisterRoutes registers the routes for Harbor repositories
func (r *Repository) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", r.index)
	router.GET("/:id", r.show)
}

// QueryParams represents the query parameters for repository listing
type QueryParams struct {
	Search *string `form:"search"`
	Sort   *string `form:"sort"`
	Page   *int    `form:"page"`
	Limit  *int    `form:"limit"`
}

// index handles the GET / endpoint for listing repositories
func (r *Repository) index(ctx *gin.Context) {
	// Handle HTML format
	if ctx.GetHeader("Accept") == "text/html" {
		// TODO: Implement HTML rendering
		// This should:
		// 1. Render the repository index template
		// 2. Return the rendered HTML
		return
	}

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
	query, err := r.harborService.NewQuery(ctx, container, &params)
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

	// Get repositories
	repositories, err := query.GetRepositories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get repositories",
			"errors":  err.Error(),
		})
		return
	}

	// Get Harbor integration details
	integration, err := r.harborService.GetIntegration(ctx, container)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get Harbor integration",
			"errors":  err.Error(),
		})
		return
	}

	// Serialize response with pagination
	serializer := model.NewHarborRepositorySerializer()
	response, err := serializer.SerializeWithPagination(ctx, repositories, &model.HarborSerializerOptions{
		URL:         integration.URL,
		ProjectName: integration.ProjectName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to serialize repositories",
			"errors":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// show handles the GET /:id endpoint for showing a repository
func (r *Repository) show(ctx *gin.Context) {
	// For frontend routing support, render the index template
	// TODO: Implement HTML rendering
	// This should:
	// 1. Render the repository index template
	// 2. Return the rendered HTML
}
