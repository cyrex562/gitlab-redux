package captcha_check

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// JsonFormatActionsSupport handles JSON form submissions with CAPTCHA checks
type JsonFormatActionsSupport struct {
	common *Common
	logger *util.Logger
}

// NewJsonFormatActionsSupport creates a new instance of JsonFormatActionsSupport
func NewJsonFormatActionsSupport(common *Common, logger *util.Logger) *JsonFormatActionsSupport {
	return &JsonFormatActionsSupport{
		common: common,
		logger: logger,
	}
}

// WithCaptchaCheckJsonFormat executes the given action with CAPTCHA check if needed for JSON format
func (j *JsonFormatActionsSupport) WithCaptchaCheckJsonFormat(ctx *gin.Context, spammable model.Spammable, action func(*gin.Context)) {
	captchaRenderFunc := func(c *gin.Context) {
		// Use 409 Conflict status code as it's appropriate for responses requiring CAPTCHA
		// https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.10
		c.JSON(409, GetSpamActionResponseFields(spammable))
	}

	j.common.WithCaptchaCheckCommon(ctx, spammable, captchaRenderFunc)
}

// GetSpamActionResponseFields returns the response fields for spam/CAPTCHA actions
func GetSpamActionResponseFields(spammable model.Spammable) gin.H {
	return gin.H{
		"spam":                spammable.IsSpam(),
		"needs_captcha_response": spammable.RenderRecaptcha(),
		"spam_log_id":         spammable.GetSpamLogID(),
		"captcha_site_key":    service.GetRecaptchaSiteKey(),
	}
}
