package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
)

var (
	ErrInvalidFile        = errors.New("invalid file")
	ErrFileNotFound       = errors.New("file not found")
	ErrInvalidUploader    = errors.New("invalid uploader")
	ErrInvalidSecret      = errors.New("invalid secret")
	ErrPathTraversal      = errors.New("path traversal attack detected")
	ErrUploadNotAuthorized = errors.New("upload not authorized")
)

// UploadService handles business logic for file uploads
type UploadService struct {
	db *sql.DB
	uploadPath string
	maxFileSize int64
}

// NewUploadService creates a new upload service
func NewUploadService(db *sql.DB, uploadPath string, maxFileSize int64) *UploadService {
	return &UploadService{
		db: db,
		uploadPath: uploadPath,
		maxFileSize: maxFileSize,
	}
}

// CreateUpload creates a new file upload
func (s *UploadService) CreateUpload(ctx context.Context, model interface{}, file *multipart.FileHeader, uploaderClass string, userID int64) (*model.Upload, error) {
	// Validate the file
	if file == nil || file.Size == 0 {
		return nil, ErrInvalidFile
	}

	// Check file size
	if file.Size > s.maxFileSize {
		return nil, ErrInvalidFile
	}

	// Generate a unique filename
	filename := s.generateFilename(file.Filename)

	// Save the file
	filePath := filepath.Join(s.uploadPath, filename)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy the file
	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Create the upload record
	upload := &model.Upload{
		ModelType:    s.getModelType(model),
		ModelID:      s.getModelID(model),
		Uploader:     uploaderClass,
		Path:         filename,
		Size:         file.Size,
		ContentType:  file.Header.Get("Content-Type"),
		UserID:       userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save the upload record
	query := `
		INSERT INTO uploads (model_type, model_id, uploader, path, size, content_type, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		RETURNING id
	`

	err = s.db.QueryRowContext(ctx, query,
		upload.ModelType,
		upload.ModelID,
		upload.Uploader,
		upload.Path,
		upload.Size,
		upload.ContentType,
		upload.UserID,
		upload.CreatedAt,
	).Scan(&upload.ID)

	if err != nil {
		return nil, err
	}

	return upload, nil
}

// FileExists checks if a file exists
func (s *UploadService) FileExists(uploader interface{}) bool {
	filePath := s.getFilePath(uploader)
	_, err := os.Stat(filePath)
	return err == nil
}

// GetFileUploader gets a file uploader for a filename
func (s *UploadService) GetFileUploader(uploader interface{}, filename string) (*model.FileUploader, error) {
	filePath := filepath.Join(s.uploadPath, filename)
	if _, err := os.Stat(filePath); err != nil {
		return nil, ErrFileNotFound
	}

	return &model.FileUploader{
		Path: filePath,
		Filename: filename,
	}, nil
}

// GetCacheSettings gets the cache settings for uploads
func (s *UploadService) GetCacheSettings() (int, map[string]string) {
	// Default to no caching
	return 0, map[string]string{
		"expires": time.Now().Add(24 * time.Hour).Format(time.RFC1123),
	}
}

// GetContentType gets the content type for a file
func (s *UploadService) GetContentType(fileUploader *model.FileUploader) string {
	// Default to application/octet-stream
	return "application/octet-stream"
}

// IsEmbeddable checks if a file is embeddable
func (s *UploadService) IsEmbeddable(fileUploader *model.FileUploader) bool {
	// Check if the file is an image
	contentType := s.GetContentType(fileUploader)
	return strings.HasPrefix(contentType, "image/")
}

// IsPDF checks if a file is a PDF
func (s *UploadService) IsPDF(fileUploader *model.FileUploader) bool {
	// Check if the file is a PDF
	contentType := s.GetContentType(fileUploader)
	return contentType == "application/pdf"
}

// AuthorizeUpload authorizes a file upload
func (s *UploadService) AuthorizeUpload(ctx context.Context, uploaderClass string) (map[string]interface{}, error) {
	// Check if the uploader class is valid
	if !s.isValidUploaderClass(uploaderClass) {
		return nil, ErrInvalidUploader
	}

	// Return the authorization data
	return map[string]interface{}{
		"upload_url": "/api/uploads",
		"max_file_size": s.maxFileSize,
	}, nil
}

// generateFilename generates a unique filename
func (s *UploadService) generateFilename(originalFilename string) string {
	// Generate a unique filename
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}

// getFilePath gets the file path for an uploader
func (s *UploadService) getFilePath(uploader interface{}) string {
	switch u := uploader.(type) {
	case *model.Upload:
		return filepath.Join(s.uploadPath, u.Path)
	case *model.FileUploader:
		return u.Path
	default:
		return ""
	}
}

// getModelType gets the model type for a model
func (s *UploadService) getModelType(model interface{}) string {
	switch model.(type) {
	case *model.Project:
		return "Project"
	case *model.User:
		return "User"
	case *model.Group:
		return "Group"
	default:
		return "Unknown"
	}
}

// getModelID gets the model ID for a model
func (s *UploadService) getModelID(model interface{}) int64 {
	switch m := model.(type) {
	case *model.Project:
		return m.ID
	case *model.User:
		return m.ID
	case *model.Group:
		return m.ID
	default:
		return 0
	}
}

// isValidUploaderClass checks if an uploader class is valid
func (s *UploadService) isValidUploaderClass(uploaderClass string) bool {
	validUploaders := []string{
		"AvatarUploader",
		"AttachmentUploader",
		"FileUploader",
		"LogoUploader",
		"PwaIconUploader",
		"HeaderLogoUploader",
		"FaviconUploader",
		"ScreenshotUploader",
	}

	for _, valid := range validUploaders {
		if uploaderClass == valid {
			return true
		}
	}

	return false
}

// HandleUpload processes the file upload for a given model
func (s *UploadService) HandleUpload(ctx *gin.Context, model interface{}) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	// Determine uploader based on model type
	uploader, err := s.getUploader(model)
	if err != nil {
		return err
	}

	// Process the upload
	return uploader.Upload(ctx, file)
}

