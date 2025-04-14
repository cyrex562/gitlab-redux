package models

import (
	"time"
)

// Member represents a group or project member
type Member struct {
	ID           int64
	GroupID      int64
	UserID       int64
	AccessLevel  string
	ExpiresAt    time.Time
	InviteEmail  string
	InviteToken  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    *User
	UpdatedBy    *User
	User         *User
	State        string
	NotifyEmail  bool
	MemberRoleID int64
}

// IsInvited checks if the member is invited
func (m *Member) IsInvited() bool {
	return m.InviteEmail != "" && m.InviteToken != ""
}

// Expires checks if the member expires
func (m *Member) Expires() bool {
	return !m.ExpiresAt.IsZero()
}

// ExpiresSoon checks if the member expires soon
func (m *Member) ExpiresSoon() bool {
	if !m.Expires() {
		return false
	}
	
	// Consider "soon" as within 7 days
	return time.Until(m.ExpiresAt) <= 7*24*time.Hour
}

// Group methods for members

// Members returns the group members
func (g *Group) Members() []*Member {
	// TODO: Implement members retrieval logic
	return nil
} 