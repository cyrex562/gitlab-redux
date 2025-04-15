package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/app/models"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	// No need for templates as Gin handles template rendering
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// GetCurrentUser retrieves the current user from the request
func (h *BaseHandler) GetCurrentUser(r *http.Request) *models.User {
	// Get user from Gin context
	if ginCtx, ok := r.Context().Value("gin-context").(*gin.Context); ok {
		if user, exists := ginCtx.Get("current_user"); exists {
			if currentUser, ok := user.(*models.User); ok {
				return currentUser
			}
		}
	}
	return nil
}

// SetCurrentUser sets the current user in the Gin context
func (h *BaseHandler) SetCurrentUser(c *gin.Context, user *models.User) {
	c.Set("current_user", user)
} 