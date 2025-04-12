package model

import (
	"time"
)

// Upload represents a file upload
type Upload struct {
	ID          int64     `json:"id"`
	ModelType   string    `json:"model_type"`
	ModelID     int64     `json:"model_id"`
	Uploader    string    `json:"uploader"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the database table name for uploads
func (u *Upload) TableName() string {
	return "uploads"
}

// ToMap converts the upload to a map
func (u *Upload) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           u.ID,
		"model_type":   u.ModelType,
		"model_id":     u.ModelID,
		"uploader":     u.Uploader,
		"path":         u.Path,
		"size":         u.Size,
		"content_type": u.ContentType,
		"user_id":      u.UserID,
		"created_at":   u.CreatedAt,
		"updated_at":   u.UpdatedAt,
	}
}

// FileUploader represents a file uploader
type FileUploader struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

// Uploadable is an interface that can be uploaded
type Uploadable interface {
	GetID() int64
	GetType() string
}
