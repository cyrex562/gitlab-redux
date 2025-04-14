package groups

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/middleware"
	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/models"
)

// DependencyProxiesHandler handles dependency proxy related requests for groups
type DependencyProxiesHandler struct {
	*BaseHandler
}

// NewDependencyProxiesHandler creates a new dependency proxies handler
func NewDependencyProxiesHandler(baseHandler *BaseHandler) *DependencyProxiesHandler {
	return &DependencyProxiesHandler{
		BaseHandler: baseHandler,
	}
}

// RegisterRoutes registers the routes for the dependency proxies handler
func (h *DependencyProxiesHandler) RegisterRoutes(router *middleware.Router) {
	// Routes will be added by child handlers
}

// VerifyDependencyProxyEnabled middleware checks if dependency proxy is enabled for the group
func (h *DependencyProxiesHandler) VerifyDependencyProxyEnabled(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		group := h.GetGroupFromContext(r)
		if group == nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		proxy, err := h.getDependencyProxy(group)
		if err != nil {
			http.Error(w, "Failed to get dependency proxy settings", http.StatusInternalServerError)
			return
		}

		if !proxy.Enabled {
			http.Error(w, "Dependency proxy not enabled", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getDependencyProxy returns the dependency proxy settings for a group
func (h *DependencyProxiesHandler) getDependencyProxy(group *models.Group) (*models.DependencyProxySetting, error) {
	// This would typically interact with your data store
	// For now, we'll return a mock implementation
	return &models.DependencyProxySetting{
		Enabled: true, // This should be fetched from the database
	}, nil
} 