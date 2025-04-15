package models

import (
	"time"
)

// JiraConnectSubscription represents a Jira Connect subscription
type JiraConnectSubscription struct {
	ID        int64     `json:"id"`
	InstallationID int64 `json:"installation_id"`
	WorkspaceID int64 `json:"workspace_id"`
	RepositoryID int64 `json:"repository_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JiraConnectWorkspace represents a Jira Connect workspace
type JiraConnectWorkspace struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JiraConnectRepository represents a Jira Connect repository
type JiraConnectRepository struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 