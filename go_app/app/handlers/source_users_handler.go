package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/gitlab-org/gitlab-redux/app/models"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// SourceUsersHandler handles source user operations
type SourceUsersHandler struct {
	*BaseHandler
	sourceUserService *services.SourceUserService
}

// NewSourceUsersHandler creates a new SourceUsersHandler
func NewSourceUsersHandler(baseHandler *BaseHandler, sourceUserService *services.SourceUserService) *SourceUsersHandler {
	return &SourceUsersHandler{
		BaseHandler:       baseHandler,
		sourceUserService: sourceUserService,
	}
}

// RegisterRoutes registers the routes for the source users handler
func (h *SourceUsersHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import/source_users/accept", h.Accept).Methods("POST")
	router.HandleFunc("/import/source_users/decline", h.Decline).Methods("POST")
	router.HandleFunc("/import/source_users/show", h.Show).Methods("GET")
}

// Accept handles accepting a source user reassignment
func (h *SourceUsersHandler) Accept(w http.ResponseWriter, r *http.Request) {
	user := h.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reassignmentToken := r.URL.Query().Get("reassignment_token")
	if reassignmentToken == "" {
		http.Error(w, "Reassignment token is required", http.StatusBadRequest)
		return
	}

	// Check if feature flag is enabled
	if !h.sourceUserService.IsFeatureEnabled(user) {
		http.Error(w, "Feature not enabled", http.StatusNotFound)
		return
	}

	// Get source user
	sourceUser, err := h.sourceUserService.GetSourceUserByToken(reassignmentToken)
	if err != nil {
		h.renderError(w, r, "invalid_invite")
		return
	}

	// Validate source user
	if !h.isSourceUserValid(sourceUser, user) {
		h.renderError(w, r, "invalid_invite")
		return
	}

	// Accept reassignment
	result, err := h.sourceUserService.AcceptReassignment(sourceUser, user, reassignmentToken)
	if err != nil {
		h.renderError(w, r, "accept_invite")
		return
	}

	if result.Success {
		h.renderBanner(w, r, "accept_invite", sourceUser)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/?alert=The+invitation+could+not+be+accepted.", http.StatusFound)
	}
}

// Decline handles declining a source user reassignment
func (h *SourceUsersHandler) Decline(w http.ResponseWriter, r *http.Request) {
	user := h.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reassignmentToken := r.URL.Query().Get("reassignment_token")
	if reassignmentToken == "" {
		http.Error(w, "Reassignment token is required", http.StatusBadRequest)
		return
	}

	// Check if feature flag is enabled
	if !h.sourceUserService.IsFeatureEnabled(user) {
		http.Error(w, "Feature not enabled", http.StatusNotFound)
		return
	}

	// Get source user
	sourceUser, err := h.sourceUserService.GetSourceUserByToken(reassignmentToken)
	if err != nil {
		h.renderError(w, r, "invalid_invite")
		return
	}

	// Validate source user
	if !h.isSourceUserValid(sourceUser, user) {
		h.renderError(w, r, "invalid_invite")
		return
	}

	// Decline reassignment
	result, err := h.sourceUserService.DeclineReassignment(sourceUser, user, reassignmentToken)
	if err != nil {
		h.renderError(w, r, "reject_invite")
		return
	}

	if result.Success {
		h.renderBanner(w, r, "reject_invite", sourceUser)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/?alert=The+invitation+could+not+be+declined.", http.StatusFound)
	}
}

// Show displays the source user page
func (h *SourceUsersHandler) Show(w http.ResponseWriter, r *http.Request) {
	user := h.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reassignmentToken := r.URL.Query().Get("reassignment_token")
	if reassignmentToken == "" {
		http.Error(w, "Reassignment token is required", http.StatusBadRequest)
		return
	}

	// Check if feature flag is enabled
	if !h.sourceUserService.IsFeatureEnabled(user) {
		http.Error(w, "Feature not enabled", http.StatusNotFound)
		return
	}

	// Get source user
	sourceUser, err := h.sourceUserService.GetSourceUserByToken(reassignmentToken)
	if err != nil {
		h.renderError(w, r, "invalid_invite")
		return
	}

	// Validate source user
	if !h.isSourceUserValid(sourceUser, user) {
		h.renderError(w, r, "invalid_invite")
		return
	}

	// Render the show page
	h.RenderTemplate(w, "import/source_users/show", map[string]interface{}{
		"source_user": sourceUser,
	})
}

// isSourceUserValid checks if the source user is valid for the current user
func (h *SourceUsersHandler) isSourceUserValid(sourceUser *models.SourceUser, currentUser *models.User) bool {
	return sourceUser != nil && 
	       sourceUser.AwaitingApproval && 
	       sourceUser.ReassignToUserID == currentUser.ID
}

// renderError renders an error banner and redirects to the root path
func (h *SourceUsersHandler) renderError(w http.ResponseWriter, r *http.Request, bannerType string) {
	h.renderBanner(w, r, bannerType, nil)
	http.Redirect(w, r, "/", http.StatusFound)
}

// renderBanner renders a banner with the given type and source user
func (h *SourceUsersHandler) renderBanner(w http.ResponseWriter, r *http.Request, bannerType string, sourceUser *models.SourceUser) {
	// This would typically set a flash message or session variable
	// For now, we'll just set a cookie with the banner type
	cookie := &http.Cookie{
		Name:  "banner_type",
		Value: bannerType,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
} 