package model

import (
	"time"
)

// Snippet represents a GitLab snippet
type Snippet struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	FileName    string    `json:"file_name"`
	Type        string    `json:"type"`
	Visibility  string    `json:"visibility"`
	AuthorID    int64     `json:"author_id"`
	ProjectID   int64     `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	WebURL      string    `json:"web_url"`
	RawURL      string    `json:"raw_url"`
	Author      *User     `json:"author"`
	Project     *Project  `json:"project"`
}

// PaginatedSnippets represents a paginated collection of snippets
type PaginatedSnippets struct {
	Items      []*Snippet `json:"items"`
	TotalItems int        `json:"total_items"`
	TotalPages int        `json:"total_pages"`
	CurrentPage int       `json:"current_page"`
	PerPage    int        `json:"per_page"`
}

// SnippetCounts represents the counts for different types of snippets
type SnippetCounts struct {
	All      int `json:"all"`
	Personal int `json:"personal"`
	Public   int `json:"public"`
	Internal int `json:"internal"`
	Private  int `json:"private"`
}
