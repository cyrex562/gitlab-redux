package glql

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/graphql"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// GlqlQueryLockedError is an error that occurs when a query is locked
type GlqlQueryLockedError struct {
	message string
}

// Error returns the error message
func (e *GlqlQueryLockedError) Error() string {
	return e.message
}

// BaseController is the base controller for all GLQL controllers
type BaseController struct {
	*graphql.GraphqlController
	rateLimiterService service.RateLimiterService
	metricsService     service.MetricsService
	contextService     service.ContextService
}

// NewBaseController creates a new BaseController
func NewBaseController(
	graphqlController *graphql.GraphqlController,
	rateLimiterService service.RateLimiterService,
	metricsService service.MetricsService,
	contextService service.ContextService,
) *BaseController {
	return &BaseController{
		GraphqlController:   graphqlController,
		rateLimiterService: rateLimiterService,
		metricsService:     metricsService,
		contextService:     contextService,
	}
}

// RegisterRoutes registers the routes for the BaseController
func (c *BaseController) RegisterRoutes(router *gin.Engine) {
	glql := router.Group("/glql")
	{
		glql.POST("/execute", c.Execute)
	}
}

// Execute handles the execute action
func (c *BaseController) Execute(ctx *gin.Context) {
	// Check rate limit
	if err := c.checkRateLimit(ctx); err != nil {
		c.handleGlqlQueryLockedError(ctx, err)
		return
	}

	// Start time for metrics
	startTime := time.Now()

	// Execute the query
	err := c.GraphqlController.Execute(ctx)

	// Handle errors
	if err != nil {
		// Increment rate limit counter if query was aborted
		if c.isQueryAbortedError(err) {
			c.incrementRateLimitCounter(ctx)
		}

		// Re-throw the error
		panic(err)
	}

	// Record metrics
	c.incrementGlqlSli(ctx, time.Since(startTime).Seconds(), c.errorTypeFrom(err))
}

// handleGlqlQueryLockedError handles the GlqlQueryLockedError
func (c *BaseController) handleGlqlQueryLockedError(ctx *gin.Context, err error) {
	// Log the exception
	c.logException(ctx, err)

	// Render error
	ctx.JSON(http.StatusForbidden, gin.H{
		"errors": []gin.H{
			{
				"message": err.Error(),
			},
		},
	})
}

// logException logs the exception
func (c *BaseController) logException(ctx *gin.Context, err error) {
	// TODO: Implement logging
}

// logs returns the logs with additional GLQL information
func (c *BaseController) logs(ctx *gin.Context) []map[string]interface{} {
	// Get base logs
	logs := c.GraphqlController.Logs(ctx)

	// Add GLQL information
	for i := range logs {
		logs[i]["glql_referer"] = ctx.GetHeader("Referer")
		logs[i]["glql_query_sha"] = c.querySha(ctx)
	}

	return logs
}

// checkRateLimit checks if the rate limit has been exceeded
func (c *BaseController) checkRateLimit(ctx *gin.Context) error {
	// Get query SHA
	querySha := c.querySha(ctx)

	// Check if rate limit has been exceeded
	if c.rateLimiterService.Peek(ctx, "glql", querySha) {
		return &GlqlQueryLockedError{
			message: "Query execution is locked due to repeated failures.",
		}
	}

	return nil
}

// incrementRateLimitCounter increments the rate limit counter
func (c *BaseController) incrementRateLimitCounter(ctx *gin.Context) bool {
	// Get query SHA
	querySha := c.querySha(ctx)

	// Increment rate limit counter
	return c.rateLimiterService.Throttled(ctx, "glql", querySha)
}

// querySha returns the SHA-256 hash of the query
func (c *BaseController) querySha(ctx *gin.Context) string {
	// Get query from request
	query := ctx.PostForm("query")
	if query == "" {
		return ""
	}

	// Calculate SHA-256 hash
	hash := sha256.Sum256([]byte(query))
	return hex.EncodeToString(hash[:])
}

// incrementGlqlSli increments the GLQL SLI metrics
func (c *BaseController) incrementGlqlSli(ctx *gin.Context, durationS float64, errorType string) {
	// Get query urgency
	queryUrgency := model.RequestUrgencyLow

	// Get labels
	labels := map[string]string{
		"endpoint_id":     c.contextService.GetCallerID(ctx),
		"feature_category": c.contextService.GetFeatureCategory(ctx),
		"query_urgency":   queryUrgency.Name,
	}

	// Record error if there is an error type
	if errorType != "" {
		labels["error_type"] = errorType
		c.metricsService.RecordGlqlError(ctx, labels, true)
		return
	}

	// Record apdex
	c.metricsService.RecordGlqlApdex(ctx, labels, durationS <= queryUrgency.Duration)
}

// errorTypeFrom returns the error type from the error
func (c *BaseController) errorTypeFrom(err error) string {
	if err == nil {
		return ""
	}

	// Check if error is a query aborted error
	if c.isQueryAbortedError(err) {
		return "query_aborted"
	}

	return "other"
}

// isQueryAbortedError checks if the error is a query aborted error
func (c *BaseController) isQueryAbortedError(err error) bool {
	// TODO: Implement check for query aborted error
	// This would typically check if the error is of a specific type
	return false
}
