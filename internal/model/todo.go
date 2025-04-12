package model

import (
	"time"
)

// Todo represents a user's todo item
type Todo struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	IssuableType string    `json:"issuable_type"`
	IssuableID   int64     `json:"issuable_id"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName returns the database table name for todos
func (t *Todo) TableName() string {
	return "todos"
}

// IsPending returns true if the todo is in a pending state
func (t *Todo) IsPending() bool {
	return t.State == "pending"
}
