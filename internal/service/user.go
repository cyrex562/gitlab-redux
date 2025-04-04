package service

import (
	"context"
	"errors"
	"time"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gorm.io/gorm"
)

// UserService handles user operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new instance of UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// List returns a paginated list of users
func (s *UserService) List(ctx context.Context, filter, searchQuery string, page int, sort string) ([]*model.User, error) {
	// TODO: Implement pagination, filtering, and sorting
	return nil, nil
}

// Get retrieves a user by ID
func (s *UserService) Get(ctx context.Context, id int64) (*model.User, error) {
	// TODO: Implement user retrieval
	return nil, nil
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, currentUser *model.User, params *model.UserParams) (*model.User, error) {
	// TODO: Implement user creation
	return nil, nil
}

// Update updates an existing user
func (s *UserService) Update(ctx context.Context, currentUser *model.User, id int64, params *model.UserParams) (*model.User, error) {
	// TODO: Implement user update
	return nil, nil
}

// Destroy removes a user
func (s *UserService) Destroy(ctx context.Context, currentUser *model.User, id int64, hardDelete bool) error {
	// TODO: Implement user deletion
	return nil
}

// Impersonate allows an admin to impersonate another user
func (s *UserService) Impersonate(ctx context.Context, currentUser *model.User, id int64) error {
	// TODO: Implement user impersonation
	return nil
}

// GetProjects returns the projects associated with a user
func (s *UserService) GetProjects(ctx context.Context, id int64) ([]*model.Project, error) {
	// TODO: Implement project retrieval
	return nil, nil
}

// GetKeys returns the SSH keys associated with a user
func (s *UserService) GetKeys(ctx context.Context, id int64) ([]*model.Key, error) {
	// TODO: Implement key retrieval
	return nil, nil
}

// GetLastAdminUser retrieves the last admin user
func (s *UserService) GetLastAdminUser(ctx context.Context) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).
		Where("admin = ?", true).
		Order("created_at DESC").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, userID uint, user *model.User) (*model.User, error) {
	// Start a transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get existing user
	var existingUser model.User
	err := tx.Where("id = ?", userID).First(&existingUser).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update user fields
	existingUser.Email = user.Email
	existingUser.Password = user.Password
	existingUser.UpdatedAt = time.Now()

	// Skip email confirmation if requested
	if user.SkipReconfirmation {
		existingUser.ConfirmedAt = time.Now()
	}

	// Save changes
	err = tx.Save(&existingUser).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &existingUser, nil
}

// CleanupNonPrimaryEmails removes non-primary emails for a user
func (s *UserService) CleanupNonPrimaryEmails(ctx context.Context, userID uint) error {
	// Start a transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete non-primary emails
	err := tx.Where("user_id = ? AND primary = ?", userID, false).Delete(&model.Email{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

// Custom errors for user operations
var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)
