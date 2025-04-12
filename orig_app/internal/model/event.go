package model

import (
	"time"
)

// Event represents a GitLab event
type Event struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	ProjectID int64     `json:"project_id"`
	ActionName string   `json:"action_name"`
	TargetID   int64    `json:"target_id"`
	TargetType string   `json:"target_type"`
	AuthorID   int64    `json:"author_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TargetTitle string   `json:"target_title"`
	Author      *User    `json:"author"`
	Project     *Project `json:"project"`
}

// RenderedEvent represents a rendered event for display
type RenderedEvent struct {
	*Event
	AuthorName     string `json:"author_name"`
	AuthorUsername string `json:"author_username"`
	AuthorAvatar   string `json:"author_avatar"`
	ProjectName    string `json:"project_name"`
	ProjectPath    string `json:"project_path"`
	TargetName     string `json:"target_name"`
	TargetPath     string `json:"target_path"`
	ActionText     string `json:"action_text"`
	TargetTypeText string `json:"target_type_text"`
}
