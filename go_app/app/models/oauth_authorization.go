package models

import (
	"time"
)

// OAuthClient represents an OAuth client
type OAuthClient struct {
	ID           uint             `json:"id"`
	Application  *OAuthApplication `json:"application"`
	RedirectURI  string           `json:"redirect_uri"`
	Scopes       string           `json:"scopes"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// PreAuthorization represents pre-authorization data
type PreAuthorization struct {
	Authorizable bool        `json:"authorizable"`
	Client       *OAuthClient `json:"client"`
	RedirectURI  string      `json:"redirect_uri"`
	Scopes       string      `json:"scopes"`
	Error        string      `json:"error,omitempty"`
}

// Authorization represents an OAuth authorization
type Authorization struct {
	RedirectURI string `json:"redirect_uri"`
	Code        string `json:"code,omitempty"`
	Token       string `json:"token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	Scopes      string `json:"scopes,omitempty"`
} 