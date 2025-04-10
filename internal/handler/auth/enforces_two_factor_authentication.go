package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// EnforcesTwoFactorAuthentication handles enforcing two-factor authentication
type EnforcesTwoFactorAuthentication struct {
	twoFactorService *service.TwoFactorService
	userService      *service.UserService
	routeService     *service.RouteService
	logger           *util.Logger
}

// NewEnforcesTwoFactorAuthentication creates a new instance of EnforcesTwoFactorAuthentication
func NewEnforcesTwoFactorAuthentication(
	twoFactorService *service.TwoFactorService,
	userService *service.UserService,
	routeService *service.RouteService,
	logger *util.Logger,
) *EnforcesTwoFactorAuthentication {
	return &EnforcesTwoFactorAuthentication{
		twoFactorService: twoFactorService,
		userService:      userService,
		routeService:     routeService,
		logger:           logger,
	}
}

// CheckTwoFactorRequirementMiddleware creates a middleware that checks two-factor authentication requirements
func (e *EnforcesTwoFactorAuthentication) CheckTwoFactorRequirementMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip for route not found
		if ctx.Request.URL.Path == "/404" {
			ctx.Next()
			return
		}

		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.Next()
			return
		}
		user := currentUser.(*model.User)

		// Check if two-factor authentication is required
		if e.IsTwoFactorAuthenticationRequired() && e.CurrentUserRequiresTwoFactor(ctx, user) {
			// Check if this is a GraphQL request
			if ctx.Request.URL.Path == "/api/graphql" {
				// Get the MFA help page URL
				mfaHelpPageURL := e.GetMfaHelpPageURL()

				// Return an error response
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": fmt.Sprintf(
						"Authentication error: enable 2FA in your profile settings to continue using GitLab: %s",
						mfaHelpPageURL,
					),
				})
				ctx.Abort()
				return
			}

			// Redirect to the two-factor authentication setup page
			ctx.Redirect(http.StatusFound, "/profile/two_factor_auth")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// IsTwoFactorAuthenticationRequired checks if two-factor authentication is required
func (e *EnforcesTwoFactorAuthentication) IsTwoFactorAuthenticationRequired() bool {
	return e.twoFactorService.IsTwoFactorAuthenticationRequired()
}

// CurrentUserRequiresTwoFactor checks if the current user needs to set up two-factor authentication
func (e *EnforcesTwoFactorAuthentication) CurrentUserRequiresTwoFactor(ctx *gin.Context, user *model.User) bool {
	return e.twoFactorService.CurrentUserNeedsToSetupTwoFactor(ctx, user) && !e.SkipTwoFactor(ctx)
}

// ExecuteActionFor2FAReason executes an action based on the two-factor authentication reason
func (e *EnforcesTwoFactorAuthentication) ExecuteActionFor2FAReason(ctx *gin.Context, actions map[string]func([]*model.Group)) error {
	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return util.NewUnauthorizedError("user not authenticated")
	}
	user := currentUser.(*model.User)

	// Get the two-factor authentication reason
	reason := e.twoFactorService.GetTwoFactorAuthenticationReason(ctx, user)

	// Get the groups enforcing two-factor authentication
	groups := e.userService.GetSourceGroupsOfTwoFactorAuthenticationRequirement(ctx, user)

	// Execute the action for the reason
	if action, ok := actions[reason]; ok {
		action(groups)
		return nil
	}

	return util.NewBadRequestError("invalid two-factor authentication reason")
}

// GetTwoFactorGracePeriod gets the two-factor authentication grace period
func (e *EnforcesTwoFactorAuthentication) GetTwoFactorGracePeriod() time.Duration {
	return e.twoFactorService.GetTwoFactorGracePeriod()
}

// IsTwoFactorGracePeriodExpired checks if the two-factor authentication grace period has expired
func (e *EnforcesTwoFactorAuthentication) IsTwoFactorGracePeriodExpired(ctx *gin.Context, user *model.User) bool {
	return e.twoFactorService.IsTwoFactorGracePeriodExpired(ctx, user)
}

// IsTwoFactorSkippable checks if two-factor authentication can be skipped
func (e *EnforcesTwoFactorAuthentication) IsTwoFactorSkippable(ctx *gin.Context, user *model.User) bool {
	return e.IsTwoFactorAuthenticationRequired() &&
		!user.IsTwoFactorEnabled() &&
		!e.IsTwoFactorGracePeriodExpired(ctx, user)
}

// SkipTwoFactor checks if two-factor authentication should be skipped
func (e *EnforcesTwoFactorAuthentication) SkipTwoFactor(ctx *gin.Context) bool {
	// Get the skip two-factor flag from the session
	skipTwoFactor, exists := ctx.Get("skip_two_factor")
	if !exists {
		return false
	}

	// Check if the skip two-factor flag is set and not expired
	skipTime, ok := skipTwoFactor.(time.Time)
	if !ok {
		return false
	}

	return skipTime.After(time.Now())
}

// GetMfaHelpPageURL gets the MFA help page URL
func (e *EnforcesTwoFactorAuthentication) GetMfaHelpPageURL() string {
	return e.routeService.GetHelpPageURL(
		"user/profile/account/two_factor_authentication.md",
		"enable-two-factor-authentication",
	)
}
