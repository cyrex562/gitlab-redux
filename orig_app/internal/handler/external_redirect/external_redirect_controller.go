package external_redirect

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// ExternalRedirectController handles external redirects
type ExternalRedirectController struct {
	configService service.ConfigService
	urlService    service.URLService
}

// NewExternalRedirectController creates a new ExternalRedirectController
func NewExternalRedirectController(
	configService service.ConfigService,
	urlService service.URLService,
) *ExternalRedirectController {
	return &ExternalRedirectController{
		configService: configService,
		urlService:    urlService,
	}
}

// RegisterRoutes registers the routes for the ExternalRedirectController
func (c *ExternalRedirectController) RegisterRoutes(router *gin.Engine) {
	router.GET("/external_redirect", c.Index)
}

// Index handles the index action
func (c *ExternalRedirectController) Index(ctx *gin.Context) {
	// Check URL parameter
	if !c.checkURLParam(ctx) {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Get URL parameter
	urlParam := c.urlParam(ctx)

	// Check if URL is known
	if c.knownURL(ctx, urlParam) {
		// Redirect to URL
		ctx.Redirect(http.StatusFound, urlParam)
		return
	}

	// Get relme keywords
	relmeKeywords := c.relmeKeywords(ctx)

	// Render the index template
	ctx.HTML(http.StatusOK, "external_redirect/index", gin.H{
		"layout": "fullscreen",
		"url":    urlParam,
		"rel":    relmeKeywords,
	})
}

// relmeKeywords gets the relme keywords from the request
func (c *ExternalRedirectController) relmeKeywords(ctx *gin.Context) string {
	rel := ctx.Query("rel")
	if rel == "" {
		return ""
	}
	return strings.TrimSpace(rel)
}

// urlParam gets the URL parameter from the request
func (c *ExternalRedirectController) urlParam(ctx *gin.Context) string {
	url := ctx.Query("url")
	if url == "" {
		return ""
	}
	return strings.TrimSpace(url)
}

// knownURL checks if the URL is known
func (c *ExternalRedirectController) knownURL(ctx *gin.Context, urlStr string) bool {
	// Parse URL
	uri, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Get GitLab URL from config
	gitlabURL := c.configService.GetGitlabURL(ctx)

	// Check if URL is from GitLab
	return uri.Host == gitlabURL.Host
}

// shouldHandleURL checks if the URL should be handled
func (c *ExternalRedirectController) shouldHandleURL(ctx *gin.Context, urlStr string) bool {
	// Check if URL is valid
	if !c.urlService.IsValidWebURL(ctx, urlStr) {
		return false
	}

	// Get request base URL and path
	baseURL := ctx.Request.URL.Scheme + "://" + ctx.Request.Host
	path := ctx.Request.URL.Path

	// Check if URL starts with base URL and path
	return !strings.HasPrefix(urlStr, baseURL+path)
}

// checkURLParam checks if the URL parameter is valid
func (c *ExternalRedirectController) checkURLParam(ctx *gin.Context) bool {
	urlParam := c.urlParam(ctx)
	return c.shouldHandleURL(ctx, urlParam)
}
