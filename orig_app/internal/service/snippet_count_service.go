package service

import (
	"context"

	"github.com/jmadden/gitlab-redux/internal/model"
)

// SnippetCountService defines the interface for snippet count operations
type SnippetCountService interface {
	// GetSnippetCounts gets the counts for different types of snippets
	GetSnippetCounts(ctx context.Context, currentUser, author *model.User) (*model.SnippetCounts, error)
}
