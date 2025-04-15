package jiraconnect

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/gitlab-redux/app/models"
)

// BranchesHandler handles Jira Connect branch-related requests
type BranchesHandler struct {
	*ApplicationHandler
}

// NewBranchesHandler creates a new BranchesHandler
func NewBranchesHandler(appHandler *ApplicationHandler) *BranchesHandler {
	return &BranchesHandler{
		ApplicationHandler: appHandler,
	}
}

// NewBranchHandler handles the new branch page request
func (h *BranchesHandler) NewBranchHandler(w http.ResponseWriter, r *http.Request) {
	// Skip JWT verification for this endpoint
	// This is equivalent to skip_before_action :verify_atlassian_jwt!, only: :new in Ruby

	// Get the new branch data
	branchData := h.newBranchData(r)

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branchData)
}

// RouteHandler handles the route request
func (h *BranchesHandler) RouteHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current Jira installation from the context
	installation := h.getCurrentJiraInstallation(r)
	if installation == nil {
		http.Error(w, "Installation not found", http.StatusNotFound)
		return
	}

	// Check if this is a proxy installation
	if h.isProxyInstallation(installation) {
		// Redirect to the create branch URL with the query string
		redirectURL := installation.CreateBranchURL + "?" + r.URL.RawQuery
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	// Otherwise, redirect to the new branch path with the query string
	newBranchPath := "/jira_connect/branches/new"
	redirectURL := newBranchPath + "?" + r.URL.RawQuery
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// isProxyInstallation checks if the installation is a proxy
func (h *BranchesHandler) isProxyInstallation(installation *models.JiraConnectInstallation) bool {
	// This would need to be implemented based on your specific logic
	// For now, we'll assume it's based on a field in the installation model
	return installation.Proxy
}

// initialBranchName generates the initial branch name from issue key and summary
func (h *BranchesHandler) initialBranchName(r *http.Request) string {
	issueKey := r.URL.Query().Get("issue_key")
	if issueKey == "" {
		return ""
	}

	issueSummary := r.URL.Query().Get("issue_summary")
	
	// Convert issue key and summary to branch name
	// This is equivalent to Issue.to_branch_name in Ruby
	return h.toBranchName(issueKey, issueSummary)
}

// toBranchName converts an issue key and summary to a branch name
func (h *BranchesHandler) toBranchName(issueKey, issueSummary string) string {
	// Sanitize the issue summary
	sanitizedSummary := h.sanitizeBranchName(issueSummary)
	
	// Combine issue key and summary
	return issueKey + "-" + sanitizedSummary
}

// sanitizeBranchName sanitizes a string to be used as a branch name
func (h *BranchesHandler) sanitizeBranchName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)
	
	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	
	// Remove any characters that aren't alphanumeric, hyphens, or underscores
	// This is a simplified version - you might need more complex sanitization
	for i := 0; i < len(name); i++ {
		c := name[i]
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			name = name[:i] + name[i+1:]
			i--
		}
	}
	
	// Limit length
	if len(name) > 100 {
		name = name[:100]
	}
	
	return name
}

// newBranchData returns the data needed for the new branch page
func (h *BranchesHandler) newBranchData(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"initial_branch_name": h.initialBranchName(r),
		"success_state_svg_path": "/assets/illustrations/empty-state/empty-merge-requests-md.svg",
	}
} 