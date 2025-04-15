package jiraconnect

import (
	"html/template"
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/app/models"
)

// BaseHandler provides common functionality for all Jira Connect handlers
type BaseHandler struct {
	templates *template.Template
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler(templates *template.Template) *BaseHandler {
	return &BaseHandler{
		templates: templates,
	}
}

// RenderTemplate renders an HTML template
func (h *BaseHandler) RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetCurrentUser retrieves the current user from the request
func (h *BaseHandler) GetCurrentUser(r *http.Request) *models.User {
	// TODO: Implement user retrieval from session/context
	return nil
} 