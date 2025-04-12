package pagination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// PageLimiter handles pagination limits
type PageLimiter struct {
	metricsService *service.MetricsService
	deviceDetector *service.DeviceDetector
	logger         *service.Logger
}

// NewPageLimiter creates a new instance of PageLimiter
func NewPageLimiter(
	metricsService *service.MetricsService,
	deviceDetector *service.DeviceDetector,
	logger *service.Logger,
) *PageLimiter {
	return &PageLimiter{
		metricsService: metricsService,
		deviceDetector: deviceDetector,
		logger:         logger,
	}
}

// SetupMiddleware sets up the middleware for page limiting
func (p *PageLimiter) SetupMiddleware(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		// This middleware doesn't do anything by default
		// It's up to the handlers to call LimitPages
		c.Next()
	})
}

// LimitPages limits the page number to a maximum value
func (p *PageLimiter) LimitPages(c *gin.Context, maxPageNumber int) error {
	// Check if maxPageNumber is valid
	if err := p.checkPageNumber(c, maxPageNumber); err != nil {
		return err
	}

	return nil
}

// checkPageNumber checks if the page number is valid
func (p *PageLimiter) checkPageNumber(c *gin.Context, maxPageNumber int) error {
	// Check if maxPageNumber is a positive integer
	if maxPageNumber <= 0 {
		return ErrPageLimitNotSensible
	}

	// Get page from query parameters
	pageStr := c.Query("page")
	if pageStr == "" {
		return nil
	}

	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return ErrPageLimitNotANumber
	}

	// Check if page exceeds maxPageNumber
	if page > maxPageNumber {
		// Record page limit interception
		p.recordPageLimitInterception(c)

		return &ErrPageOutOfBounds{
			MaxPageNumber: maxPageNumber,
		}
	}

	return nil
}

// DefaultPageOutOfBoundsResponse returns a default response for page out of bounds errors
func (p *PageLimiter) DefaultPageOutOfBoundsResponse(c *gin.Context, err error) {
	c.Status(http.StatusBadRequest)
}

// recordPageLimitInterception records the page limit being hit in metrics
func (p *PageLimiter) recordPageLimitInterception(c *gin.Context) {
	// Get user agent
	userAgent := c.GetHeader("User-Agent")

	// Check if user agent is from a bot
	isBot := p.deviceDetector.IsBot(userAgent)

	// Get controller and action from context
	controller := c.GetString("controller")
	action := c.GetString("action")

	// Record metric
	p.metricsService.IncrementCounter("gitlab_page_out_of_bounds", map[string]string{
		"controller": controller,
		"action":     action,
		"bot":        strconv.FormatBool(isBot),
	})
}

// Errors
var (
	ErrPageLimitNotANumber  = &service.Error{
		Code:    "page_limit_not_a_number",
		Message: "Page limit must be a number",
	}
	ErrPageLimitNotSensible = &service.Error{
		Code:    "page_limit_not_sensible",
		Message: "Page limit must be a positive number",
	}
)

// ErrPageOutOfBounds is returned when the page number exceeds the maximum
type ErrPageOutOfBounds struct {
	MaxPageNumber int
}

// Error returns the error message
func (e *ErrPageOutOfBounds) Error() string {
	return "Page number exceeds the maximum of " + strconv.Itoa(e.MaxPageNumber)
}

// Code returns the error code
func (e *ErrPageOutOfBounds) Code() string {
	return "page_out_of_bounds"
}
