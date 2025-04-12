package workhorse

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// WorkhorseAuthorizationHandler handles Workhorse authorization
type WorkhorseAuthorizationHandler struct {
	uploaderService *service.UploaderService
	workhorseService *service.WorkhorseService
}

// NewWorkhorseAuthorizationHandler creates a new Workhorse authorization handler
func NewWorkhorseAuthorizationHandler(
	uploaderService *service.UploaderService,
	workhorseService *service.WorkhorseService,
) *WorkhorseAuthorizationHandler {
	return &WorkhorseAuthorizationHandler{
		uploaderService: uploaderService,
		workhorseService: workhorseService,
	}
}

// RegisterRoutes registers the Workhorse routes
func (h *WorkhorseAuthorizationHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Skip CSRF protection for the authorize endpoint
	router.POST("/authorize", h.authorize)
}

// authorize handles the Workhorse authorization request
func (h *WorkhorseAuthorizationHandler) authorize(c *gin.Context) {
	// Set the content type for Workhorse internal API
	h.setWorkhorseInternalAPIContentType(c)

	// Get the uploader class
	uploaderClass, err := h.getUploaderClass(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the maximum size
	maximumSize, err := h.getMaximumSize(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Authorize the upload
	authorized, err := h.uploaderService.WorkhorseAuthorize(c, uploaderClass, false, maximumSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file"})
		return
	}

	c.JSON(http.StatusOK, authorized)
}

// setWorkhorseInternalAPIContentType sets the content type for Workhorse internal API
func (h *WorkhorseAuthorizationHandler) setWorkhorseInternalAPIContentType(c *gin.Context) {
	c.Header("Content-Type", "application/json")
}

// isFileValid checks if a file is valid
func (h *WorkhorseAuthorizationHandler) isFileValid(file *model.UploadedFile) bool {
	if file == nil {
		return false
	}

	// Get the file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.OriginalFilename), "."))

	// Check if the extension is in the allowlist
	return h.isExtensionAllowed(ext)
}

// isExtensionAllowed checks if a file extension is allowed
func (h *WorkhorseAuthorizationHandler) isExtensionAllowed(ext string) bool {
	// Get the extension allowlist
	allowlist := h.getFileExtensionAllowlist()

	// Check if the extension is in the allowlist
	for _, allowed := range allowlist {
		if allowed == ext {
			return true
		}
	}

	return false
}

// getUploaderClass gets the uploader class for the current request
func (h *WorkhorseAuthorizationHandler) getUploaderClass(c *gin.Context) (string, error) {
	// This should be implemented by the specific handler
	return "", errors.New("getUploaderClass not implemented")
}

// getMaximumSize gets the maximum size for the current request
func (h *WorkhorseAuthorizationHandler) getMaximumSize(c *gin.Context) (int64, error) {
	// This should be implemented by the specific handler
	return 0, errors.New("getMaximumSize not implemented")
}

// getFileExtensionAllowlist gets the file extension allowlist
func (h *WorkhorseAuthorizationHandler) getFileExtensionAllowlist() []string {
	// Default to the ImportExportUploader allowlist
	return []string{
		"gz", "bz2", "tar", "zip", "rar", "7z", "xz",
		"json", "yaml", "yml", "xml", "txt", "md", "markdown",
		"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx",
		"jpg", "jpeg", "png", "gif", "svg", "ico",
		"mp3", "mp4", "avi", "mov", "wmv", "flv", "mkv",
		"sql", "db", "sqlite", "sqlite3",
		"bak", "backup", "old",
	}
}
