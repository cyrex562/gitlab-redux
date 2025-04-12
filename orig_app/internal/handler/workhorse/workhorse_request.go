package workhorse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// WorkhorseRequestHandler handles Workhorse request verification
type WorkhorseRequestHandler struct {
	workhorseService *service.WorkhorseService
}

// NewWorkhorseRequestHandler creates a new Workhorse request handler
func NewWorkhorseRequestHandler(workhorseService *service.WorkhorseService) *WorkhorseRequestHandler {
	return &WorkhorseRequestHandler{
		workhorseService: workhorseService,
	}
}

// VerifyWorkhorseAPI verifies that the request is coming from GitLab Workhorse
func (h *WorkhorseRequestHandler) VerifyWorkhorseAPI(c *gin.Context) {
	// Get the request headers
	headers := c.Request.Header

	// Verify the request
	err := h.workhorseService.VerifyAPIRequest(headers)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Workhorse request"})
		c.Abort()
		return
	}

	// Continue to the next handler
	c.Next()
}

// RegisterMiddleware registers the Workhorse request verification middleware
func (h *WorkhorseRequestHandler) RegisterMiddleware(router *gin.Engine) {
	router.Use(h.VerifyWorkhorseAPI)
}
