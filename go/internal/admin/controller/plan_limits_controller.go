package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// PlanLimitsController handles admin plan limits management
type PlanLimitsController struct {
	baseController
	planLimitsService *service.PlanLimitsService
}

// NewPlanLimitsController creates a new plan limits controller
func NewPlanLimitsController(apiClient *api.Client) *PlanLimitsController {
	return &PlanLimitsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		planLimitsService: service.NewPlanLimitsService(apiClient),
	}
}

// Create updates plan limits
func (c *PlanLimitsController) Create(w http.ResponseWriter, r *http.Request) {
	var params struct {
		PlanID                      int64 `json:"plan_id"`
		ConanMaxFileSize           int64 `json:"conan_max_file_size"`
		HelmMaxFileSize            int64 `json:"helm_max_file_size"`
		MavenMaxFileSize           int64 `json:"maven_max_file_size"`
		NpmMaxFileSize             int64 `json:"npm_max_file_size"`
		NugetMaxFileSize           int64 `json:"nuget_max_file_size"`
		PypiMaxFileSize            int64 `json:"pypi_max_file_size"`
		TerraformModuleMaxFileSize int64 `json:"terraform_module_max_file_size"`
		GenericPackagesMaxFileSize int64 `json:"generic_packages_max_file_size"`
		CiInstanceLevelVariables   int64 `json:"ci_instance_level_variables"`
		CiPipelineSize            int64 `json:"ci_pipeline_size"`
		CiActiveJobs              int64 `json:"ci_active_jobs"`
		CiProjectSubscriptions    int64 `json:"ci_project_subscriptions"`
		CiPipelineSchedules       int64 `json:"ci_pipeline_schedules"`
		CiNeedsSizeLimit          int64 `json:"ci_needs_size_limit"`
		CiRegisteredGroupRunners  int64 `json:"ci_registered_group_runners"`
		CiRegisteredProjectRunners int64 `json:"ci_registered_project_runners"`
		DotenvSize                int64 `json:"dotenv_size"`
		DotenvVariables           int64 `json:"dotenv_variables"`
		PipelineHierarchySize     int64 `json:"pipeline_hierarchy_size"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get the referer path for redirection
	redirectPath := r.Referer()
	if redirectPath == "" {
		redirectPath = "/admin/application_settings/general"
	}

	// Update plan limits
	if err := c.planLimitsService.UpdatePlanLimits(params.PlanID, params); err != nil {
		format := r.URL.Query().Get("format")
		if format == "json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		helper.RespondWithError(w, http.StatusUnprocessableEntity, "Failed to update plan limits")
		return
	}

	// Handle different response formats
	format := r.URL.Query().Get("format")
	if format == "json" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Default to HTML response with redirect
	http.Redirect(w, r, redirectPath, http.StatusFound)
}
