package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// VerifiesWithEmailHandler handles HTTP requests for email verification
type VerifiesWithEmailHandler struct {
	emailVerificationService *service.EmailVerificationService
	rateLimitService        *service.RateLimitService
	authService             *service.AuthService
}

// NewVerifiesWithEmailHandler creates a new handler instance
func NewVerifiesWithEmailHandler(
	emailVerificationService *service.EmailVerificationService,
	rateLimitService *service.RateLimitService,
	authService *service.AuthService,
) *VerifiesWithEmailHandler {
	return &VerifiesWithEmailHandler{
		emailVerificationService: emailVerificationService,
		rateLimitService:        rateLimitService,
		authService:             authService,
	}
}

// RegisterRoutes registers the handler routes
func (h *VerifiesWithEmailHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/verify", h.verifyWithEmail)
			auth.POST("/resend-verification", h.resendVerificationCode)
			auth.POST("/update-email", h.updateEmail)
			auth.GET("/successful-verification", h.successfulVerification)
		}
	}
}

// verifyWithEmail handles POST /api/auth/verify
func (h *VerifiesWithEmailHandler) verifyWithEmail(c *gin.Context) {
	// Skip verification if two-factor is enabled or if it's a QA request
	if h.shouldSkipVerification(c) {
		c.Next()
		return
	}

	// Find the user
	user, err := h.findUser(c)
	if err != nil || user == nil {
		c.Next()
		return
	}

	// Check if user is active
	if !user.IsActive() {
		c.Next()
		return
	}

	// Check if verification token is provided
	verificationToken := c.PostForm("verification_token")
	if verificationToken != "" {
		// Verify the token
		result, err := h.emailVerificationService.ValidateToken(c.Request.Context(), user, verificationToken)
		if err != nil {
			h.handleVerificationFailure(c, user, "invalid_token", "Invalid verification token")
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid verification token"})
			return
		}

		if result.Status == "success" {
			h.handleVerificationSuccess(c, user)
			c.JSON(http.StatusOK, gin.H{"status": "success", "redirect_path": "/users/successful-verification"})
			return
		} else {
			h.handleVerificationFailure(c, user, result.Reason, result.Message)
			c.JSON(http.StatusBadRequest, result)
			return
		}
	}

	// Check if email verification is required
	if h.requireEmailVerificationEnabled(user) {
		// Check rate limit
		if h.rateLimitService.IsRateLimited(c.Request.Context(), "user_sign_in", user.ID) {
			interval := h.rateLimitService.GetRateLimitInterval("user_sign_in")
			message := fmt.Sprintf("Maximum login attempts exceeded. Wait %s and try again.", interval)
			c.Redirect(http.StatusFound, "/users/sign_in?alert="+message)
			return
		}

		// Verify password
		password := c.PostForm("password")
		if h.authService.ValidatePassword(user, password) {
			h.verifyEmail(c, user)
		}
	}
}

