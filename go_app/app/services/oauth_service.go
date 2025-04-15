package services

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"net/http"
	"time"

	"gitlab.com/gitlab-org/gitlab-redux/app/models"
	"gorm.io/gorm"
)

var (
	ErrInvalidUserCode = errors.New("invalid user code")
	ErrExpiredCode     = errors.New("code has expired")
)

// OAuthService handles OAuth-related operations
type OAuthService struct {
	db            *gorm.DB
	deviceCodes   map[string]DeviceCodeInfo
}

type DeviceCodeInfo struct {
	UserCode   string
	ExpiresAt  time.Time
	IsVerified bool
}

// NewOAuthService creates a new OAuthService
func NewOAuthService(db *gorm.DB) *OAuthService {
	return &OAuthService{
		db:            db,
		deviceCodes:   make(map[string]DeviceCodeInfo),
	}
}

// GetUserApplications retrieves all OAuth applications for a user
func (s *OAuthService) GetUserApplications(userID uint) ([]*models.OAuthApplication, error) {
	// TODO: Implement database query to get user's applications
	return []*models.OAuthApplication{}, nil
}

// GetUserAuthorizedTokens retrieves all authorized tokens for a user
func (s *OAuthService) GetUserAuthorizedTokens(userID uint) ([]*models.OAuthToken, error) {
	// TODO: Implement database query to get user's authorized tokens
	return []*models.OAuthToken{}, nil
}

// GetUserApplication retrieves a specific OAuth application for a user
func (s *OAuthService) GetUserApplication(userID, appID uint) (*models.OAuthApplication, error) {
	// TODO: Implement database query to get specific application
	return nil, errors.New("not implemented")
}

// CreateApplication creates a new OAuth application
func (s *OAuthService) CreateApplication(app *models.OAuthApplication) (*models.OAuthApplication, error) {
	// TODO: Implement database insert
	return app, nil
}

// RenewApplicationSecret generates a new secret for an OAuth application
func (s *OAuthService) RenewApplicationSecret(app *models.OAuthApplication) (string, error) {
	// TODO: Implement secret generation and database update
	return "new_secret", nil
}

// GetPreAuthorization retrieves pre-authorization data from the request
func (s *OAuthService) GetPreAuthorization(r *http.Request) (*models.PreAuthorization, error) {
	// TODO: Implement pre-authorization retrieval
	return &models.PreAuthorization{
		Authorizable: true,
		Client: &models.OAuthClient{
			Application: &models.OAuthApplication{},
		},
	}, nil
}

// Authorize authorizes an OAuth application for a user
func (s *OAuthService) Authorize(preAuth *models.PreAuthorization, user *models.User) (*models.Authorization, error) {
	// TODO: Implement authorization
	return &models.Authorization{
		RedirectURI: preAuth.RedirectURI,
	}, nil
}

// DowngradeScopes downgrades scopes for login
func (s *OAuthService) DowngradeScopes(preAuth *models.PreAuthorization, authType string) error {
	// TODO: Implement scope downgrading
	return nil
}

// EnsureReadUserScope ensures the application has read_user scope
func (s *OAuthService) EnsureReadUserScope(app *models.OAuthApplication) error {
	// TODO: Implement read_user scope check
	return nil
}

// HasReadUserScope checks if the application has read_user scope
func (s *OAuthService) HasReadUserScope(app *models.OAuthApplication) bool {
	// TODO: Implement read_user scope check
	return false
}

// HasAPIScope checks if the application has API scope
func (s *OAuthService) HasAPIScope(app *models.OAuthApplication) bool {
	// TODO: Implement API scope check
	return false
}

// AddReadUserScope adds read_user scope to the application
func (s *OAuthService) AddReadUserScope(app *models.OAuthApplication) error {
	// TODO: Implement adding read_user scope
	return nil
}

// HasDangerousScopes checks if the application has dangerous scopes
func (s *OAuthService) HasDangerousScopes(app *models.OAuthApplication) bool {
	// TODO: Implement dangerous scopes check
	return false
}

