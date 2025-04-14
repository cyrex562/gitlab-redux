package models

import "time"

// SourceUser represents a user from an external source
type SourceUser struct {
	ID        int64
	GroupID   int64
	UserID    int64
	Source    string
	Username  string
	Email     string
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsReassigned checks if the source user is reassigned
func (u *SourceUser) IsReassigned() bool {
	return u.UserID != 0
}

// IsAwaitingReassignment checks if the source user is awaiting reassignment
func (u *SourceUser) IsAwaitingReassignment() bool {
	return u.UserID == 0
} 