package models

import "time"

// Release represents a GitLab release
type Release struct {
	ID          int64     `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GroupID     int64     `json:"group_id"`
	AuthorID    int64     `json:"author_id"`
} 