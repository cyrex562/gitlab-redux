package import_controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/bulk_imports"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/rate_limiter"
)

const (
	PollingInterval = 3000
)

// BulkImportsController handles bulk import functionality
type BulkImportsController struct {
	bulkImportsService *bulk_imports.Service
	rateLimiter       *rate_limiter.Service
}

// NewBulkImportsController creates a new bulk imports controller
func NewBulkImportsController(
	bulkImportsService *bulk_imports.Service,
	rateLimiter *rate_limiter.Service,
) *BulkImportsController {
	return &BulkImportsController{
		bulkImportsService: bulkImportsService,
		rateLimiter:       rateLimiter,
	}
}

// RegisterRoutes registers the routes for the bulk imports controller
func (c *BulkImportsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/configure", c.Configure)
	router.GET("/status", c.Status)
	router.GET("/history", c.History)
	router.GET("/failures", c.Failures)
	router.POST("/create", c.Create)
	router.GET("/realtime_changes", c.RealtimeChanges)
}

// Configure handles bulk import configuration
func (c *BulkImportsController) Configure(ctx *gin.Context) {
	ctx.Set("bulk_import_gitlab_access_token", strings.TrimSpace(ctx.PostForm("bulk_import_gitlab_access_token")))
	ctx.Set("bulk_import_gitlab_url", ctx.PostForm("bulk_import_gitlab_url"))

	if err := c.verifyBlockedURI(ctx); err != nil {
		c.clearSessionData(ctx)
		ctx.Redirect(http.StatusFound, "/new_group#import-group-pane")
		ctx.Set("alert", c.safeFormat("Specified URL cannot be used: \"%v\"", err))
		return
	}

	if err := c.validateConfigureParams(ctx); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, c.getStatusURL(ctx))
}

// Status handles GET requests for bulk import status
func (c *BulkImportsController) Status(ctx *gin.Context) {
	format := ctx.DefaultQuery("format", "json")
	if format == "json" {
		data, err := c.bulkImportsService.GetImportableData(ctx, c.getQueryParams(ctx), c.getCredentials(ctx))
		if err != nil {
			c.handleConnectionError(ctx, err)
			return
		}

		for _, header := range c.getPaginationHeaders() {
			ctx.Header(header, data.Response.Headers[header])
		}

		ctx.JSON(http.StatusOK, gin.H{
			"importable_data":      c.serializedData(data.Response.ParsedResponse),
			"version_validation":   data.VersionValidation,
		})
		return
	}

	// HTML format
	namespaceID := ctx.Query("namespace_id")
	if namespaceID != "" {
		namespace, err := c.bulkImportsService.FindNamespace(ctx, namespaceID)
		if err != nil {
			ctx.HTML(http.StatusNotFound, "errors/404", nil)
			return
		}

		user := ctx.MustGet("current_user").(*models.User)
		if !user.CanCreateSubgroup(namespace) {
			ctx.HTML(http.StatusNotFound, "errors/404", nil)
			return
		}
	}

	ctx.Set("source_url", ctx.GetString("bulk_import_gitlab_url"))
	ctx.HTML(http.StatusOK, "import/bulk_imports/status", nil)
}

// History handles GET requests for bulk import history
func (c *BulkImportsController) History(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "import/bulk_imports/history", nil)
}

// Failures handles GET requests for bulk import failures
func (c *BulkImportsController) Failures(ctx *gin.Context) {
	bulkImport, err := c.getBulkImport(ctx)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	entity, err := c.getBulkImportEntity(ctx, bulkImport)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	ctx.Set("bulk_import_entity", entity)
	ctx.HTML(http.StatusOK, "import/bulk_imports/failures", nil)
}

