package snippets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// SendBlob handles sending snippet blobs
type SendBlob struct {
	snippetService *service.SnippetService
	logger        *util.Logger
}

// NewSendBlob creates a new instance of SendBlob
func NewSendBlob(snippetService *service.SnippetService, logger *util.Logger) *SendBlob {
	return &SendBlob{
		snippetService: snippetService,
		logger:        logger,
	}
}

// SendSnippetBlob sends a blob from a snippet
func (s *SendBlob) SendSnippetBlob(ctx *gin.Context, snippet *model.Snippet, blob *model.Blob) {
	// Set workhorse content type
	s.setWorkhorseContentType(ctx)

	// Get content disposition
	inline := s.getContentDisposition(ctx) == "inline"

	// Check if anonymous users can cache the blob
	allowCaching := s.canAnonymousCacheBlob(ctx, snippet)

	// Send the blob
	s.sendBlob(ctx, snippet, blob, inline, allowCaching)
}

// setWorkhorseContentType sets the workhorse content type header
func (s *SendBlob) setWorkhorseContentType(ctx *gin.Context) {
	ctx.Header("Gitlab-Workhorse-Send-Data", "git-blob")
}

// getContentDisposition determines the content disposition based on the inline parameter
func (s *SendBlob) getContentDisposition(ctx *gin.Context) string {
	if ctx.Query("inline") == "false" {
		return "attachment"
	}
	return "inline"
}

// canAnonymousCacheBlob checks if anonymous users can cache the blob
func (s *SendBlob) canAnonymousCacheBlob(ctx *gin.Context, snippet *model.Snippet) bool {
	// TODO: Implement anonymous user permission check
	// This should check if anonymous users have permission to cache blobs for this snippet
	return false
}

// sendBlob sends the blob with the specified parameters
func (s *SendBlob) sendBlob(ctx *gin.Context, snippet *model.Snippet, blob *model.Blob, inline bool, allowCaching bool) {
	// Set content type
	ctx.Header("Content-Type", blob.GetContentType())

	// Set content disposition
	disposition := "attachment"
	if inline {
		disposition = "inline"
	}
	ctx.Header("Content-Disposition", disposition)

	// Set cache control headers
	if allowCaching {
		ctx.Header("Cache-Control", "public")
	} else {
		ctx.Header("Cache-Control", "private")
	}

	// Write blob content
	ctx.Data(http.StatusOK, blob.GetContentType(), blob.GetContent())
}
