package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// HooksHandler handles system hook requests
type HooksHandler struct {
	hookService *service.HookService
}

// NewHooksHandler creates a new HooksHandler instance
func NewHooksHandler(hookService *service.HookService) *HooksHandler {
	return &HooksHandler{
		hookService: hookService,
	}
}

// Test handles the POST request to test a system hook
func (h *HooksHandler) Test(c *gin.Context) {
	// Check if system hooks are enabled
	if !h.isSystemHooksEnabled() {
		c.JSON(http.StatusNotFound, gin.H{"error": "System hooks are not available"})
		return
	}

	hookID := c.Param("id")
	trigger := c.Query("trigger")

	result, err := h.hookService.TestHook(c, hookID, trigger)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to test web hook"})
		return
	}

	// Set flash message based on result
	h.setHookExecutionNotice(c, result)

	// Redirect back or to default path
	redirectPath := c.DefaultQuery("redirect", "/admin/hooks")
	c.Redirect(http.StatusFound, redirectPath)
}

// isSystemHooksEnabled checks if system hooks are available
func (h *HooksHandler) isSystemHooksEnabled() bool {
	// TODO: Implement proper check based on your system configuration
	// This should check if the system is not GitLab.com
	return true
}

// setHookExecutionNotice sets a flash message based on the hook test result
func (h *HooksHandler) setHookExecutionNotice(c *gin.Context, result *service.HookTestResult) {
	if result.Success {
		c.SetFlash("notice", "Hook executed successfully")
	} else {
		c.SetFlash("alert", "Hook execution failed: "+result.Error)
	}
}
