package controller

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab/workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab/workhorse/internal/service"
)

const (
	buildsPerPage = 30
)

// JobsController handles admin CI jobs management
type JobsController struct {
	baseController
	jobsService *service.JobsService
}

// NewJobsController creates a new jobs controller
func NewJobsController(apiClient *api.Client) *JobsController {
	return &JobsController{
		baseController: baseController{
			apiClient: apiClient,
		},
		jobsService: service.NewJobsService(apiClient),
	}
}

// Index handles the jobs listing page
func (c *JobsController) Index(w http.ResponseWriter, r *http.Request) {
	// Set feature flag for admin jobs filter runner type
	helper.SetFeatureFlag(w, "admin_jobs_filter_runner_type", "ops")

	// TODO: Implement jobs listing view
	// This would typically:
	// 1. Fetch jobs with pagination
	// 2. Apply any filters
	// 3. Render the jobs view template
	helper.RespondWithError(w, http.StatusNotImplemented, "Jobs listing not yet implemented")
}

// CancelAll cancels all running or pending jobs
func (c *JobsController) CancelAll(w http.ResponseWriter, r *http.Request) {
	if err := c.jobsService.CancelAllRunningOrPendingJobs(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to cancel jobs")
		return
	}

	// Redirect to jobs index page
	http.Redirect(w, r, "/admin/jobs", http.StatusSeeOther)
}
