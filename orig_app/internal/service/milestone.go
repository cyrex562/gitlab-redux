package service

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/repository"
)

// MilestoneService defines the interface for milestone-related operations
type MilestoneService interface {
	// FindMilestones finds milestones based on search parameters
	FindMilestones(ctx context.Context, params *MilestoneSearchParams) ([]*model.Milestone, error)

	// GetMilestoneStateCount gets the count of milestones by state
	GetMilestoneStateCount(ctx context.Context, projectIDs []int64, groupIDs []int64) (*model.MilestoneStateCount, error)

	// GetMilestoneJSON gets milestones in JSON format
	GetMilestoneJSON(ctx context.Context, params *MilestoneSearchParams) ([]*model.MilestoneJSON, error)
}

// MilestoneSearchParams represents the parameters for searching milestones
type MilestoneSearchParams struct {
	State       *model.MilestoneState
	SearchTitle string
	GroupIDs    []int64
	ProjectIDs  []int64
	Page        int
	PerPage     int
}

// milestoneService implements the MilestoneService interface
type milestoneService struct {
	db            *sql.DB
	milestoneRepo repository.MilestoneRepository
	groupRepo     repository.GroupRepository
}

// NewMilestoneService creates a new MilestoneService
func NewMilestoneService(db *sql.DB, milestoneRepo repository.MilestoneRepository, groupRepo repository.GroupRepository) MilestoneService {
	return &milestoneService{
		db:            db,
		milestoneRepo: milestoneRepo,
		groupRepo:     groupRepo,
	}
}

// FindMilestones finds milestones based on search parameters
func (s *milestoneService) FindMilestones(ctx context.Context, params *MilestoneSearchParams) ([]*model.Milestone, error) {
	// Find milestones by search parameters
	milestones, err := s.milestoneRepo.FindByParams(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find milestones: %w", err)
	}

	return milestones, nil
}

// GetMilestoneStateCount gets the count of milestones by state
func (s *milestoneService) GetMilestoneStateCount(ctx context.Context, projectIDs []int64, groupIDs []int64) (*model.MilestoneStateCount, error) {
	// Get milestone state count
	count, err := s.milestoneRepo.GetStateCount(ctx, projectIDs, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone state count: %w", err)
	}

	return count, nil
}

// GetMilestoneJSON gets milestones in JSON format
func (s *milestoneService) GetMilestoneJSON(ctx context.Context, params *MilestoneSearchParams) ([]*model.MilestoneJSON, error) {
	// Find milestones
	milestones, err := s.FindMilestones(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert to JSON format
	result := make([]*model.MilestoneJSON, 0, len(milestones))
	for _, milestone := range milestones {
		json := &model.MilestoneJSON{
			ID:      milestone.ID,
			Title:   milestone.Title,
			DueDate: milestone.DueDate,
			Name:    milestone.Title, // In the Ruby code, name is the same as title
		}
		result = append(result, json)
	}

	return result, nil
}