// getUploader returns the appropriate uploader for the model type
func (s *UploadService) getUploader(model interface{}) (Uploader, error) {
	switch model.(type) {
	case *model.Project:
		return &FileUploader{config: s.config}, nil
	case *model.Group:
		return &NamespaceFileUploader{config: s.config}, nil
	default:
		return nil, fmt.Errorf("invalid model type")
	}
}

// GetProject retrieves a project by ID
func (s *UploadService) GetProject(ctx context.Context, id string) (*model.Project, error) {
	// TODO: Implement project retrieval from database
	return nil, fmt.Errorf("not implemented")
}

// GetGroup retrieves a group by ID
func (s *UploadService) GetGroup(ctx context.Context, id string) (*model.Group, error) {
	// TODO: Implement group retrieval from database
	return nil, fmt.Errorf("not implemented")
}

// ShouldBypassAuthChecks checks if authentication should be bypassed
func (s *UploadService) ShouldBypassAuthChecks() bool {
	return s.config.BypassAuthChecks
}

// Uploader interface defines the contract for file uploaders
type Uploader interface {
	Upload(ctx *gin.Context, file *multipart.FileHeader) error
}

// FileUploader handles project file uploads
type FileUploader struct {
	config *Config
}

// Upload implements the Uploader interface for project files
func (u *FileUploader) Upload(ctx *gin.Context, file *multipart.FileHeader) error {
	// TODO: Implement project file upload logic
	return fmt.Errorf("not implemented")
}

// NamespaceFileUploader handles group file uploads
type NamespaceFileUploader struct {
	config *Config
}

// Upload implements the Uploader interface for group files
func (u *NamespaceFileUploader) Upload(ctx *gin.Context, file *multipart.FileHeader) error {
	// TODO: Implement group file upload logic
	return fmt.Errorf("not implemented")
}

// Config holds configuration for the UploadService
type Config struct {
	UploadPath        string
	BypassAuthChecks  bool
	MaxFileSize       int64
	AllowedFileTypes []string
}
