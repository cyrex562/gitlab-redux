package issuable

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// IssuableCollections handles collections of issuable items
type IssuableCollections struct {
	userService        *service.UserService
	issuableService    *service.IssuableService
	metadataService    *service.IssuableMetadataService
	sortingService     *service.SortingService
	rateLimitService   *service.RateLimitService
	logger             *service.Logger
}

// NewIssuableCollections creates a new instance of IssuableCollections
func NewIssuableCollections(
	userService *service.UserService,
	issuableService *service.IssuableService,
	metadataService *service.IssuableMetadataService,
	sortingService *service.SortingService,
	rateLimitService *service.RateLimitService,
	logger *service.Logger,
) *IssuableCollections {
	return &IssuableCollections{
		userService:      userService,
		issuableService:  issuableService,
		metadataService:  metadataService,
		sortingService:   sortingService,
		rateLimitService: rateLimitService,
		logger:           logger,
	}
}

// SetIssuablesIndex sets up the issuables index with pagination
func (i *IssuableCollections) SetIssuablesIndex(ctx *gin.Context) error {
	// Get the issuables collection
	issuables, err := i.GetIssuablesCollection(ctx)
	if err != nil {
		return err
	}

	// Set pagination
	totalPages, err := i.SetPagination(ctx, issuables)
	if err != nil {
		return err
	}

	// Check if we need to redirect out of range
	if err := i.RedirectOutOfRange(ctx, issuables, totalPages); err != nil {
		return err
	}

	return nil
}

// SetPagination sets up pagination for the issuables
func (i *IssuableCollections) SetPagination(ctx *gin.Context, issuables []model.Issuable) (int, error) {
	// Get the format (atom or html)
	format := ctx.DefaultQuery("format", "html")

	// Get the row count based on format
	rowCount := -1
	if format != "atom" {
		rowCount = len(issuables)
	}

	// Get the page parameter
	page := ctx.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return 0, fmt.Errorf("invalid page number: %w", err)
	}

	// Get the sort parameter
	sort := ctx.DefaultQuery("sort", "")

	// Apply pagination
	perPage := 20
	if sort == "relative_position" {
		perPage = 100
	}

	// Calculate pagination
	start := (pageNum - 1) * perPage
	end := start + perPage
	if end > len(issuables) {
		end = len(issuables)
	}

	// Get the current user
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		return 0, err
	}

	// Get the paginated issuables
	paginatedIssuables := issuables[start:end]

	// Get the issuable metadata
	metadata, err := i.metadataService.GetIssuableMetadata(currentUser, paginatedIssuables)
	if err != nil {
		return 0, err
	}

	// Set the data in the context
	ctx.Set("issuables", paginatedIssuables)
	ctx.Set("issuable_meta_data", metadata)
	ctx.Set("total_pages", i.CalculateTotalPages(paginatedIssuables, rowCount))

	return i.CalculateTotalPages(paginatedIssuables, rowCount), nil
}

// GetIssuablesCollection gets the collection of issuables
func (i *IssuableCollections) GetIssuablesCollection(ctx *gin.Context) ([]model.Issuable, error) {
	// Get the finder options
	options, err := i.GetFinderOptions(ctx)
	if err != nil {
		return nil, err
	}

	// Execute the finder
	issuables, err := i.issuableService.ExecuteFinder(ctx, options)
	if err != nil {
		return nil, err
	}

	// Preload the collection
	return i.PreloadForCollection(ctx, issuables)
}

// CalculateTotalPages calculates the total number of pages
func (i *IssuableCollections) CalculateTotalPages(relation []model.Issuable, rowCount int) int {
	limit := 20 // Default page size

	if limit == 0 {
		return 1
	}

	if rowCount == -1 {
		page, _ := strconv.Atoi(i.getPageFromContext())
		return page + 1
	}

	return (rowCount + limit - 1) / limit
}

// GetPerPageForRelativePosition gets the per page value for relative position sorting
func (i *IssuableCollections) GetPerPageForRelativePosition() int {
	return 100
}

// GetIssuableFinderFor gets the finder for the given finder class
func (i *IssuableCollections) GetIssuableFinderFor(finderClass string) (*service.IssuableFinder, error) {
	// Get the finder options
	options, err := i.GetFinderOptions(nil)
	if err != nil {
		return nil, err
	}

	// Create the finder
	return service.NewIssuableFinder(options), nil
}

