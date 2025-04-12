package service

import (
	"context"

	"github.com/jmadden/gitlab-redux/internal/model"
)

// SnippetService defines the interface for snippet-related operations
type SnippetService interface {
	// FindSnippets finds snippets based on the given parameters
	FindSnippets(ctx context.Context, currentUser, author *model.User, scope, sort string, page, perPage int) (*model.PaginatedSnippets, error)

	// GetSnippetByID gets a snippet by its ID
	GetSnippetByID(ctx context.Context, id int64, currentUser *model.User) (*model.Snippet, error)

	// CreateSnippet creates a new snippet
	CreateSnippet(ctx context.Context, snippet *model.Snippet, currentUser *model.User) (*model.Snippet, error)

	// UpdateSnippet updates an existing snippet
	UpdateSnippet(ctx context.Context, snippet *model.Snippet, currentUser *model.User) (*model.Snippet, error)

	// DeleteSnippet deletes a snippet
	DeleteSnippet(ctx context.Context, id int64, currentUser *model.User) error
}
