package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

// ApplicationSettings represents the application-wide settings
type ApplicationSettings struct {
	ID                              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
	SignupEnabled                   bool      `json:"signup_enabled"`
	SigninEnabled                   bool      `json:"signin_enabled"`
	GitlabEnabled                   bool      `json:"gitlab_enabled"`
	GithubEnabled                   bool      `json:"github_enabled"`
	GoogleOAuth2Enabled             bool      `json:"google_oauth2_enabled"`
	GoogleOAuth2ClientID           string    `json:"google_oauth2_client_id"`
	GoogleOAuth2ClientSecret       string    `json:"google_oauth2_client_secret"`
	RegistrationToken              string    `json:"registration_token"`
	HealthCheckAccessToken         string    `json:"health_check_access_token"`
	ErrorTrackingAccessToken       string    `json:"error_tracking_access_token"`
	RepositoryCheckStatesEnabled   bool      `json:"repository_check_states_enabled"`
	RepositoryCheckStatesInterval  int       `json:"repository_check_states_interval"`
	UsageDataEnabled              bool      `json:"usage_data_enabled"`
	UsageDataLastReportedAt       time.Time `json:"usage_data_last_reported_at"`
	UsageDataReportingInterval    int       `json:"usage_data_reporting_interval"`
	MetricsEnabled                bool      `json:"metrics_enabled"`
	MetricsPort                   int       `json:"metrics_port"`
	ProfilingEnabled              bool      `json:"profiling_enabled"`
	ProfilingPort                 int       `json:"profiling_port"`
	NetworkTimeout                int       `json:"network_timeout"`
	NetworkRetryCount             int       `json:"network_retry_count"`
	NetworkRetryInterval          int       `json:"network_retry_interval"`
	DefaultTheme                  string    `json:"default_theme"`
	DefaultLocale                 string    `json:"default_locale"`
	DefaultTimezone               string    `json:"default_timezone"`
	DefaultDateFormat             string    `json:"default_date_format"`
	DefaultTimeFormat             string    `json:"default_time_format"`
}

// ApplicationSettingParams represents the parameters for updating application settings
type ApplicationSettingParams struct {
	SignupEnabled                   *bool   `json:"signup_enabled"`
	SigninEnabled                   *bool   `json:"signin_enabled"`
	GitlabEnabled                   *bool   `json:"gitlab_enabled"`
	GithubEnabled                   *bool   `json:"github_enabled"`
	GoogleOAuth2Enabled             *bool   `json:"google_oauth2_enabled"`
	GoogleOAuth2ClientID           *string `json:"google_oauth2_client_id"`
	GoogleOAuth2ClientSecret       *string `json:"google_oauth2_client_secret"`
	RepositoryCheckStatesEnabled   *bool   `json:"repository_check_states_enabled"`
	RepositoryCheckStatesInterval  *int    `json:"repository_check_states_interval"`
	UsageDataEnabled              *bool   `json:"usage_data_enabled"`
	UsageDataReportingInterval    *int    `json:"usage_data_reporting_interval"`
	MetricsEnabled                *bool   `json:"metrics_enabled"`
	MetricsPort                   *int    `json:"metrics_port"`
	ProfilingEnabled              *bool   `json:"profiling_enabled"`
	ProfilingPort                 *int    `json:"profiling_port"`
	NetworkTimeout                *int    `json:"network_timeout"`
	NetworkRetryCount             *int    `json:"network_retry_count"`
	NetworkRetryInterval          *int    `json:"network_retry_interval"`
	DefaultTheme                  *string `json:"default_theme"`
	DefaultLocale                 *string `json:"default_locale"`
	DefaultTimezone               *string `json:"default_timezone"`
	DefaultDateFormat             *string `json:"default_date_format"`
	DefaultTimeFormat             *string `json:"default_time_format"`
}

var currentSettings *ApplicationSettings

// GetCurrentApplicationSettings returns the current application settings
func GetCurrentApplicationSettings() *ApplicationSettings {
	if currentSettings == nil {
		// TODO: Load settings from database
		currentSettings = &ApplicationSettings{
			SignupEnabled:                   true,
			SigninEnabled:                   true,
			GitlabEnabled:                   true,
			GithubEnabled:                   true,
			GoogleOAuth2Enabled:             false,
			RepositoryCheckStatesEnabled:   true,
			RepositoryCheckStatesInterval:  3600,
			UsageDataEnabled:              true,
			UsageDataReportingInterval:    86400,
			MetricsEnabled:                false,
			MetricsPort:                   9090,
			ProfilingEnabled:              false,
			ProfilingPort:                 6060,
			NetworkTimeout:                30,
			NetworkRetryCount:             3,
			NetworkRetryInterval:          5,
			DefaultTheme:                  "light",
			DefaultLocale:                 "en",
			DefaultTimezone:               "UTC",
			DefaultDateFormat:             "YYYY-MM-DD",
			DefaultTimeFormat:             "HH:mm:ss",
		}
	}
	return currentSettings
}

