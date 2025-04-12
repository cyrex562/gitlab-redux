package tree

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// RedirectsForMissingPath handles redirects for missing paths in tree views
type RedirectsForMissingPath struct {
	projectService *service.ProjectService
	i18nService    *service.I18nService
	logger         *service.Logger
}

// NewRedirectsForMissingPath creates a new instance of RedirectsForMissingPath
func NewRedirectsForMissingPath(
	projectService *service.ProjectService,
	i18nService *service.I18nService,
	logger *service.Logger,
) *RedirectsForMissingPath {
	return &RedirectsForMissingPath{
		projectService: projectService,
		i18nService:   i18nService,
		logger:        logger,
	}
}

// RedirectToTreeRootForMissingPath redirects to the tree root when a path is missing
func (r *RedirectsForMissingPath) RedirectToTreeRootForMissingPath(c *gin.Context, project interface{}, ref string, path string, refType string) error {
	// Generate project tree path
	treePath, err := r.projectService.GetProjectTreePath(project, ref, refType)
	if err != nil {
		return err
	}

	// Get missing path notice
	notice := r.missingPathOnRef(path, ref)

	// Redirect with notice
	c.Redirect(302, treePath)
	c.Set("notice", notice)

	return nil
}

// MissingPathOnRef generates the message for a missing path on a ref
func (r *RedirectsForMissingPath) missingPathOnRef(path string, ref string) string {
	// Truncate path
	truncatedPath := r.truncatePath(path)

	// Format message using i18n service
	return r.i18nService.Format(`"%s" did not exist on "%s"`, truncatedPath, ref)
}

// TruncatePath truncates a path with a separator
func (r *RedirectsForMissingPath) truncatePath(path string) string {
	// Reverse the path
	runes := []rune(path)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	reversedPath := string(runes)

	// Truncate with separator
	truncated := r.truncateWithSeparator(reversedPath, 60, "/")

	// Reverse back
	runes = []rune(truncated)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// TruncateWithSeparator truncates a string at a separator
func (r *RedirectsForMissingPath) truncateWithSeparator(s string, maxLen int, separator string) string {
	if len(s) <= maxLen {
		return s
	}

	// Find the last occurrence of the separator within the maxLen
	truncateIdx := strings.LastIndex(s[:maxLen], separator)
	if truncateIdx == -1 {
		return s[:maxLen]
	}

	return s[:truncateIdx]
}
