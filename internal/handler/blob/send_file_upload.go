package blob

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SendFileUpload handles sending file uploads to clients
type SendFileUpload struct {
	workhorseClient *service.WorkhorseClient
	cdnService      *service.CDNService
}

// NewSendFileUpload creates a new instance of SendFileUpload
func NewSendFileUpload(
	workhorseClient *service.WorkhorseClient,
	cdnService *service.CDNService,
) *SendFileUpload {
	return &SendFileUpload{
		workhorseClient: workhorseClient,
		cdnService:      cdnService,
	}
}

// SendUpload sends a file upload to the client
func (s *SendFileUpload) SendUpload(
	c *gin.Context,
	fileUpload *model.FileUpload,
	sendParams map[string]interface{},
	redirectParams map[string]interface{},
	attachment string,
	proxy bool,
	disposition string,
) {
	// Default disposition to "attachment" if not provided
	if disposition == "" {
		disposition = "attachment"
	}

	// Get content type for the attachment
	contentType := s.contentTypeFor(attachment)

	// Handle attachment-specific parameters
	if attachment != "" {
		// Format the content disposition
		responseDisposition := fmt.Sprintf("%s; filename=\"%s\"", disposition, attachment)

		// Set response headers for cloud storage
		if redirectParams == nil {
			redirectParams = make(map[string]interface{})
		}
		if redirectParams["query"] == nil {
			redirectParams["query"] = make(map[string]string)
		}
		queryParams := redirectParams["query"].(map[string]string)
		queryParams["response-content-disposition"] = responseDisposition
		queryParams["response-content-type"] = contentType

		// Handle JavaScript files to avoid cross-origin protection
		if filepath.Ext(attachment) == ".js" {
			if sendParams == nil {
				sendParams = make(map[string]interface{})
			}
			sendParams["content_type"] = "text/plain"
		}

		// Set filename and disposition in send params
		if sendParams == nil {
			sendParams = make(map[string]interface{})
		}
		sendParams["filename"] = attachment
		sendParams["disposition"] = disposition
	}

	// Handle different types of file uploads
	if s.isImageScalingRequest(fileUpload, c) {
		// Handle image scaling request
		location := fileUpload.Path
		if !fileUpload.IsFileStorage() {
			location = fileUpload.URL
		}

		width := c.Query("width")
		widthInt := 0
		if width != "" {
			fmt.Sscanf(width, "%d", &widthInt)
		}

		// Send scaled image through workhorse
		headers := s.workhorseClient.SendScaledImage(location, widthInt, contentType)
		for key, value := range headers {
			c.Header(key, value)
		}
		c.Status(http.StatusOK)
	} else if fileUpload.IsFileStorage() {
		// Send file directly
		c.File(fileUpload.Path)
	} else if fileUpload.IsProxyDownloadEnabled() || proxy {
		// Send URL through workhorse
		headers := s.workhorseClient.SendURL(fileUpload.URL, redirectParams)
		for key, value := range headers {
			c.Header(key, value)
		}
		c.Status(http.StatusOK)
	} else {
		// Redirect to CDN URL
		clientIP := c.ClientIP()
		fileURL := s.cdnService.GetFileURL(fileUpload, clientIP, redirectParams)
		c.Redirect(http.StatusFound, fileURL)
	}
}

// ContentTypeFor returns the content type for an attachment
func (s *SendFileUpload) ContentTypeFor(attachment string) string {
	if attachment == "" {
		return ""
	}
	return s.guessContentType(attachment)
}

// GuessContentType guesses the content type for a filename
func (s *SendFileUpload) GuessContentType(filename string) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}

// IsImageScalingRequest checks if the request is for image scaling
func (s *SendFileUpload) IsImageScalingRequest(fileUpload *model.FileUpload, c *gin.Context) bool {
	return s.isAvatarSafeForScaling(fileUpload, c) || s.isPwaIconSafeForScaling(fileUpload, c)
}

