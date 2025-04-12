package repository

import (
	"context"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// LabelRepository defines the interface for label-related database operations
type LabelRepository interface {
	// FindDistinctByProjects finds distinct labels based on project IDs
	FindDistinctByProjects(ctx context.Context, projectIDs []int64) ([]*model.Label, error)
}

// ProjectRepository defines the interface for project-related database operations
type ProjectRepository interface {
	// UserHasAccess checks if a user has access to a project
	UserHasAccess(ctx context.Context, userID, projectID int64) (bool, error)
}
