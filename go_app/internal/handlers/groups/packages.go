package groups

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/middleware"
)

// PackagesHandler handles group package operations
type PackagesHandler struct {
	*BaseHandler
}

// NewPackagesHandler creates a new packages handler
func NewPackagesHandler(baseHandler *BaseHandler) *PackagesHandler {
	return &PackagesHandler{
		BaseHandler: baseHandler,
	}
}

// RegisterRoutes registers the routes for the packages handler
func (h *PackagesHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/{group_id}/packages", h.Index).Methods("GET")
	router.HandleFunc("/groups/{group_id}/packages/{id}", h.Show).Methods("GET")
}

// Index handles the packages index page
func (h *PackagesHandler) Index(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !h.verifyPackagesEnabled(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.respondWithHTML(w, "packages/index", map[string]interface{}{
		"group": group,
	})
}

// Show handles displaying a package
// This renders the index template to allow frontend routing to work on page refresh
func (h *PackagesHandler) Show(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !h.verifyPackagesEnabled(r) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.respondWithHTML(w, "packages/index", map[string]interface{}{
		"group": group,
	})
}

// Helper methods

func (h *PackagesHandler) verifyPackagesEnabled(r *http.Request) bool {
	group := h.GetGroupFromContext(r)
	return group.PackagesFeatureEnabled()
} 