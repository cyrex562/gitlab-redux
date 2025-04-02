package services

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
)

// UpdateApplicationSettingsService handles updating application settings
type UpdateApplicationSettingsService struct {
	settings *models.ApplicationSettings
	user     *models.User
	params   *models.ApplicationSettingParams
}

// NewUpdateApplicationSettingsService creates a new instance of UpdateApplicationSettingsService
func NewUpdateApplicationSettingsService(settings *models.ApplicationSettings, user *models.User, params *models.ApplicationSettingParams) *UpdateApplicationSettingsService {
	return &UpdateApplicationSettingsService{
		settings: settings,
		user:     user,
		params:   params,
	}
}

// Execute performs the update operation
func (s *UpdateApplicationSettingsService) Execute() (bool, error) {
	if !s.user.IsAdmin() {
		return false, ErrUnauthorized
	}

	if err := s.settings.Update(s.params); err != nil {
		return false, err
	}

	return true, nil
}
