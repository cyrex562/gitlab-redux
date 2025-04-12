package ratelimit

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// SearchRateLimitable handles rate limiting for search requests
type SearchRateLimitable struct {
	rateLimiter *service.ApplicationRateLimiter
	logger      *util.Logger
	settings    *service.ApplicationSettings
}

// NewSearchRateLimitable creates a new instance of SearchRateLimitable
func NewSearchRateLimitable(
	rateLimiter *service.ApplicationRateLimiter,
	logger *util.Logger,
	settings *service.ApplicationSettings,
) *SearchRateLimitable {
	return &SearchRateLimitable{
		rateLimiter: rateLimiter,
		logger:      logger,
		settings:    settings,
	}
}

// CheckSearchRateLimit applies rate limiting to search requests
func (s *SearchRateLimitable) CheckSearchRateLimit(c *gin.Context) error {
	// Get the current user from the context
	currentUser, exists := c.Get("current_user")
	if !exists {
		currentUser = nil
	}

	// Get the client IP address
	clientIP := c.ClientIP()

	if currentUser != nil {
		// For authenticated users, apply rate limit with user and search scope
		// Because every search in the UI typically runs concurrent searches with different
		// scopes to get counts, we apply rate limits on the search scope if it is present.
		//
		// If abusive search is detected, we have stricter limits and ignore the search scope.
		searchScope := s.safeSearchScope(c)

		// Create scope array with user and search scope (if present)
		scope := []interface{}{currentUser}
		if searchScope != "" {
			scope = append(scope, searchScope)
		}

		// Get the allowlist from application settings
		allowlist := s.settings.GetSearchRateLimitAllowlist()

		// Check rate limit with the appropriate scope and allowlist
		return s.rateLimiter.CheckRateLimit(
			c,
			"search_rate_limit",
			scope,
			false,
			map[string]interface{}{
				"users_allowlist": allowlist,
			},
			nil,
		)
	} else {
		// For unauthenticated users, apply rate limit with IP address
		return s.rateLimiter.CheckRateLimit(
			c,
			"search_rate_limit_unauthenticated",
			[]interface{}{clientIP},
			false,
			nil,
			nil,
		)
	}
}

// SafeSearchScope safely extracts the search scope parameter, guarding against abusive values
func (s *SearchRateLimitable) safeSearchScope(c *gin.Context) string {
	// Get the scope parameter from the request
	scope := c.Query("scope")

	// Check if the scope is abusive
	if s.isAbusiveSearch(scope) {
		return ""
	}

	return scope
}

// IsAbusiveSearch checks if the search scope is abusive
func (s *SearchRateLimitable) isAbusiveSearch(scope string) bool {
	// This is a placeholder implementation
	// In a real implementation, this would check if the scope is abusive
	// based on length, content, or other criteria

	// For example, if the scope is too long or contains invalid keywords
	if len(scope) > 1000 {
		return true
	}

	return false
}
