package sorting

import (
	"strings"

	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/auth"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/cookies"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/pagination"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/sorting"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// SortingPreference provides functionality for handling sorting preferences
type SortingPreference struct {
	sortingHelper    *sorting.SortingHelper
	cookiesHelper    *cookies.CookiesHelper
	paginationParams *pagination.PaginationParams
	authService      *auth.AuthService
	userService      *model.UserService
}

// NewSortingPreference creates a new instance of SortingPreference
func NewSortingPreference(
	sortingHelper *sorting.SortingHelper,
	cookiesHelper *cookies.CookiesHelper,
	paginationParams *pagination.PaginationParams,
	authService *auth.AuthService,
	userService *model.UserService,
) *SortingPreference {
	return &SortingPreference{
		sortingHelper:    sortingHelper,
		cookiesHelper:    cookiesHelper,
		paginationParams: paginationParams,
		authService:      authService,
		userService:      userService,
	}
}

// SetSortOrder sets the sort order based on user preferences, cookies, or pagination params
func (s *SortingPreference) SetSortOrder(field string, defaultOrder string) string {
	sortOrder := s.setSortOrderFromUserPreference(field)
	if sortOrder == "" {
		sortOrder = s.setSortOrderFromCookie(field)
	}
	if sortOrder == "" {
		sortOrder = s.paginationParams.Sort
	}

	// Some types of sorting might not be available on the dashboard
	if !s.validSortOrder(sortOrder) {
		return defaultOrder
	}

	return sortOrder
}

// SortingField returns the field to store the sorting parameter
// This should be implemented by controllers that use this module
func (s *SortingPreference) SortingField() string {
	return ""
}

// DefaultSortOrder returns the default sort order
// This should be implemented by controllers that use this module
func (s *SortingPreference) DefaultSortOrder() string {
	return ""
}

// LegacySortCookieName returns the legacy sort cookie name
// This should be implemented by controllers that use this module
func (s *SortingPreference) LegacySortCookieName() string {
	return ""
}

// setSortOrderFromUserPreference sets the sort order from user preferences
func (s *SortingPreference) setSortOrderFromUserPreference(field string) string {
	currentUser := s.authService.CurrentUser()
	if currentUser == nil || field == "" {
		return ""
	}

	userPreference := currentUser.UserPreference
	if userPreference == nil {
		return ""
	}

	sortParam := s.paginationParams.Sort
	if sortParam == "" {
		sortParam = userPreference.GetPreference(field)
	}

	// In read-only mode, just return the sort param without updating
	if s.userService.IsDatabaseReadOnly() {
		return sortParam
	}

	// Update user preference if it has changed
	if userPreference.GetPreference(field) != sortParam {
		userPreference.SetPreference(field, sortParam)
	}

	return sortParam
}

// setSortOrderFromCookie sets the sort order from cookies
func (s *SortingPreference) setSortOrderFromCookie(field string) string {
	legacyCookieName := s.LegacySortCookieName()
	if legacyCookieName == "" {
		return ""
	}

	sortParam := ""
	if s.paginationParams.Sort != "" {
		sortParam = s.paginationParams.Sort
	} else {
		// Fallback to legacy cookie value for backward compatibility
		sortParam = s.cookiesHelper.GetCookie(legacyCookieName)
		if sortParam == "" {
			sortParam = s.cookiesHelper.GetCookie(s.rememberSortingKey(field))
		}
	}

	sortValue := s.updateCookieValue(sortParam)
	s.cookiesHelper.SetSecureCookie(s.rememberSortingKey(field), sortValue)
	return sortValue
}

// rememberSortingKey converts sorting_field to legacy cookie name for backwards compatibility
// :merge_requests_sort => 'mergerequest_sort'
// :issues_sort => 'issue_sort'
func (s *SortingPreference) rememberSortingKey(field string) string {
	if field == "" {
		return ""
	}

	parts := strings.Split(field, "_")
	if len(parts) < 2 {
		return ""
	}

	// Take all parts except the last one (which is "sort")
	baseParts := parts[:len(parts)-1]

	// Singularize each part (simplified version)
	for i, part := range baseParts {
		if strings.HasSuffix(part, "ies") {
			baseParts[i] = strings.TrimSuffix(part, "ies") + "y"
		} else if strings.HasSuffix(part, "es") {
			baseParts[i] = strings.TrimSuffix(part, "es")
		} else if strings.HasSuffix(part, "s") {
			baseParts[i] = strings.TrimSuffix(part, "s")
		}
	}

	return strings.Join(baseParts, "") + "_sort"
}

// updateCookieValue updates old values to the actual ones
func (s *SortingPreference) updateCookieValue(value string) string {
	switch value {
	case "id_asc":
		return s.sortingHelper.SortValueOldestCreated()
	case "id_desc":
		return s.sortingHelper.SortValueRecentlyCreated()
	case "downvotes_asc", "downvotes_desc":
		return s.sortingHelper.SortValuePopularity()
	default:
		return value
	}
}

// validSortOrder checks if the sort order is valid
func (s *SortingPreference) validSortOrder(sortOrder string) bool {
	if sortOrder == "" {
		return false
	}

	if strings.Contains(sortOrder, "weight") {
		return s.sortingHelper.CanSortByIssueWeight()
	}

	if strings.Contains(sortOrder, "merged_at") {
		return s.sortingHelper.CanSortByMergedDate()
	}

	return true
}
