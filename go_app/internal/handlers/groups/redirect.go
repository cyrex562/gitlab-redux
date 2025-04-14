package groups

import (
	"net/http"
	"strconv"

	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/middleware"
	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/services"
)

// RedirectHandler handles group redirect operations
type RedirectHandler struct {
	groupService *services.GroupService
	authorizer   *services.Authorizer
}

// NewRedirectHandler creates a new redirect handler
func NewRedirectHandler(groupService *services.GroupService, authorizer *services.Authorizer) *RedirectHandler {
	return &RedirectHandler{
		groupService: groupService,
		authorizer:   authorizer,
	}
}

// RegisterRoutes registers the routes for the redirect handler
func (h *RedirectHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/redirect/{id}", h.RedirectFromID).Methods("GET")
}

// RedirectFromID handles redirecting from a group ID to the group's page
func (h *RedirectHandler) RedirectFromID(w http.ResponseWriter, r *http.Request) {
	// Get group ID from URL parameters
	groupID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Get group from service
	group, err := h.groupService.GetGroup(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Get user from context
	user := h.getUserFromContext(r)
	if user == nil {
		// No user in context, which is fine as this controller skips authentication
	}

	// Check if user can read the group
	if !h.authorizer.CanReadGroup(user, group) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Redirect to the group's page
	http.Redirect(w, r, "/groups/"+strconv.FormatInt(group.ID, 10), http.StatusSeeOther)
}

// Helper methods

func (h *RedirectHandler) getUserFromContext(r *http.Request) *models.User {
	// Implementation depends on your context management
	// This is a placeholder that should be replaced with your actual implementation
	return nil
} 