// IsPwaIconSafeForScaling checks if a PWA icon is safe for scaling
func (s *SendFileUpload) IsPwaIconSafeForScaling(fileUpload *model.FileUpload, c *gin.Context) bool {
	width := c.Query("width")
	widthInt := 0
	if width != "" {
		fmt.Sscanf(width, "%d", &widthInt)
	}

	return fileUpload.IsImageSafeForScaling() &&
		s.isMountedAsPwaIcon(fileUpload) &&
		s.isValidImageScalingWidth(widthInt, []int{192, 512}) // Example allowed widths
}

// IsAvatarSafeForScaling checks if an avatar is safe for scaling
func (s *SendFileUpload) IsAvatarSafeForScaling(fileUpload *model.FileUpload, c *gin.Context) bool {
	width := c.Query("width")
	widthInt := 0
	if width != "" {
		fmt.Sscanf(width, "%d", &widthInt)
	}

	return fileUpload.IsImageSafeForScaling() &&
		s.isMountedAsAvatar(fileUpload) &&
		s.isValidImageScalingWidth(widthInt, []int{32, 64, 128}) // Example allowed widths
}

// IsMountedAsAvatar checks if a file upload is mounted as an avatar
func (s *SendFileUpload) IsMountedAsAvatar(fileUpload *model.FileUpload) bool {
	return fileUpload.MountedAs == "avatar"
}

// IsMountedAsPwaIcon checks if a file upload is mounted as a PWA icon
func (s *SendFileUpload) IsMountedAsPwaIcon(fileUpload *model.FileUpload) bool {
	return fileUpload.MountedAs == "pwa_icon"
}

// IsValidImageScalingWidth checks if the requested width is valid for scaling
func (s *SendFileUpload) IsValidImageScalingWidth(width int, allowedScalarWidths []int) bool {
	for _, allowedWidth := range allowedScalarWidths {
		if allowedWidth == width {
			return true
		}
	}
	return false
}

// Private helper methods
func (s *SendFileUpload) contentTypeFor(attachment string) string {
	if attachment == "" {
		return ""
	}
	return s.guessContentType(attachment)
}

func (s *SendFileUpload) guessContentType(filename string) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}

func (s *SendFileUpload) isImageScalingRequest(fileUpload *model.FileUpload, c *gin.Context) bool {
	return s.isAvatarSafeForScaling(fileUpload, c) || s.isPwaIconSafeForScaling(fileUpload, c)
}

func (s *SendFileUpload) isPwaIconSafeForScaling(fileUpload *model.FileUpload, c *gin.Context) bool {
	width := c.Query("width")
	widthInt := 0
	if width != "" {
		fmt.Sscanf(width, "%d", &widthInt)
	}

	return fileUpload.IsImageSafeForScaling() &&
		s.isMountedAsPwaIcon(fileUpload) &&
		s.isValidImageScalingWidth(widthInt, []int{192, 512}) // Example allowed widths
}

func (s *SendFileUpload) isAvatarSafeForScaling(fileUpload *model.FileUpload, c *gin.Context) bool {
	width := c.Query("width")
	widthInt := 0
	if width != "" {
		fmt.Sscanf(width, "%d", &widthInt)
	}

	return fileUpload.IsImageSafeForScaling() &&
		s.isMountedAsAvatar(fileUpload) &&
		s.isValidImageScalingWidth(widthInt, []int{32, 64, 128}) // Example allowed widths
}

func (s *SendFileUpload) isMountedAsAvatar(fileUpload *model.FileUpload) bool {
	return fileUpload.MountedAs == "avatar"
}

func (s *SendFileUpload) isMountedAsPwaIcon(fileUpload *model.FileUpload) bool {
	return fileUpload.MountedAs == "pwa_icon"
}

func (s *SendFileUpload) isValidImageScalingWidth(width int, allowedScalarWidths []int) bool {
	for _, allowedWidth := range allowedScalarWidths {
		if allowedWidth == width {
			return true
		}
	}
	return false
}
