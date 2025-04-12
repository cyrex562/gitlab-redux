package import

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// ImportUrlParams handles processing and sanitizing import URLs
type ImportUrlParams struct {
	urlSanitizer *service.UrlSanitizer
}

// NewImportUrlParams creates a new instance of ImportUrlParams
func NewImportUrlParams(urlSanitizer *service.UrlSanitizer) *ImportUrlParams {
	return &ImportUrlParams{
		urlSanitizer: urlSanitizer,
	}
}

// GetImportUrlParams returns the processed import URL parameters
func (i *ImportUrlParams) GetImportUrlParams(ctx *gin.Context) map[string]string {
	// Check if import_url is present in the project parameters
	importURL, exists := ctx.GetPostForm("project[import_url]")
	if !exists || importURL == "" {
		return map[string]string{}
	}

	// Get the import URL credentials if present
	importURLUser, _ := ctx.GetPostForm("project[import_url_user]")
	importURLPassword, _ := ctx.GetPostForm("project[import_url_password]")

	// Process the import URL with credentials
	fullURL := i.importParamsToFullURL(importURL, importURLUser, importURLPassword)

	// Return the processed parameters
	return map[string]string{
		"import_url":  fullURL,
		"import_type": "git", // Always set import_type to 'git' to prevent stale values
	}
}

// importParamsToFullURL processes the import URL with credentials
func (i *ImportUrlParams) importParamsToFullURL(importURL, importURLUser, importURLPassword string) string {
	// Create credentials map
	credentials := map[string]string{
		"user":     importURLUser,
		"password": importURLPassword,
	}

	// Use the URL sanitizer to get the full URL
	return i.urlSanitizer.SanitizeURL(importURL, credentials)
}
