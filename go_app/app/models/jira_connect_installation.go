package models

import (
	"time"
)

// JiraConnectInstallation represents a Jira Connect installation
type JiraConnectInstallation struct {
	ID        int64     `json:"id"`
	ClientKey string    `json:"client_key"`
	SharedSecret string `json:"shared_secret"`
	BaseURL   string    `json:"base_url"`
	Proxy     bool      `json:"proxy"`
	CreateBranchURL string `json:"create_branch_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JiraUser represents a Jira user
type JiraUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
} 