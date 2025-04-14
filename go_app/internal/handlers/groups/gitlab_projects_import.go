package groups

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"gitlab.com/gitlab-org/gitlab/go_app/internal/middleware"
	"gitlab.com/gitlab-org/gitlab/go_app/internal/services/groups"
)

// GitlabProjectsImportHandler handles GitLab project import operations
type GitlabProjectsImportHandler struct {
	*BaseHandler
	groupService *groups.Service
	templates    *template.Template
}

// NewGitlabProjectsImportHandler creates a new GitLab projects import handler
func NewGitlabProjectsImportHandler(baseHandler *BaseHandler, groupService *groups.Service, templates *template.Template) *GitlabProjectsImportHandler {
	return &GitlabProjectsImportHandler{
		BaseHandler:  baseHandler,
		groupService: groupService,
		templates:    templates,
	}
}

// RegisterRoutes registers the routes for the GitLab projects import handler
func (h *GitlabProjectsImportHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/import/gitlab_projects/new", h.New).Methods("GET")
	router.HandleFunc("/import/gitlab_projects/create", h.Create).Methods("POST")
}

// New handles displaying the new import form
func (h *GitlabProjectsImportHandler) New(w http.ResponseWriter, r *http.Request) {
	// Get namespace from query params
	namespaceID := r.URL.Query().Get("namespace_id")
	if namespaceID == "" {
		http.Error(w, "Namespace ID is required", http.StatusBadRequest)
		return
	}

	// Get namespace
	namespace, err := h.groupService.GetGroup(r.Context(), namespaceID)
	if err != nil {
		http.Error(w, "Namespace not found", http.StatusNotFound)
		return
	}

	// Check if user can import projects
	user := h.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !h.groupService.CanImportProjects(r.Context(), user, namespace) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Render template
	h.RenderTemplate(w, "import/gitlab_projects/new", map[string]interface{}{
		"namespace": namespace,
		"path":      r.URL.Query().Get("path"),
	})
}

// Create handles creating a new project import
func (h *GitlabProjectsImportHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.RedirectWithFlash(w, r, "/import/gitlab_projects/new", "alert", "Unable to process project import file")
		return
	}

	// Get form values
	params := &groups.CreateProjectParams{
		Name:        r.FormValue("name"),
		Path:        r.FormValue("path"),
		NamespaceID: r.FormValue("namespace_id"),
	}

	// Get file
	file, header, err := r.FormFile("file")
	if err != nil {
		h.RedirectWithFlash(w, r, "/import/gitlab_projects/new", "alert", "Unable to process project import file")
		return
	}
	defer file.Close()

	// Validate file
	if !h.isValidImportFile(header.Filename) {
		h.RedirectWithFlash(w, r, "/import/gitlab_projects/new", "alert", "You need to upload a GitLab project export archive (ending in .gz)")
		return
	}

	// Create project
	user := h.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := h.groupService.CreateProject(r.Context(), params, file, user)
	if err != nil {
		h.RedirectWithFlash(w, r, "/import/gitlab_projects/new", "alert", 
			fmt.Sprintf("Project could not be imported: %s", err.Error()))
		return
	}

	h.RedirectWithFlash(w, r, fmt.Sprintf("/projects/%s", result.Project.Path), "notice", 
		fmt.Sprintf("Project '%s' is being imported.", result.Project.Name))
}

// Helper methods

func (h *GitlabProjectsImportHandler) isValidImportFile(filename string) bool {
	return filepath.Ext(filename) == ".gz"
}

// RenderTemplate renders a template with the given data
func (h *GitlabProjectsImportHandler) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	return h.templates.ExecuteTemplate(w, name, data)
} 