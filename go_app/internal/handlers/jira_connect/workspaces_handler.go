package jira_connect

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cyrex562/gitlab-redux/internal/middleware"
	"github.com/cyrex562/gitlab-redux/internal/services"
)

// WorkspacesHandler handles Jira Connect workspace operations
type WorkspacesHandler struct {
	jiraService *services.JiraService
}

// NewWorkspacesHandler creates a new workspaces handler
func NewWorkspacesHandler(jiraService *services.JiraService) *WorkspacesHandler {
	return &WorkspacesHandler{
		jiraService: jiraService,
	}
}

// RegisterRoutes registers the routes for the workspaces handler
func (h *WorkspacesHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/api/v4/jira/connect/workspaces/search", h.Search).Methods("GET")
}

// Search handles searching for workspaces (namespaces) with Jira installations
func (h *WorkspacesHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Get search query from request parameters
	searchQuery := r.URL.Query().Get("searchQuery")
	
	// Sanitize the search query (equivalent to ActionController::Base.helpers.sanitize)
	sanitizedQuery := sanitizeString(searchQuery)
	
	// Get the current Jira installation
	installation, err := h.jiraService.GetCurrentInstallation(r)
	if err != nil {
		http.Error(w, "Failed to get Jira installation", http.StatusInternalServerError)
		return
	}
	
	// Get available namespaces (equivalent to Namespace.without_project_namespaces.with_jira_installation)
	namespaces, err := h.jiraService.GetAvailableNamespaces(installation.ID, sanitizedQuery)
	if err != nil {
		http.Error(w, "Failed to get namespaces", http.StatusInternalServerError)
		return
	}
	
	// Return the workspaces as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Create response structure
	response := map[string]interface{}{
		"workspaces": namespaces,
	}
	
	// Encode and write the response
	json.NewEncoder(w).Encode(response)
}

// sanitizeString sanitizes a string to prevent XSS attacks
// This is a simplified version of ActionController::Base.helpers.sanitize
func sanitizeString(s string) string {
	// Basic sanitization - remove HTML tags
	// In a real implementation, you would use a proper HTML sanitizer
	return strings.TrimSpace(s)
} 