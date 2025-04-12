package models

// Group represents a GitLab group
type Group struct {
	ID                    int64
	Path                  string
	DependencyProxySetting *DependencyProxySetting
}

// DependencyProxySetting represents group dependency proxy settings
type DependencyProxySetting struct {
	ID      int64
	GroupID int64
	Enabled bool
} 