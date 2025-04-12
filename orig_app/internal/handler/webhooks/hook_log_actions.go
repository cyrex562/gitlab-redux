package webhooks

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// HookLogActions handles actions related to webhook logs
type HookLogActions struct {
	resendService *service.WebHooksEventsResendService
	executionNotice *HookExecutionNotice
	logger *util.Logger
}

// NewHookLogActions creates a new instance of HookLogActions
func NewHookLogActions(resendService *service.WebHooksEventsResendService, executionNotice *HookExecutionNotice, logger *util.Logger) *HookLogActions {
	return &HookLogActions{
		resendService: resendService,
		executionNotice: executionNotice,
		logger: logger,
	}
}

// Show displays a webhook log
func (h *HookLogActions) Show(ctx *gin.Context, hook model.WebHook) {
	// Hide search settings
	ctx.Set("hide_search_settings", true)

	// Get the hook log ID from the URL parameters
	hookLogID := ctx.Param("id")

	// Find the hook log
	hookLog, err := h.getHookLog(hook, hookLogID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Hook log not found"})
		return
	}

	// Render the show template
	ctx.HTML(200, "webhooks/logs/show.html", gin.H{
		"hook": hook,
		"hookLog": hookLog,
		"hide_search_settings": true,
	})
}

// Retry retries a webhook execution
func (h *HookLogActions) Retry(ctx *gin.Context, hook model.WebHook) {
	// Get the hook log ID from the URL parameters
	hookLogID := ctx.Param("id")

	// Find the hook log
	hookLog, err := h.getHookLog(hook, hookLogID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Hook log not found"})
		return
	}

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	// Execute the hook
	result := h.executeHook(ctx, hookLog, user)

	if result.Success {
		// Redirect to the after retry path
		ctx.Redirect(302, h.getAfterRetryRedirectPath(ctx, hook))
	} else {
		// Set a warning flash message
		ctx.Set("flash_warning", result.Message)

		// Redirect back with a fallback
		ctx.Redirect(302, h.getAfterRetryRedirectPath(ctx, hook))
	}
}

// getHookLog retrieves a webhook log by ID
func (h *HookLogActions) getHookLog(hook model.WebHook, hookLogID string) (model.WebHookLog, error) {
	return hook.GetWebHookLogs().Find(hookLogID)
}

// executeHook executes a webhook and sets the execution notice
func (h *HookLogActions) executeHook(ctx *gin.Context, hookLog model.WebHookLog, user model.User) service.ServiceResponse {
	// Execute the hook
	result := h.resendService.Execute(hookLog, user)

	// Set the hook execution notice
	h.executionNotice.SetHookExecutionNotice(ctx, result.ToMap())

	return result
}

// getAfterRetryRedirectPath returns the path to redirect to after a retry
func (h *HookLogActions) getAfterRetryRedirectPath(ctx *gin.Context, hook model.WebHook) string {
	// This would typically be a method on the hook model or a helper function
	// For now, we'll return a default path
	return "/webhooks/" + hook.ID() + "/logs"
}