// GetFinderOptions gets the finder options
func (i *IssuableCollections) GetFinderOptions(ctx *gin.Context) (map[string]interface{}, error) {
	options := make(map[string]interface{})

	// Get the state parameter
	state := i.getStateFromContext(ctx)
	if state == "" {
		state = i.GetDefaultState()
	}
	options["state"] = state

	// Get the scope parameter
	scope := i.getScopeFromContext(ctx)
	if scope != "" {
		options["scope"] = scope
	}

	// Get the confidential parameter
	confidential := i.getConfidentialFromContext(ctx)
	options["confidential"] = confidential

	// Get the sort parameter
	sort := i.getSortFromContext(ctx)
	if sort == "" {
		sort = i.GetDefaultSortOrder(state)
	}
	options["sort"] = sort

	// Check for exact IID search
	search := i.getSearchFromContext(ctx)
	if search != "" {
		re := regexp.MustCompile(`^#(?P<iid>\d+)$`)
		matches := re.FindStringSubmatch(search)
		if len(matches) > 1 {
			options["iids"] = matches[1]
			options["search"] = nil
		} else {
			options["search"] = search
		}
	}

	// Add project or group options
	if projectID := i.getProjectIDFromContext(ctx); projectID != "" {
		options["project_id"] = projectID
		options["attempt_project_search_optimizations"] = true
	} else if groupID := i.getGroupIDFromContext(ctx); groupID != "" {
		options["group_id"] = groupID
		options["include_subgroups"] = true
		options["attempt_group_search_optimizations"] = true
	}

	return options, nil
}

// GetDefaultState gets the default state
func (i *IssuableCollections) GetDefaultState() string {
	return "opened"
}

// GetLegacySortCookieName gets the legacy sort cookie name
func (i *IssuableCollections) GetLegacySortCookieName() string {
	return "issuable_sort"
}

// GetDefaultSortOrder gets the default sort order based on state
func (i *IssuableCollections) GetDefaultSortOrder(state string) string {
	switch state {
	case "opened", "all":
		return "created_desc"
	case "merged", "closed":
		return "updated_desc"
	default:
		return "created_desc"
	}
}

// GetFinder gets the finder
func (i *IssuableCollections) GetFinder(ctx *gin.Context) (*service.IssuableFinder, error) {
	// Get the finder type
	finderType := i.GetFinderType(ctx)

	// Get the finder
	return i.GetIssuableFinderFor(finderType)
}

// GetCollectionType gets the collection type
func (i *IssuableCollections) GetCollectionType(ctx *gin.Context) string {
	finderType := i.GetFinderType(ctx)

	switch finderType {
	case "IssuesFinder":
		return "Issue"
	case "MergeRequestsFinder":
		return "MergeRequest"
	default:
		return ""
	}
}

// PreloadForCollection preloads the collection
func (i *IssuableCollections) PreloadForCollection(ctx *gin.Context, issuables []model.Issuable) ([]model.Issuable, error) {
	collectionType := i.GetCollectionType(ctx)

	// Common attributes to preload
	commonAttributes := []string{"author", "assignees", "labels", "milestone"}

	switch collectionType {
	case "Issue":
		// Preload issue-specific attributes
		attributes := append(commonAttributes, "work_item_type", "project", "project.namespace")
		return i.issuableService.PreloadAttributes(ctx, issuables, attributes)

	case "MergeRequest":
		// Preload merge request-specific attributes
		attributes := append(commonAttributes,
			"target_project",
			"latest_merge_request_diff",
			"approvals",
			"approved_by_users",
			"reviewers",
			"source_project.route",
			"head_pipeline.project",
			"target_project.namespace")
		return i.issuableService.PreloadAttributes(ctx, issuables, attributes)

	default:
		return issuables, nil
	}
}

// RedirectOutOfRange redirects if the page is out of range
func (i *IssuableCollections) RedirectOutOfRange(ctx *gin.Context, issuables []model.Issuable, totalPages int) error {
	page, err := strconv.Atoi(i.getPageFromContext())
	if err != nil {
		return err
	}

	if page > totalPages && totalPages > 0 {
		// Redirect to the last page
		ctx.Redirect(http.StatusFound, fmt.Sprintf("%s?page=%d", ctx.Request.URL.Path, totalPages))
		return fmt.Errorf("redirected to last page")
	}

	return nil
}

// Helper methods to get values from context
func (i *IssuableCollections) getPageFromContext() string {
	return "1" // Default to page 1
}

func (i *IssuableCollections) getStateFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.DefaultQuery("state", "")
}

func (i *IssuableCollections) getScopeFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.DefaultQuery("scope", "")
}

func (i *IssuableCollections) getConfidentialFromContext(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}
	confidential := ctx.DefaultQuery("confidential", "")
	return confidential == "true"
}

func (i *IssuableCollections) getSortFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.DefaultQuery("sort", "")
}

func (i *IssuableCollections) getSearchFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.DefaultQuery("search", "")
}

func (i *IssuableCollections) getProjectIDFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.DefaultQuery("project_id", "")
}

func (i *IssuableCollections) getGroupIDFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.DefaultQuery("group_id", "")
}

// GetFinderType gets the finder type based on the context
func (i *IssuableCollections) GetFinderType(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}

	// This would be determined by the route or controller
	// For now, we'll return a default
	return "IssuesFinder"
}
