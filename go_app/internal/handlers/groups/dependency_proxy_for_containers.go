package groups

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/middleware"
	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/go_app/internal/services/dependency_proxy"
)

const (
	// DistributionAPIVersion is the Docker Distribution API version
	DistributionAPIVersion = "registry/2.0"
	// ManifestDigestHeader is the header for manifest digest
	ManifestDigestHeader = "Docker-Content-Digest"
	// DependencyContentTypeHeader is the header for dependency content type
	DependencyContentTypeHeader = "Gitlab-Workhorse-Dependency-Content-Type"
)

// DependencyProxyForContainersHandler handles container dependency proxy operations
type DependencyProxyForContainersHandler struct {
	*DependencyProxiesHandler
	registryService *dependency_proxy.RegistryService
	fileUploader   *dependency_proxy.FileUploader
}

// NewDependencyProxyForContainersHandler creates a new dependency proxy for containers handler
func NewDependencyProxyForContainersHandler(
	baseHandler *DependencyProxiesHandler,
	registryService *dependency_proxy.RegistryService,
	fileUploader *dependency_proxy.FileUploader,
) *DependencyProxyForContainersHandler {
	return &DependencyProxyForContainersHandler{
		DependencyProxiesHandler: baseHandler,
		registryService:         registryService,
		fileUploader:           fileUploader,
	}
}

// RegisterRoutes registers the routes for the dependency proxy for containers handler
func (h *DependencyProxyForContainersHandler) RegisterRoutes(router *middleware.Router) {
	router.HandleFunc("/groups/:group_id/dependency_proxy/containers/manifest", h.Manifest).Methods("GET")
	router.HandleFunc("/groups/:group_id/dependency_proxy/containers/blob", h.Blob).Methods("GET")
	router.HandleFunc("/groups/:group_id/dependency_proxy/containers/authorize_upload_blob", h.AuthorizeUploadBlob).Methods("POST")
	router.HandleFunc("/groups/:group_id/dependency_proxy/containers/upload_blob", h.UploadBlob).Methods("POST")
	router.HandleFunc("/groups/:group_id/dependency_proxy/containers/authorize_upload_manifest", h.AuthorizeUploadManifest).Methods("POST")
	router.HandleFunc("/groups/:group_id/dependency_proxy/containers/upload_manifest", h.UploadManifest).Methods("POST")
}

