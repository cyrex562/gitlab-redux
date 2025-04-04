package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// HookLogsHandler handles web hook log requests
type HookLogsHandler struct {
	hookLogService *service.HookLogService
}

// NewHookLogsHandler creates a new HookLogsHandler instance
func NewHookLogsHandler(hookLogService *service.HookLogService) *HookLogsHandler {
	return &HookLogsHandler{
		hookLogService: hookLogService,
	}
}

// Retry handles the POST request to retry a failed web hook
func (h *HookLogsHandler) Retry(c *gin.Context) {
	hookID := c.Param("hook_id")
	logID := c.Param("log_id")

	err := h.hookLogService.RetryHook(c, hookID, logID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retry web hook"})
		return
	}

	// Redirect to the hook edit page after retry
	c.Redirect(http.StatusFound, "/admin/hooks/"+hookID+"/edit")
}

// Show handles the GET request to display a web hook log
func (h *HookLogsHandler) Show(c *gin.Context) {
	hookID := c.Param("hook_id")
	logID := c.Param("log_id")

	log, err := h.hookLogService.GetHookLog(c, hookID, logID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch web hook log"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// getHook retrieves the system hook by ID
func (h *HookLogsHandler) getHook(c *gin.Context) (*model.SystemHook, error) {
	hookID := c.Param("hook_id")
	return h.hookLogService.GetSystemHook(c, hookID)
}
