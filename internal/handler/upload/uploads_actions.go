package upload

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// UploadsActionsHandler handles HTTP requests for file uploads
type UploadsActionsHandler struct {
	uploadService *service.UploadService
}

// NewUploadsActionsHandler creates a new handler instance
func NewUploadsActionsHandler(uploadService *service.UploadService) *UploadsActionsHandler {
	return &UploadsActionsHandler{
		uploadService: uploadService,
	}
}

// RegisterRoutes registers the handler routes
func (h *UploadsActionsHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		uploads := api.Group("/uploads")
		{
			uploads.POST("/", h.create)
			uploads.GET("/:filename", h.show)
			uploads.POST("/authorize", h.authorize)
		}
	}
}

// create handles POST /api/uploads
func (h *UploadsActionsHandler) create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the model from the context (set by middleware)
	model, exists := c.Get("model")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "model not found"})
		return
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Get the uploader class from the context (set by middleware)
	uploaderClass, exists := c.Get("uploader_class")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uploader class not found"})
		return
	}

	// Create the upload
	upload, err := h.uploadService.CreateUpload(c.Request.Context(), model, file, uploaderClass.(string), userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"link": upload.ToMap()})
}

// show handles GET /api/uploads/:filename
func (h *UploadsActionsHandler) show(c *gin.Context) {
	filename := c.Param("filename")

	// Check for path traversal
	if strings.Contains(filename, "..") {
		c.Status(http.StatusBadRequest)
		return
	}

	// Get the uploader from the context (set by middleware)
	uploader, exists := c.Get("uploader")
	if !exists {
		c.Status(http.StatusNotFound)
		return
	}

	// Check if the file exists
	if !h.uploadService.FileExists(uploader) {
		c.Status(http.StatusNotFound)
		return
	}

	// Get the file uploader
	fileUploader, err := h.uploadService.GetFileUploader(uploader, filename)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Set cache headers
	ttl, directives := h.uploadService.GetCacheSettings()
	if ttl > 0 {
		c.Header("Cache-Control", fmt.Sprintf("private, must-revalidate, max-age=%d", ttl))
		c.Header("Expires", directives["expires"])
	}

	// Set content type
	c.Header("Content-Type", h.uploadService.GetContentType(fileUploader))

	// Set content disposition
	disposition := "attachment"
	if h.uploadService.IsEmbeddable(fileUploader) || h.uploadService.IsPDF(fileUploader) {
		disposition = "inline"
	}
	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=%s", disposition, filepath.Base(filename)))

	// Send the file
	c.File(fileUploader.Path)
}

// authorize handles POST /api/uploads/authorize
func (h *UploadsActionsHandler) authorize(c *gin.Context) {
	// Set content type
	c.Header("Content-Type", "application/json")

	// Get the uploader class from the context (set by middleware)
	uploaderClass, exists := c.Get("uploader_class")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uploader class not found"})
		return
	}

	// Authorize the upload
	authorized, err := h.uploadService.AuthorizeUpload(c.Request.Context(), uploaderClass.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error uploading file"})
		return
	}

	c.JSON(http.StatusOK, authorized)
}
