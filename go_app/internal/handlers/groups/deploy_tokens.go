package groups

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/internal/middleware"
	"github.com/cyrex562/gitlab-redux/internal/services/deploy_tokens"
	"gitlab.com/gitlab-org/gitlab/internal/models"
)

// DeployTokensHandler handles deploy token operations for groups
type DeployTokensHandler struct {
	*BaseHandler
	revokeService *deploy_tokens.RevokeService
}

// NewDeployTokensHandler creates a new deploy tokens handler
func NewDeployTokensHandler(baseHandler *BaseHandler, revokeService *deploy_tokens.RevokeService) *DeployTokensHandler {
	return &DeployTokensHandler{
		BaseHandler:   baseHandler,
		revokeService: revokeService,
	}
}

// RegisterRoutes registers the routes for the deploy tokens handler
func (h *DeployTokensHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/:group_id/deploy_tokens/revoke", h.Revoke).Methods("POST")
}

// Revoke handles revoking a deploy token
func (h *DeployTokensHandler) Revoke(w http.ResponseWriter, r *http.Request) {
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

	// Check if user has permission to revoke deploy tokens
	if !h.authorizeDestroyDeployToken(group, user) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get token ID from request
	tokenID := r.FormValue("token_id")
	if tokenID == "" {
		http.Error(w, "Token ID is required", http.StatusBadRequest)
		return
	}

	// Revoke the token
	err := h.revokeService.Execute(group, user, tokenID)
	if err != nil {
		http.Error(w, "Failed to revoke token", http.StatusInternalServerError)
		return
	}

	// Redirect to the repository settings page with the deploy tokens section
	http.Redirect(w, r, "/groups/"+group.Path+"/settings/repository#js-deploy-tokens", http.StatusSeeOther)
}

// authorizeDestroyDeployToken checks if the user has permission to revoke deploy tokens
func (h *DeployTokensHandler) authorizeDestroyDeployToken(group *models.Group, user *models.User) bool {
	// TODO: Implement authorization logic
	// This should check if the user has the appropriate permissions to revoke deploy tokens
	return true
} 