package models

import "time"

// Group represents a GitLab group
type Group struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AuthorID    int64     `json:"author_id"`
	Variables   []*GroupVariable `json:"variables,omitempty"`
	Errors      []string  `json:"errors,omitempty"`
	// Add other fields as needed
}

// UsageQuotasEnabled checks if usage quotas are enabled for the group
func (g *Group) UsageQuotasEnabled() bool {
	// TODO: Implement the actual logic
	return true
} 