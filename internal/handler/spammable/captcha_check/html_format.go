package captcha_check

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// HtmlFormatActionsSupport handles HTML form submissions with CAPTCHA checks
type HtmlFormatActionsSupport struct {
	common *Common
	logger *util.Logger
}

// NewHtmlFormatActionsSupport creates a new instance of HtmlFormatActionsSupport
func NewHtmlFormatActionsSupport(common *Common, logger *util.Logger) *HtmlFormatActionsSupport {
	return &HtmlFormatActionsSupport{
		common: common,
		logger: logger,
	}
}

// WithCaptchaCheckHtmlFormat executes the given action with CAPTCHA check if needed for HTML format
func (h *HtmlFormatActionsSupport) WithCaptchaCheckHtmlFormat(ctx *gin.Context, spammable model.Spammable, action func(*gin.Context)) {
	captchaRenderFunc := func(c *gin.Context) {
		c.HTML(200, "captcha_check.html", gin.H{
			"spammable": spammable,
		})
	}

	h.common.WithCaptchaCheckCommon(ctx, spammable, captchaRenderFunc)
}

// ConvertHtmlSpamParamsToHeaders converts spam/CAPTCHA values from form field params to headers
func (h *HtmlFormatActionsSupport) ConvertHtmlSpamParamsToHeaders(ctx *gin.Context) {
	// Check if we have either g-recaptcha-response or spam_log_id in the form data
	gRecaptchaResponse := ctx.PostForm("g-recaptcha-response")
	spamLogID := ctx.PostForm("spam_log_id")

	if gRecaptchaResponse == "" && spamLogID == "" {
		return
	}

	// Set the headers if the values are present
	if gRecaptchaResponse != "" {
		ctx.Request.Header.Set("X-GitLab-Captcha-Response", gRecaptchaResponse)
	}
	if spamLogID != "" {
		ctx.Request.Header.Set("X-GitLab-Spam-Log-Id", spamLogID)
	}

	// Reset the spam params on the request context since they have changed mid-request
	// Note: This assumes you have a RequestContext implementation similar to the Ruby version
	// You'll need to implement this part based on your Go application's context management
	// requestContext := GetRequestContext(ctx)
	// requestContext.SetSpamParams(NewSpamParamsFromRequest(ctx.Request))
}