// Manifest handles manifest requests
func (h *DependencyProxyForContainersHandler) Manifest(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	image := r.URL.Query().Get("image")
	tag := r.URL.Query().Get("tag")

	token, err := h.ensureTokenGranted(group, image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	result, err := h.registryService.FindCachedManifest(group, image, tag, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.Manifest != nil {
		h.sendManifest(w, r, result.Manifest, result.FromCache)
	} else {
		h.sendDependency(w, r, h.manifestHeader(token), h.registryService.ManifestURL(image, tag), h.manifestFileName(image, tag))
	}
}

// Blob handles blob requests
func (h *DependencyProxyForContainersHandler) Blob(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	sha := r.URL.Query().Get("sha")
	blobFileName := h.blobFileName(sha)

	blob, err := group.FindDependencyProxyBlobByFileName(blobFileName)
	if err == nil && blob != nil {
		// Track event
		eventName := h.trackingEventName("blob", true)
		h.trackPackageEvent(eventName, "dependency_proxy", group, h.GetUserFromContext(r))

		h.sendUpload(w, r, blob.File)
	} else {
		token, err := h.ensureTokenGranted(group, r.URL.Query().Get("image"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		h.sendDependency(w, r, h.tokenHeader(token), h.registryService.BlobURL(r.URL.Query().Get("image"), sha), blobFileName)
	}
}

// AuthorizeUploadBlob handles blob upload authorization
func (h *DependencyProxyForContainersHandler) AuthorizeUploadBlob(w http.ResponseWriter, r *http.Request) {
	h.setWorkhorseInternalAPIContentType(w)
	
	auth, err := h.fileUploader.WorkhorseAuthorize(false, models.MaxBlobFileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(auth)
}

// UploadBlob handles blob uploads
func (h *DependencyProxyForContainersHandler) UploadBlob(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	sha := r.FormValue("sha")
	blobFileName := h.blobFileName(sha)

	blob := &models.DependencyProxyBlob{
		FileName: blobFileName,
		Size:     header.Size,
	}

	if err := blob.SaveFile(file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	if err := group.CreateDependencyProxyBlob(blob); err != nil {
		http.Error(w, "Failed to create blob", http.StatusInternalServerError)
		return
	}

	// Track event
	eventName := h.trackingEventName("blob", false)
	h.trackPackageEvent(eventName, "dependency_proxy", group, h.GetUserFromContext(r))

	w.WriteHeader(http.StatusOK)
}

// AuthorizeUploadManifest handles manifest upload authorization
func (h *DependencyProxyForContainersHandler) AuthorizeUploadManifest(w http.ResponseWriter, r *http.Request) {
	h.setWorkhorseInternalAPIContentType(w)
	
	auth, err := h.fileUploader.WorkhorseAuthorize(false, models.MaxManifestFileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(auth)
}

// UploadManifest handles manifest uploads
func (h *DependencyProxyForContainersHandler) UploadManifest(w http.ResponseWriter, r *http.Request) {
	group := h.GetGroupFromContext(r)
	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	image := r.FormValue("image")
	tag := r.FormValue("tag")
	manifestFileName := h.manifestFileName(image, tag)

	manifest := &models.DependencyProxyManifest{
		FileName:    manifestFileName,
		ContentType: r.Header.Get(DependencyContentTypeHeader),
		Digest:      r.Header.Get(ManifestDigestHeader),
		Size:        header.Size,
	}

	if err := manifest.SaveFile(file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	existingManifest, err := group.FindActiveDependencyProxyManifestByFileName(manifestFileName)
	if err == nil && existingManifest != nil {
		if err := existingManifest.Update(manifest); err != nil {
			http.Error(w, "Failed to update manifest", http.StatusInternalServerError)
			return
		}
	} else {
		if err := group.CreateDependencyProxyManifest(manifest); err != nil {
			http.Error(w, "Failed to create manifest", http.StatusInternalServerError)
			return
		}
	}

	// Track event
	eventName := h.trackingEventName("manifest", false)
	h.trackPackageEvent(eventName, "dependency_proxy", group, h.GetUserFromContext(r))

	w.WriteHeader(http.StatusOK)
}

// Helper methods

func (h *DependencyProxyForContainersHandler) sendManifest(w http.ResponseWriter, r *http.Request, manifest *models.DependencyProxyManifest, fromCache bool) {
	w.Header().Set(ManifestDigestHeader, manifest.Digest)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", manifest.Size))
	w.Header().Set("Docker-Distribution-Api-Version", DistributionAPIVersion)
	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", manifest.Digest))

	// Track event
	eventName := h.trackingEventName("manifest", fromCache)
	h.trackPackageEvent(eventName, "dependency_proxy", h.GetGroupFromContext(r), h.GetUserFromContext(r))

	h.sendUpload(w, r, manifest.File)
}

func (h *DependencyProxyForContainersHandler) blobFileName(sha string) string {
	return fmt.Sprintf("%s.gz", strings.TrimPrefix(sha, "sha256:"))
}

func (h *DependencyProxyForContainersHandler) manifestFileName(image, tag string) string {
	// TODO: Implement path traversal check
	return fmt.Sprintf("%s:%s.json", image, tag)
}

func (h *DependencyProxyForContainersHandler) ensureTokenGranted(group *models.Group, image string) (string, error) {
	token, err := h.registryService.RequestToken(image, group.DependencyProxySetting)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (h *DependencyProxyForContainersHandler) tokenHeader(token string) map[string][]string {
	return map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}
}

func (h *DependencyProxyForContainersHandler) manifestHeader(token string) map[string][]string {
	headers := h.tokenHeader(token)
	headers["Accept"] = models.AcceptedManifestTypes
	return headers
}

func (h *DependencyProxyForContainersHandler) trackingEventName(objectType string, fromCache bool) string {
	eventName := fmt.Sprintf("pull_%s", objectType)
	if fromCache {
		eventName = fmt.Sprintf("%s_from_cache", eventName)
	}
	return eventName
}

func (h *DependencyProxyForContainersHandler) setWorkhorseInternalAPIContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (h *DependencyProxyForContainersHandler) sendUpload(w http.ResponseWriter, r *http.Request, file io.Reader) {
	// TODO: Implement file upload sending logic
	io.Copy(w, file)
}

func (h *DependencyProxyForContainersHandler) sendDependency(w http.ResponseWriter, r *http.Request, headers map[string][]string, url string, fileName string) {
	// TODO: Implement dependency sending logic
	// This would typically involve making a request to the external registry
	// and streaming the response back to the client
} 