// resendVerificationCode handles POST /api/auth/resend-verification
func (h *VerifiesWithEmailHandler) resendVerificationCode(c *gin.Context) {
	// Find the verification user
	user, err := h.findVerificationUser(c)
	if err != nil || user == nil {
		c.Next()
		return
	}

	// Check rate limit
	if h.rateLimitService.IsRateLimited(c.Request.Context(), "email_verification_code_send", user.ID) {
		interval := h.rateLimitService.GetRateLimitInterval("email_verification_code_send")
		message := fmt.Sprintf("You've reached the maximum amount of resends. Wait %s and try again.", interval)
		c.JSON(http.StatusTooManyRequests, gin.H{"status": "failure", "message": message})
		return
	}

	// Get email from request
	email := c.PostForm("email")
	if email != "" {
		// Check if it's a secondary email
		secondaryEmail := h.emailVerificationService.GetSecondaryEmail(user, email)
		if secondaryEmail != "" {
			h.sendVerificationInstructions(c, user, secondaryEmail)
		}
	} else {
		h.sendVerificationInstructions(c, user, "")
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// updateEmail handles POST /api/auth/update-email
func (h *VerifiesWithEmailHandler) updateEmail(c *gin.Context) {
	// Find the verification user
	user, err := h.findVerificationUser(c)
	if err != nil || user == nil {
		c.Next()
		return
	}

	// Log verification
	h.logVerification(c, user, "email_update_requested")

	// Get email from request
	email := c.PostForm("email")

	// Update email
	result, err := h.emailVerificationService.UpdateEmail(c.Request.Context(), user, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update email"})
		return
	}

	if result.Status == "success" {
		h.sendVerificationInstructions(c, user, "")
	} else {
		h.handleVerificationFailure(c, user, result.Reason, result.Message)
	}

	c.JSON(http.StatusOK, result)
}

// successfulVerification handles GET /api/auth/successful-verification
func (h *VerifiesWithEmailHandler) successfulVerification(c *gin.Context) {
	// Clear verification user ID from session
	c.SetCookie("verification_user_id", "", -1, "/", "", false, true)

	// Get redirect URL
	redirectURL := h.authService.GetAfterSignInPath(c.GetInt64("user_id"))

	// Render the page
	c.HTML(http.StatusOK, "auth/successful_verification.html", gin.H{
		"redirect_url": redirectURL,
	})
}

// Private helper methods

func (h *VerifiesWithEmailHandler) shouldSkipVerification(c *gin.Context) bool {
	// Check if two-factor is enabled
	if h.authService.IsTwoFactorEnabled(c.GetInt64("user_id")) {
		return true
	}

	// Check if it's a QA request
	userAgent := c.GetHeader("User-Agent")
	return h.authService.IsQaRequest(userAgent)
}

func (h *VerifiesWithEmailHandler) findUser(c *gin.Context) (*model.User, error) {
	// Try to find user from session
	userID := c.GetInt64("user_id")
	if userID > 0 {
		return h.authService.GetUserByID(c.Request.Context(), userID)
	}

	// Try to find user from verification session
	return h.findVerificationUser(c)
}

func (h *VerifiesWithEmailHandler) findVerificationUser(c *gin.Context) (*model.User, error) {
	// Get verification user ID from session
	verificationUserID, err := c.Cookie("verification_user_id")
	if err != nil {
		return nil, err
	}

	// Convert to int64
	userID, err := strconv.ParseInt(verificationUserID, 10, 64)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	return h.authService.GetUserByID(c.Request.Context(), userID)
}

func (h *VerifiesWithEmailHandler) sendVerificationInstructions(c *gin.Context, user *model.User, secondaryEmail string) {
	// Generate token
	token, err := h.emailVerificationService.GenerateToken(c.Request.Context(), user)
	if err != nil {
		return
	}

	// Lock access
	err = h.authService.LockAccess(c.Request.Context(), user, false, "")
	if err != nil {
		return
	}

	// Send email
	email := secondaryEmail
	if email == "" {
		email = h.emailVerificationService.GetVerificationEmail(user)
	}

	err = h.emailVerificationService.SendVerificationEmail(c.Request.Context(), email, token)
	if err != nil {
		return
	}

	// Log verification
	h.logVerification(c, user, "instructions_sent")
}

func (h *VerifiesWithEmailHandler) verifyEmail(c *gin.Context, user *model.User) {
	// Check if unlock token exists
	if user.UnlockToken != "" {
		// Check if token is expired
		if h.emailVerificationService.IsTokenExpired(user) {
			h.sendVerificationInstructions(c, user, "")
		}

		// Prompt for email verification
		h.promptForEmailVerification(c, user)
	} else if user.IsAccessLocked() || !h.authService.IsTrustedIPAddress(c.Request.Context(), user, c.ClientIP()) {
		// Require email verification if:
		// - their account has been locked because of too many failed login attempts, or
		// - they have logged in before, but never from the current ip address
		reason := ""
		if !user.IsAccessLocked() {
			reason = "sign in from untrusted IP address"
		}

		// Check rate limit
		if !h.rateLimitService.IsRateLimited(c.Request.Context(), "email_verification_code_send", user.ID) {
			h.sendVerificationInstructions(c, user, reason)
		}

		// Prompt for email verification
		h.promptForEmailVerification(c, user)
	}
}

func (h *VerifiesWithEmailHandler) promptForEmailVerification(c *gin.Context, user *model.User) {
	// Set verification user ID in session
	c.SetCookie("verification_user_id", fmt.Sprintf("%d", user.ID), 3600, "/", "", false, true)

	// Set user in context
	c.Set("resource", user)

	// Render the page
	c.HTML(http.StatusOK, "auth/email_verification.html", gin.H{
		"user": user,
	})
}

func (h *VerifiesWithEmailHandler) handleVerificationFailure(c *gin.Context, user *model.User, reason, message string) {
	// Add error to user
	user.AddError("base", message)

	// Log verification
	h.logVerification(c, user, "failed_attempt", reason)
}

func (h *VerifiesWithEmailHandler) handleVerificationSuccess(c *gin.Context, user *model.User) {
	// Confirm user if unconfirmed
	if h.emailVerificationService.IsUnconfirmedVerificationEmail(user) {
		user.Confirm()
	}

	// Set email reset offered at if nil
	if user.EmailResetOfferedAt.IsZero() {
		user.EmailResetOfferedAt = time.Now()
	}

	// Unlock access
	err := h.authService.UnlockAccess(c.Request.Context(), user)
	if err != nil {
		return
	}

	// Log verification
	h.logVerification(c, user, "successful")

	// Sign in user
	h.authService.SignIn(c, user)

	// Log audit event
	h.authService.LogAuditEvent(c, c.GetInt64("user_id"), user, "email_verification")

	// Log user activity
	h.authService.LogUserActivity(c, user)

	// Verify known sign in
	h.authService.VerifyKnownSignIn(c, user)
}

func (h *VerifiesWithEmailHandler) logVerification(c *gin.Context, user *model.User, event string, reason string) {
	// Log verification
	h.authService.LogVerification(c, user, event, reason)
}

func (h *VerifiesWithEmailHandler) requireEmailVerificationEnabled(user *model.User) bool {
	// Check if email verification is required
	return h.emailVerificationService.IsEmailVerificationRequired(user) &&
		!h.emailVerificationService.IsSkipEmailVerificationEnabled(user)
}
