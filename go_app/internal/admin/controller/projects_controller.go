package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

// ProjectsController handles admin project management
type ProjectsController struct {
	baseController
	projectsService *service.ProjectsService
}

// NewProjectsController creates a new projects controller
func NewProjectsController(apiClient *api.Client) *ProjectsController {
	return &ProjectsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		projectsService: service.NewProjectsService(apiClient),
	}
}

// Index displays a list of projects
func (c *ProjectsController) Index(w http.ResponseWriter, r *http.Request) {
	// Set default sort if not provided
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "latest_activity_desc"
	}

	// Handle archived parameter
	archived := r.URL.Query().Get("archived")
	if r.URL.Query().Get("last_repository_check_failed") != "" && archived == "" {
		archived = "true"
	}

	// Get projects with filters
	projects, err := c.projectsService.GetProjects(r.URL.Query())
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	// Handle different response formats
	format := r.URL.Query().Get("format")
	if format == "json" {
		// TODO: Implement HTML template rendering for projects list
		response := struct {
			HTML string `json:"html"`
		}{
			HTML: "Projects list HTML template",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Default to HTML response
	// TODO: Implement HTML template rendering for projects list
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Show displays project details
func (c *ProjectsController) Show(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	projectID := r.URL.Query().Get("id")

	if namespaceID == "" || projectID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	project, err := c.projectsService.GetProject(namespaceID, projectID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// Get group members if project belongs to a group
	if project.Group != nil {
		groupMembers, err := c.projectsService.GetGroupMembers(project.Group.ID, r.URL.Query().Get("group_members_page"))
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch group members")
			return
		}
		// TODO: Add group members to response
	}

	// Get project members
	projectMembers, err := c.projectsService.GetProjectMembers(project.ID, r.URL.Query().Get("project_members_page"))
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch project members")
		return
	}

	// Get access requesters
	requesters, err := c.projectsService.GetAccessRequesters(project.ID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch access requesters")
		return
	}

	// TODO: Implement HTML template rendering for project details
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Destroy removes a project
func (c *ProjectsController) Destroy(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	projectID := r.URL.Query().Get("id")

	if namespaceID == "" || projectID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	project, err := c.projectsService.GetProject(namespaceID, projectID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	if err := c.projectsService.DestroyProject(project.ID); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	http.Redirect(w, r, "/admin/projects", http.StatusFound)
}

// Transfer moves a project to a new namespace
func (c *ProjectsController) Transfer(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	projectID := r.URL.Query().Get("id")
	newNamespaceID := r.URL.Query().Get("new_namespace_id")

	if namespaceID == "" || projectID == "" || newNamespaceID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	project, err := c.projectsService.GetProject(namespaceID, projectID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	if err := c.projectsService.TransferProject(project.ID, newNamespaceID); err != nil {
		helper.RespondWithError(w, http.StatusUnprocessableEntity, "Failed to transfer project")
		return
	}

	http.Redirect(w, r, "/admin/projects/"+projectID, http.StatusFound)
}

// Edit displays the form for editing a project
func (c *ProjectsController) Edit(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	projectID := r.URL.Query().Get("id")

	if namespaceID == "" || projectID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	project, err := c.projectsService.GetProject(namespaceID, projectID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// TODO: Implement HTML template rendering for edit project form
	helper.RespondWithError(w, http.StatusNotImplemented, "HTML response not yet implemented")
}

// Update modifies a project
func (c *ProjectsController) Update(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	projectID := r.URL.Query().Get("id")

	if namespaceID == "" || projectID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	var params struct {
		Description              string `json:"description"`
		Name                    string `json:"name"`
		RunnerRegistrationEnabled bool   `json:"runner_registration_enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := c.projectsService.GetProject(namespaceID, projectID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	if err := c.projectsService.UpdateProject(project.ID, params); err != nil {
		helper.RespondWithError(w, http.StatusUnprocessableEntity, "Failed to update project")
		return
	}

	// Reset runner registration token if disabled
	if !params.RunnerRegistrationEnabled {
		if err := c.projectsService.ResetRunnerRegistrationToken(project.ID); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to reset runner registration token")
			return
		}
	}

	http.Redirect(w, r, "/admin/projects/"+projectID, http.StatusFound)
}

// RepositoryCheck triggers a repository check
func (c *ProjectsController) RepositoryCheck(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	projectID := r.URL.Query().Get("id")

	if namespaceID == "" || projectID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	project, err := c.projectsService.GetProject(namespaceID, projectID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	if err := c.projectsService.TriggerRepositoryCheck(project.ID); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to trigger repository check")
		return
	}

	http.Redirect(w, r, "/admin/projects/"+projectID, http.StatusFound)
}
