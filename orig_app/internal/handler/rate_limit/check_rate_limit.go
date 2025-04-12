package ratelimit

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// CheckRateLimit handles rate limiting for controller actions
type CheckRateLimit struct {
	rateLimiter *service.ApplicationRateLimiter
	logger      *util.Logger
}

// NewCheckRateLimit creates a new instance of CheckRateLimit
func NewCheckRateLimit(
	rateLimiter *service.ApplicationRateLimiter,
	logger *util.Logger,
) *CheckRateLimit {
	return &CheckRateLimit{
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// CheckRateLimit checks if the rate limit for a given action is throttled
func (c *CheckRateLimit) CheckRateLimit(
	ctx *gin.Context,
	key string,
	scope string,
	redirectBack bool,
	options map[string]interface{},
	block func() error,
) error {
	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		currentUser = nil
	}

	// Check if the request is throttled
	isThrottled, err := c.rateLimiter.ThrottledRequest(ctx.Request, currentUser, key, scope, options)
	if err != nil {
		c.logger.Error("Failed to check rate limit", err)
		return err
	}

	// If not throttled, execute the block if provided
	if !isThrottled {
		if block != nil {
			return block()
		}
		return nil
	}

	// If throttled, handle the response
	message := "This endpoint has been requested too many times. Try again later."

	if redirectBack {
		// Redirect back with an alert message
		ctx.Redirect(302, ctx.Request.Referer())
		ctx.Set("flash_alert", message)
		return nil
	}

	// Return a too many requests response
	ctx.String(429, message)
	return nil
}
