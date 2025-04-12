package oauth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// OauthApplications handles OAuth application functionality
type OauthApplications struct {
	configService *service.ConfigService
	sessionService *service.SessionService
	authService   *service.AuthService
	logger        *service.Logger
}

// NewOauthApplications creates a new instance of OauthApplications
func NewOauthApplications(
	configService *service.ConfigService,
	sessionService *service.SessionService,
	authService *service.AuthService,
	logger *service.Logger,
) *OauthApplications {
	return &OauthApplications{
		configService:  configService,
		sessionService: sessionService,
		authService:    authService,
		logger:         logger,
	}
}

// SetupMiddleware sets up the middleware for OAuth applications
func (o *OauthApplications) SetupMiddleware(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		// Prepare scopes for create and update actions
		if (c.Request.URL.Path == "/oauth/applications" && c.Request.Method == http.MethodPost) ||
			(c.Request.URL.Path == "/oauth/applications/:id" && c.Request.Method == http.MethodPut) {
			if err := o.prepareScopes(c); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Failed to prepare scopes",
				})
				return
			}
		}

		c.Next()
	})
}

// PrepareScopes prepares the scopes for OAuth applications
func (o *OauthApplications) prepareScopes(ctx *gin.Context) error {
	// Get application params from request
	var appParams map[string]interface{}
	if err := ctx.ShouldBindJSON(&appParams); err != nil {
		return err
	}

	// Get scopes from params
	scopes, ok := appParams["scopes"].([]interface{})
	if !ok {
		return nil
	}

	// Join scopes with space
	scopesStr := make([]string, len(scopes))
	for i, scope := range scopes {
		scopesStr[i] = scope.(string)
	}
	appParams["scopes"] = strings.Join(scopesStr, " ")

	// Set application params back to request
	ctx.Set("doorkeeper_application", appParams)

	return nil
}

// SetCreatedSession sets the created session flag
func (o *OauthApplications) SetCreatedSession(ctx *gin.Context) error {
	// Set created session flag
	return o.sessionService.Set(ctx, "oauth_applications_created", true)
}

// GetCreatedSession gets the created session flag
func (o *OauthApplications) GetCreatedSession(ctx *gin.Context) (bool, error) {
	// Get created session flag
	created, err := o.sessionService.Get(ctx, "oauth_applications_created")
	if err != nil {
		return false, err
	}

	// Delete created session flag
	if err := o.sessionService.Delete(ctx, "oauth_applications_created"); err != nil {
		return false, err
	}

	return created.(bool), nil
}

// LoadScopes loads the available scopes
func (o *OauthApplications) LoadScopes(ctx *gin.Context) ([]string, error) {
	// Get configuration
	config, err := o.configService.GetConfiguration()
	if err != nil {
		return nil, err
	}

	// Get scopes from configuration
	scopes := config.OAuth.Scopes

	// Filter out AI workflow, dynamic user, and self rotate scopes
	filteredScopes := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		if scope != "ai_workflow" && scope != "dynamic_user" && scope != "self_rotate" {
			filteredScopes = append(filteredScopes, scope)
		}
	}

	return filteredScopes, nil
}

// PermittedParams returns the permitted parameters for OAuth applications
func (o *OauthApplications) PermittedParams() []string {
	return []string{"name", "redirect_uri", "scopes", "confidential"}
}

// ApplicationParams returns the application parameters from the request
func (o *OauthApplications) ApplicationParams(ctx *gin.Context) (map[string]interface{}, error) {
	// Get application params from request
	var appParams map[string]interface{}
	if err := ctx.ShouldBindJSON(&appParams); err != nil {
		return nil, err
	}

	// Filter permitted params
	filteredParams := make(map[string]interface{})
	for _, param := range o.PermittedParams() {
		if value, ok := appParams[param]; ok {
			filteredParams[param] = value
		}
	}

	return filteredParams, nil
}
