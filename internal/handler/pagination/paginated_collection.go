package pagination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// PaginatedCollection handles paginated collections
type PaginatedCollection struct {
	urlService *service.URLService
	logger     *service.Logger
}

// NewPaginatedCollection creates a new instance of PaginatedCollection
func NewPaginatedCollection(
	urlService *service.URLService,
	logger *service.Logger,
) *PaginatedCollection {
	return &PaginatedCollection{
		urlService: urlService,
		logger:     logger,
	}
}

// RedirectOutOfRange redirects to the last page if the current page is out of range
func (p *PaginatedCollection) RedirectOutOfRange(c *gin.Context, collection *service.PaginatedCollection, totalPages int) bool {
	// If total pages is 0, return false
	if totalPages == 0 {
		return false
	}

	// Get current page from query parameters
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}

	// Parse current page
	currentPage, err := strconv.Atoi(pageStr)
	if err != nil {
		currentPage = 1
	}

	// Check if current page is out of range
	outOfRange := currentPage > totalPages

	// If out of range, redirect to the last page
	if outOfRange {
		// Get safe parameters
		safeParams := p.getSafeParams(c)

		// Add page parameter
		safeParams["page"] = strconv.Itoa(totalPages)

		// Generate URL
		url, err := p.urlService.GenerateURL(c, safeParams)
		if err != nil {
			p.logger.Error("Failed to generate URL", err)
			return outOfRange
		}

		// Redirect to the last page
		c.Redirect(http.StatusFound, url)
	}

	return outOfRange
}

// getSafeParams gets safe parameters from the request
func (p *PaginatedCollection) getSafeParams(c *gin.Context) map[string]string {
	// Get all query parameters
	params := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// Remove page parameter
	delete(params, "page")

	return params
}
