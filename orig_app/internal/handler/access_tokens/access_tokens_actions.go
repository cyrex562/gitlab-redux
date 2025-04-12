package access_tokens

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// AccessTokensActions handles actions related to access tokens
type AccessTokensActions struct {
	createService *service.ResourceAccessTokensCreateService
	revokeService *service.ResourceAccessTokensRevokeService
	rotateService *service.ResourceAccessTokensRotateService
	authService   *service.AuthService
	logger        *util.Logger
}

// NewAccessTokensActions creates a new instance of AccessTokensActions
func NewAccessTokensActions(
	createService *service.ResourceAccessTokensCreateService,
	revokeService *service.ResourceAccessTokensRevokeService,
	rotateService *service.ResourceAccessTokensRotateService,
	authService *service.AuthService,
	logger *util.Logger,
) *AccessTokensActions {
	return &AccessTokensActions{
		createService: createService,
		revokeService: revokeService,
		rotateService: rotateService,
		authService:   authService,
		logger:        logger,
	}
}

// Index displays a list of access tokens
func (a *AccessTokensActions) Index(ctx *gin.Context, resource model.Resource) {
	// Check permission
	if !a.checkPermission(ctx, "read_resource_access_tokens", resource) {
		return
	}

	// Create a new personal access token
	resourceAccessToken := model.NewPersonalAccessToken()

	// Set index variables
	scopes, activeTokens, activeTokensSize, inactiveTokensSize := a.setIndexVars(ctx, resource)

	// Handle different response formats
	if ctx.GetHeader("Accept") == "application/json" {
		ctx.JSON(200, gin.H{
			"active_access_tokens": activeTokens,
		})
	} else {
		ctx.HTML(200, "access_tokens/index.html", gin.H{
			"resource_access_token": resourceAccessToken,
			"scopes":               scopes,
			"active_access_tokens": activeTokens,
			"active_access_tokens_size": activeTokensSize,
			"inactive_access_tokens_size": inactiveTokensSize,
		})
	}
}

// Create creates a new access token
func (a *AccessTokensActions) Create(ctx *gin.Context, resource model.Resource) {
	// Check permission
	if !a.checkPermission(ctx, "create_resource_access_tokens", resource) {
		return
	}

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	// Get create parameters
	createParams := a.getCreateParams(ctx)

	// Create the token
	tokenResponse := a.createService.Execute(user, resource, createParams)

	if tokenResponse.Success {
		resourceAccessToken := tokenResponse.Payload["access_token"].(model.PersonalAccessToken)
		tokens, size := a.getActiveAccessTokens(resource)

		ctx.JSON(200, gin.H{
			"new_token":           resourceAccessToken.GetToken(),
			"active_access_tokens": tokens,
			"total":               size,
		})
	} else {
		ctx.JSON(422, gin.H{
			"errors": tokenResponse.Errors,
		})
	}
}

// Revoke revokes an access token
func (a *AccessTokensActions) Revoke(ctx *gin.Context, resource model.Resource) {
	// Check permission
	if !a.checkPermission(ctx, "destroy_resource_access_tokens", resource) {
		return
	}

	// Get the token ID from the URL parameters
	tokenID := ctx.Param("id")

	// Find the token
	resourceAccessToken := a.finder().Find(tokenID)

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	// Revoke the token
	revokedResponse := a.revokeService.Execute(user, resource, resourceAccessToken)

	if revokedResponse.Success {
		// Set a success flash message
		ctx.Set("flash_notice", fmt.Sprintf("Revoked access token %s!", resourceAccessToken.GetName()))
	} else {
		// Set an error flash message
		ctx.Set("flash_alert", fmt.Sprintf("Could not revoke access token %s.", resourceAccessToken.GetName()))
	}

	// Redirect to the access tokens path
	ctx.Redirect(302, a.getResourceAccessTokensPath(resource))
}

// Rotate rotates an access token
func (a *AccessTokensActions) Rotate(ctx *gin.Context, resource model.Resource) {
	// Check permission
	if !a.checkPermission(ctx, "manage_resource_access_tokens", resource) {
		return
	}

	// Get rotate parameters
	rotateParams := a.getRotateParams(ctx)

	// Find the token
	token := a.finder().Find(rotateParams["id"].(string))

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	// Rotate the token
	result := a.rotateService.New(user, token, resource, true).Execute()
	resourceAccessToken := result.Payload["personal_access_token"].(model.PersonalAccessToken)

	if result.Success {
		tokens, size := a.getActiveAccessTokens(resource)

		ctx.JSON(200, gin.H{
			"new_token":           resourceAccessToken.GetToken(),
			"active_access_tokens": tokens,
			"total":               size,
		})
	} else {
		ctx.JSON(422, gin.H{
			"message": result.Message,
		})
	}
}

