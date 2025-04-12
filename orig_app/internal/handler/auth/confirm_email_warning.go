package auth

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// ConfirmEmailWarning handles displaying a warning to users who need to confirm their email
type ConfirmEmailWarning struct {
	settingsService *service.SettingsService
	logger          *util.Logger
	email           string
}

// NewConfirmEmailWarning creates a new instance of ConfirmEmailWarning
func NewConfirmEmailWarning(
	settingsService *service.SettingsService,
	logger *util.Logger,
) *ConfirmEmailWarning {
	return &ConfirmEmailWarning{
		settingsService: settingsService,
		logger:          logger,
	}
}

// SetConfirmWarning sets a warning message if the user needs to confirm their email
func (c *ConfirmEmailWarning) SetConfirmWarning(ctx *gin.Context) {
	// Check if we should show the warning
	if !c.shouldShowConfirmWarning(ctx) {
		return
	}

	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return
	}
	user := currentUser.(*model.User)

	// Return if the user is already confirmed
	if user.Confirmed() {
		return
	}

	// Get the email to display
	email := c.getEmail(user)

	// Create the warning message
	message := c.formatConfirmWarningMessage(
		email,
		c.createResendLink(email),
		c.createUpdateLink(),
	)

	// Set the warning message in the flash
	ctx.Set("flash_warning", template.HTML(message))
}

// shouldShowConfirmWarning checks if we should show the confirm warning
func (c *ConfirmEmailWarning) shouldShowConfirmWarning(ctx *gin.Context) bool {
	// Check if it's an HTML request
	if ctx.GetHeader("Accept") != "text/html" {
		return false
	}

	// Check if it's a GET request
	if ctx.Request.Method != http.MethodGet {
		return false
	}

	// Check if email confirmation is set to soft
	return c.settingsService.EmailConfirmationSettingSoft()
}

// getEmail gets the email to display
func (c *ConfirmEmailWarning) getEmail(user *model.User) string {
	if c.email != "" {
		return c.email
	}

	if user.UnconfirmedEmail != "" {
		c.email = user.UnconfirmedEmail
	} else {
		c.email = user.Email
	}

	return c.email
}

// formatConfirmWarningMessage formats the confirm warning message
func (c *ConfirmEmailWarning) formatConfirmWarningMessage(email, resendLink, updateLink string) string {
	return fmt.Sprintf(
		"Please check your email (%s) to verify that you own this address and unlock the power of CI/CD. "+
			"Didn't receive it? %s. Wrong email address? %s.",
		template.HTMLEscapeString(email),
		resendLink,
		updateLink,
	)
}

// createResendLink creates a link to resend the confirmation email
func (c *ConfirmEmailWarning) createResendLink(email string) string {
	return fmt.Sprintf(
		`<a href="/users/confirmation?user[email]=%s" data-method="post">Resend it</a>`,
		template.HTMLEscapeString(email),
	)
}

// createUpdateLink creates a link to update the email address
func (c *ConfirmEmailWarning) createUpdateLink() string {
	return `<a href="/-/profile">Update it</a>`
}
