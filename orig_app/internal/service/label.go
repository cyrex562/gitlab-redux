package service

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/repository"
)

// LabelService defines the interface for label-related operations
type LabelService interface {
	// FindDistinctLabelsByProjects finds distinct labels based on project IDs
	FindDistinctLabelsByProjects(ctx context.Context, user *model.User, projectIDs []int64) ([]*model.Label, error)
}

// labelService implements the LabelService interface
type labelService struct {
	db           *sql.DB
	labelRepo    repository.LabelRepository
	projectRepo  repository.ProjectRepository
}

// NewLabelService creates a new LabelService
func NewLabelService(db *sql.DB, labelRepo repository.LabelRepository, projectRepo repository.ProjectRepository) LabelService {
	return &labelService{
		db:           db,
		labelRepo:    labelRepo,
		projectRepo:  projectRepo,
	}
}

// FindDistinctLabelsByProjects finds distinct labels based on project IDs
func (s *labelService) FindDistinctLabelsByProjects(ctx context.Context, user *model.User, projectIDs []int64) ([]*model.Label, error) {
	// Verify user has access to all projects
	for _, projectID := range projectIDs {
		hasAccess, err := s.projectRepo.UserHasAccess(ctx, user.ID, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to check project access: %w", err)
		}
		if !hasAccess {
			return nil, fmt.Errorf("user does not have access to project %d", projectID)
		}
	}

	// Find distinct labels by project IDs
	labels, err := s.labelRepo.FindDistinctByProjects(ctx, projectIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find distinct labels: %w", err)
	}

	return labels, nil
}
