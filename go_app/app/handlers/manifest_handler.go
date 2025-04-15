package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/gitlab-org/gitlab-redux/app/models"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

const (
	maxManifestSizeMB = 1
)

// ManifestHandler handles manifest import operations
type ManifestHandler struct {
	*BaseHandler
	manifestService *services.ManifestService
}

// NewManifestHandler creates a new ManifestHandler
func NewManifestHandler(baseHandler *BaseHandler, manifestService *services.ManifestService) *ManifestHandler {
	return &ManifestHandler{
		BaseHandler:     baseHandler,
		manifestService: manifestService,
	}
}

// RegisterRoutes registers the routes for the manifest handler
func (h *ManifestHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import/manifest/new", h.New).Methods("GET")
	router.HandleFunc("/import/manifest/status", h.Status).Methods("GET")
	router.HandleFunc("/import/manifest/upload", h.Upload).Methods("POST")
	router.HandleFunc("/import/manifest/create", h.Create).Methods("POST")
}

// New renders the new manifest import form
func (h *ManifestHandler) New(w http.ResponseWriter, r *http.Request) {
	if !h.manifestService.IsImportEnabled() {
		http.Error(w, "Manifest import is not enabled", http.StatusNotFound)
		return
	}

	h.RenderTemplate(w, "import/manifest/new", nil)
}

// Status returns the status of manifest imports
func (h *ManifestHandler) Status(w http.ResponseWriter, r *http.Request) {
	user := h.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	metadata := h.manifestService.GetMetadata(user)
	if metadata == nil {
		http.Redirect(w, r, "/import/manifest/new", http.StatusFound)
		return
	}

	response := map[string]interface{}{
		"repositories": metadata.Repositories,
		"group_id":    metadata.GroupID,
	}

	json.NewEncoder(w).Encode(response)
}

// Upload handles manifest file upload
func (h *ManifestHandler) Upload(w http.ResponseWriter, r *http.Request) {
	user := h.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(maxManifestSizeMB << 20) // Convert MB to bytes
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	// Get the uploaded file
	file, header, err := r.FormFile("manifest")
	if err != nil {
		http.Error(w, "Failed to get manifest file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > maxManifestSizeMB<<20 {
		h.RenderTemplate(w, "import/manifest/new", map[string]interface{}{
			"errors": []string{fmt.Sprintf("Import manifest files cannot exceed %d MB", maxManifestSizeMB)},
		})
		return
	}

	// Get group ID from form
	groupID := r.FormValue("group_id")
	group, err := h.manifestService.GetGroup(groupID)
	if err != nil {
		h.RenderTemplate(w, "import/manifest/new", map[string]interface{}{
			"errors": []string{"Invalid group selected"},
		})
		return
	}

	// Check permissions
	if !h.manifestService.CanImportProjects(user, group) {
		h.RenderTemplate(w, "import/manifest/new", map[string]interface{}{
			"errors": []string{"You don't have enough permissions to import projects in the selected group"},
		})
		return
	}

	// Process manifest
	manifest, err := h.manifestService.ProcessManifest(file)
	if err != nil {
		h.RenderTemplate(w, "import/manifest/new", map[string]interface{}{
			"errors": []string{err.Error()},
		})
		return
	}

	// Save manifest metadata
	err = h.manifestService.SaveMetadata(user, manifest.Projects, group.ID)
	if err != nil {
		h.RenderTemplate(w, "import/manifest/new", map[string]interface{}{
			"errors": []string{"Failed to save manifest metadata"},
		})
		return
	}

	http.Redirect(w, r, "/import/manifest/status", http.StatusFound)
}

// Create creates a new project from manifest
func (h *ManifestHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := h.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		RepoID string `json:"repo_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	metadata := h.manifestService.GetMetadata(user)
	if metadata == nil {
		http.Error(w, "No manifest metadata found", http.StatusBadRequest)
		return
	}

	// Find repository in manifest
	var repo *models.Repository
	for _, r := range metadata.Repositories {
		if r.ID == req.RepoID {
			repo = &r
			break
		}
	}

	if repo == nil {
		http.Error(w, "Repository not found in manifest", http.StatusNotFound)
		return
	}

	// Create project
	project, err := h.manifestService.CreateProject(repo, metadata.GroupID, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(project)
} 