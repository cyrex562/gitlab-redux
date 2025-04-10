package issuable

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// IssuableCollectionsAction handles collections of issuable items
type IssuableCollectionsAction struct {
	issuableCollections *IssuableCollections
	issuesCalendar *IssuesCalendar
	rateLimitService *service.RateLimitService
	issuableService *service.IssuableService
	userService *service.UserService
	metadataService *service.IssuableMetadataService
	logger *service.Logger
}

// NewIssuableCollectionsAction creates a new instance of IssuableCollectionsAction
func NewIssuableCollectionsAction(
	issuableCollections *IssuableCollections,
	issuesCalendar *IssuesCalendar,
	rateLimitService *service.RateLimitService,
	issuableService *service.IssuableService,
	userService *service.UserService,
	metadataService *service.IssuableMetadataService,
	logger *service.Logger,
) *IssuableCollectionsAction {
	return &IssuableCollectionsAction{
		issuableCollections: issuableCollections,
		issuesCalendar: issuesCalendar,
		rateLimitService: rateLimitService,
		issuableService: issuableService,
		userService: userService,
		metadataService: metadataService,
		logger: logger,
	}
}

// SetupRoutes sets up the routes for the IssuableCollectionsAction
func (i *IssuableCollectionsAction) SetupRoutes(router *gin.RouterGroup) {
	// Check search rate limit middleware
	router.Use(i.checkSearchRateLimitMiddleware())

	// Set up routes
	router.GET("/issues", i.Issues)
	router.GET("/merge_requests", i.MergeRequests)
	router.GET("/issues_calendar", i.IssuesCalendar)
}

