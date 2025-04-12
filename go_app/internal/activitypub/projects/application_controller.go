package projects

import (
	"net/http"
	"path/filepath"

	"gitlab.com/gitlab-org/gitlab/go/internal/activitypub"
	"gitlab.com/gitlab-org/gitlab/go/internal/feature"
	"gitlab.com/gitlab-org/gitlab/go/internal/models"
)

// ApplicationController is the base controller for ActivityPub project endpoints
type ApplicationController struct {
	activitypub.ApplicationController
	project *models.Project
}

// BeforeAction performs setup before handling requests
func (c *ApplicationController) BeforeAction(w http.ResponseWriter, r *http.Request) error {
	if err := c.setProject(w, r); err != nil {
		return err
	}

	if err := c.ensureProjectFeatureFlag(); err != nil {
		return err
	}

	return nil
}

// setProject finds and sets the project based on request parameters
func (c *ApplicationController) setProject(w http.ResponseWriter, r *http.Request) error {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		projectID = r.URL.Query().Get("id")
	}

	if projectID == "" {
		return nil
	}

	namespaceID := r.URL.Query().Get("namespace_id")
	path := filepath.Join(namespaceID, projectID)

	project, err := c.findRoutable(models.Project{}, path, r.URL.Path, func(p *models.Project) bool {
		return p.Public && !p.PendingDelete
	})

	if err != nil {
		return err
	}

	c.project = project
	return nil
}

// ensureProjectFeatureFlag checks if the ActivityPub project feature is enabled
func (c *ApplicationController) ensureProjectFeatureFlag() error {
	if !feature.Enabled("activity_pub_project", c.project) {
		return activitypub.ErrNotFound
	}
	return nil
}

// GetProject returns the current project
func (c *ApplicationController) GetProject() *models.Project {
	return c.project
}
