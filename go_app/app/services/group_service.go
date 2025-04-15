package services

import (
	"errors"

	"gorm.io/gorm"
)

// Group represents a group in the system
type Group struct {
	ID             uint   `gorm:"primaryKey"`
	OrganizationID uint   `gorm:"not null"`
	Name           string `gorm:"not null"`
	Path           string `gorm:"not null"`
	FullPath       string `gorm:"not null"`
}

// GroupParams represents parameters for creating a group
type GroupParams struct {
	Name           string `json:"name" binding:"required"`
	Path           string `json:"path" binding:"required"`
	OrganizationID uint   `json:"organization_id"`
}

// GroupService handles group-related operations
type GroupService struct {
	db *gorm.DB
}

// NewGroupService creates a new GroupService
func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{
		db: db,
	}
}

// FindByID finds a group by ID in an organization
func (s *GroupService) FindByID(organizationID uint, groupID string) (*Group, error) {
	var group Group
	if err := s.db.Where("organization_id = ? AND id = ?", organizationID, groupID).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(userID uint, params GroupParams) (*Group, error) {
	// Create the group
	group := &Group{
		OrganizationID: params.OrganizationID,
		Name:           params.Name,
		Path:           params.Path,
		FullPath:       params.Path, // This would be more complex in a real implementation
	}

	if err := s.db.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

// DestroyGroup deletes a group
func (s *GroupService) DestroyGroup(userID uint, groupID uint) error {
	// Find the group
	var group Group
	if err := s.db.First(&group, groupID).Error; err != nil {
		return err
	}

	// Delete the group
	if err := s.db.Delete(&group).Error; err != nil {
		return err
	}

	return nil
}

// CanViewEditPage checks if a user can view the edit page for a group
func (s *GroupService) CanViewEditPage(userID uint, groupID uint) bool {
	// TODO: Implement proper permission check
	return true
}

// CanRemoveGroup checks if a user can remove a group
func (s *GroupService) CanRemoveGroup(userID uint, groupID uint) bool {
	// TODO: Implement proper permission check
	return true
} 