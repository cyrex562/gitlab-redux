package services

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
)

// ResetRegistrationTokenService handles resetting registration tokens
type ResetRegistrationTokenService struct {
	settings *models.ApplicationSettings
	user     *models.User
}

// NewResetRegistrationTokenService creates a new instance of ResetRegistrationTokenService
func NewResetRegistrationTokenService(settings *models.ApplicationSettings, user *models.User) *ResetRegistrationTokenService {
	return &ResetRegistrationTokenService{
		settings: settings,
		user:     user,
	}
}

// Execute performs the reset operation
func (s *ResetRegistrationTokenService) Execute() error {
	if !s.user.IsAdmin() {
		return ErrUnauthorized
	}

	return s.settings.ResetRegistrationToken()
}
