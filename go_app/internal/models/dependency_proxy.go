package models

import (
	"mime/multipart"
)

const (
	DistributionAPIVersion = "registry/2.0"
	BlobMaxFileSize       = 5 * 1024 * 1024 * 1024 // 5GB
	ManifestMaxFileSize   = 10 * 1024 * 1024       // 10MB
)

// DependencyProxyBlob represents a cached container blob
type DependencyProxyBlob struct {
	ID        int64
	GroupID   int64
	FileName  string
	Size      int64
	File      *multipart.File
}

// DependencyProxyManifest represents a cached container manifest
type DependencyProxyManifest struct {
	ID          int64
	GroupID     int64
	FileName    string
	Size        int64
	ContentType string
	Digest      string
	File        *multipart.File
	Active      bool
}

// Registry provides URLs for the container registry
type Registry struct{}

// ManifestURL returns the URL for a container manifest
func (r *Registry) ManifestURL(image, tag string) string {
	// TODO: Implement actual registry URL construction
	return ""
}

// BlobURL returns the URL for a container blob
func (r *Registry) BlobURL(image, sha string) string {
	// TODO: Implement actual registry URL construction
	return ""
}

// FileUploader handles file uploads for dependency proxy
type FileUploader struct {
	HasLength    bool
	MaximumSize  int64
}

// WorkhorseAuthorize generates authorization for Workhorse file uploads
func (u *FileUploader) WorkhorseAuthorize() map[string]interface{} {
	return map[string]interface{}{
		"MaximumSize": u.MaximumSize,
		"HasLength":   u.HasLength,
	}
} 