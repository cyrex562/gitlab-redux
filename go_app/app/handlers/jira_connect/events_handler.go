package jiraconnect

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"gitlab.com/gitlab-org/gitlab-redux/app/models"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

// EventsHandler handles Jira Connect event-related requests
type EventsHandler struct {
	*ApplicationHandler
	jiraService *services.JiraService
}

// NewEventsHandler creates a new EventsHandler
func NewEventsHandler(appHandler *ApplicationHandler, jiraService *services.JiraService) *EventsHandler {
	return &EventsHandler{
		ApplicationHandler: appHandler,
		jiraService:       jiraService,
	}
}

// InstalledHandler handles the installed event
func (h *EventsHandler) InstalledHandler(w http.ResponseWriter, r *http.Request) {
	// Verify asymmetric JWT
	installation, err := h.verifyAsymmetricJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Transform params to snake_case
	transformedParams := h.transformKeysToSnakeCase(params)

	// Determine if we need to create or update the installation
	var success bool
	if installation != nil {
		// Update existing installation
		success = h.updateInstallation(installation, transformedParams)
	} else {
		// Create new installation
		success = h.createInstallation(transformedParams)
	}

	// Return appropriate response
	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

// UninstalledHandler handles the uninstalled event
func (h *EventsHandler) UninstalledHandler(w http.ResponseWriter, r *http.Request) {
	// Verify asymmetric JWT
	installation, err := h.verifyAsymmetricJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Destroy the installation
	success := h.destroyInstallation(installation)

	// Return appropriate response
	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

// verifyAsymmetricJWT verifies the asymmetric JWT token
func (h *EventsHandler) verifyAsymmetricJWT(r *http.Request) (*models.JiraConnectInstallation, error) {
	// Get the auth token
	authToken := h.getAuthToken(r)
	if authToken == "" {
		return nil, nil
	}

	// Parse request body to get client key
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, err
	}

	// Transform params to snake_case
	transformedParams := h.transformKeysToSnakeCase(params)
	clientKey, ok := transformedParams["client_key"].(string)
	if !ok {
		return nil, nil
	}

	// Get JWT verification claims
	claims := h.getJWTVerificationClaims(r, clientKey)

	// Verify JWT
	if !h.jiraService.VerifyAsymmetricJWT(authToken, claims) {
		return nil, nil
	}

	// Find installation by client key
	return h.jiraService.FindInstallationByClientKey(clientKey), nil
}

// getJWTVerificationClaims gets the JWT verification claims
func (h *EventsHandler) getJWTVerificationClaims(r *http.Request, clientKey string) map[string]interface{} {
	// Calculate audiences
	audiences := h.calculateAudiences()

	// Create QSH
	qsh := h.jiraService.CreateQueryStringHash(r.URL.String(), r.Method, h.jiraService.GetJiraConnectBaseURL(r))

	// Return claims
	return map[string]interface{}{
		"aud": audiences,
		"iss": clientKey,
		"qsh": qsh,
	}
}

// calculateAudiences calculates the JWT audiences
func (h *EventsHandler) calculateAudiences() []string {
	// Get base URL
	baseURL := h.jiraService.GetJiraConnectBaseURL(nil)

	// Check if we need to enforce HTTPS
	var audiences []string
	if h.jiraService.EnforceJiraBaseURLHTTPS() {
		// Replace http with https in base URL
		httpsURL := strings.Replace(baseURL, "http://", "https://", 1)
		audiences = append(audiences, httpsURL)
	} else {
		audiences = append(audiences, baseURL)
	}

	// Check for additional audience URL
	additionalURL := h.jiraService.GetJiraConnectAdditionalAudienceURL()
	if additionalURL != "" {
		// Append path to additional URL
		additionalURLWithPath := h.jiraService.AppendPath(additionalURL, "-/jira_connect")
		audiences = append(audiences, additionalURLWithPath)
	}

	return audiences
}

// createInstallation creates a new Jira Connect installation
func (h *EventsHandler) createInstallation(params map[string]interface{}) bool {
	// Extract required fields
	clientKey, _ := params["client_key"].(string)
	sharedSecret, _ := params["shared_secret"].(string)
	baseURL, _ := params["base_url"].(string)

	// Create installation
	installation := &models.JiraConnectInstallation{
		ClientKey:    clientKey,
		SharedSecret: sharedSecret,
		BaseURL:      baseURL,
	}

	// Save installation
	return h.jiraService.CreateInstallation(installation)
}

// updateInstallation updates an existing Jira Connect installation
func (h *EventsHandler) updateInstallation(installation *models.JiraConnectInstallation, params map[string]interface{}) bool {
	// Extract required fields
	sharedSecret, _ := params["shared_secret"].(string)
	baseURL, _ := params["base_url"].(string)

	// Update installation
	installation.SharedSecret = sharedSecret
	installation.BaseURL = baseURL

	// Save installation
	return h.jiraService.UpdateInstallation(installation)
}

// destroyInstallation destroys a Jira Connect installation
func (h *EventsHandler) destroyInstallation(installation *models.JiraConnectInstallation) bool {
	// Destroy installation
	return h.jiraService.DestroyInstallation(installation)
}

// transformKeysToSnakeCase transforms map keys from camelCase to snake_case
func (h *EventsHandler) transformKeysToSnakeCase(params map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range params {
		// Convert camelCase to snake_case
		snakeKey := h.camelToSnake(key)
		result[snakeKey] = value
	}

	return result
}

// camelToSnake converts a camelCase string to snake_case
func (h *EventsHandler) camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
} 