package repository

import (
	"context"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// MilestoneRepository defines the interface for milestone-related database operations
type MilestoneRepository interface {
	// FindByParams finds milestones based on search parameters
	FindByParams(ctx context.Context, params *service.MilestoneSearchParams) ([]*model.Milestone, error)

	// GetStateCount gets the count of milestones by state
	GetStateCount(ctx context.Context, projectIDs []int64, groupIDs []int64) (*model.MilestoneStateCount, error)
}

// GroupRepository defines the interface for group-related database operations
type GroupRepository interface {
	// FindByUser finds groups for a user
	FindByUser(ctx context.Context, userID int64, allAvailable bool) ([]*model.Group, error)
}
