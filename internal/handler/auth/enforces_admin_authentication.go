package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// EnforcesAdminAuthentication handles enforcing admin authentication
type EnforcesAdminAuthentication struct {
	authService *service.AuthService
	userService *service.UserService
	logger      *util.Logger
}

// NewEnforcesAdminAuthentication creates a new instance of EnforcesAdminAuthentication
func NewEnforcesAdminAuthentication(
	authService *service.AuthService,
	userService *service.UserService,
	logger *util.Logger,
) *EnforcesAdminAuthentication {
	return &EnforcesAdminAuthentication{
		authService: authService,
		userService: userService,
		logger:      logger,
	}
}

// AuthenticateAdminMiddleware creates a middleware that enforces admin authentication
func (e *EnforcesAdminAuthentication) AuthenticateAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			ctx.Abort()
			return
		}
		user := currentUser.(*model.User)

		// Check if the user is an admin
		if !user.IsAdmin() {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			ctx.Abort()
			return
		}

		// Check if admin mode is enabled
		if !e.authService.IsAdminModeEnabled() {
			ctx.Next()
			return
		}

		// Check if the user is in admin mode
		userMode := e.userService.GetUserMode(ctx, user)
		if !userMode.IsAdminMode() {
			// Request admin mode
			e.userService.RequestAdminMode(ctx, user)

			// Store the current location if it's storable
			if e.isStorableLocation(ctx) {
				e.authService.StoreLocation(ctx, "redirect", ctx.Request.URL.Path)
			}

			// Redirect to the admin session path
			ctx.Redirect(http.StatusFound, "/admin/session/new")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// AuthorizeAbilityMiddleware creates a middleware that authorizes a specific ability
func (e *EnforcesAdminAuthentication) AuthorizeAbilityMiddleware(ability string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			ctx.Abort()
			return
		}
		user := currentUser.(*model.User)

		// Check if the user is an admin
		if user.IsAdmin() {
			// If the user is an admin, enforce admin authentication
			e.AuthenticateAdminMiddleware()(ctx)
			return
		}

		// Check if the user has the required ability
		if !e.userService.Can(ctx, user, ability) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Permission denied",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// isStorableLocation checks if the current location is storable
func (e *EnforcesAdminAuthentication) isStorableLocation(ctx *gin.Context) bool {
	return ctx.Request.URL.Path != "/admin/session/new"
}
