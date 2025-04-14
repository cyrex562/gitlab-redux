package models

import (
	"time"
)

// Group represents a group in the system
type Group struct {
	ID          string
	Name        string
	Path        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User represents a user in the system
type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session represents a user session
type Session struct {
	ID           string
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Project represents a project in the system
type Project struct {
	ID             string
	Name           string
	Path           string
	Description    string
	ImportState    *ImportState
	ImportFinished bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ImportState represents the state of an import
type ImportState struct {
	Status     string
	LastError  string
	FinishedAt *time.Time
}

// IsFinished checks if the import is finished
func (s *ImportState) IsFinished() bool {
	return s.Status == "finished" || s.FinishedAt != nil
}

// IsFailed checks if the import failed
func (s *ImportState) IsFailed() bool {
	return s.Status == "failed" || s.LastError != ""
}

// GitHubRepo represents a GitHub repository
type GitHubRepo struct {
	ID          string
	Name        string
	FullName    string
	Description string
	Private     bool
	Fork        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ImportFailure represents an import failure
type ImportFailure struct {
	ID          string
	ProjectID   string
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// FeatureFlags represents feature flags for the system
type FeatureFlags struct {
	Flags map[string]bool
}

// Member represents a member of a group or project
type Member struct {
	ID        string
	UserID    string
	GroupID   string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SourceUser represents a user from an external source
type SourceUser struct {
	ID        string
	Username  string
	Email     string
	Source    string
	CreatedAt time.Time
	UpdatedAt time.Time
} 