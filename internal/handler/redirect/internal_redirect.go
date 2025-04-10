package redirect

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// InternalRedirect handles safe redirects within the application
type InternalRedirect struct {
	allowedHosts []string
}

// NewInternalRedirect creates a new instance of InternalRedirect
func NewInternalRedirect(allowedHosts []string) *InternalRedirect {
	return &InternalRedirect{
		allowedHosts: allowedHosts,
	}
}

// SafeRedirectPath checks if a path is safe to redirect to
func (i *InternalRedirect) SafeRedirectPath(path string) string {
	if path == "" {
		return ""
	}

	// Verify that the string starts with a '/' and a known route character
	validPathRegex := regexp.MustCompile(`^/[-\w].*$`)
	if !validPathRegex.MatchString(path) {
		return ""
	}

	// Parse the URI
	uri, err := url.Parse(path)
	if err != nil {
		return ""
	}

	// Return the full path for the URI
	return i.FullPathForURI(uri)
}

// SafeRedirectPathForURL checks if a URL is safe to redirect to
func (i *InternalRedirect) SafeRedirectPathForURL(urlStr string, ctx *gin.Context) string {
	if urlStr == "" {
		return ""
	}

	// Parse the URL
	uri, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	// Check if the host is allowed
	if i.HostAllowed(uri, ctx) {
		return i.SafeRedirectPath(i.FullPathForURI(uri))
	}

	return ""
}

// SanitizeRedirect sanitizes a redirect URL or path
func (i *InternalRedirect) SanitizeRedirect(urlOrPath string, ctx *gin.Context) string {
	// Try to sanitize as a path first
	safePath := i.SafeRedirectPath(urlOrPath)
	if safePath != "" {
		return safePath
	}

	// If not a safe path, try to sanitize as a URL
	return i.SafeRedirectPathForURL(urlOrPath, ctx)
}

// HostAllowed checks if a URI's host is allowed
func (i *InternalRedirect) HostAllowed(uri *url.URL, ctx *gin.Context) bool {
	// Get the request host and port
	requestHost := ctx.Request.Host
	hostParts := strings.Split(requestHost, ":")
	requestHostName := hostParts[0]
	requestPort := "80" // Default HTTP port
	if len(hostParts) > 1 {
		requestPort = hostParts[1]
	}

	// Check if the URI host matches the request host
	if uri.Host != requestHostName {
		return false
	}

	// Check if the URI port matches the request port
	uriPort := uri.Port()
	if uriPort == "" {
		// If no port is specified, use the default port based on the scheme
		if uri.Scheme == "https" {
			uriPort = "443"
		} else {
			uriPort = "80"
		}
	}

	return uriPort == requestPort
}

// FullPathForURI returns the full path for a URI
func (i *InternalRedirect) FullPathForURI(uri *url.URL) string {
	// Combine the path and query
	pathWithQuery := uri.Path
	if uri.RawQuery != "" {
		pathWithQuery += "?" + uri.RawQuery
	}

	// Add the fragment if present
	if uri.Fragment != "" {
		pathWithQuery += "#" + uri.Fragment
	}

	return pathWithQuery
}

// RefererPath returns the path from the referer header
func (i *InternalRedirect) RefererPath(ctx *gin.Context) string {
	referer := ctx.Request.Referer()
	if referer == "" {
		return ""
	}

	uri, err := url.Parse(referer)
	if err != nil {
		return ""
	}

	return uri.Path
}
