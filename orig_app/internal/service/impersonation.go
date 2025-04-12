package service

import (
	"context"
	"errors"
	"time"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gorm.io/gorm"
)

// TokenFinderOptions represents options for finding tokens
type TokenFinderOptions struct {
	UserID         uint
	OrganizationID uint
	State         string
	Sort          string
}

// ImpersonationService handles business logic for impersonation tokens
type ImpersonationService struct {
	db *gorm.DB
}

// NewImpersonationService creates a new ImpersonationService instance
func NewImpersonationService(db *gorm.DB) *ImpersonationService {
	return &ImpersonationService{
		db: db,
	}
}

// GetUserByUsername retrieves a user by their username
func (s *ImpersonationService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAvailableScopes retrieves available scopes for the current user
func (s *ImpersonationService) GetAvailableScopes(ctx context.Context) ([]string, error) {
	// Get all available scopes
	scopes := []string{
		"read_user",
		"write_user",
		"read_repository",
		"write_repository",
		"read_api",
		"write_api",
		"read_registry",
		"write_registry",
	}

	// TODO: Implement virtual registry filtering
	// This should filter scopes based on:
	// 1. User permissions
	// 2. Virtual registry settings
	// 3. Feature flags

	return scopes, nil
}

// GetActiveTokens retrieves active impersonation tokens for a user
func (s *ImpersonationService) GetActiveTokens(ctx context.Context, userID uint) ([]model.ImpersonationToken, error) {
	options := TokenFinderOptions{
		UserID: userID,
		State:  "active",
		Sort:   "expires_asc",
	}
	return s.findTokens(ctx, options)
}

// findTokens finds tokens based on the provided options
func (s *ImpersonationService) findTokens(ctx context.Context, options TokenFinderOptions) ([]model.ImpersonationToken, error) {
	query := s.db.WithContext(ctx)

	// Apply filters
	if options.UserID > 0 {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.OrganizationID > 0 {
		query = query.Where("organization_id = ?", options.OrganizationID)
	}

	// Apply state filter
	switch options.State {
	case "active":
		query = query.Where("revoked = ? AND expires_at > ?", false, time.Now())
	case "revoked":
		query = query.Where("revoked = ?", true)
	case "expired":
		query = query.Where("expires_at <= ?", time.Now())
	}

	// Apply sorting
	switch options.Sort {
	case "expires_asc":
		query = query.Order("expires_at ASC")
	case "expires_desc":
		query = query.Order("expires_at DESC")
	case "created_asc":
		query = query.Order("created_at ASC")
	case "created_desc":
		query = query.Order("created_at DESC")
	}

	var tokens []model.ImpersonationToken
	err := query.Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// CreateToken creates a new impersonation token
func (s *ImpersonationService) CreateToken(ctx context.Context, userID uint, token *model.ImpersonationToken) (*model.ImpersonationToken, error) {
	// Set user ID and generate token
	token.UserID = userID
	token.Token = generateSecureToken()
	token.CreatedAt = time.Now()
	token.LastUsedAt = time.Now()

	// Validate token data
	if err := s.validateToken(token); err != nil {
		return nil, err
	}

	err := s.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}

// RevokeToken revokes an impersonation token
func (s *ImpersonationService) RevokeToken(ctx context.Context, userID uint, tokenID string) (*model.ImpersonationToken, error) {
	var token model.ImpersonationToken
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, tokenID).
		First(&token).Error
	if err != nil {
		return nil, err
	}

	token.Revoked = true
	token.RevokedAt = time.Now()

	err = s.db.WithContext(ctx).Save(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// RotateToken rotates an impersonation token
func (s *ImpersonationService) RotateToken(ctx context.Context, userID uint, tokenID string) (*model.ImpersonationToken, error) {
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

	// Get the original token
	var oldToken model.ImpersonationToken
	err := tx.Where("user_id = ? AND id = ?", userID, tokenID).First(&oldToken).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create new token with same properties
	newToken := model.ImpersonationToken{
		UserID:          oldToken.UserID,
		Name:            oldToken.Name,
		Description:     oldToken.Description,
		Scopes:          oldToken.Scopes,
		ExpiresAt:       oldToken.ExpiresAt,
		OrganizationID:  oldToken.OrganizationID,
		Token:           generateSecureToken(),
		CreatedAt:       time.Now(),
		LastUsedAt:      time.Now(),
	}

	// Validate new token
	if err := s.validateToken(&newToken); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Create(&newToken).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Revoke the old token
	oldToken.Revoked = true
	oldToken.RevokedAt = time.Now()
	err = tx.Save(&oldToken).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &newToken, nil
}

// StopImpersonation stops the current impersonation session and returns the original user
func (s *ImpersonationService) StopImpersonation(ctx context.Context) (*model.User, error) {
	// Get current user from context
	currentUser, exists := ctx.Value("current_user").(*model.User)
	if !exists {
		return nil, ErrNoCurrentUser
	}

	// Get impersonated user from context
	impersonatedUser, exists := ctx.Value("impersonated_user").(*model.User)
	if !exists {
		return nil, ErrNoImpersonatedUser
	}

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

	// Update impersonation session status
	err := tx.Exec("UPDATE impersonation_sessions SET ended_at = ? WHERE impersonator_id = ? AND impersonated_id = ? AND ended_at IS NULL",
		time.Now(), currentUser.ID, impersonatedUser.ID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return impersonatedUser, nil
}

// validateToken validates token data before creation
func (s *ImpersonationService) validateToken(token *model.ImpersonationToken) error {
	if token.Name == "" {
		return ErrInvalidTokenName
	}
	if token.Description == "" {
		return ErrInvalidTokenDescription
	}
	if len(token.Scopes) == 0 {
		return ErrInvalidTokenScopes
	}
	if token.ExpiresAt.Before(time.Now()) {
		return ErrInvalidTokenExpiration
	}
	return nil
}

// generateSecureToken generates a secure random token
func generateSecureToken() string {
	// TODO: Implement proper secure token generation
	// This should use a cryptographically secure random number generator
	return "generated_token_here"
}

// Custom errors for token validation and impersonation
var (
	ErrInvalidTokenName        = errors.New("token name cannot be empty")
	ErrInvalidTokenDescription = errors.New("token description cannot be empty")
	ErrInvalidTokenScopes      = errors.New("token must have at least one scope")
	ErrInvalidTokenExpiration  = errors.New("token expiration must be in the future")
	ErrNoCurrentUser          = errors.New("no current user found in context")
	ErrNoImpersonatedUser     = errors.New("no impersonated user found in context")
)
