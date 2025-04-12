package controller

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// OrganizationsController handles admin organization management
type OrganizationsController struct {
	baseController
	organizationsService *service.OrganizationsService
}

// NewOrganizationsController creates a new organizations controller
func NewOrganizationsController(apiClient *api.Client) *OrganizationsController {
	return &OrganizationsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		organizationsService: service.NewOrganizationsService(apiClient),
	}
}

// Index displays the organizations list
func (c *OrganizationsController) Index(w http.ResponseWriter, r *http.Request) {
	// Check if organizations feature is enabled
	if !c.organizationsService.IsOrganizationsEnabled() {
		helper.RespondWithError(w, http.StatusForbidden, "Organizations feature is not enabled")
		return
	}

	// Set feature flag for organization creation
	helper.SetFeatureFlag(w, "allow_organization_creation", "ops")

	// TODO: Implement HTML template rendering for organizations list
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}
