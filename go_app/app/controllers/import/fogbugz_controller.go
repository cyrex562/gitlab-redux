package import_controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/fogbugz"
	"gitlab.com/gitlab-org/gitlab-redux/internal/services/rate_limiter"
)

// FogbugzController handles Fogbugz import functionality
type FogbugzController struct {
	*BaseController
	fogbugzService *fogbugz.Service
	rateLimiter    *rate_limiter.Service
}

// NewFogbugzController creates a new Fogbugz controller
func NewFogbugzController(
	baseController *BaseController,
	fogbugzService *fogbugz.Service,
	rateLimiter *rate_limiter.Service,
) *FogbugzController {
	return &FogbugzController{
		BaseController: baseController,
		fogbugzService: fogbugzService,
		rateLimiter:    rateLimiter,
	}
}

// RegisterRoutes registers the routes for the Fogbugz controller
func (c *FogbugzController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/new", c.New)
	router.POST("/callback", c.Callback)
	router.GET("/new_user_map", c.NewUserMap)
	router.POST("/create_user_map", c.CreateUserMap)
	router.GET("/status", c.Status)
	router.POST("/create", c.Create)
}

// New handles the new import page
func (c *FogbugzController) New(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "import/fogbugz/new", nil)
}

// Callback handles the Fogbugz callback
func (c *FogbugzController) Callback(ctx *gin.Context) {
	if c.isThrottled(ctx) {
		ctx.Redirect(http.StatusFound, c.getNewURL(ctx))
		ctx.Set("alert", "Rate limit exceeded")
		return
	}

	if err := c.verifyBlockedURI(ctx); err != nil {
		ctx.Redirect(http.StatusFound, c.getNewURL(ctx))
		ctx.Set("alert", c.safeFormat("Specified URL cannot be used: \"%v\"", err))
		return
	}

	client, err := c.fogbugzService.NewClient(ctx, c.getImportParams(ctx))
	if err != nil {
		ctx.Redirect(http.StatusFound, c.getNewURL(ctx))
		ctx.Set("alert", "Could not connect to FogBugz, check your URL")
		return
	}

	token, err := client.GetToken(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, c.getNewURL(ctx))
		ctx.Set("alert", "Could not authenticate with FogBugz")
		return
	}

	ctx.Set("fogbugz_token", token)
	ctx.Set("fogbugz_uri", ctx.PostForm("uri"))

	ctx.Redirect(http.StatusFound, c.getNewUserMapURL(ctx))
}

// NewUserMap handles the new user map page
func (c *FogbugzController) NewUserMap(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "import/fogbugz/new_user_map", nil)
}

// CreateUserMap handles user map creation
func (c *FogbugzController) CreateUserMap(ctx *gin.Context) {
	userMap := c.getUserMapParams(ctx)
	if !c.validUserMap(userMap) {
		ctx.Set("alert", "All users must have a name.")
		ctx.HTML(http.StatusUnprocessableEntity, "import/fogbugz/new_user_map", nil)
		return
	}

	ctx.Set("fogbugz_user_map", userMap)
	ctx.Set("notice", "The user map has been saved. Continue by selecting the projects you want to import.")

	ctx.Redirect(http.StatusFound, c.getStatusURL(ctx))
}

// Create handles project creation from Fogbugz
func (c *FogbugzController) Create(ctx *gin.Context) {
	credentials := c.getCredentials(ctx)
	serviceParams := c.getServiceParams(ctx)

	result, err := c.fogbugzService.ImportProject(ctx, c.getClient(ctx), ctx.MustGet("current_user").(*models.User), serviceParams, credentials)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result.Project)
}

// ImportableRepos returns the list of importable Fogbugz repositories
func (c *FogbugzController) ImportableRepos(ctx *gin.Context) ([]interface{}, error) {
	return c.getClient(ctx).GetRepos(ctx)
}

// IncompatibleRepos returns an empty list as Fogbugz has no incompatible repos
func (c *FogbugzController) IncompatibleRepos(ctx *gin.Context) ([]interface{}, error) {
	return []interface{}{}, nil
}

// ProviderName returns "fogbugz" as the provider name
func (c *FogbugzController) ProviderName(ctx *gin.Context) string {
	return "fogbugz"
}

// ProviderURL returns the Fogbugz URI from session
func (c *FogbugzController) ProviderURL(ctx *gin.Context) string {
	return ctx.GetString("fogbugz_uri")
}

// Private helper methods

func (c *FogbugzController) fogbugzImportEnabled() bool {
	// TODO: Implement check for Fogbugz import feature flag
	return true
}

func (c *FogbugzController) isThrottled(ctx *gin.Context) bool {
	user := ctx.MustGet("current_user").(*models.User)
	return c.rateLimiter.IsThrottled(ctx, user, "fogbugz_import")
}

func (c *FogbugzController) verifyBlockedURI(ctx *gin.Context) error {
	uri := ctx.PostForm("uri")
	return c.fogbugzService.ValidateURL(ctx, uri, c.allowLocalRequests(ctx))
}

func (c *FogbugzController) allowLocalRequests(ctx *gin.Context) bool {
	// TODO: Implement settings check
	return true
}

func (c *FogbugzController) getClient(ctx *gin.Context) *fogbugz.Client {
	return c.fogbugzService.NewClient(ctx, c.getCredentials(ctx))
}

func (c *FogbugzController) getUserMap(ctx *gin.Context) map[string]interface{} {
	client := c.getClient(ctx)
	userMap := client.GetUserMap(ctx)

	storedUserMap := ctx.GetString("fogbugz_user_map")
	if storedUserMap != "" {
		// TODO: Implement user map update
	}

	return userMap
}

func (c *FogbugzController) handleUnauthorized(ctx *gin.Context, err error) {
	ctx.Redirect(http.StatusFound, c.getNewURL(ctx))
	ctx.Set("alert", err.Error())
}

func (c *FogbugzController) getImportParams(ctx *gin.Context) map[string]string {
	return map[string]string{
		"uri":     ctx.PostForm("uri"),
		"email":   ctx.PostForm("email"),
		"password": ctx.PostForm("password"),
	}
}

func (c *FogbugzController) getUserMapParams(ctx *gin.Context) map[string]interface{} {
	// TODO: Implement parameter extraction
	return map[string]interface{}{}
}

func (c *FogbugzController) validUserMap(userMap map[string]interface{}) bool {
	if userMap == nil {
		return false
	}

	for _, user := range userMap {
		if userMap, ok := user.(map[string]interface{}); ok {
			if name, ok := userMap["name"].(string); !ok || name == "" {
				return false
			}
		}
	}
	return true
}

func (c *FogbugzController) getCredentials(ctx *gin.Context) *fogbugz.Credentials {
	return &fogbugz.Credentials{
		URI:   ctx.GetString("fogbugz_uri"),
		Token: ctx.GetString("fogbugz_token"),
	}
}

func (c *FogbugzController) getServiceParams(ctx *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range ctx.Request.PostForm {
		params[k] = v[0]
	}

	params["umap"] = c.getUserMap(ctx)
	// TODO: Add organization_id

	return params
}

func (c *FogbugzController) getNewURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.fogbugzService.GetNewURL(ctx, "fogbugz", namespaceID)
}

func (c *FogbugzController) getNewUserMapURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.fogbugzService.GetNewUserMapURL(ctx, "fogbugz", namespaceID)
}

func (c *FogbugzController) getStatusURL(ctx *gin.Context) string {
	namespaceID := ctx.Query("namespace_id")
	return c.fogbugzService.GetStatusURL(ctx, "fogbugz", namespaceID)
}

func (c *FogbugzController) safeFormat(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
} 