// RevokeToken revokes an OAuth token
func (s *OAuthService) RevokeToken(token string, tokenTypeHint string) error {
	// Find token in database
	var oauthToken OAuthToken
	if err := s.db.Where("token = ?", token).First(&oauthToken).Error; err != nil {
		return err
	}

	// Delete the token
	if err := s.db.Delete(&oauthToken).Error; err != nil {
		return err
	}

	return nil
}

// RevokeApplicationTokens revokes all tokens for an application
func (s *OAuthService) RevokeApplicationTokens(userID, appID uint) error {
	// TODO: Implement application token revocation
	return nil
}

// FindDeviceGrantByUserCode finds a device grant by user code
func (s *OAuthService) FindDeviceGrantByUserCode(userCode string) (*models.DeviceGrant, error) {
	// TODO: Implement device grant lookup
	return &models.DeviceGrant{
		UserCode: userCode,
		Scopes:   "",
	}, nil
}

// GenerateDeviceCode generates a new device code and user code
func (s *OAuthService) GenerateDeviceCode() (string, string, error) {
	// Generate a random device code
	deviceCodeBytes := make([]byte, 20)
	if _, err := rand.Read(deviceCodeBytes); err != nil {
		return "", "", err
	}
	deviceCode := base32.StdEncoding.EncodeToString(deviceCodeBytes)

	// Generate a user-friendly code (8 characters)
	userCodeBytes := make([]byte, 5)
	if _, err := rand.Read(userCodeBytes); err != nil {
		return "", "", err
	}
	userCode := base32.StdEncoding.EncodeToString(userCodeBytes)[:8]

	// Store the codes with expiration
	s.deviceCodes[deviceCode] = DeviceCodeInfo{
		UserCode:   userCode,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
		IsVerified: false,
	}

	return deviceCode, userCode, nil
}

// VerifyUserCode verifies a user code and returns the associated device grant
func (s *OAuthService) VerifyUserCode(userCode string) (bool, error) {
	// Find the device code associated with this user code
	var foundDeviceCode string
	for deviceCode, info := range s.deviceCodes {
		if info.UserCode == userCode {
			foundDeviceCode = deviceCode
			break
		}
	}

	if foundDeviceCode == "" {
		return false, errors.New("invalid user code")
	}

	info := s.deviceCodes[foundDeviceCode]
	if time.Now().After(info.ExpiresAt) {
		delete(s.deviceCodes, foundDeviceCode)
		return false, errors.New("code expired")
	}

	// Mark the code as verified
	info.IsVerified = true
	s.deviceCodes[foundDeviceCode] = info

	return true, nil
}

// generateRandomCode generates a random code of specified length
func generateRandomCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes), nil
}

// OAuthToken represents an OAuth access token
type OAuthToken struct {
	Token         string    `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null"`
	ApplicationID uint      `gorm:"not null"`
	Scopes        []string  `gorm:"type:text"`
	ExpiresIn     int64     `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
}

// TokenInfo represents information about an OAuth token
type TokenInfo struct {
	UserID        uint      `json:"resource_owner_id"`
	ApplicationID uint      `json:"application_id"`
	Scopes        []string  `json:"scopes"`
	ExpiresIn     int64     `json:"expires_in"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetTokenInfo retrieves information about an OAuth token
func (s *OAuthService) GetTokenInfo(token string) (*TokenInfo, error) {
	// Find token in database
	var oauthToken OAuthToken
	if err := s.db.Where("token = ?", token).First(&oauthToken).Error; err != nil {
		return nil, err
	}

	// Check if token is expired
	if oauthToken.ExpiresIn > 0 && time.Now().After(oauthToken.CreatedAt.Add(time.Duration(oauthToken.ExpiresIn)*time.Second)) {
		return nil, errors.New("token expired")
	}

	return &TokenInfo{
		UserID:        oauthToken.UserID,
		ApplicationID: oauthToken.ApplicationID,
		Scopes:        oauthToken.Scopes,
		ExpiresIn:     oauthToken.ExpiresIn,
		CreatedAt:     oauthToken.CreatedAt,
	}, nil
} 