// ResetRegistrationToken generates a new registration token
func (s *ApplicationSettings) ResetRegistrationToken() error {
	token, err := generateRandomToken(32)
	if err != nil {
		return err
	}
	s.RegistrationToken = token
	// TODO: Save to database
	return nil
}

// ResetHealthCheckAccessToken generates a new health check access token
func (s *ApplicationSettings) ResetHealthCheckAccessToken() error {
	token, err := generateRandomToken(32)
	if err != nil {
		return err
	}
	s.HealthCheckAccessToken = token
	// TODO: Save to database
	return nil
}

// ResetErrorTrackingAccessToken generates a new error tracking access token
func (s *ApplicationSettings) ResetErrorTrackingAccessToken() error {
	token, err := generateRandomToken(32)
	if err != nil {
		return err
	}
	s.ErrorTrackingAccessToken = token
	// TODO: Save to database
	return nil
}

// Update updates the application settings with the provided parameters
func (s *ApplicationSettings) Update(params *ApplicationSettingParams) error {
	if params == nil {
		return errors.New("params cannot be nil")
	}

	// Update boolean fields
	if params.SignupEnabled != nil {
		s.SignupEnabled = *params.SignupEnabled
	}
	if params.SigninEnabled != nil {
		s.SigninEnabled = *params.SigninEnabled
	}
	if params.GitlabEnabled != nil {
		s.GitlabEnabled = *params.GitlabEnabled
	}
	if params.GithubEnabled != nil {
		s.GithubEnabled = *params.GithubEnabled
	}
	if params.GoogleOAuth2Enabled != nil {
		s.GoogleOAuth2Enabled = *params.GoogleOAuth2Enabled
	}
	if params.RepositoryCheckStatesEnabled != nil {
		s.RepositoryCheckStatesEnabled = *params.RepositoryCheckStatesEnabled
	}
	if params.UsageDataEnabled != nil {
		s.UsageDataEnabled = *params.UsageDataEnabled
	}
	if params.MetricsEnabled != nil {
		s.MetricsEnabled = *params.MetricsEnabled
	}
	if params.ProfilingEnabled != nil {
		s.ProfilingEnabled = *params.ProfilingEnabled
	}

	// Update string fields
	if params.GoogleOAuth2ClientID != nil {
		s.GoogleOAuth2ClientID = *params.GoogleOAuth2ClientID
	}
	if params.GoogleOAuth2ClientSecret != nil {
		s.GoogleOAuth2ClientSecret = *params.GoogleOAuth2ClientSecret
	}
	if params.DefaultTheme != nil {
		s.DefaultTheme = *params.DefaultTheme
	}
	if params.DefaultLocale != nil {
		s.DefaultLocale = *params.DefaultLocale
	}
	if params.DefaultTimezone != nil {
		s.DefaultTimezone = *params.DefaultTimezone
	}
	if params.DefaultDateFormat != nil {
		s.DefaultDateFormat = *params.DefaultDateFormat
	}
	if params.DefaultTimeFormat != nil {
		s.DefaultTimeFormat = *params.DefaultTimeFormat
	}

	// Update integer fields
	if params.RepositoryCheckStatesInterval != nil {
		s.RepositoryCheckStatesInterval = *params.RepositoryCheckStatesInterval
	}
	if params.UsageDataReportingInterval != nil {
		s.UsageDataReportingInterval = *params.UsageDataReportingInterval
	}
	if params.MetricsPort != nil {
		s.MetricsPort = *params.MetricsPort
	}
	if params.ProfilingPort != nil {
		s.ProfilingPort = *params.ProfilingPort
	}
	if params.NetworkTimeout != nil {
		s.NetworkTimeout = *params.NetworkTimeout
	}
	if params.NetworkRetryCount != nil {
		s.NetworkRetryCount = *params.NetworkRetryCount
	}
	if params.NetworkRetryInterval != nil {
		s.NetworkRetryInterval = *params.NetworkRetryInterval
	}

	// TODO: Save to database
	return nil
}

// generateRandomToken generates a random token of specified length
func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
