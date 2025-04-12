package services

import (
	"errors"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gorm.io/gorm"
)

// ApplicationService handles business logic for OAuth applications
type ApplicationService struct {
	db *gorm.DB
}

// NewApplicationService creates a new instance of ApplicationService
func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{
		db: db,
	}
}

// FindAll returns all OAuth applications with pagination
func (s *ApplicationService) FindAll(cursor string) ([]models.Application, int64, error) {
	var applications []models.Application
	var totalCount int64

	query := s.db.Model(&models.Application{})
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if cursor != "" {
		query = query.Where("id > ?", cursor)
	}

	if err := query.Order("id ASC").Limit(20).Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, totalCount, nil
}

// FindByID finds an application by its ID
func (s *ApplicationService) FindByID(id uint64) (*models.Application, error) {
	var application models.Application
	if err := s.db.First(&application, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		return nil, err
	}
	return &application, nil
}

// GetAvailableScopes returns the list of available OAuth scopes
func (s *ApplicationService) GetAvailableScopes() []string {
	// TODO: Implement scope filtering based on configuration
	return []string{
		"api",
		"read_user",
		"read_repository",
		"write_repository",
		"read_registry",
		"write_registry",
	}
}

// Create creates a new OAuth application
func (s *ApplicationService) Create(params models.ApplicationParams, r *http.Request) (*models.Application, error) {
	application := &models.Application{
		Name:         params.Name,
		RedirectURI:  params.RedirectURI,
		Scopes:       strings.Join(params.Scopes, " "),
		Confidential: params.Confidential,
		Trusted:      params.Trusted,
	}

	if err := s.db.Create(application).Error; err != nil {
		return nil, err
	}

	return application, nil
}

// Update updates an existing OAuth application
func (s *ApplicationService) Update(id uint64, params models.ApplicationParams) (*models.Application, error) {
	application, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	application.Name = params.Name
	application.RedirectURI = params.RedirectURI
	application.Scopes = strings.Join(params.Scopes, " ")
	application.Confidential = params.Confidential
	application.Trusted = params.Trusted

	if err := s.db.Save(application).Error; err != nil {
		return nil, err
	}

	return application, nil
}

// RenewSecret generates a new secret for an OAuth application
func (s *ApplicationService) RenewSecret(id uint64) (string, error) {
	application, err := s.FindByID(id)
	if err != nil {
		return "", err
	}

	// Generate a new secret (implementation depends on your requirements)
	newSecret := generateRandomString(32)
	application.Secret = newSecret

	if err := s.db.Save(application).Error; err != nil {
		return "", err
	}

	return newSecret, nil
}

// Delete removes an OAuth application
func (s *ApplicationService) Delete(id uint64) error {
	return s.db.Delete(&models.Application{}, id).Error
}

// ResetWebIdeOAuthApplicationSettings resets the Web IDE OAuth application settings
func (s *ApplicationService) ResetWebIdeOAuthApplicationSettings() bool {
	// TODO: Implement Web IDE OAuth application settings reset
	return true
}

// Helper function to generate random strings
func generateRandomString(length int) string {
	// TODO: Implement secure random string generation
	return "random-secret-string"
}
