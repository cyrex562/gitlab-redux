package service

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gorm.io/gorm"
)

// GroupService handles business logic for groups
type GroupService struct {
	db *gorm.DB
}

// NewGroupService creates a new GroupService instance
func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{
		db: db,
	}
}

// GetGroups retrieves a paginated list of groups with optional filtering
func (s *GroupService) GetGroups(ctx context.Context, page int, sort string, name string) ([]model.Group, error) {
	var groups []model.Group
	query := s.db.WithContext(ctx)

	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	// Apply sorting
	switch sort {
	case "name":
		query = query.Order("name ASC")
	case "created_at":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("name ASC")
	}

	// Apply pagination
	limit := 20 // Default page size
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	err := query.Find(&groups).Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupWithDetails retrieves a group with its members, projects, and statistics
func (s *GroupService) GetGroupWithDetails(ctx context.Context, groupID string, membersPage, projectsPage int) (*model.GroupDetails, error) {
	var group model.Group
	err := s.db.WithContext(ctx).First(&group, groupID).Error
	if err != nil {
		return nil, err
	}

	// Get members with pagination
	var members []model.GroupMember
	membersLimit := 20
	membersOffset := (membersPage - 1) * membersLimit
	err = s.db.WithContext(ctx).
		Where("group_id = ?", group.ID).
		Order("access_level DESC").
		Offset(membersOffset).
		Limit(membersLimit).
		Find(&members).Error
	if err != nil {
		return nil, err
	}

	// Get projects with pagination
	var projects []model.Project
	projectsLimit := 20
	projectsOffset := (projectsPage - 1) * projectsLimit
	err = s.db.WithContext(ctx).
		Where("group_id = ?", group.ID).
		Offset(projectsOffset).
		Limit(projectsLimit).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return &model.GroupDetails{
		Group:     group,
		Members:   members,
		Projects:  projects,
	}, nil
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(ctx context.Context, group *model.Group) (*model.Group, error) {
	err := s.db.WithContext(ctx).Create(group).Error
	if err != nil {
		return nil, err
	}
	return group, nil
}

// UpdateGroup updates an existing group
func (s *GroupService) UpdateGroup(ctx context.Context, groupID string, group *model.Group) (*model.Group, error) {
	result := s.db.WithContext(ctx).Where("id = ?", groupID).Updates(group)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("group not found")
	}
	return group, nil
}

// DeleteGroup deletes a group
func (s *GroupService) DeleteGroup(ctx context.Context, groupID string) error {
	result := s.db.WithContext(ctx).Delete(&model.Group{}, groupID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("group not found")
	}
	return nil
}

// IsDependencyProxyAvailable checks if dependency proxy is available for a group
func (g *GroupService) IsDependencyProxyAvailable(ctx context.Context, group *model.Group) (bool, error) {
	// TODO: Implement dependency proxy availability check
	// This should:
	// 1. Check if the group has dependency proxy feature enabled
	// 2. Return the result
	return false, nil
}

// GetDependencyProxyPolicySubject returns the dependency proxy policy subject for a group
func (g *GroupService) GetDependencyProxyPolicySubject(ctx context.Context, group *model.Group) (*model.DependencyProxyPolicySubject, error) {
	// TODO: Implement policy subject retrieval
	// This should:
	// 1. Get the dependency proxy policy subject for the group
	// 2. Return the result
	return nil, nil
}

// Config holds configuration for the GroupService
type Config struct {
	// Add configuration options as needed
}
