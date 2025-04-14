package models

import (
	"time"
)

// Group represents a group in the system
type Group struct {
	ID              string
	Name            string
	Path            string
	Description     string
	VisibilityLevel int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DependencyProxySetting *DependencyProxySetting
	Features              map[string]bool
}

// DependencyProxySetting represents group dependency proxy settings
type DependencyProxySetting struct {
	ID      int64
	GroupID int64
	Enabled bool
}

// PackagesFeatureEnabled checks if packages feature is enabled for the group
func (g *Group) PackagesFeatureEnabled() bool {
	if g.Features == nil {
		return false
	}
	return g.Features["packages"]
} 