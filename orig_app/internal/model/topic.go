package model

import "time"

// Topic represents a GitLab project topic
type Topic struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	Name           string    `json:"name"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Avatar         string    `json:"avatar,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TopicParams represents the parameters for creating or updating a topic
type TopicParams struct {
	OrganizationID int64  `json:"organization_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Avatar         string `json:"avatar,omitempty"`
}

// TitleOrName returns the topic's title if available, otherwise returns the name
func (t *Topic) TitleOrName() string {
	if t.Title != "" {
		return t.Title
	}
	return t.Name
}
