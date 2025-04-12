package member

import (
	"strings"
	"unicode"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/reflect"
)

// RendersMemberAccess handles rendering of member access in GitLab
type RendersMemberAccess struct {
	userService UserService
}

// UserService defines the interface for user-related operations
type UserService interface {
	MaxMemberAccessForIDs(entityType string, ids []int64) map[int64]int
}

// NewRendersMemberAccess creates a new instance of RendersMemberAccess
func NewRendersMemberAccess(userService UserService) *RendersMemberAccess {
	return &RendersMemberAccess{
		userService: userService,
	}
}

// PrepareGroupsForRendering prepares groups for rendering by preloading max member access
func (r *RendersMemberAccess) PrepareGroupsForRendering(groups []*model.Group) []*model.Group {
	r.preloadMaxMemberAccessForCollection("group", groups)
	return groups
}

// preloadMaxMemberAccessForCollection preloads the maximum member access for a collection of objects
func (r *RendersMemberAccess) preloadMaxMemberAccessForCollection(entityType string, collection interface{}) {
	// Check if collection is empty
	if collection == nil || reflect.IsEmpty(collection) {
		return
	}

	// Get collection IDs
	var ids []int64
	switch c := collection.(type) {
	case []*model.Group:
		ids = make([]int64, 0, len(c))
		for _, item := range c {
			ids = append(ids, item.ID)
		}
	case []*model.Project:
		ids = make([]int64, 0, len(c))
		for _, item := range c {
			ids = append(ids, item.ID)
		}
	default:
		// Handle other collection types if needed
		return
	}

	// Convert entity type to snake case for method name
	entityTypeSnake := toSnakeCase(entityType)

	// Call the appropriate method on the user service
	r.userService.MaxMemberAccessForIDs(entityTypeSnake, ids)
}

// toSnakeCase converts a string from camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
