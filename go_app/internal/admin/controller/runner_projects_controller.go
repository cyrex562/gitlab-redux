package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// RunnerProjectsController handles admin runner-project assignment management
type RunnerProjectsController struct {
	baseController
	runnerProjectsService *service.RunnerProjectsService
}

// NewRunnerProjectsController creates a new runner projects controller
func NewRunnerProjectsController(apiClient *api.Client) *RunnerProjectsController {
	return &RunnerProjectsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		runnerProjectsService: service.NewRunnerProjectsService(apiClient),
	}
}

// Create assigns a runner to a project
func (c *RunnerProjectsController) Create(w http.ResponseWriter, r *http.Request) {
	var params struct {
		NamespaceID string `json:"namespace_id"`
		ProjectID   string `json:"project_id"`
		RunnerProject struct {
			RunnerID int64 `json:"runner_id"`
		} `json:"runner_project"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get project by full path
	project, err := c.runnerProjectsService.GetProject(params.NamespaceID, params.ProjectID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// Get runner by ID
	runner, err := c.runnerProjectsService.GetRunner(params.RunnerProject.RunnerID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Runner not found")
		return
	}

	// Assign runner to project
	if err := c.runnerProjectsService.AssignRunner(runner.ID, project.ID); err != nil {
		http.Redirect(w, r, "/admin/runners/"+runner.ID+"/edit", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/admin/runners/"+runner.ID+"/edit", http.StatusFound)
}

// Destroy unassigns a runner from a project
func (c *RunnerProjectsController) Destroy(w http.ResponseWriter, r *http.Request) {
	var params struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get runner project by ID
	runnerProject, err := c.runnerProjectsService.GetRunnerProject(params.ID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Runner project not found")
		return
	}

	// Unassign runner from project
	if err := c.runnerProjectsService.UnassignRunner(runnerProject.ID); err != nil {
		http.Redirect(w, r, "/admin/runners/"+runnerProject.RunnerID+"/edit", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/admin/runners/"+runnerProject.RunnerID+"/edit", http.StatusFound)
}
