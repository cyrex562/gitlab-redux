package groups

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/middleware"
)

// InfrastructureRegistryHandler handles infrastructure registry operations
type InfrastructureRegistryHandler struct {
	*BaseHandler
}

// NewInfrastructureRegistryHandler creates a new infrastructure registry handler
func NewInfrastructureRegistryHandler(baseHandler *BaseHandler) *InfrastructureRegistryHandler {
	return &InfrastructureRegistryHandler{
		BaseHandler: baseHandler,
	}
}

// RegisterRoutes registers the routes for the infrastructure registry handler
func (h *InfrastructureRegistryHandler) RegisterRoutes(router *middleware.Router) {
	// Routes will be added by child handlers
	// This base controller just provides the packages verification middleware
}

// VerifyPackagesEnabled middleware verifies that packages are enabled for the group
func (h *InfrastructureRegistryHandler) VerifyPackagesEnabled(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		group := h.GetGroupFromContext(r)
		if group == nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		if !group.PackagesFeatureEnabled() {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
} 