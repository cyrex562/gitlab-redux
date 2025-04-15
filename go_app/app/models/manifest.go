package models

// Manifest represents a manifest file for importing projects
type Manifest struct {
	Projects []Repository `json:"projects"`
}

// ManifestMetadata stores metadata about a manifest import
type ManifestMetadata struct {
	GroupID     string       `json:"group_id"`
	Repositories []Repository `json:"repositories"`
}

// Repository represents a repository in a manifest
type Repository struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// Group represents a GitLab group
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// User represents a GitLab user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Project represents a GitLab project
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	GroupID     string `json:"group_id"`
} 