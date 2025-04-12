package performance

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gitlab-org/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/config"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// WithPerformanceBarHandler handles performance bar functionality
type WithPerformanceBarHandler struct {
	config *config.Config
	authService *service.AuthService
	performanceService *service.PerformanceService
}

// NewWithPerformanceBarHandler creates a new performance bar handler
func NewWithPerformanceBarHandler(
	config *config.Config,
	authService *service.AuthService,
	performanceService *service.PerformanceService,
) *WithPerformanceBarHandler {
	return &WithPerformanceBarHandler{
		config: config,
		authService: authService,
		performanceService: performanceService,
	}
}

// SetPeekEnabledForCurrentRequest sets the peek enabled flag for the current request
func (h *WithPerformanceBarHandler) SetPeekEnabledForCurrentRequest(c *gin.Context) {
	// Get the current user
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		return
	}

	// Check if the performance bar is enabled for the current request
	enabled := h.isPeekEnabled(c, user)

	// Store the result in the request context
	c.Set("peek_enabled", enabled)
}

// IsPeekEnabled checks if the performance bar is enabled for the current request
func (h *WithPerformanceBarHandler) IsPeekEnabled(c *gin.Context) bool {
	// Get the stored value from the context
	enabled, exists := c.Get("peek_enabled")
	if !exists {
		return false
	}

	return enabled.(bool)
}

// isPeekEnabled determines if the performance bar should be enabled
func (h *WithPerformanceBarHandler) isPeekEnabled(c *gin.Context, user *model.User) bool {
	// Check if we're in development mode
	isDevelopment := os.Getenv("GIN_MODE") != "release"

	// Get the cookie value
	cookieEnabled := false
	cookie, err := c.Cookie("perf_bar_enabled")
	if err == nil && cookie == "true" {
		cookieEnabled = true
	}

	// Set the cookie in development mode if it's not set
	if isDevelopment && cookie == "" {
		c.SetCookie("perf_bar_enabled", "true", 0, "/", "", false, true)
	}

	// Check if the user is allowed to use the performance bar
	userAllowed := h.performanceService.IsAllowedForUser(user)

	// Return true only if both conditions are met
	return cookieEnabled && userAllowed
}
