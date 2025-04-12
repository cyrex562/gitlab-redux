package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SkipsAlreadySignedInMessage is a middleware that skips showing an "already signed in"
// warning message on registrations and logins
type SkipsAlreadySignedInMessage struct {
	authService *service.AuthService
	i18nService *service.I18nService
}

// NewSkipsAlreadySignedInMessage creates a new instance of SkipsAlreadySignedInMessage
func NewSkipsAlreadySignedInMessage(
	authService *service.AuthService,
	i18nService *service.I18nService,
) *SkipsAlreadySignedInMessage {
	return &SkipsAlreadySignedInMessage{
		authService: authService,
		i18nService: i18nService,
	}
}

// RegisterMiddleware registers the middleware for the specified routes
func (s *SkipsAlreadySignedInMessage) RegisterMiddleware(router *gin.Engine) {
	// Register middleware for new and create actions
	// In Go/Gin, we would typically use router groups or middleware for specific routes
	// This is equivalent to the Ruby skip_before_action and before_action
	router.Use(s.requireNoAuthenticationWithoutFlash)
}

// RequireNoAuthenticationWithoutFlash is a middleware that requires no authentication
// but skips the "already signed in" flash message
func (s *SkipsAlreadySignedInMessage) RequireNoAuthenticationWithoutFlash() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, require no authentication (this would typically redirect if user is authenticated)
		s.authService.RequireNoAuthentication(c)

		// Check if the flash message is "already authenticated"
		flashMessage := c.GetString("flash_alert")
		alreadyAuthenticatedMessage := s.i18nService.T("devise.failure.already_authenticated")

		if flashMessage == alreadyAuthenticatedMessage {
			// Clear the flash message
			c.Set("flash_alert", "")
		}

		c.Next()
	}
}

// RequireNoAuthenticationWithoutFlashForRoutes applies the middleware to specific routes
func (s *SkipsAlreadySignedInMessage) RequireNoAuthenticationWithoutFlashForRoutes(
	router *gin.Engine,
	routes []string,
) {
	for _, route := range routes {
		router.GET(route, s.RequireNoAuthenticationWithoutFlash())
		router.POST(route, s.RequireNoAuthenticationWithoutFlash())
	}
}
