package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// TwoFactorAdmin handles two-factor authentication for admin mode
type TwoFactorAdmin struct {
	webauthnService *service.WebauthnAuthenticateService
	twoFactorAuth   *TwoFactorAuth
	logger          *util.Logger
}

// NewTwoFactorAdmin creates a new instance of TwoFactorAdmin
func NewTwoFactorAdmin(
	webauthnService *service.WebauthnAuthenticateService,
	twoFactorAuth *TwoFactorAuth,
	logger *util.Logger,
) *TwoFactorAdmin {
	return &TwoFactorAdmin{
		webauthnService: webauthnService,
		twoFactorAuth:   twoFactorAuth,
		logger:          logger,
	}
}

// AdminModePromptForTwoFactor prompts the user for two-factor authentication in admin mode
func (t *TwoFactorAdmin) AdminModePromptForTwoFactor(ctx *gin.Context, user model.User) {
	// Set the user in the context for admin views
	ctx.Set("user", user)

	// Check if the user can log in
	if !user.Can("log_in") {
		t.handleLockedUser(ctx, user)
		return
	}

	// Set the OTP user ID in the session
	ctx.Set("otp_user_id", user.GetID())

	// Set up WebAuthn authentication
	t.setupWebauthnAuthentication(ctx, user)

	// Render the two-factor authentication page
	ctx.HTML(200, "admin/sessions/two_factor.html", gin.H{
		"user": user,
	})
}

// AdminModeAuthenticateWithTwoFactor authenticates the user with two-factor authentication in admin mode
func (t *TwoFactorAdmin) AdminModeAuthenticateWithTwoFactor(ctx *gin.Context) {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	// Check if the user can log in
	if !user.Can("log_in") {
		t.handleLockedUser(ctx, user)
		return
	}

	// Get the user parameters from the request
	userParams := t.getUserParams(ctx)

	// Check if the OTP attempt is present and the OTP user ID is set in the session
	if userParams["otp_attempt"] != nil && ctx.Get("otp_user_id") != nil {
		t.adminModeAuthenticateWithTwoFactorViaOTP(ctx, user)
	} else if userParams["device_response"] != nil && ctx.Get("otp_user_id") != nil {
		t.adminModeAuthenticateWithTwoFactorViaWebauthn(ctx, user)
	} else if user != nil && user.ValidPassword(userParams["password"].(string)) {
		t.AdminModePromptForTwoFactor(ctx, user)
	} else {
		t.invalidLoginRedirect(ctx)
	}
}

// adminModeAuthenticateWithTwoFactorViaOTP authenticates the user with two-factor authentication via OTP in admin mode
func (t *TwoFactorAdmin) adminModeAuthenticateWithTwoFactorViaOTP(ctx *gin.Context, user model.User) {
	// Get the user parameters from the request
	userParams := t.getUserParams(ctx)

	// Check if the OTP attempt is valid
	if t.validOtpAttempt(ctx, user, userParams["otp_attempt"].(string)) {
		// Remove any lingering user data from login
		ctx.Set("otp_user_id", nil)

		// Save the user unless the database is read-only
		if !t.isDatabaseReadOnly() {
			user.Save()
		}

		// Enable admin mode
		t.enableAdminMode(ctx)
	} else {
		t.adminHandleTwoFactorFailure(ctx, user, "OTP", "Invalid two-factor code.")
	}
}

// adminModeAuthenticateWithTwoFactorViaWebauthn authenticates the user with two-factor authentication via WebAuthn in admin mode
func (t *TwoFactorAdmin) adminModeAuthenticateWithTwoFactorViaWebauthn(ctx *gin.Context, user model.User) {
	// Get the user parameters from the request
	userParams := t.getUserParams(ctx)

	// Get the challenge from the session
	challenge, _ := ctx.Get("challenge")

	// Authenticate with WebAuthn
	if t.webauthnService.Execute(user, userParams["device_response"].(string), challenge.(string)) {
		t.adminHandleTwoFactorSuccess(ctx)
	} else {
		t.adminHandleTwoFactorFailure(ctx, user, "WebAuthn", "Authentication via WebAuthn device failed.")
	}
}

