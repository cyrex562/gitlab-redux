package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// SessionsController handles admin mode authentication
type SessionsController struct {
	baseController
	sessionsService *service.SessionsService
}

// NewSessionsController creates a new sessions controller
func NewSessionsController(apiClient *api.Client) *SessionsController {
	return &SessionsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		sessionsService: service.NewSessionsService(apiClient),
	}
}

// userIsAdmin checks if the current user has admin access
func (c *SessionsController) userIsAdmin(r *http.Request) bool {
	return c.sessionsService.CanAccessAdminArea(r)
}

// New displays the admin mode login form
func (c *SessionsController) New(w http.ResponseWriter, r *http.Request) {
	if !c.userIsAdmin(r) {
		helper.RespondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	// Check if user is already in admin mode
	if c.sessionsService.IsAdminMode(r) {
		redirectPath := c.getRedirectPath(r)
		http.Redirect(w, r, redirectPath, http.StatusFound)
		return
	}

	// Request admin mode if not already requested
	if !c.sessionsService.IsAdminModeRequested(r) {
		if err := c.sessionsService.RequestAdminMode(r); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to request admin mode")
			return
		}
	}

	// Store redirect location
	if err := c.sessionsService.StoreRedirectLocation(r); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to store redirect location")
		return
	}

	// TODO: Implement HTML template rendering for admin mode login form
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Create handles admin mode authentication
func (c *SessionsController) Create(w http.ResponseWriter, r *http.Request) {
	if !c.userIsAdmin(r) {
		helper.RespondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	// Check if admin mode was requested
	if !c.sessionsService.IsAdminModeRequested(r) {
		http.Redirect(w, r, "/admin/sessions/new", http.StatusFound)
		return
	}

	var params struct {
		User struct {
			Password      string `json:"password"`
			OTPAttempt    string `json:"otp_attempt"`
			DeviceResponse string `json:"device_response"`
		} `json:"user"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if two-factor authentication is enabled
	if c.sessionsService.IsTwoFactorEnabled(r) {
		// Handle two-factor authentication
		if err := c.sessionsService.AuthenticateWithTwoFactor(r, params.User.OTPAttempt, params.User.DeviceResponse); err != nil {
			// TODO: Implement HTML template rendering for two-factor authentication form
			helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
			return
		}
	} else {
		// Handle password authentication
		if err := c.sessionsService.EnableAdminMode(r, params.User.Password); err != nil {
			// TODO: Implement HTML template rendering for login form with error
			helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
			return
		}
	}

	// Redirect to stored location or admin root
	redirectPath := c.getRedirectPath(r)
	http.Redirect(w, r, redirectPath, http.StatusFound)
}

// Destroy disables admin mode
func (c *SessionsController) Destroy(w http.ResponseWriter, r *http.Request) {
	if !c.userIsAdmin(r) {
		helper.RespondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	if err := c.sessionsService.DisableAdminMode(r); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to disable admin mode")
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// getRedirectPath determines the redirect path after admin mode authentication
func (c *SessionsController) getRedirectPath(r *http.Request) string {
	redirectPath := c.sessionsService.GetStoredRedirectLocation(r)
	if redirectPath == "" {
		redirectPath = "/admin"
	}

	// Check if redirect path is excluded
	if c.isExcludedRedirectPath(redirectPath) {
		return "/admin"
	}

	return redirectPath
}

// isExcludedRedirectPath checks if the path is in the list of excluded paths
func (c *SessionsController) isExcludedRedirectPath(path string) bool {
	excludedPaths := []string{
		"/admin/sessions/new",
		"/admin/sessions",
	}

	for _, excluded := range excludedPaths {
		if path == excluded {
			return true
		}
	}

	return false
}
