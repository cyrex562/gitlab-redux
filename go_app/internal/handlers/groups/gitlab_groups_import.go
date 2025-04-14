package groups

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/go_app/internal/middleware"
	"gitlab.com/gitlab-org/gitlab/go_app/internal/models"
	"gitlab.com/gitlab-org/gitlab/go_app/internal/services/groups"
)

// GitlabGroupsImportHandler handles GitLab group import operations
type GitlabGroupsImportHandler struct {
	*BaseHandler
	groupService *groups.Service
}

// NewGitlabGroupsImportHandler creates a new GitLab groups import handler
func NewGitlabGroupsImportHandler(baseHandler *BaseHandler, groupService *groups.Service) *GitlabGroupsImportHandler {
	return &GitlabGroupsImportHandler{
		BaseHandler:  baseHandler,
		groupService: groupService,
	}
}

// RegisterRoutes registers the routes for the GitLab groups import handler
func (h *GitlabGroupsImportHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/import/gitlab_groups/create", h.Create).Methods("POST")
}

// Create handles creating a new group import
func (h *GitlabGroupsImportHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Check rate limit
	if err := h.checkImportRateLimit(r); err != nil {
		h.RedirectWithFlash(w, r, "/groups/new#import-group-pane", "alert", "This endpoint has been requested too many times. Try again later.")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.RedirectWithFlash(w, r, "/groups/new#import-group-pane", "alert", "Unable to process group import file")
		return
	}

	// Get form values
	params := &groups.CreateParams{
		Path:     r.FormValue("path"),
		Name:     r.FormValue("name"),
		ParentID: r.FormValue("parent_id"),
	}

	// Get file
	file, _, err := r.FormFile("file")
	if err != nil {
		h.RedirectWithFlash(w, r, "/groups/new#import-group-pane", "alert", "Unable to process group import file")
		return
	}
	defer file.Close()

	// Set visibility level
	params.VisibilityLevel = h.getClosestAllowedVisibilityLevel(params.ParentID)

	// Create group
	user := h.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := h.groupService.Create(r.Context(), params, file, user)
	if err != nil {
		h.RedirectWithFlash(w, r, "/groups/new#import-group-pane", "alert", 
			fmt.Sprintf("Group could not be imported: %s", err.Error()))
		return
	}

	// Start import
	if err := h.groupService.StartImport(r.Context(), result.Group, user); err != nil {
		h.RedirectWithFlash(w, r, result.Group.Path, "alert", "Group import could not be scheduled")
		return
	}

	h.RedirectWithFlash(w, r, result.Group.Path, "notice", 
		fmt.Sprintf("Group '%s' is being imported.", result.Group.Name))
}

// Helper methods

func (h *GitlabGroupsImportHandler) checkImportRateLimit(r *http.Request) error {
	user := h.GetUserFromContext(r)
	if user == nil {
		return fmt.Errorf("unauthorized")
	}

	return h.groupService.CheckImportRateLimit(r.Context(), user)
}

func (h *GitlabGroupsImportHandler) getClosestAllowedVisibilityLevel(parentID string) int {
	if parentID == "" {
		return models.VisibilityPrivate
	}

	parent, err := h.groupService.GetGroup(context.Background(), parentID)
	if err != nil {
		return models.VisibilityPrivate
	}

	return h.groupService.GetClosestAllowedVisibilityLevel(parent.VisibilityLevel)
} 