package groups

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/middleware"
)

// DependencyProxyAuthHandler handles authentication for dependency proxy
type DependencyProxyAuthHandler struct {
	*DependencyProxiesHandler
}

// NewDependencyProxyAuthHandler creates a new dependency proxy auth handler
func NewDependencyProxyAuthHandler(baseHandler *DependencyProxiesHandler) *DependencyProxyAuthHandler {
	return &DependencyProxyAuthHandler{
		DependencyProxiesHandler: baseHandler,
	}
}

// RegisterRoutes registers the routes for the dependency proxy auth handler
func (h *DependencyProxyAuthHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/:group_id/dependency_proxy_auth", h.Authenticate).Methods("GET")
}

// Authenticate handles the authentication endpoint
func (h *DependencyProxyAuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Return empty response with 200 OK status, matching the Ruby implementation
	w.WriteHeader(http.StatusOK)
} 