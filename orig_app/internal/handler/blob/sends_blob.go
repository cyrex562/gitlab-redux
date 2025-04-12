package blob

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SendsBlob handles sending blob data to clients
type SendsBlob struct {
	blobHelper     *BlobHelper
	sendFileUpload *SendFileUpload
	lfsService     *service.LFSService
}

// NewSendsBlob creates a new instance of SendsBlob
func NewSendsBlob(
	blobHelper *BlobHelper,
	sendFileUpload *SendFileUpload,
	lfsService *service.LFSService,
) *SendsBlob {
	return &SendsBlob{
		blobHelper:     blobHelper,
		sendFileUpload: sendFileUpload,
		lfsService:     lfsService,
	}
}

// SendBlob sends a blob to the client
func (s *SendsBlob) SendBlob(
	c *gin.Context,
	repository *model.Repository,
	blob *model.Blob,
	inline bool,
	allowCaching bool,
) {
	if blob == nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Set security header
	c.Header("X-Content-Type-Options", "nosniff")

	// Check if blob is already cached
	if s.isCachedBlob(c, blob, allowCaching) {
		return
	}

	// Handle different blob types
	if blob.IsStoredExternally() {
		s.sendLFSObject(c, blob, repository.Project)
	} else {
		s.sendGitBlob(c, repository, blob, inline)
	}
}

// IsCachedBlob checks if the blob is already cached
func (s *SendsBlob) IsCachedBlob(
	c *gin.Context,
	blob *model.Blob,
	allowCaching bool,
) bool {
	// Check if the response is stale
	stale := s.isStale(c, blob.ID)

	// Determine cache time based on ref and commit
	maxAge := s.determineCacheTime(c)

	// Set cache headers
	s.setCacheHeaders(c, maxAge, allowCaching)

	return !stale
}

// SendLFSObject sends an LFS object to the client
func (s *SendsBlob) SendLFSObject(
	c *gin.Context,
	blob *model.Blob,
	project *model.Project,
) {
	lfsObject := s.findLFSObject(blob)

	if lfsObject != nil && s.lfsService.IsProjectAllowedAccess(lfsObject, project) {
		s.sendFileUpload.SendUpload(
			c,
			lfsObject.File,
			nil,
			nil,
			blob.Name,
			false,
			"attachment",
		)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// SendGitBlob sends a Git blob to the client
func (s *SendsBlob) SendGitBlob(
	c *gin.Context,
	repository *model.Repository,
	blob *model.Blob,
	inline bool,
) {
	// Implementation depends on how Git blobs are stored and served
	// This would typically involve reading the blob from the repository
	// and sending it to the client with appropriate headers
	// ...
}

// FindLFSObject finds an LFS object by its OID
func (s *SendsBlob) FindLFSObject(blob *model.Blob) *model.LFSObject {
	lfsObject := s.lfsService.FindByOID(blob.LFSOID)
	if lfsObject != nil && s.lfsService.FileExists(lfsObject) {
		return lfsObject
	}
	return nil
}

// Private helper methods

func (s *SendsBlob) isCachedBlob(
	c *gin.Context,
	blob *model.Blob,
	allowCaching bool,
) bool {
	return s.IsCachedBlob(c, blob, allowCaching)
}

func (s *SendsBlob) isStale(c *gin.Context, blobID string) bool {
	// Check if the response is stale based on the blob ID
	// This would typically involve checking the If-None-Match header
	// against the blob ID
	// ...
	return false
}

func (s *SendsBlob) determineCacheTime(c *gin.Context) time.Duration {
	// Get ref and commit from context
	ref := c.GetString("ref")
	commit := c.GetString("commit")

	// If ref and commit match, the blob is immutable
	if ref != "" && commit != "" && ref == commit {
		return model.BlobCacheTimeImmutable
	}

	// Otherwise, the blob may change over time
	return model.BlobCacheTime
}

func (s *SendsBlob) setCacheHeaders(
	c *gin.Context,
	maxAge time.Duration,
	allowCaching bool,
) {
	// Set cache headers
	c.Header("Cache-Control", s.buildCacheControlHeader(maxAge, allowCaching))
	c.Header("Expires", time.Now().Add(maxAge).Format(time.RFC1123))
}

func (s *SendsBlob) buildCacheControlHeader(
	maxAge time.Duration,
	allowCaching bool,
) string {
	// Build Cache-Control header
	// ...
	return ""
}

func (s *SendsBlob) sendLFSObject(
	c *gin.Context,
	blob *model.Blob,
	project *model.Project,
) {
	s.SendLFSObject(c, blob, project)
}

func (s *SendsBlob) sendGitBlob(
	c *gin.Context,
	repository *model.Repository,
	blob *model.Blob,
	inline bool,
) {
	s.SendGitBlob(c, repository, blob, inline)
}

func (s *SendsBlob) findLFSObject(blob *model.Blob) *model.LFSObject {
	return s.FindLFSObject(blob)
}
