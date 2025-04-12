package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/gitlab/internal/middleware"
	"gitlab.com/gitlab-org/gitlab/internal/models"
	"gitlab.com/gitlab-org/gitlab/internal/services"
)

const (
	sendDependencyContentTypeHeader = "X-GitLab-Send-Dependency-Content-Type"
	manifestDigestHeader           = "Docker-Content-Digest"
)

// DependencyProxyContainersController handles container registry proxy requests
type DependencyProxyContainersController struct {
	group *models.Group
	token string
}

// NewDependencyProxyContainersController creates a new containers controller
func NewDependencyProxyContainersController(group *models.Group) *DependencyProxyContainersController {
	return &DependencyProxyContainersController{
		group: group,
	}
}

// RegisterRoutes registers the container proxy routes
func (c *DependencyProxyContainersController) RegisterRoutes(r *http.ServeMux) {
	// Manifest routes
	r.HandleFunc("/groups/"+c.group.Path+"/dependency_proxy/containers/manifests", middleware.Chain(
		c.verifyDependencyProxyEnabled,
		c.ensureTokenGranted,
	)(c.handleManifest))

	// Blob routes
	r.HandleFunc("/groups/"+c.group.Path+"/dependency_proxy/containers/blobs", middleware.Chain(
		c.verifyDependencyProxyEnabled,
		c.ensureTokenGranted,
	)(c.handleBlob))

	// Upload routes
	r.HandleFunc("/groups/"+c.group.Path+"/dependency_proxy/containers/blobs/uploads", middleware.Chain(
		c.verifyDependencyProxyEnabled,
		c.verifyWorkhorseAPI,
	)(c.handleBlobUpload))

	r.HandleFunc("/groups/"+c.group.Path+"/dependency_proxy/containers/manifests/uploads", middleware.Chain(
		c.verifyDependencyProxyEnabled,
		c.verifyWorkhorseAPI,
	)(c.handleManifestUpload))
}

// verifyDependencyProxyEnabled middleware checks if dependency proxy is enabled
func (c *DependencyProxyContainersController) verifyDependencyProxyEnabled(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setting := c.getDependencyProxySetting()
		if setting == nil || !setting.Enabled {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next(w, r)
	}
}

// ensureTokenGranted middleware ensures a valid token is present
func (c *DependencyProxyContainersController) ensureTokenGranted(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image := r.URL.Query().Get("image")
		result := services.NewRequestTokenService(image, c.getDependencyProxySetting()).Execute()

		if result.Status == "success" {
			c.token = result.Token
			next(w, r)
		} else {
			http.Error(w, result.Message, result.HTTPStatus)
		}
	}
}

// verifyWorkhorseAPI middleware verifies Workhorse API requests
func (c *DependencyProxyContainersController) verifyWorkhorseAPI(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement Workhorse API verification
		next(w, r)
	}
}

func (c *DependencyProxyContainersController) handleManifest(w http.ResponseWriter, r *http.Request) {
	image := r.URL.Query().Get("image")
	tag := r.URL.Query().Get("tag")

	result := services.NewFindCachedManifestService(c.group, image, tag, c.token).Execute()

	if result.Status == "success" {
		if result.Manifest != nil {
			c.sendManifest(w, r, result.Manifest, result.FromCache)
		} else {
			// TODO: Implement proxy to registry
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	} else {
		http.Error(w, result.Message, result.HTTPStatus)
	}
}

func (c *DependencyProxyContainersController) handleBlob(w http.ResponseWriter, r *http.Request) {
	sha := r.URL.Query().Get("sha")
	fileName := strings.TrimPrefix(sha, "sha256:") + ".gz"

	// Placeholder for blob lookup
	_ = c.group.DependencyProxySetting.GroupID // Using GroupID to demonstrate we'll need it for the lookup
	_ = fileName // Will be used to look up the blob

	// TODO: Implement blob handling - look up blob by fileName in database
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (c *DependencyProxyContainersController) handleBlobUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.authorizeUploadBlob(w, r)
	} else if r.Method == http.MethodPut {
		c.uploadBlob(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *DependencyProxyContainersController) handleManifestUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.authorizeUploadManifest(w, r)
	} else if r.Method == http.MethodPut {
		c.uploadManifest(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *DependencyProxyContainersController) authorizeUploadBlob(w http.ResponseWriter, r *http.Request) {
	uploader := &models.FileUploader{
		HasLength:   false,
		MaximumSize: models.BlobMaxFileSize,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploader.WorkhorseAuthorize())
}

func (c *DependencyProxyContainersController) authorizeUploadManifest(w http.ResponseWriter, r *http.Request) {
	uploader := &models.FileUploader{
		HasLength:   false,
		MaximumSize: models.ManifestMaxFileSize,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploader.WorkhorseAuthorize())
}

func (c *DependencyProxyContainersController) uploadBlob(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement blob upload
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (c *DependencyProxyContainersController) uploadManifest(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement manifest upload
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (c *DependencyProxyContainersController) sendManifest(w http.ResponseWriter, r *http.Request, manifest *models.DependencyProxyManifest, fromCache bool) {
	w.Header().Set(manifestDigestHeader, manifest.Digest)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", manifest.Size))
	w.Header().Set("Docker-Distribution-Api-Version", models.DistributionAPIVersion)
	w.Header().Set("Etag", fmt.Sprintf(`"%s"`, manifest.Digest))
	w.Header().Set("Content-Type", manifest.ContentType)

	// TODO: Implement actual file sending and event tracking
	w.WriteHeader(http.StatusOK)
}

func (c *DependencyProxyContainersController) getDependencyProxySetting() *models.DependencyProxySetting {
	return c.group.DependencyProxySetting
} 