// Create handles POST requests for creating bulk imports
func (c *BulkImportsController) Create(ctx *gin.Context) {
	if c.isThrottled(ctx) {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"success": false})
		return
	}

	params := c.getCreateParams(ctx)
	if !c.validCreateParams(params) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"success": false})
		return
	}

	responses := make([]gin.H, 0)
	for _, entry := range params {
		if entry["destination_name"] != "" {
			if entry["destination_slug"] == "" {
				entry["destination_slug"] = entry["destination_name"]
			}
			delete(entry, "destination_name")
		}

		result, err := c.bulkImportsService.Create(ctx, ctx.MustGet("current_user").(*models.User), entry, c.getCredentials(ctx))
		if err != nil {
			responses = append(responses, gin.H{
				"success": false,
				"message": err.Error(),
			})
			continue
		}

		responses = append(responses, gin.H{
			"success": true,
			"id":      result.ID,
			"message": result.Message,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}

// RealtimeChanges handles GET requests for realtime changes
func (c *BulkImportsController) RealtimeChanges(ctx *gin.Context) {
	ctx.Header("X-Poll-Interval", fmt.Sprintf("%d", PollingInterval))

	imports := c.getCurrentUserBulkImports(ctx)
	ctx.JSON(http.StatusOK, imports)
}

// Private helper methods

func (c *BulkImportsController) getBulkImport(ctx *gin.Context) (*models.BulkImport, error) {
	id := ctx.Query("id")
	if id == "" {
		return nil, fmt.Errorf("missing bulk import id")
	}

	return c.bulkImportsService.FindBulkImport(ctx, id)
}

func (c *BulkImportsController) getBulkImportEntity(ctx *gin.Context, bulkImport *models.BulkImport) (*models.BulkImportEntity, error) {
	entityID := ctx.Query("entity_id")
	if entityID == "" {
		return nil, fmt.Errorf("missing entity id")
	}

	return c.bulkImportsService.FindBulkImportEntity(ctx, bulkImport, entityID)
}

func (c *BulkImportsController) getPaginationHeaders() []string {
	return []string{
		"x-next-page",
		"x-page",
		"x-per-page",
		"x-prev-page",
		"x-total",
		"x-total-pages",
	}
}

func (c *BulkImportsController) serializedData(data interface{}) interface{} {
	// TODO: Implement data serialization
	return data
}

func (c *BulkImportsController) getQueryParams(ctx *gin.Context) map[string]interface{} {
	params := map[string]interface{}{
		"top_level_only":    true,
		"min_access_level":  models.AccessLevelOwner,
	}

	filter := c.sanitizedFilterParam(ctx)
	if filter != "" {
		params["search"] = filter
	}

	return params
}

func (c *BulkImportsController) validateConfigureParams(ctx *gin.Context) error {
	client := c.bulkImportsService.NewHTTPClient(
		c.getCredentials(ctx).URL,
		c.getCredentials(ctx).AccessToken,
	)

	if err := client.ValidateInstanceVersion(ctx); err != nil {
		return err
	}

	return client.ValidateImportScopes(ctx)
}

func (c *BulkImportsController) getCreateParams(ctx *gin.Context) []map[string]interface{} {
	// TODO: Implement parameter extraction
	return []map[string]interface{}{}
}

func (c *BulkImportsController) validCreateParams(params []map[string]interface{}) bool {
	for _, param := range params {
		if param["source_type"] != "group_entity" {
			return false
		}
	}
	return true
}

func (c *BulkImportsController) ensureBulkImportEnabled(ctx *gin.Context) bool {
	// TODO: Implement feature flag check
	return true
}

func (c *BulkImportsController) verifyBlockedURI(ctx *gin.Context) error {
	url := ctx.GetString("bulk_import_gitlab_url")
	return c.bulkImportsService.ValidateURL(ctx, url, c.allowLocalRequests(ctx))
}

func (c *BulkImportsController) allowLocalRequests(ctx *gin.Context) bool {
	// TODO: Implement settings check
	return true
}

func (c *BulkImportsController) handleConnectionError(ctx *gin.Context, err error) {
	c.clearSessionData(ctx)

	errorMsg := c.safeFormat("Unable to connect to server: %v", err)
	ctx.Set("alert", errorMsg)

	format := ctx.DefaultQuery("format", "json")
	if format == "json" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"message":  errorMsg,
				"redirect": "/new_group",
			},
		})
		return
	}

	ctx.Redirect(http.StatusFound, "/new_group#import-group-pane")
}

func (c *BulkImportsController) clearSessionData(ctx *gin.Context) {
	ctx.Set("bulk_import_gitlab_url", "")
	ctx.Set("bulk_import_gitlab_access_token", "")
}

func (c *BulkImportsController) getCredentials(ctx *gin.Context) *bulk_imports.Credentials {
	return &bulk_imports.Credentials{
		URL:         ctx.GetString("bulk_import_gitlab_url"),
		AccessToken: ctx.GetString("bulk_import_gitlab_access_token"),
	}
}

func (c *BulkImportsController) sanitizedFilterParam(ctx *gin.Context) string {
	filter := ctx.Query("filter")
	if filter == "" {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(filter))
}

func (c *BulkImportsController) getCurrentUserBulkImports(ctx *gin.Context) []interface{} {
	user := ctx.MustGet("current_user").(*models.User)
	return c.bulkImportsService.GetUserBulkImports(ctx, user)
}

func (c *BulkImportsController) isThrottled(ctx *gin.Context) bool {
	user := ctx.MustGet("current_user").(*models.User)
	return c.rateLimiter.IsThrottled(ctx, user, "bulk_import")
}

func (c *BulkImportsController) getStatusURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.bulkImportsService.GetStatusURL(ctx, "bulk_imports", namespaceID)
}

func (c *BulkImportsController) safeFormat(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
} 