package snippets

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/blob"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/notes"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/pagination"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SnippetsActions provides actions for handling snippets
type SnippetsActions struct {
	rendersNotes      *notes.RendersNotes
	rendersBlob       *blob.RendersBlob
	paginatedCollection *pagination.PaginatedCollection
	noteableMetadata  *service.NoteableMetadataService
	sendBlob          *SendBlob
	snippetsSort      *SnippetsSort
	analyticsTracking *service.AnalyticsTrackingService
	workhorseService  *service.WorkhorseService
}

// NewSnippetsActions creates a new instance of SnippetsActions
func NewSnippetsActions(
	rendersNotes *notes.RendersNotes,
	rendersBlob *blob.RendersBlob,
	paginatedCollection *pagination.PaginatedCollection,
	noteableMetadata *service.NoteableMetadataService,
	sendBlob *SendBlob,
	snippetsSort *SnippetsSort,
	analyticsTracking *service.AnalyticsTrackingService,
	workhorseService *service.WorkhorseService,
) *SnippetsActions {
	return &SnippetsActions{
		rendersNotes:      rendersNotes,
		rendersBlob:       rendersBlob,
		paginatedCollection: paginatedCollection,
		noteableMetadata:  noteableMetadata,
		sendBlob:          sendBlob,
		snippetsSort:      snippetsSort,
		analyticsTracking: analyticsTracking,
		workhorseService:  workhorseService,
	}
}

// RegisterRoutes registers the routes for snippets actions
func (s *SnippetsActions) RegisterRoutes(router *gin.Engine) {
	// Register routes for snippets actions
	// In Go/Gin, we would typically use router groups
	snippetsGroup := router.Group("/snippets")
	{
		snippetsGroup.GET("/:id/edit", s.Edit)
		snippetsGroup.GET("/:id/raw", s.Raw)
		snippetsGroup.GET("/:id", s.Show)
	}
}

// Edit handles the edit action for a snippet
func (s *SnippetsActions) Edit(c *gin.Context) {
	// In Go, we would typically render a template or return JSON
	// This is a placeholder for the actual implementation
	c.HTML(http.StatusOK, "snippets/edit", gin.H{
		"snippet": c.MustGet("snippet"),
	})
}

// Raw handles the raw action for a snippet
// This endpoint is being replaced by Snippets::BlobController#raw
// Support for old raw links will be maintained via this action but
// it will only return the first blob found
func (s *SnippetsActions) Raw(c *gin.Context) {
	snippet := c.MustGet("snippet").(*model.Snippet)
	blob := s.getBlob(snippet)

	// Set content type for workhorse
	s.workhorseService.SetContentType(c)

	// Check if the blob has a snippet (old format)
	if blob.Snippet != nil {
		// Send the data with appropriate headers
		content := s.convertLineEndings(blob.Data, c.Query("line_ending"))
		filename := model.SanitizeFileName(blob.Name)

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.String(http.StatusOK, content)
	} else {
		// Send the snippet blob using the SendBlob service
		s.sendBlob.SendSnippetBlob(c, snippet, blob)
	}
}

// Show handles the show action for a snippet
func (s *SnippetsActions) Show(c *gin.Context) {
	snippet := c.MustGet("snippet").(*model.Snippet)

	// Track analytics event
	s.analyticsTracking.TrackEvent(c, "i_snippets_show", nil)

	// Handle different formats
	format := c.GetHeader("Accept")
	if strings.Contains(format, "application/javascript") || c.Query("format") == "js" {
		s.showJS(c, snippet)
	} else {
		s.showHTML(c, snippet)
	}
}

// IsJSRequest checks if the request is a JavaScript request
func (s *SnippetsActions) IsJSRequest(c *gin.Context) bool {
	format := c.GetHeader("Accept")
	return strings.Contains(format, "application/javascript") || c.Query("format") == "js"
}

// Private helper methods

func (s *SnippetsActions) showHTML(c *gin.Context, snippet *model.Snippet) {
	// Create a new note for the snippet
	note := &model.Note{
		Noteable: snippet,
		Project:  snippet.Project,
	}

	// Get discussions and notes
	discussions := snippet.Discussions
	notes := s.prepareNotesForRendering(discussions)

	// Render the template
	c.HTML(http.StatusOK, "snippets/show", gin.H{
		"snippet":     snippet,
		"note":        note,
		"noteable":    snippet,
		"discussions": discussions,
		"notes":       notes,
	})
}

func (s *SnippetsActions) showJS(c *gin.Context, snippet *model.Snippet) {
	if !snippet.IsEmbeddable() {
		c.Status(http.StatusNotFound)
		return
	}

	// Conditionally expand blobs
	blobs := s.getBlobs(snippet)
	s.conditionallyExpandBlobs(blobs)

	// Render the template
	c.HTML(http.StatusOK, "shared/snippets/show", gin.H{
		"snippet": snippet,
		"blobs":   blobs,
	})
}

func (s *SnippetsActions) getBlob(snippet *model.Snippet) *model.Blob {
	blobs := s.getBlobs(snippet)
	if len(blobs) > 0 {
		return blobs[0]
	}
	return nil
}

func (s *SnippetsActions) getBlobs(snippet *model.Snippet) []*model.Blob {
	if snippet.IsEmptyRepo() {
		return []*model.Blob{snippet.Blob}
	}
	return snippet.Blobs
}

func (s *SnippetsActions) convertLineEndings(content string, lineEnding string) string {
	if lineEnding == "raw" {
		return content
	}
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func (s *SnippetsActions) prepareNotesForRendering(discussions []*model.Discussion) []*model.Note {
	var notes []*model.Note
	for _, discussion := range discussions {
		notes = append(notes, discussion.Notes...)
	}
	return s.rendersNotes.PrepareNotesForRendering(notes)
}

func (s *SnippetsActions) conditionallyExpandBlobs(blobs []*model.Blob) {
	// Implementation would depend on how blobs are expanded
	// This is a placeholder for the actual implementation
}
