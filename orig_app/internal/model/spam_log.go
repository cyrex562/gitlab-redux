package model

import "time"

// SpamLog represents a spam detection log entry
type SpamLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents a GitLab user
type User struct {
	ID                      int64  `json:"id"`
	Username               string `json:"username"`
	TrustedWithSpamAttribute bool  `json:"trusted_with_spam_attribute"`
	// TODO: Add other necessary user fields
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	// TODO: Implement admin check logic
	return false
}