// Inactive displays a list of inactive access tokens
func (a *AccessTokensActions) Inactive(ctx *gin.Context, resource model.Resource) {
	// Check permission
	if !a.checkPermission(ctx, "read_resource_access_tokens", resource) {
		return
	}

	// Get the page number from the query parameters
	page := ctx.DefaultQuery("page", "1")

	// Get inactive tokens
	tokens := a.getInactiveAccessTokens(resource).Page(page)

	// Add pagination headers
	a.addPaginationHeaders(ctx, tokens)

	// Render the tokens
	ctx.JSON(200, a.represent(tokens))
}

// checkPermission checks if the current user has the specified permission
func (a *AccessTokensActions) checkPermission(ctx *gin.Context, action string, resource model.Resource) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")
	user, _ := currentUser.(model.User)

	// Check if the user has the permission
	if !a.authService.Can(user, action, resource) {
		ctx.JSON(404, gin.H{"error": "Not found"})
		return false
	}

	return true
}

// getCreateParams gets the parameters for creating a token
func (a *AccessTokensActions) getCreateParams(ctx *gin.Context) map[string]interface{} {
	var params struct {
		Name        string   `json:"name" binding:"required"`
		ExpiresAt   string   `json:"expires_at"`
		Description string   `json:"description"`
		AccessLevel string   `json:"access_level"`
		Scopes      []string `json:"scopes"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		return nil
	}

	return map[string]interface{}{
		"name":        params.Name,
		"expires_at":  params.ExpiresAt,
		"description": params.Description,
		"access_level": params.AccessLevel,
		"scopes":      params.Scopes,
	}
}

// getRotateParams gets the parameters for rotating a token
func (a *AccessTokensActions) getRotateParams(ctx *gin.Context) map[string]interface{} {
	return map[string]interface{}{
		"id": ctx.Param("id"),
	}
}

// setIndexVars sets the variables for the index action
func (a *AccessTokensActions) setIndexVars(ctx *gin.Context, resource model.Resource) ([]string, []model.PersonalAccessToken, int, int) {
	// Load resource members
	resource.LoadMembers()

	// Get available scopes
	scopes := a.authService.GetAvailableScopesFor(resource)

	// Filter token scopes
	scopes = a.authService.FilterTokenScopes(scopes, ctx)

	// Get active and inactive tokens
	activeTokens, activeTokensSize := a.getActiveAccessTokens(resource)
	inactiveTokensSize := a.getInactiveAccessTokens(resource).Size()

	return scopes, activeTokens, activeTokensSize, inactiveTokensSize
}

// finder creates a new personal access tokens finder
func (a *AccessTokensActions) finder(options map[string]interface{}) *model.PersonalAccessTokensFinder {
	if options == nil {
		options = make(map[string]interface{})
	}

	options["user"] = a.getBotUsers(resource)
	options["impersonation"] = false

	return model.NewPersonalAccessTokensFinder(options)
}

// getBotUsers gets the bot users for a resource
func (a *AccessTokensActions) getBotUsers(resource model.Resource) []model.User {
	return resource.GetBots()
}

// getActiveAccessTokens gets the active access tokens for a resource
func (a *AccessTokensActions) getActiveAccessTokens(resource model.Resource) ([]model.PersonalAccessToken, int) {
	// This would typically be a method on the resource model
	// For now, we'll return an empty slice and 0
	return []model.PersonalAccessToken{}, 0
}

// getInactiveAccessTokens gets the inactive access tokens for a resource
func (a *AccessTokensActions) getInactiveAccessTokens(resource model.Resource) *model.PaginatedCollection {
	// This would typically be a method on the resource model
	// For now, we'll return an empty collection
	return model.NewPaginatedCollection([]model.PersonalAccessToken{})
}

// addPaginationHeaders adds pagination headers to the response
func (a *AccessTokensActions) addPaginationHeaders(ctx *gin.Context, collection *model.PaginatedCollection) {
	ctx.Header("X-Total", fmt.Sprintf("%d", collection.Total()))
	ctx.Header("X-Total-Pages", fmt.Sprintf("%d", collection.TotalPages()))
	ctx.Header("X-Per-Page", fmt.Sprintf("%d", collection.PerPage()))
	ctx.Header("X-Page", fmt.Sprintf("%d", collection.Page()))
	ctx.Header("X-Next-Page", fmt.Sprintf("%d", collection.NextPage()))
	ctx.Header("X-Prev-Page", fmt.Sprintf("%d", collection.PrevPage()))
}

// represent represents a collection of tokens
func (a *AccessTokensActions) represent(collection *model.PaginatedCollection) interface{} {
	// This would typically be a method on the collection model
	// For now, we'll return the collection as is
	return collection
}

// getResourceAccessTokensPath gets the path to the access tokens for a resource
func (a *AccessTokensActions) getResourceAccessTokensPath(resource model.Resource) string {
	// This would typically be a method on the resource model
	// For now, we'll return a default path
	return fmt.Sprintf("/%s/%s/access_tokens", resource.GetType(), resource.GetID())
}