// enableAdminMode enables admin mode for the current user
func (t *TwoFactorAdmin) enableAdminMode(ctx *gin.Context) {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	// Get the current user mode
	currentUserMode := user.GetCurrentUserMode()

	// Enable admin mode
	if currentUserMode.EnableAdminMode(true) {
		// Redirect to the redirect path
		ctx.Redirect(302, t.getRedirectPath(ctx))
		ctx.Set("flash_notice", "Admin mode enabled")
	} else {
		t.invalidLoginRedirect(ctx)
	}
}

// invalidLoginRedirect redirects to the login page with an invalid login message
func (t *TwoFactorAdmin) invalidLoginRedirect(ctx *gin.Context) {
	ctx.Set("flash_alert", "Invalid login or password")
	ctx.HTML(200, "admin/sessions/new.html", gin.H{})
}

// adminHandleTwoFactorSuccess handles a successful two-factor authentication in admin mode
func (t *TwoFactorAdmin) adminHandleTwoFactorSuccess(ctx *gin.Context) {
	// Remove any lingering user data from login
	ctx.Set("otp_user_id", nil)
	ctx.Set("challenge", nil)

	// Enable admin mode
	t.enableAdminMode(ctx)
}

// adminHandleTwoFactorFailure handles a failed two-factor authentication in admin mode
func (t *TwoFactorAdmin) adminHandleTwoFactorFailure(ctx *gin.Context, user model.User, method string, message string) {
	// Increment the failed attempts counter
	user.IncrementFailedAttempts()

	// Log the failed two-factor authentication
	t.logFailedTwoFactor(ctx, user, method)

	// Log the failed admin mode login
	t.logger.Info(fmt.Sprintf("Failed Admin Mode Login: user=%s ip=%s method=%s", user.GetUsername(), ctx.ClientIP(), method))

	// Set the flash alert
	ctx.Set("flash_alert", message)

	// Prompt for two-factor authentication
	t.AdminModePromptForTwoFactor(ctx, user)
}

// handleLockedUser handles a locked user
func (t *TwoFactorAdmin) handleLockedUser(ctx *gin.Context, user model.User) {
	// This would typically be a method on the TwoFactorAuth struct
	// For now, we'll just redirect to the login page
	t.invalidLoginRedirect(ctx)
}

// setupWebauthnAuthentication sets up WebAuthn authentication
func (t *TwoFactorAdmin) setupWebauthnAuthentication(ctx *gin.Context, user model.User) {
	// This would typically be a method on the TwoFactorAuth struct
	// For now, we'll just set a challenge in the session
	ctx.Set("challenge", "challenge")
}

// validOtpAttempt checks if the OTP attempt is valid
func (t *TwoFactorAdmin) validOtpAttempt(ctx *gin.Context, user model.User, otpAttempt string) bool {
	// This would typically be a method on the TwoFactorAuth struct
	// For now, we'll just return true
	return true
}

// isDatabaseReadOnly checks if the database is read-only
func (t *TwoFactorAdmin) isDatabaseReadOnly() bool {
	// This would typically be a method on the database service
	// For now, we'll just return false
	return false
}

// getUserParams gets the user parameters from the request
func (t *TwoFactorAdmin) getUserParams(ctx *gin.Context) map[string]interface{} {
	var params struct {
		Password      string `form:"password"`
		OtpAttempt    string `form:"otp_attempt"`
		DeviceResponse string `form:"device_response"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		return make(map[string]interface{})
	}

	return map[string]interface{}{
		"password":       params.Password,
		"otp_attempt":    params.OtpAttempt,
		"device_response": params.DeviceResponse,
	}
}

// logFailedTwoFactor logs a failed two-factor authentication
func (t *TwoFactorAdmin) logFailedTwoFactor(ctx *gin.Context, user model.User, method string) {
	// This would typically be a method on the TwoFactorAuth struct
	// For now, we'll just log the failed attempt
	t.logger.Info(fmt.Sprintf("Failed Two-Factor Authentication: user=%s method=%s", user.GetUsername(), method))
}

// getRedirectPath gets the redirect path
func (t *TwoFactorAdmin) getRedirectPath(ctx *gin.Context) string {
	// This would typically be a method on the auth service
	// For now, we'll just return a default path
	return "/admin"
}
