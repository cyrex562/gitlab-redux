package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

const (
	tagsLimit = 20
)

// RunnersController handles admin CI runner management
type RunnersController struct {
	baseController
	runnersService *service.RunnersService
}

// NewRunnersController creates a new runners controller
func NewRunnersController(apiClient *api.Client) *RunnersController {
	return &RunnersController{
		baseController: baseController{
			apiClient: apiClient,
		},
		runnersService: service.NewRunnersService(apiClient),
	}
}

// Index displays a list of runners
func (c *RunnersController) Index(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement HTML template rendering for runners list
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Show displays runner details
func (c *RunnersController) Show(w http.ResponseWriter, r *http.Request) {
	runnerID := r.URL.Query().Get("id")
	if runnerID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing runner ID")
		return
	}

	runner, err := c.runnersService.GetRunner(runnerID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Runner not found")
		return
	}

	// TODO: Implement HTML template rendering for runner details
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Edit displays the form for editing a runner
func (c *RunnersController) Edit(w http.ResponseWriter, r *http.Request) {
	runnerID := r.URL.Query().Get("id")
	if runnerID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing runner ID")
		return
	}

	runner, err := c.runnersService.GetRunner(runnerID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Runner not found")
		return
	}

	// Get available projects for assignment
	search := r.URL.Query().Get("search")
	page := r.URL.Query().Get("page")
	projects, err := c.runnersService.GetAvailableProjects(runner.ID, search, page)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch available projects")
		return
	}

	// TODO: Implement HTML template rendering for edit runner form
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// New displays the form for creating a new runner
func (c *RunnersController) New(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement HTML template rendering for new runner form
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Register handles runner registration
func (c *RunnersController) Register(w http.ResponseWriter, r *http.Request) {
	runnerID := r.URL.Query().Get("id")
	if runnerID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing runner ID")
		return
	}

	runner, err := c.runnersService.GetRunner(runnerID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Runner not found")
		return
	}

	if !runner.RegistrationAvailable {
		helper.RespondWithError(w, http.StatusNotFound, "Runner registration is not available")
		return
	}

	// TODO: Implement HTML template rendering for runner registration
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Update modifies a runner
func (c *RunnersController) Update(w http.ResponseWriter, r *http.Request) {
	runnerID := r.URL.Query().Get("id")
	if runnerID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing runner ID")
		return
	}

	var params struct {
		Description string `json:"description"`
		TagList     []string `json:"tag_list"`
		// Add other editable fields as needed
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	runner, err := c.runnersService.GetRunner(runnerID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Runner not found")
		return
	}

	if err := c.runnersService.UpdateRunner(runner.ID, params); err != nil {
		// Get available projects for assignment
		search := r.URL.Query().Get("search")
		page := r.URL.Query().Get("page")
		projects, err := c.runnersService.GetAvailableProjects(runner.ID, search, page)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch available projects")
			return
		}

		// TODO: Implement HTML template rendering for show runner view
		helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
		return
	}

	http.Redirect(w, r, "/admin/runners/"+runnerID+"/edit", http.StatusFound)
}

// TagList retrieves a list of runner tags
func (c *RunnersController) TagList(w http.ResponseWriter, r *http.Request) {
	tags, err := c.runnersService.GetTags(r.URL.Query())
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch tags")
		return
	}

	// Limit the number of tags
	if len(tags) > tagsLimit {
		tags = tags[:tagsLimit]
	}

	json.NewEncoder(w).Encode(tags)
}

// RunnerSetupScripts provides runner setup scripts
func (c *RunnersController) RunnerSetupScripts(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement runner setup scripts
	helper.RespondWithError(w, http.StatusNotImplemented, "Runner setup scripts not yet implemented")
}
