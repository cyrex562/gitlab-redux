package captcha_check

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// RestApiActionsSupport handles CAPTCHA checks for REST API actions
type RestApiActionsSupport struct {
	common *Common
	logger *util.Logger
}

// NewRestApiActionsSupport creates a new instance of RestApiActionsSupport
func NewRestApiActionsSupport(common *Common, logger *util.Logger) *RestApiActionsSupport {
	return &RestApiActionsSupport{
		common: common,
		logger: logger,
	}
}

// WithCaptchaCheckRestApi executes the given action with CAPTCHA check if needed for REST API
func (r *RestApiActionsSupport) WithCaptchaCheckRestApi(ctx *gin.Context, spammable model.Spammable, action func(*gin.Context)) {
	captchaRenderFunc := func(c *gin.Context) {
		// Get the spam action response fields
		fields := GetSpamActionResponseFields(spammable)

		// Remove the spam field as it's not needed in the response
		delete(fields, "spam")

		// Use 409 Conflict status code as it's appropriate for responses requiring CAPTCHA
		// https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.10
		status := 409

		// Add error message in the format expected by the API
		// Note: This nested 'error' key may not be consistent with all other API error responses
		// as mentioned in the original Ruby code
		fields["message"] = gin.H{
			"error": spammable.GetFullErrorMessages(),
		}

		// Render the structured API error
		c.JSON(status, fields)
	}

	r.common.WithCaptchaCheckCommon(ctx, spammable, captchaRenderFunc)
}
