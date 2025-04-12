package security

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/config"
)

// WebIdeCSPHandler handles Content Security Policy for Web IDE
type WebIdeCSPHandler struct {
	config *config.Config
}

// NewWebIdeCSPHandler creates a new Web IDE CSP handler
func NewWebIdeCSPHandler(config *config.Config) *WebIdeCSPHandler {
	return &WebIdeCSPHandler{
		config: config,
	}
}

// IncludeWebIdeCSP adds Web IDE CSP directives to the response
func (h *WebIdeCSPHandler) IncludeWebIdeCSP(c *gin.Context) {
	// Get existing CSP directives
	csp := c.GetHeader("Content-Security-Policy")
	if csp == "" {
		return
	}

	// Parse the base URL
	baseURL, err := url.Parse(c.Request.URL.String())
	if err != nil {
		return
	}

	// Set the base path
	baseURL.Path = h.config.GitLab.RelativeURLRoot
	if baseURL.Path == "" {
		baseURL.Path = "/"
	}

	// Add webpack path
	baseURL.Path = strings.TrimRight(baseURL.Path, "/") + "/assets/webpack/"
	baseURL.RawQuery = "" // Remove query parameters

	// Get default-src directive
	defaultSrc := h.getDirectiveValues(csp, "default-src")
	if defaultSrc == "" {
		return
	}

	// Update frame-src directive
	frameSrc := h.getDirectiveValues(csp, "frame-src")
	if frameSrc == "" {
		frameSrc = defaultSrc
	}
	frameSrc += " " + baseURL.String() + " https://*.web-ide.gitlab-static.net/"
	csp = h.updateDirective(csp, "frame-src", frameSrc)

	// Update worker-src directive
	workerSrc := h.getDirectiveValues(csp, "worker-src")
	if workerSrc == "" {
		workerSrc = defaultSrc
	}
	workerSrc += " " + baseURL.String()
	csp = h.updateDirective(csp, "worker-src", workerSrc)

	// Set the updated CSP header
	c.Header("Content-Security-Policy", csp)
}

// getDirectiveValues extracts values from a CSP directive
func (h *WebIdeCSPHandler) getDirectiveValues(csp, directive string) string {
	parts := strings.Split(csp, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, directive+" ") {
			return strings.TrimPrefix(part, directive+" ")
		}
	}
	return ""
}

// updateDirective updates a CSP directive with new values
func (h *WebIdeCSPHandler) updateDirective(csp, directive, values string) string {
	parts := strings.Split(csp, ";")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, directive+" ") {
			parts[i] = directive + " " + values
			return strings.Join(parts, "; ")
		}
	}
	return csp + "; " + directive + " " + values
}
