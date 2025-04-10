package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// TwoFactorAuth handles two-factor authentication
type TwoFactorAuth struct {
	webauthnService *service.WebauthnAuthenticateService
	webauthnConfig  *service.WebauthnConfiguration
	authService     *service.AuthService
	notificationService *service.NotificationService
	logger          *util.Logger
}

// NewTwoFactorAuth creates a new instance of TwoFactorAuth
func NewTwoFactorAuth(
	webauthnService *service.WebauthnAuthenticateService,
	webauthnConfig *service.WebauthnConfiguration,
	authService *service.AuthService,
	notificationService *service.NotificationService,
	logger *util.Logger,
) *TwoFactorAuth {
	return &TwoFactorAuth{
		webauthnService: webauthnService,
		webauthnConfig:  webauthnConfig,
		authService:     authService,
		notificationService: notificationService,
		logger:          logger,
	}
}

// PromptForTwoFactor prompts the user for two-factor authentication
func (t *TwoFactorAuth) PromptForTwoFactor(ctx *gin.Context, user model.User) {
	// Set the user in the context for Devise views
	ctx.Set("user", user)

	// Check if the user can log in
	if !user.Can("log_in") {
		t.handleLockedUser(ctx, user)
		return
	}

	// Store the user's ID in the session
	ctx.Set("otp_user_id", user.GetID())

	// Store the user's password hash in the session
	passwordHash := sha256.Sum256([]byte(user.GetEncryptedPassword()))
	ctx.Set("user_password_hash", hex.EncodeToString(passwordHash[:]))

	// Add GON variables
	t.addGonVariables(ctx)

	// Set up WebAuthn authentication
	t.setupWebauthnAuthentication(ctx, user)

	// Render the two-factor authentication page
	ctx.HTML(200, "devise/sessions/two_factor.html", gin.H{
		"user": user,
	})
}

// HandleLockedUser handles a locked user
func (t *TwoFactorAuth) HandleLockedUser(ctx *gin.Context, user model.User) {
	// Clear the two-factor attempt
	t.clearTwoFactorAttempt(ctx)

	// Redirect to the login page
	t.lockedUserRedirect(ctx, user)
}

// LockedUserRedirect redirects a locked user to the login page
func (t *TwoFactorAuth) LockedUserRedirect(ctx *gin.Context, user model.User) {
	ctx.Redirect(302, "/users/sign_in")
	ctx.Set("flash_alert", t.lockedUserRedirectAlert(user))
}

// AuthenticateWithTwoFactor authenticates the user with two-factor authentication
func (t *TwoFactorAuth) AuthenticateWithTwoFactor(ctx *gin.Context) {
	// Find the user
	user := t.findUser(ctx)

	// Check if the user can log in
	if !user.Can("log_in") {
		t.handleLockedUser(ctx, user)
		return
	}

	// Check if the user's password has changed
	if t.userPasswordChanged(ctx, user) {
		t.handleChangedUser(ctx, user)
		return
	}

	// Get the user parameters from the request
	userParams := t.getUserParams(ctx)

	// Check if the OTP attempt is present and the OTP user ID is set in the session
	if userParams["otp_attempt"] != nil && ctx.Get("otp_user_id") != nil {
		t.authenticateWithTwoFactorViaOTP(ctx, user)
	} else if userParams["device_response"] != nil && ctx.Get("otp_user_id") != nil {
		t.authenticateWithTwoFactorViaWebauthn(ctx, user)
	} else if user != nil && user.ValidPassword(userParams["password"].(string)) {
		t.PromptForTwoFactor(ctx, user)
	}
}

// LockedUserRedirectAlert returns the alert message for a locked user
func (t *TwoFactorAuth) lockedUserRedirectAlert(user model.User) string {
	if user.IsAccessLocked() {
		return "Your account is locked."
	} else if !user.IsConfirmed() {
		return "You have to confirm your email address before continuing."
	} else {
		return "Invalid login or password"
	}
}

