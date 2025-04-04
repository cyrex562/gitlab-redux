package captcha_check

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// Common handles common CAPTCHA check functionality
type Common struct {
	recaptchaService *service.RecaptchaService
	logger          *util.Logger
}

// NewCommon creates a new instance of Common
func NewCommon(recaptchaService *service.RecaptchaService, logger *util.Logger) *Common {
	return &Common{
		recaptchaService: recaptchaService,
		logger:          logger,
	}
}

// WithCaptchaCheckCommon executes the given action with CAPTCHA check if needed
func (c *Common) WithCaptchaCheckCommon(ctx *gin.Context, spammable model.Spammable, captchaRenderFunc func(*gin.Context)) {
	// If CAPTCHA is not necessary, execute the action directly
	if !spammable.RenderRecaptcha() {
		captchaRenderFunc(ctx)
		return
	}

	// Load recaptcha configurations
	if err := c.recaptchaService.LoadConfigurations(); err != nil {
		c.logger.Error("Failed to load recaptcha configurations", "error", err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// Render the CAPTCHA
	captchaRenderFunc(ctx)
}

// Spammable interface defines the methods required for CAPTCHA checking
type Spammable interface {
	// RenderRecaptcha returns true if a CAPTCHA should be rendered
	RenderRecaptcha() bool
}
