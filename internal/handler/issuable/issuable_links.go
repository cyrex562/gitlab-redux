package issuable

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// IssuableLinks handles links between issuable items
type IssuableLinks struct {
	issuableService *service.IssuableService
	logger          *service.Logger
}

// NewIssuableLinks creates a new instance of IssuableLinks
func NewIssuableLinks(
	issuableService *service.IssuableService,
	logger *service.Logger,
) *IssuableLinks {
	return &IssuableLinks{
		issuableService: issuableService,
		logger:          logger,
	}
}

// Index handles the index action for issuable links
func (i *IssuableLinks) Index(ctx *gin.Context) {
	// Get the issuables
	issuables, err := i.GetIssuables(ctx)
	if err != nil {
		i.logger.Error("Failed to get issuables", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issuables"})
		return
	}

	// Return the issuables as JSON
	ctx.JSON(http.StatusOK, issuables)
}

// Create handles the create action for issuable links
func (i *IssuableLinks) Create(ctx *gin.Context) {
	// Get the create parameters
	params, err := i.GetCreateParams(ctx)
	if err != nil {
		i.logger.Error("Failed to get create parameters", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Execute the create service
	result, err := i.GetCreateService(ctx).Execute(ctx, params)
	if err != nil {
		i.logger.Error("Failed to create issuable links", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create issuable links"})
		return
	}

	// Get the issuables
	issuables, err := i.GetIssuables(ctx)
	if err != nil {
		i.logger.Error("Failed to get issuables", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issuables"})
		return
	}

	// Return the result and issuables as JSON
	ctx.JSON(result.HTTPStatus, gin.H{
		"message":   result.Message,
		"issuables": issuables,
	})
}

// Destroy handles the destroy action for issuable links
func (i *IssuableLinks) Destroy(ctx *gin.Context) {
	// Execute the destroy service
	result, err := i.GetDestroyService(ctx).Execute(ctx)
	if err != nil {
		i.logger.Error("Failed to destroy issuable links", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to destroy issuable links"})
		return
	}

	// Get the issuables
	issuables, err := i.GetIssuables(ctx)
	if err != nil {
		i.logger.Error("Failed to get issuables", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issuables"})
		return
	}

	// Return the result and issuables as JSON
	ctx.JSON(result.HTTPStatus, gin.H{
		"issuables": issuables,
	})
}

// GetIssuables gets the list of issuables
func (i *IssuableLinks) GetIssuables(ctx *gin.Context) ([]model.Issuable, error) {
	// Get the list service
	listService := i.GetListService(ctx)
	if listService == nil {
		return nil, fmt.Errorf("list service not implemented")
	}

	// Execute the list service
	return listService.Execute(ctx)
}

// GetListService gets the list service
// This method should be overridden by the implementing struct
func (i *IssuableLinks) GetListService(ctx *gin.Context) *service.IssuableListService {
	// This is a placeholder that should be overridden
	return nil
}

// GetCreateParams gets the create parameters
func (i *IssuableLinks) GetCreateParams(ctx *gin.Context) (map[string]interface{}, error) {
	// Get the issuable references from the request
	var request struct {
		IssuableReferences []string `json:"issuable_references"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		return nil, err
	}

	// Return the parameters
	return map[string]interface{}{
		"issuable_references": request.IssuableReferences,
	}, nil
}

// GetCreateService gets the create service
// This method should be overridden by the implementing struct
func (i *IssuableLinks) GetCreateService(ctx *gin.Context) *service.IssuableCreateService {
	// This is a placeholder that should be overridden
	return nil
}

// GetDestroyService gets the destroy service
// This method should be overridden by the implementing struct
func (i *IssuableLinks) GetDestroyService(ctx *gin.Context) *service.IssuableDestroyService {
	// This is a placeholder that should be overridden
	return nil
}
