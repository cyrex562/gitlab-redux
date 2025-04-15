package models

import "time"

// Project represents a GitLab project
type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GroupID     int64     `json:"group_id"`
	AuthorID    int64     `json:"author_id"`
	Archived    bool      `json:"archived"`
	// Add other fields as needed
} 