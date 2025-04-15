package models

import (
	"time"
)

// OAuthApplication represents an OAuth application
type OAuthApplication struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	RedirectURI  string    `json:"redirect_uri"`
	Scopes       string    `json:"scopes"`
	Secret       string    `json:"-"`
	OwnerID      uint      `json:"owner_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Confidential bool      `json:"confidential"`
	Trusted      bool      `json:"trusted"`
}

// OAuthToken represents an OAuth token
type OAuthToken struct {
	ID            uint      `json:"id"`
	ApplicationID uint      `json:"application_id"`
	UserID        uint      `json:"user_id"`
	Token         string    `json:"-"`
	RefreshToken  string    `json:"-"`
	ExpiresIn     int       `json:"expires_in"`
	Scopes        string    `json:"scopes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
} 