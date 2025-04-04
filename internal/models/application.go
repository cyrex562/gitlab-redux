package models

import (
	"time"
)

// Application represents an OAuth application
type Application struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	UID          string    `json:"uid"`
	Secret       string    `json:"secret"`
	RedirectURI  string    `json:"redirect_uri"`
	Scopes       string    `json:"scopes"`
	Confidential bool      `json:"confidential"`
	Trusted      bool      `json:"trusted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ApplicationParams represents the parameters for creating or updating an application
type ApplicationParams struct {
	Name         string   `json:"name" binding:"required"`
	RedirectURI  string   `json:"redirect_uri" binding:"required"`
	Scopes       []string `json:"scopes"`
	Confidential bool     `json:"confidential"`
	Trusted      bool     `json:"trusted"`
}

// TableName returns the table name for the Application model
func (Application) TableName() string {
	return "oauth_applications"
}
