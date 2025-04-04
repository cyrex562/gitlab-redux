package webhooks

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// HookExecutionNotice handles setting flash notices based on webhook execution results
type HookExecutionNotice struct{}

// NewHookExecutionNotice creates a new instance of HookExecutionNotice
func NewHookExecutionNotice() *HookExecutionNotice {
	return &HookExecutionNotice{}
}

// SetHookExecutionNotice sets a flash notice based on the webhook execution result
func (h *HookExecutionNotice) SetHookExecutionNotice(ctx *gin.Context, result map[string]interface{}) {
	// Extract HTTP status and message from the result
	httpStatus, hasStatus := result["http_status"].(int)
	message, hasMessage := result["message"].(string)

	if !hasMessage {
		message = "Unknown error"
	}

	// Set appropriate flash notice based on HTTP status
	if hasStatus && httpStatus >= 200 && httpStatus < 400 {
		// Success case
		ctx.Set("flash_notice", fmt.Sprintf("Hook executed successfully: HTTP %d", httpStatus))
	} else if hasStatus {
		// HTTP error case
		ctx.Set("flash_alert", fmt.Sprintf("Hook executed successfully but returned HTTP %d %s", httpStatus, message))
	} else {
		// Execution failure case
		ctx.Set("flash_alert", fmt.Sprintf("Hook execution failed: %s", message))
	}
}
