package jiraconnect

import (
	"context"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/gitlab-redux/app/models"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// ApplicationHandler serves as the base controller for Jira Connect handlers
type ApplicationHandler struct {
	*BaseHandler
	jiraService *services.JiraService
}

// NewApplicationHandler creates a new ApplicationHandler
func NewApplicationHandler(baseHandler *BaseHandler, jiraService *services.JiraService) *ApplicationHandler {
	return &ApplicationHandler{
		BaseHandler: baseHandler,
		jiraService: jiraService,
	}
}

// VerifyAtlassianJWTMiddleware creates a middleware that verifies the Atlassian JWT token
func (h *ApplicationHandler) VerifyAtlassianJWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip verification for app descriptor
		if r.URL.Path == "/jira_connect/app_descriptor" {
			next(w, r)
			return
		}

		// Verify JWT token
		if !h.verifyAtlassianJWT(r) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Store current Jira installation in request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "current_jira_installation", h.getCurrentJiraInstallation(r))
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// verifyAtlassianJWT verifies the Atlassian JWT token
func (h *ApplicationHandler) verifyAtlassianJWT(r *http.Request) bool {
	// Get installation from JWT
	installation := h.installationFromJWT(r)
	if installation == nil {
		return false
	}

	// Verify JWT signature with stored shared secret
	return h.jiraService.VerifyJWT(r, installation.SharedSecret)
}

// verifyQSHClaim verifies the QSH claim in the JWT token
func (h *ApplicationHandler) verifyQSHClaim(r *http.Request) bool {
	// Skip verification for JSON requests with context-qsh claim
	if r.Header.Get("Content-Type") == "application/json" && h.jiraService.VerifyContextQSHClaim(r) {
		return true
	}

	// Verify QSH claim matches the current request
	return h.jiraService.VerifyQSHClaim(r, r.URL.String(), r.Method, h.jiraService.GetJiraConnectBaseURL(r))
}

// installationFromJWT retrieves the Jira installation from the JWT token
func (h *ApplicationHandler) installationFromJWT(r *http.Request) *models.JiraConnectInstallation {
	// Get issuer claim from JWT
	issClaim := h.jiraService.GetIssuerClaim(r)
	if issClaim == "" {
		return nil
	}

	// Find installation by client key
	return h.jiraService.FindInstallationByClientKey(issClaim)
}

// getCurrentJiraInstallation retrieves the current Jira installation from the request context
func (h *ApplicationHandler) getCurrentJiraInstallation(r *http.Request) *models.JiraConnectInstallation {
	installation, ok := r.Context().Value("current_jira_installation").(*models.JiraConnectInstallation)
	if !ok {
		return nil
	}
	return installation
}

// getJiraUser retrieves the Jira user from the JWT token
func (h *ApplicationHandler) getJiraUser(r *http.Request) *models.JiraUser {
	// Get installation from JWT
	installation := h.installationFromJWT(r)
	if installation == nil {
		return nil
	}

	// Get subject claim from JWT
	subClaim := h.jiraService.GetSubjectClaim(r)
	if subClaim == "" {
		return nil
	}

	// Get user info from Jira
	return h.jiraService.GetUserInfo(installation, subClaim)
}

// getAuthToken retrieves the authentication token from the request
func (h *ApplicationHandler) getAuthToken(r *http.Request) string {
	// Check for JWT parameter
	jwt := r.URL.Query().Get("jwt")
	if jwt != "" {
		return jwt
	}

	// Check for Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 {
			return parts[1]
		}
	}

	return ""
} 