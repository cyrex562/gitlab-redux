package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jmadden/gitlab-redux/internal/model"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

// EmailVerificationService handles business logic for email verification
type EmailVerificationService struct {
	db *sql.DB
	emailService *EmailService
	featureService *FeatureService
}

// NewEmailVerificationService creates a new email verification service
func NewEmailVerificationService(db *sql.DB, emailService *EmailService, featureService *FeatureService) *EmailVerificationService {
	return &EmailVerificationService{
		db: db,
		emailService: emailService,
		featureService: featureService,
	}
}

// GenerateToken generates a verification token for a user
func (s *EmailVerificationService) GenerateToken(ctx context.Context, user *model.User) (string, error) {
	// Generate a random token
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	// Encode the token
	encodedToken := base64.URLEncoding.EncodeToString(token)

	// Hash the token
	hashedToken, err := s.hashToken(encodedToken)
	if err != nil {
		return "", err
	}

	// Update the user's unlock token
	query := `
		UPDATE users
		SET unlock_token = $1, locked_at = $2
		WHERE id = $3
	`
	_, err = s.db.ExecContext(ctx, query, hashedToken, time.Now(), user.ID)
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}

// ValidateToken validates a verification token for a user
func (s *EmailVerificationService) ValidateToken(ctx context.Context, user *model.User, token string) (*model.VerificationResult, error) {
	// Hash the token
	hashedToken, err := s.hashToken(token)
	if err != nil {
		return nil, err
	}

	// Check if the token matches
	if user.UnlockToken != hashedToken {
		return &model.VerificationResult{
			Status:  "error",
			Reason:  "invalid_token",
			Message: "Invalid verification token",
		}, nil
	}

	// Check if the token is expired
	if s.IsTokenExpired(user) {
		return &model.VerificationResult{
			Status:  "error",
			Reason:  "token_expired",
			Message: "Verification token has expired",
		}, nil
	}

	return &model.VerificationResult{
		Status: "success",
	}, nil
}

// UpdateEmail updates a user's email address
func (s *EmailVerificationService) UpdateEmail(ctx context.Context, user *model.User, email string) (*model.VerificationResult, error) {
	// Validate the email
	if !s.emailService.IsValidEmail(email) {
		return &model.VerificationResult{
			Status:  "error",
			Reason:  "invalid_email",
			Message: "Invalid email address",
		}, nil
	}

	// Check if the email is already in use
	exists, err := s.emailService.IsEmailInUse(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return &model.VerificationResult{
			Status:  "error",
			Reason:  "email_in_use",
			Message: "Email address is already in use",
		}, nil
	}

	// Update the user's email
	query := `
		UPDATE users
		SET email = $1, updated_at = $2
		WHERE id = $3
	`
	_, err = s.db.ExecContext(ctx, query, email, time.Now(), user.ID)
	if err != nil {
		return nil, err
	}

	return &model.VerificationResult{
		Status: "success",
	}, nil
}

// GetSecondaryEmail gets a user's secondary email
func (s *EmailVerificationService) GetSecondaryEmail(user *model.User, email string) string {
	// Check if the email is a secondary email
	for _, secondaryEmail := range user.SecondaryEmails {
		if secondaryEmail.Email == email && secondaryEmail.Confirmed {
			return email
		}
	}
	return ""
}

// GetVerificationEmail gets the email to send verification instructions to
func (s *EmailVerificationService) GetVerificationEmail(user *model.User) string {
	return user.Email
}

// IsTokenExpired checks if a user's unlock token is expired
func (s *EmailVerificationService) IsTokenExpired(user *model.User) bool {
	if user.LockedAt.IsZero() {
		return false
	}

	// Token is valid for 30 minutes
	return user.LockedAt.Add(30 * time.Minute).Before(time.Now())
}

// IsUnconfirmedVerificationEmail checks if a user's verification email is unconfirmed
func (s *EmailVerificationService) IsUnconfirmedVerificationEmail(user *model.User) bool {
	return !user.Confirmed
}

// IsEmailVerificationRequired checks if email verification is required for a user
func (s *EmailVerificationService) IsEmailVerificationRequired(user *model.User) bool {
	return s.featureService.IsEnabled("require_email_verification", user)
}

// IsSkipEmailVerificationEnabled checks if skipping email verification is enabled for a user
func (s *EmailVerificationService) IsSkipEmailVerificationEnabled(user *model.User) bool {
	return s.featureService.IsEnabled("skip_require_email_verification", user)
}

// hashToken hashes a token using bcrypt
func (s *EmailVerificationService) hashToken(token string) (string, error) {
	// TODO: Implement bcrypt hashing
	return token, nil
}
