package model

import (
	"time"
)

// Label represents a label in the system
type Label struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProjectID   int64     `json:"project_id"`
	GroupID     int64     `json:"group_id"`
}

// LabelAppearance represents the appearance of a label for serialization
type LabelAppearance struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Type        string `json:"type"`
	TextColor   string `json:"text_color"`
}
