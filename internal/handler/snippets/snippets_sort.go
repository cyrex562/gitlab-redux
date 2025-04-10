package snippets

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/pagination"
)

// SnippetsSort provides sorting functionality for snippets
type SnippetsSort struct {
	paginationParams *pagination.PaginationParams
}

// NewSnippetsSort creates a new instance of SnippetsSort
func NewSnippetsSort(paginationParams *pagination.PaginationParams) *SnippetsSort {
	return &SnippetsSort{
		paginationParams: paginationParams,
	}
}

// SortParam returns the sort parameter for snippets
// If no sort parameter is provided, it defaults to 'updated_desc'
func (s *SnippetsSort) SortParam() string {
	sort := s.paginationParams.Sort
	if sort == "" {
		return "updated_desc"
	}
	return sort
}
