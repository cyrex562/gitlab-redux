package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// InvisibleCaptchaOnSignup handles invisible CAPTCHA functionality for signup forms
type InvisibleCaptchaOnSignup struct {
	configService *service.ConfigService
	metricsService *service.MetricsService
	logger *service.Logger
	invisibleCaptchaService *service.InvisibleCaptchaService
}

// NewInvisibleCaptchaOnSignup creates a new instance of InvisibleCaptchaOnSignup
func NewInvisibleCaptchaOnSignup(
	configService *service.ConfigService,
	metricsService *service.MetricsService,
	logger *service.Logger,
	invisibleCaptchaService *service.InvisibleCaptchaService,
) *InvisibleCaptchaOnSignup {
	return &InvisibleCaptchaOnSignup{
		configService: configService,
		metricsService: metricsService,
		logger: logger,
		invisibleCaptchaService: invisibleCaptchaService,
	}
}

// SetupMiddleware sets up the invisible CAPTCHA middleware
func (i *InvisibleCaptchaOnSignup) SetupMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Only apply to POST requests to the signup endpoint
		if ctx.Request.Method != http.MethodPost || ctx.Request.URL.Path != "/users" {
			ctx.Next()
			return
		}

		// Check if invisible CAPTCHA is enabled
		if !i.configService.IsInvisibleCaptchaEnabled() {
			ctx.Next()
			return
		}

		// Check for honeypot spam
		if i.invisibleCaptchaService.IsHoneypotSpam(ctx) {
			i.OnHoneypotSpamCallback(ctx)
			ctx.Abort()
			return
		}

		// Check for timestamp spam
		if i.invisibleCaptchaService.IsTimestampSpam(ctx) {
			i.OnTimestampSpamCallback(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// OnHoneypotSpamCallback handles honeypot spam detection
func (i *InvisibleCaptchaOnSignup) OnHoneypotSpamCallback(ctx *gin.Context) {
	// Check if invisible CAPTCHA is enabled
	if !i.configService.IsInvisibleCaptchaEnabled() {
		return
	}

	// Increment the honeypot counter
	i.invisibleCaptchaService.IncrementHoneypotCounter()

	// Log the request
	i.LogRequest(ctx, "Invisible_Captcha_Honeypot_Request")

	// Return a 200 OK response
	ctx.Status(http.StatusOK)
}

// OnTimestampSpamCallback handles timestamp spam detection
func (i *InvisibleCaptchaOnSignup) OnTimestampSpamCallback(ctx *gin.Context) {
	// Check if invisible CAPTCHA is enabled
	if !i.configService.IsInvisibleCaptchaEnabled() {
		return
	}

	// Increment the timestamp counter
	i.invisibleCaptchaService.IncrementTimestampCounter()

	// Log the request
	i.LogRequest(ctx, "Invisible_Captcha_Timestamp_Request")

	// Redirect to the login page with an error message
	ctx.Redirect(http.StatusFound, "/users/sign_in")
	ctx.Set("flash_message", i.invisibleCaptchaService.GetTimestampErrorMessage())
}

// LogRequest logs the request information
func (i *InvisibleCaptchaOnSignup) LogRequest(ctx *gin.Context, message string) {
	requestInfo := map[string]interface{}{
		"message":        message,
		"env":           "invisible_captcha_signup_bot_detected",
		"remote_ip":     ctx.ClientIP(),
		"request_method": ctx.Request.Method,
		"path":          ctx.Request.URL.Path,
	}

	i.logger.Error(requestInfo)
}
