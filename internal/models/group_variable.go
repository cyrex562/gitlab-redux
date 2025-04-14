package models

import "time"

// GroupVariable represents a GitLab group CI/CD variable
type GroupVariable struct {
	ID           int64     `json:"id"`
	VariableType string    `json:"variable_type"`
	Key          string    `json:"key"`
	Description  string    `json:"description,omitempty"`
	SecretValue  string    `json:"secret_value,omitempty"`
	Protected    bool      `json:"protected,omitempty"`
	Masked       bool      `json:"masked,omitempty"`
	Hidden       bool      `json:"hidden,omitempty"`
	Raw          bool      `json:"raw,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	GroupID      int64     `json:"group_id"`
} 