// ClearTwoFactorAttempt clears the two-factor attempt data from the session
func (t *TwoFactorAuth) clearTwoFactorAttempt(ctx *gin.Context) {
	ctx.Set("otp_user_id", nil)
	ctx.Set("user_password_hash", nil)
	ctx.Set("challenge", nil)
}

// AuthenticateWithTwoFactorViaOTP authenticates the user with two-factor authentication via OTP
func (t *TwoFactorAuth) authenticateWithTwoFactorViaOTP(ctx *gin.Context, user model.User) {
	// Get the user parameters from the request
	userParams := t.getUserParams(ctx)

	// Check if the OTP attempt is valid
	if t.validOtpAttempt(ctx, user, userParams["otp_attempt"].(string)) {
		// Clear the two-factor attempt
		t.clearTwoFactorAttempt(ctx)

		// Remember the user if requested
		if userParams["remember_me"] == "1" {
			t.rememberMe(ctx, user)
		}

		// Save the user
		user.Save()

		// Sign in the user
		t.signIn(ctx, user, "two_factor_authenticated", "authentication")
	} else {
		// Send a notification about the failed attempt
		t.sendTwoFactorOtpAttemptFailedEmail(ctx, user)

		// Handle the two-factor failure
		t.handleTwoFactorFailure(ctx, user, "OTP", "Invalid two-factor code.")
	}
}

// AuthenticateWithTwoFactorViaWebauthn authenticates the user with two-factor authentication via WebAuthn
func (t *TwoFactorAuth) authenticateWithTwoFactorViaWebauthn(ctx *gin.Context, user model.User) {
	// Get the user parameters from the request
	userParams := t.getUserParams(ctx)

	// Get the challenge from the session
	challenge, _ := ctx.Get("challenge")

	// Authenticate with WebAuthn
	if t.webauthnService.Execute(user, userParams["device_response"].(string), challenge.(string)) {
		t.handleTwoFactorSuccess(ctx, user)
	} else {
		t.handleTwoFactorFailure(ctx, user, "WebAuthn", "Authentication via WebAuthn device failed.")
	}
}

// SetupWebauthnAuthentication sets up WebAuthn authentication
func (t *TwoFactorAuth) setupWebauthnAuthentication(ctx *gin.Context, user model.User) {
	// Check if the user has WebAuthn registrations
	if len(user.GetWebauthnRegistrations()) > 0 {
		// Get the WebAuthn registration IDs
		webauthnRegistrationIDs := make([]string, 0)
		for _, registration := range user.GetWebauthnRegistrations() {
			webauthnRegistrationIDs = append(webauthnRegistrationIDs, registration.GetCredentialXid())
		}

		// Get the WebAuthn options
		getOptions := t.webauthnConfig.GetCredentialOptionsForGet(
			webauthnRegistrationIDs,
			"discouraged",
			map[string]interface{}{
				"appid": t.webauthnConfig.GetOrigin(),
			},
		)

		// Store the challenge in the session
		ctx.Set("challenge", getOptions.Challenge)

		// Add the WebAuthn options to the GON variables
		ctx.Set("gon_webauthn_options", getOptions)
	}
}

// HandleTwoFactorSuccess handles a successful two-factor authentication
func (t *TwoFactorAuth) handleTwoFactorSuccess(ctx *gin.Context, user model.User) {
	// Clear the two-factor attempt
	t.clearTwoFactorAttempt(ctx)

	// Get the user parameters from the request
	userParams := t.getUserParams(ctx)

	// Remember the user if requested
	if userParams["remember_me"] == "1" {
		t.rememberMe(ctx, user)
	}

	// Sign in the user
	t.signIn(ctx, user, "two_factor_authenticated", "authentication")
}

