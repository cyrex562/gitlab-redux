package service

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GithubOAuthConfig represents the GitHub OAuth configuration
type GithubOAuthConfig struct {
	AppID     string
	AppSecret string
	VerifySSL bool
	RedirectURL string
}

// GithubOAuthService handles GitHub OAuth operations
type GithubOAuthService struct {
	config *GithubOAuthConfig
}

// NewGithubOAuthService creates a new instance of GithubOAuthService
func NewGithubOAuthService(config *GithubOAuthConfig) *GithubOAuthService {
	return &GithubOAuthService{
		config: config,
	}
}

// GetOAuthConfig returns the OAuth configuration for GitHub
func (s *GithubOAuthService) GetOAuthConfig() (*oauth2.Config, error) {
	if s.config == nil {
		return nil, &OAuthConfigMissingError{}
	}

	// Create transport with SSL verification option
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !s.config.VerifySSL,
		},
	}

	// Create HTTP client with custom transport
	client := &http.Client{
		Transport: transport,
	}

	// Create OAuth config
	config := &oauth2.Config{
		ClientID:     s.config.AppID,
		ClientSecret: s.config.AppSecret,
		RedirectURL:  s.config.RedirectURL,
		Scopes:       []string{"repo", "read:org"},
		Endpoint:     github.Endpoint,
	}

	// Set custom HTTP client
	ctx := oauth2.WithClient(oauth2.BackgroundContext(), client)
	config = config.WithContext(ctx)

	return config, nil
}

// OAuthConfigMissingError represents an error when OAuth configuration is missing
type OAuthConfigMissingError struct{}

func (e *OAuthConfigMissingError) Error() string {
	return "Missing OAuth configuration for GitHub"
}
