package groups

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cyrex562/gitlab-redux/internal/middleware"
	"github.com/cyrex562/gitlab-redux/internal/models"
	"github.com/cyrex562/gitlab-redux/internal/services/group_links"
)

// GroupLinksHandler handles group link operations
type GroupLinksHandler struct {
	*BaseHandler
	updateService  *group_links.UpdateService
	destroyService *group_links.DestroyService
}

// NewGroupLinksHandler creates a new group links handler
func NewGroupLinksHandler(
	baseHandler *BaseHandler,
	updateService *group_links.UpdateService,
	destroyService *group_links.DestroyService,
) *GroupLinksHandler {
	return &GroupLinksHandler{
		BaseHandler:    baseHandler,
		updateService:  updateService,
		destroyService: destroyService,
	}
}

// RegisterRoutes registers the routes for the group links handler
func (h *GroupLinksHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/:group_id/group_links/:id", h.Update).Methods("PUT")
	router.HandleFunc("/groups/:group_id/group_links/:id", h.Destroy).Methods("DELETE")
}

// Update handles updating a group link
func (h *GroupLinksHandler) Update(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	user := h.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check if user has admin permissions
	if !h.authorizeAdminGroup(group, user) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get group link ID from URL
	groupLinkID := r.URL.Query().Get(":id")
	if groupLinkID == "" {
		http.Error(w, "Group link ID is required", http.StatusBadRequest)
		return
	}

	// Find the group link
	groupLink, err := group.FindSharedWithGroupLink(groupLinkID)
	if err != nil {
		http.Error(w, "Group link not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var params struct {
		GroupAccess  string    `json:"group_access"`
		ExpiresAt    time.Time `json:"expires_at"`
		MemberRoleID int64     `json:"member_role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the group link
	err = h.updateService.Execute(groupLink, user, params.GroupAccess, params.ExpiresAt, params.MemberRoleID)
	if err != nil {
		http.Error(w, "Failed to update group link", http.StatusInternalServerError)
		return
	}

	// Return response based on whether the link expires
	if groupLink.Expires() {
		response := map[string]interface{}{
			"expires_in":   h.timeAgoWithTooltip(groupLink.ExpiresAt),
			"expires_soon": groupLink.ExpiresSoon(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Destroy handles destroying a group link
func (h *GroupLinksHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	user := h.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check if user has admin permissions
	if !h.authorizeAdminGroup(group, user) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get group link ID from URL
	groupLinkID := r.URL.Query().Get(":id")
	if groupLinkID == "" {
		http.Error(w, "Group link ID is required", http.StatusBadRequest)
		return
	}

	// Find the group link
	groupLink, err := group.FindSharedWithGroupLink(groupLinkID)
	if err != nil {
		http.Error(w, "Group link not found", http.StatusNotFound)
		return
	}

	// Destroy the group link
	err = h.destroyService.Execute(group, user, groupLink)
	if err != nil {
		http.Error(w, "Failed to destroy group link", http.StatusInternalServerError)
		return
	}

	// Handle different response formats
	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, "/groups/"+group.Path+"/group_members", http.StatusFound)
	}
}

// authorizeAdminGroup checks if the user has admin permissions for the group
func (h *GroupLinksHandler) authorizeAdminGroup(group *models.Group, user *models.User) bool {
	// TODO: Implement authorization logic
	// This should check if the user has admin permissions for the group
	return true
}

// timeAgoWithTooltip returns a formatted time ago string with tooltip
func (h *GroupLinksHandler) timeAgoWithTooltip(t time.Time) string {
	// TODO: Implement time ago with tooltip formatting
	return ""
} 