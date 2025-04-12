package redirect

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// ContinueParams handles processing and validating continue parameters for redirects
type ContinueParams struct {
	redirectService *service.RedirectService
	logger          *util.Logger
}

// NewContinueParams creates a new instance of ContinueParams
func NewContinueParams(
	redirectService *service.RedirectService,
	logger *util.Logger,
) *ContinueParams {
	return &ContinueParams{
		redirectService: redirectService,
		logger:          logger,
	}
}

// GetContinueParams gets and validates the continue parameters from the request
func (c *ContinueParams) GetContinueParams(ctx *gin.Context) map[string]interface{} {
	// Get the continue parameters from the request
	continueParams, exists := ctx.GetQuery("continue")
	if !exists {
		return make(map[string]interface{})
	}

	// Parse the continue parameters
	params := make(map[string]interface{})

	// Get the redirect path
	to, exists := ctx.GetQuery("continue[to]")
	if exists {
		// Validate and sanitize the redirect path
		params["to"] = c.redirectService.SafeRedirectPath(ctx, to)
	}

	// Get the notice
	notice, exists := ctx.GetQuery("continue[notice]")
	if exists {
		params["notice"] = notice
	}

	// Get the notice_now
	noticeNow, exists := ctx.GetQuery("continue[notice_now]")
	if exists {
		params["notice_now"] = noticeNow
	}

	return params
}
