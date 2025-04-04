package harbor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// Artifact provides Harbor artifact listing and querying functionality
type Artifact struct {
	harborService *service.HarborService
}

// NewArtifact creates a new instance of Artifact
func NewArtifact(harborService *service.HarborService) *Artifact {
	return &Artifact{
		harborService: harborService,
	}
}

// RegisterRoutes registers the routes for Harbor artifacts
func (a *Artifact) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/", a.index)
}

// QueryParams represents the query parameters for artifact listing
type QueryParams struct {
	RepositoryID *int64  `form:"repository_id"`
	Search       *string `form:"search"`
	Sort         *string `form:"sort"`
	Page         *int    `form:"page"`
	Limit        *int    `form:"limit"`
}

// index handles the GET / endpoint for listing artifacts
func (a *Artifact) index(ctx *gin.Context) {
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
	query, err := a.harborService.NewQuery(ctx, container, &params)
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

	// Get artifacts
	artifacts, err := query.GetArtifacts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get artifacts",
			"errors":  err.Error(),
		})
		return
	}

	// Serialize response with pagination
	serializer := model.NewHarborArtifactSerializer()
	response, err := serializer.SerializeWithPagination(ctx, artifacts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to serialize artifacts",
			"errors":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
