package service

import (
	"context"
	"errors"
	"net/http"
)

var (
	// ErrRegistryConnection represents a registry connection error
	ErrRegistryConnection = errors.New("failed to connect to container registry")
)

// RegistryService handles registry operations
type RegistryService struct {
	client *http.Client
	config *RegistryConfig
}

// RegistryConfig represents the registry configuration
type RegistryConfig struct {
	URL      string
	Username string
	Password string
}

// NewRegistryService creates a new instance of RegistryService
func NewRegistryService(config *RegistryConfig) *RegistryService {
	return &RegistryService{
		client: &http.Client{},
		config: config,
	}
}

// PingRegistry pings the container registry to check connectivity
func (s *RegistryService) PingRegistry(ctx context.Context) error {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", s.config.URL+"/v2/", nil)
	if err != nil {
		return err
	}

	// Add authentication if provided
	if s.config.Username != "" && s.config.Password != "" {
		req.SetBasicAuth(s.config.Username, s.config.Password)
	}

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return ErrRegistryConnection
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		return ErrRegistryConnection
	}

	return nil
}

// GetRegistryInfo gets information about the container registry
func (s *RegistryService) GetRegistryInfo(ctx context.Context) (map[string]interface{}, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", s.config.URL+"/v2/", nil)
	if err != nil {
		return nil, err
	}

	// Add authentication if provided
	if s.config.Username != "" && s.config.Password != "" {
		req.SetBasicAuth(s.config.Username, s.config.Password)
	}

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, ErrRegistryConnection
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, ErrRegistryConnection
	}

	// Return registry info
	return map[string]interface{}{
		"version": "2.0",
		"url":     s.config.URL,
	}, nil
}
