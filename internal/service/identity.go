package service

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gorm.io/gorm"
)

// IdentityService handles business logic for user identities
type IdentityService struct {
	db *gorm.DB
}

// NewIdentityService creates a new IdentityService instance
func NewIdentityService(db *gorm.DB) *IdentityService {
	return &IdentityService{
		db: db,
	}
}

// GetUserByUsername retrieves a user by their username
func (s *IdentityService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserIdentities retrieves all identities for a user
func (s *IdentityService) GetUserIdentities(ctx context.Context, userID uint) ([]model.Identity, error) {
	var identities []model.Identity
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&identities).Error
	if err != nil {
		return nil, err
	}
	return identities, nil
}

// CreateIdentity creates a new user identity
func (s *IdentityService) CreateIdentity(ctx context.Context, identity *model.Identity) (*model.Identity, error) {
	err := s.db.WithContext(ctx).Create(identity).Error
	if err != nil {
		return nil, err
	}
	return identity, nil
}

// UpdateIdentity updates an existing user identity
func (s *IdentityService) UpdateIdentity(ctx context.Context, userID uint, identityID string, identity *model.Identity) (*model.Identity, error) {
	result := s.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, identityID).
		Updates(identity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("identity not found")
	}
	return identity, nil
}

// DeleteIdentity removes a user identity
func (s *IdentityService) DeleteIdentity(ctx context.Context, userID uint, identityID string) error {
	result := s.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, identityID).
		Delete(&model.Identity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("identity not found")
	}
	return nil
}

// RepairLdapBlocked repairs the LDAP blocked status for a user
func (s *IdentityService) RepairLdapBlocked(ctx context.Context, userID uint) error {
	// TODO: Implement LDAP blocked status repair
	// This would typically involve:
	// 1. Checking user's identities
	// 2. Updating the user's LDAP blocked status based on their identities
	return nil
}
