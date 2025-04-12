package model

import (
	"time"
)

// BroadcastMessageDismissal represents a user's dismissal of a broadcast message
type BroadcastMessageDismissal struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	MessageID int64     `json:"message_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CookieKey returns the cookie key for this dismissal
func (d *BroadcastMessageDismissal) CookieKey() string {
	return "broadcast_message_dismissed_" + string(d.MessageID)
}

// IsExpired returns true if the dismissal has expired
func (d *BroadcastMessageDismissal) IsExpired() bool {
	return time.Now().After(d.ExpiresAt)
}
