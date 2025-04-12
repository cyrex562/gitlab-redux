package repository

import (
	"net/http"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// RepositorySettingsRedirect handles redirecting to repository settings
type RepositorySettingsRedirect struct {
	router Router
}

// Router defines the interface for routing operations
type Router interface {
	ProjectSettingsRepositoryPath(project *model.Project, anchor string) string
}

// NewRepositorySettingsRedirect creates a new instance of RepositorySettingsRedirect
func NewRepositorySettingsRedirect(router Router) *RepositorySettingsRedirect {
	return &RepositorySettingsRedirect{
		router: router,
	}
}

// RedirectToRepositorySettings redirects to the project's repository settings page
func (r *RepositorySettingsRedirect) RedirectToRepositorySettings(w http.ResponseWriter, project *model.Project, anchor string) {
	path := r.router.ProjectSettingsRepositoryPath(project, anchor)
	http.Redirect(w, nil, path, http.StatusFound)
}
