package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// AuthMiddleware handles authentication for routes
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware ensures the user is authenticated
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header or cookie
		token := m.getToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate token and get user
		user, err := m.authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("current_user", user)
		c.Set("gin-context", c)

		c.Next()
	}
}

// getToken gets the token from header or cookie
func (m *AuthMiddleware) getToken(c *gin.Context) string {
	// Try to get token from Authorization header
	token := c.GetHeader("Authorization")
	if token != "" {
		return token
	}

	// Try to get token from cookie
	cookie, err := c.Cookie("auth_token")
	if err == nil {
		return cookie
	}

	return ""
} 