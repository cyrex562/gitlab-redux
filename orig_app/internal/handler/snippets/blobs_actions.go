package snippets

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

var (
	// ErrNoRepository represents a missing repository error
	ErrNoRepository = errors.New("snippet repository does not exist")
	// ErrNoBlob represents a missing blob error
	ErrNoBlob = errors.New("blob not found")
)

// BlobsActions provides functionality for handling snippet blob actions
type BlobsActions struct {
	snippetService *service.SnippetService
	logger        *util.Logger
}

// NewBlobsActions creates a new instance of BlobsActions
func NewBlobsActions(snippetService *service.SnippetService, logger *util.Logger) *BlobsActions {
	return &BlobsActions{
		snippetService: snippetService,
		logger:        logger,
	}
}

// RegisterRoutes registers the routes for blob actions
func (b *BlobsActions) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/raw", b.authorizeReadSnippet, b.ensureRepository, b.ensureBlob, b.raw)
}

// raw handles the GET /raw endpoint for retrieving raw blob content
func (b *BlobsActions) raw(ctx *gin.Context) {
	snippet := b.getSnippet(ctx)
	blob := b.getBlob(ctx)

	// Send the blob
	b.sendSnippetBlob(ctx, snippet, blob)
}

// authorizeReadSnippet middleware ensures the user can read the snippet
func (b *BlobsActions) authorizeReadSnippet(ctx *gin.Context) {
	snippet := b.getSnippet(ctx)
	if !b.snippetService.CanRead(ctx, snippet) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Next()
}

// ensureRepository middleware ensures the snippet repository exists
func (b *BlobsActions) ensureRepository(ctx *gin.Context) {
	snippet := b.getSnippet(ctx)
	if !snippet.RepoExists() {
		b.logger.Error("Snippet raw blob attempt with no repo", map[string]interface{}{
			"snippet_id": snippet.ID,
		})
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	ctx.Next()
}

// ensureBlob middleware ensures the blob exists
func (b *BlobsActions) ensureBlob(ctx *gin.Context) {
	blob := b.getBlob(ctx)
	if blob == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.Next()
}

// getSnippet gets the snippet from the context
func (b *BlobsActions) getSnippet(ctx *gin.Context) *model.Snippet {
	snippet, exists := ctx.Get("snippet")
	if !exists {
		return nil
	}
	return snippet.(*model.Snippet)
}

// getBlob gets the blob from the context or loads it if not present
func (b *BlobsActions) getBlob(ctx *gin.Context) *model.Blob {
	// Check if blob is already in context
	if blob, exists := ctx.Get("blob"); exists {
		return blob.(*model.Blob)
	}

	// Get snippet and ref parameters
	snippet := b.getSnippet(ctx)
	params := &model.RefParams{
		ID:      ctx.Param("id"),
		Ref:     ctx.Query("ref"),
		Path:    ctx.Query("path"),
		RefType: ctx.Query("ref_type"),
	}

	// Extract ref and get blob
	commit, path, err := b.snippetService.ExtractRef(ctx, snippet, params)
	if err != nil {
		return nil
	}

	// Get blob from repository
	blob, err := b.snippetService.GetBlobAt(ctx, snippet, commit.ID, path)
	if err != nil {
		return nil
	}

	// Store blob in context
	ctx.Set("blob", blob)
	return blob
}

// sendSnippetBlob sends the blob content as response
func (b *BlobsActions) sendSnippetBlob(ctx *gin.Context, snippet *model.Snippet, blob *model.Blob) {
	// Set content type header
	ctx.Header("Content-Type", blob.GetContentType())

	// Set content disposition header
	ctx.Header("Content-Disposition", "inline")

	// Set cache control headers
	ctx.Header("Cache-Control", "private")

	// Write blob content
	ctx.Data(http.StatusOK, blob.GetContentType(), blob.GetContent())
}