// CheckSearchRateLimitMiddleware creates a middleware to check the search rate limit
func (i *IssuableCollectionsAction) checkSearchRateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the search parameter is present
		if ctx.DefaultQuery("search", "") != "" {
			// Check the search rate limit
			if err := i.rateLimitService.CheckSearchRateLimit(ctx); err != nil {
				// Handle the error
				i.handleError(ctx, err)
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

// Issues handles the issues action
func (i *IssuableCollectionsAction) Issues(ctx *gin.Context) {
	// Check if the request is for HTML or Atom
	format := ctx.DefaultQuery("format", "html")
	if format == "html" {
		// Render the HTML response
		ctx.HTML(http.StatusOK, "issues", nil)
	} else if format == "atom" {
		// Get the current user from the context
		currentUser, err := i.userService.GetCurrentUser(ctx)
		if err != nil {
			i.handleError(ctx, err)
			return
		}

		// Get the issuables collection
		issuables, err := i.issuableCollections.GetIssuablesCollection(ctx)
		if err != nil {
			i.handleError(ctx, err)
			return
		}

		// Filter out archived issuables
		nonArchivedIssuables := i.filterNonArchived(issuables)

		// Get the page parameter
		page := ctx.DefaultQuery("page", "1")

		// Paginate the issuables
		paginatedIssuables, err := i.paginateIssuables(nonArchivedIssuables, page)
		if err != nil {
			i.handleError(ctx, err)
			return
		}

		// Get the issuable metadata
		issuableMetadata, err := i.metadataService.GetIssuableMetadata(currentUser, paginatedIssuables)
		if err != nil {
			i.handleError(ctx, err)
			return
		}

		// Set the issuables and metadata in the context
		ctx.Set("issues", paginatedIssuables)
		ctx.Set("issuable_meta_data", issuableMetadata)

		// Render the Atom response
		ctx.HTML(http.StatusOK, "xml", nil)
	}
}

// MergeRequests handles the merge requests action
func (i *IssuableCollectionsAction) MergeRequests(ctx *gin.Context) {
	// Render the merge requests
	i.renderMergeRequests(ctx)
}

// IssuesCalendar handles the issues calendar action
func (i *IssuableCollectionsAction) IssuesCalendar(ctx *gin.Context) {
	// Get the issuables collection
	issuables, err := i.issuableCollections.GetIssuablesCollection(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Render the issues calendar
	i.issuesCalendar.RenderIssuesCalendar(ctx, issuables)
}

// GetSortingField gets the sorting field based on the action name
func (i *IssuableCollectionsAction) getSortingField(ctx *gin.Context) string {
	// Get the action name
	actionName := ctx.Request.URL.Path
	parts := strings.Split(actionName, "/")
	if len(parts) > 0 {
		actionName = parts[len(parts)-1]
	}

	// Check the action name
	switch actionName {
	case "issues":
		return model.IssueSortingPreferenceField
	case "merge_requests", "search_merge_requests":
		return model.MergeRequestSortingPreferenceField
	default:
		return ""
	}
}

// GetFinderType gets the finder type based on the action name
func (i *IssuableCollectionsAction) getFinderType(ctx *gin.Context) string {
	// Get the action name
	actionName := ctx.Request.URL.Path
	parts := strings.Split(actionName, "/")
	if len(parts) > 0 {
		actionName = parts[len(parts)-1]
	}

	// Check the action name
	switch actionName {
	case "issues", "issues_calendar":
		return "IssuesFinder"
	case "merge_requests", "search_merge_requests":
		return "MergeRequestsFinder"
	default:
		return ""
	}
}

// GetFinderOptions gets the finder options
func (i *IssuableCollectionsAction) getFinderOptions(ctx *gin.Context) map[string]interface{} {
	// Get the base finder options
	finderOptions := i.issuableCollections.GetFinderOptions(ctx)

	// Get the issue types
	issueTypes := model.IssueTypesForList

	// Add the non-archived and issue types options
	finderOptions["non_archived"] = true
	finderOptions["issue_types"] = issueTypes

	return finderOptions
}

// RenderMergeRequests renders the merge requests
func (i *IssuableCollectionsAction) renderMergeRequests(ctx *gin.Context) {
	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the issuables collection
	issuables, err := i.issuableCollections.GetIssuablesCollection(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the page parameter
	page := ctx.DefaultQuery("page", "1")

	// Paginate the issuables
	paginatedIssuables, err := i.paginateIssuables(issuables, page)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the issuable metadata
	issuableMetadata, err := i.metadataService.GetIssuableMetadata(currentUser, paginatedIssuables)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Set the merge requests and metadata in the context
	ctx.Set("merge_requests", paginatedIssuables)
	ctx.Set("issuable_meta_data", issuableMetadata)

	// Render the merge requests
	ctx.HTML(http.StatusOK, "merge_requests", nil)
}

// FilterNonArchived filters out archived issuables
func (i *IssuableCollectionsAction) filterNonArchived(issuables []model.Issuable) []model.Issuable {
	// Create a new slice for the non-archived issuables
	nonArchivedIssuables := make([]model.Issuable, 0, len(issuables))

	// Iterate over the issuables
	for _, issuable := range issuables {
		// Check if the issuable is not archived
		if !issuable.IsArchived() {
			// Add the issuable to the non-archived issuables
			nonArchivedIssuables = append(nonArchivedIssuables, issuable)
		}
	}

	return nonArchivedIssuables
}

// PaginateIssuables paginates the issuables
func (i *IssuableCollectionsAction) paginateIssuables(issuables []model.Issuable, page string) ([]model.Issuable, error) {
	// Parse the page parameter
	pageNum, err := i.parseInt(page)
	if err != nil {
		return nil, err
	}

	// Get the page size
	pageSize := 20

	// Calculate the start and end indices
	startIndex := (pageNum - 1) * pageSize
	endIndex := startIndex + pageSize

	// Check if the start index is out of bounds
	if startIndex >= len(issuables) {
		return []model.Issuable{}, nil
	}

	// Check if the end index is out of bounds
	if endIndex > len(issuables) {
		endIndex = len(issuables)
	}

	// Return the paginated issuables
	return issuables[startIndex:endIndex], nil
}

// HandleError handles an error
func (i *IssuableCollectionsAction) handleError(ctx *gin.Context, err error) {
	// Log the error
	i.logger.Error(err)

	// Check if the request is for HTML or JSON
	format := ctx.DefaultQuery("format", "html")
	if format == "html" {
		ctx.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// ParseInt parses an integer
func (i *IssuableCollectionsAction) parseInt(s string) (int, error) {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return 0, fmt.Errorf("not implemented")
}
