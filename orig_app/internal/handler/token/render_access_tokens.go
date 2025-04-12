package token

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// RenderAccessTokens handles rendering of access tokens
type RenderAccessTokens struct {
	tokenService    *service.TokenService
	paginationService *service.PaginationService
	calendarService *service.CalendarService
	logger          *service.Logger
}

// NewRenderAccessTokens creates a new instance of RenderAccessTokens
func NewRenderAccessTokens(
	tokenService *service.TokenService,
	paginationService *service.PaginationService,
	calendarService *service.CalendarService,
	logger *service.Logger,
) *RenderAccessTokens {
	return &RenderAccessTokens{
		tokenService:     tokenService,
		paginationService: paginationService,
		calendarService:  calendarService,
		logger:           logger,
	}
}

// ActiveAccessTokens gets active access tokens with pagination
func (r *RenderAccessTokens) ActiveAccessTokens(c *gin.Context) (interface{}, int, error) {
	// Create finder options
	options := map[string]interface{}{
		"state": "active",
		"sort":  "expires_asc",
	}

	// Get tokens with users preloaded
	tokens, err := r.tokenService.FindTokensWithUsers(options)
	if err != nil {
		return nil, 0, err
	}

	// Get total size
	size := len(tokens)

	// Get page number
	page := r.getPage(c)

	// Paginate tokens
	paginatedTokens, err := r.paginationService.PaginateTokens(tokens, page)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination headers
	err = r.addPaginationHeaders(c, paginatedTokens)
	if err != nil {
		return nil, 0, err
	}

	// Represent tokens
	representedTokens, err := r.tokenService.RepresentTokens(paginatedTokens)
	if err != nil {
		return nil, 0, err
	}

	return representedTokens, size, nil
}

// InactiveAccessTokens gets inactive access tokens
func (r *RenderAccessTokens) InactiveAccessTokens(c *gin.Context) (interface{}, error) {
	// Create finder options
	options := map[string]interface{}{
		"state": "inactive",
		"sort":  "updated_at_desc",
	}

	// Get tokens with users preloaded
	tokens, err := r.tokenService.FindTokensWithUsers(options)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// AddPaginationHeaders adds pagination headers to the response
func (r *RenderAccessTokens) addPaginationHeaders(c *gin.Context, relation interface{}) error {
	// Get pagination parameters
	perPage, err := r.paginationService.GetPerPage(relation)
	if err != nil {
		return err
	}

	currentPage, err := r.paginationService.GetCurrentPage(relation)
	if err != nil {
		return err
	}

	nextPage, err := r.paginationService.GetNextPage(relation)
	if err != nil {
		return err
	}

	prevPage, err := r.paginationService.GetPrevPage(relation)
	if err != nil {
		return err
	}

	totalCount, err := r.paginationService.GetTotalCount(relation)
	if err != nil {
		return err
	}

	// Get permitted params
	params := map[string]interface{}{
		"page":     c.Query("page"),
		"per_page": c.Query("per_page"),
	}

	// Build and execute pagination headers
	err = r.paginationService.BuildOffsetHeaders(c, perPage, currentPage, nextPage, prevPage, totalCount, params)
	if err != nil {
		return err
	}

	return nil
}

// GetPage gets the page number from request parameters
func (r *RenderAccessTokens) getPage(c *gin.Context) int {
	page := c.DefaultQuery("page", "1")
	pageNum := 1

	if pageVal, err := r.paginationService.ParsePageNumber(page); err == nil {
		pageNum = pageVal
	}

	return pageNum
}

// ExpiryICS generates an ICS calendar for token expiry dates
func (r *RenderAccessTokens) ExpiryICS(tokens []map[string]interface{}) (string, error) {
	// Create new calendar
	calendar, err := r.calendarService.NewCalendar()
	if err != nil {
		return "", err
	}

	// Add events for each token
	for _, token := range tokens {
		expiresAt, ok := token["expires_at"].(string)
		if !ok {
			continue
		}

		name, ok := token["name"].(string)
		if !ok {
			continue
		}

		// Parse expires_at date
		expiresAtDate := expiresAt
		for _, char := range "-" {
			expiresAtDate = r.calendarService.RemoveChar(expiresAtDate, char)
		}

		// Create event
		event := map[string]interface{}{
			"dtstart":  expiresAtDate,
			"dtend":    expiresAtDate,
			"summary":  "Token " + name + " expires today",
		}

		err = r.calendarService.AddEvent(calendar, event)
		if err != nil {
			return "", err
		}
	}

	// Convert calendar to ICS format
	ics, err := r.calendarService.ToICS(calendar)
	if err != nil {
		return "", err
	}

	return ics, nil
}
