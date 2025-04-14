package models

import (
	"time"
)

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