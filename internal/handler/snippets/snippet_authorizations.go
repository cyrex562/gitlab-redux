package snippets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SnippetAuthorizations provides authorization methods for snippets
type SnippetAuthorizations struct {
	authService *service.AuthService
}

// NewSnippetAuthorizations creates a new instance of SnippetAuthorizations
func NewSnippetAuthorizations(authService *service.AuthService) *SnippetAuthorizations {
	return &SnippetAuthorizations{
		authService: authService,
	}
}

// AuthorizeReadSnippet checks if the current user can read the snippet
// Returns true if authorized, false otherwise
func (s *SnippetAuthorizations) AuthorizeReadSnippet(c *gin.Context, snippet *model.Snippet) bool {
	currentUser := s.getCurrentUser(c)
	if currentUser == nil {
		c.Status(http.StatusNotFound)
		return false
	}

	if !s.authService.Can(currentUser, "read_snippet", snippet) {
		c.Status(http.StatusNotFound)
		return false
	}

	return true
}

// AuthorizeUpdateSnippet checks if the current user can update the snippet
// Returns true if authorized, false otherwise
func (s *SnippetAuthorizations) AuthorizeUpdateSnippet(c *gin.Context, snippet *model.Snippet) bool {
	currentUser := s.getCurrentUser(c)
	if currentUser == nil {
		c.Status(http.StatusNotFound)
		return false
	}

	if !s.authService.Can(currentUser, "update_snippet", snippet) {
		c.Status(http.StatusNotFound)
		return false
	}

	return true
}

// AuthorizeAdminSnippet checks if the current user can admin the snippet
// Returns true if authorized, false otherwise
func (s *SnippetAuthorizations) AuthorizeAdminSnippet(c *gin.Context, snippet *model.Snippet) bool {
	currentUser := s.getCurrentUser(c)
	if currentUser == nil {
		c.Status(http.StatusNotFound)
		return false
	}

	if !s.authService.Can(currentUser, "admin_snippet", snippet) {
		c.Status(http.StatusNotFound)
		return false
	}

	return true
}

// AuthorizeCreateSnippet checks if the current user can create a snippet
// Returns true if authorized, false otherwise
func (s *SnippetAuthorizations) AuthorizeCreateSnippet(c *gin.Context) bool {
	currentUser := s.getCurrentUser(c)
	if currentUser == nil {
		c.Status(http.StatusNotFound)
		return false
	}

	if !s.authService.Can(currentUser, "create_snippet", nil) {
		c.Status(http.StatusNotFound)
		return false
	}

	return true
}

// Private helper methods

func (s *SnippetAuthorizations) getCurrentUser(c *gin.Context) *model.User {
	// Get the current user from the context
	user, exists := c.Get("current_user")
	if !exists {
		return nil
	}
	return user.(*model.User)
}
