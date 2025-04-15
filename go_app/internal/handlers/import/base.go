package importhandler

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab/go_app/internal/middleware"
)

// BaseHandler provides common functionality for import handlers
type BaseHandler struct {
	router *middleware.Router
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(router *middleware.Router) *BaseHandler {
	return &BaseHandler{
		router: router,
	}
}

// GetUserFromContext gets the current user from the request context
func (h *BaseHandler) GetUserFromContext(r *http.Request) interface{} {
	// TODO: Implement user retrieval from context
	return nil
}

// RedirectWithFlash redirects to a URL with a flash message
func (h *BaseHandler) RedirectWithFlash(w http.ResponseWriter, r *http.Request, url, flashType, message string) {
	// TODO: Implement flash message handling
	http.Redirect(w, r, url, http.StatusSeeOther)
} 