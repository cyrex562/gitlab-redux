package models

import (
	"time"
)

// GroupLink represents a link between groups
type GroupLink struct {
	ID           int64
	GroupID      int64
	SharedGroupID int64
	GroupAccess  string
	ExpiresAt    time.Time
	MemberRoleID int64
}

// Expires checks if the group link expires
func (l *GroupLink) Expires() bool {
	return !l.ExpiresAt.IsZero()
}

// ExpiresSoon checks if the group link expires soon
func (l *GroupLink) ExpiresSoon() bool {
	if !l.Expires() {
		return false
	}
	
	// Consider "soon" as within 7 days
	return time.Until(l.ExpiresAt) <= 7*24*time.Hour
}

// Group methods for group links

// FindSharedWithGroupLink finds a shared with group link by ID
func (g *Group) FindSharedWithGroupLink(id string) (*GroupLink, error) {
	// TODO: Implement group link finding logic
	return nil, nil
} 