// HandleTwoFactorFailure handles a failed two-factor authentication
func (t *TwoFactorAuth) handleTwoFactorFailure(ctx *gin.Context, user model.User, method string, message string) {
	// Increment the failed attempts counter
	user.IncrementFailedAttempts()

	// Log the failed two-factor authentication
	t.logFailedTwoFactor(ctx, user, method)

	// Log the failed login
	t.logger.Info(fmt.Sprintf("Failed Login: user=%s ip=%s method=%s", user.GetUsername(), ctx.ClientIP(), method))

	// Set the flash alert
	ctx.Set("flash_alert", message)

	// Prompt for two-factor authentication
	t.PromptForTwoFactor(ctx, user)
}

// SendTwoFactorOtpAttemptFailedEmail sends a notification about a failed OTP attempt
func (t *TwoFactorAuth) sendTwoFactorOtpAttemptFailedEmail(ctx *gin.Context, user model.User) {
	t.notificationService.TwoFactorOtpAttemptFailed(user, ctx.ClientIP())
}

// LogFailedTwoFactor logs a failed two-factor authentication
func (t *TwoFactorAuth) logFailedTwoFactor(ctx *gin.Context, user model.User, method string) {
	// This is a no-op by default
	// It can be overridden in EE
}

// HandleChangedUser handles a user whose password has changed
func (t *TwoFactorAuth) handleChangedUser(ctx *gin.Context, user model.User) {
	// Clear the two-factor attempt
	t.clearTwoFactorAttempt(ctx)

	// Redirect to the login page
	ctx.Redirect(302, "/users/sign_in")
	ctx.Set("flash_alert", "An error occurred. Please sign in again.")
}

// UserPasswordChanged checks if the user's password has changed
func (t *TwoFactorAuth) userPasswordChanged(ctx *gin.Context, user model.User) bool {
	// Get the user's password hash from the session
	userPasswordHash, exists := ctx.Get("user_password_hash")
	if !exists {
		return false
	}

	// Calculate the current password hash
	currentPasswordHash := sha256.Sum256([]byte(user.GetEncryptedPassword()))
	currentPasswordHashHex := hex.EncodeToString(currentPasswordHash[:])

	// Compare the password hashes
	return userPasswordHash.(string) != currentPasswordHashHex
}

// FindUser finds the user from the request
func (t *TwoFactorAuth) findUser(ctx *gin.Context) model.User {
	// This would typically be a method on the auth service
	// For now, we'll just return a dummy user
	return model.NewUser()
}

// AddGonVariables adds GON variables to the context
func (t *TwoFactorAuth) addGonVariables(ctx *gin.Context) {
	// This would typically be a method on the auth service
	// For now, we'll just set a dummy variable
	ctx.Set("gon_variables", map[string]interface{}{
		"webauthn": map[string]interface{}{
			"options": "{}",
		},
	})
}

// ValidOtpAttempt checks if the OTP attempt is valid
func (t *TwoFactorAuth) validOtpAttempt(ctx *gin.Context, user model.User, otpAttempt string) bool {
	// This would typically be a method on the auth service
	// For now, we'll just return true
	return true
}

// RememberMe remembers the user
func (t *TwoFactorAuth) rememberMe(ctx *gin.Context, user model.User) {
	// This would typically be a method on the auth service
	// For now, we'll just set a cookie
	ctx.SetCookie("remember_user_token", "token", 30*24*60*60, "/", "", false, true)
}

// SignIn signs in the user
func (t *TwoFactorAuth) signIn(ctx *gin.Context, user model.User, message string, event string) {
	// This would typically be a method on the auth service
	// For now, we'll just set the current user in the context
	ctx.Set("current_user", user)
	ctx.Set("flash_notice", "Signed in successfully.")
}

// GetUserParams gets the user parameters from the request
func (t *TwoFactorAuth) getUserParams(ctx *gin.Context) map[string]interface{} {
	var params struct {
		Password       string `form:"password"`
		OtpAttempt     string `form:"otp_attempt"`
		DeviceResponse string `form:"device_response"`
		RememberMe     string `form:"remember_me"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		return make(map[string]interface{})
	}

	return map[string]interface{}{
		"password":        params.Password,
		"otp_attempt":     params.OtpAttempt,
		"device_response": params.DeviceResponse,
		"remember_me":     params.RememberMe,
	}
}
