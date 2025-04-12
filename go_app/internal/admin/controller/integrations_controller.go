package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// IntegrationsController handles admin integration-related requests
type IntegrationsController struct {
	baseController
	integrationService *service.IntegrationService
}

// NewIntegrationsController creates a new integrations controller
func NewIntegrationsController(apiClient *api.Client) *IntegrationsController {
	return &IntegrationsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		integrationService: service.NewIntegrationService(apiClient),
	}
}

// Overrides handles the overrides endpoint
func (c *IntegrationsController) Overrides(w http.ResponseWriter, r *http.Request) {
	if !c.integrationService.IsInstanceLevelEnabled() {
		helper.RespondWithError(w, http.StatusNotFound, "Instance level integrations are not enabled")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "json" {
		c.handleJSONOverrides(w, r)
		return
	}

	// Default to HTML response
	c.handleHTMLOverrides(w, r)
}

func (c *IntegrationsController) handleJSONOverrides(w http.ResponseWriter, r *http.Request) {
	projects, err := c.integrationService.GetProjectsWithActiveIntegration()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	response := struct {
		Projects []interface{} `json:"projects"`
	}{
		Projects: projects,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *IntegrationsController) handleHTMLOverrides(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement HTML template rendering
	// This would typically use a template engine to render the overrides view
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}
