package service

import (
	"context"

	"github.com/jmadden/gitlab-redux/internal/model"
)

// ProjectService defines the interface for project-related operations
type ProjectService interface {
	// FindProjects finds projects based on the given parameters and current user
	FindProjects(ctx context.Context, params map[string]interface{}, currentUser *model.User) ([]*model.Project, error)

	// GetProjectByID gets a project by its ID
	GetProjectByID(ctx context.Context, id int64, currentUser *model.User) (*model.Project, error)

	// GetStarredProjects gets the starred projects for the current user
	GetStarredProjects(ctx context.Context, currentUser *model.User) ([]*model.Project, error)

	// StarProject stars a project for the current user
	StarProject(ctx context.Context, projectID int64, currentUser *model.User) error

	// UnstarProject unstars a project for the current user
	UnstarProject(ctx context.Context, projectID int64, currentUser *model.User) error
}
