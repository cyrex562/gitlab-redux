package issuable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// IssuableActions handles common actions for issuable items
type IssuableActions struct {
	authService *service.AuthService
	issuableService *service.IssuableService
	serializerService *service.SerializerService
	discussionService *service.DiscussionService
	userService *service.UserService
	projectService *service.ProjectService
	markdownService *service.MarkdownService
	errorTrackingService *service.ErrorTrackingService
	logger *service.Logger
}

// NewIssuableActions creates a new instance of IssuableActions
func NewIssuableActions(
	authService *service.AuthService,
	issuableService *service.IssuableService,
	serializerService *service.SerializerService,
	discussionService *service.DiscussionService,
	userService *service.UserService,
	projectService *service.ProjectService,
	markdownService *service.MarkdownService,
	errorTrackingService *service.ErrorTrackingService,
	logger *service.Logger,
) *IssuableActions {
	return &IssuableActions{
		authService: authService,
		issuableService: issuableService,
		serializerService: serializerService,
		discussionService: discussionService,
		userService: userService,
		projectService: projectService,
		markdownService: markdownService,
		errorTrackingService: errorTrackingService,
		logger: logger,
	}
}

// Show displays an issuable item
func (i *IssuableActions) Show(ctx *gin.Context) {
	// Get the issuable from the context
	issuable, err := i.getIssuable(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Set application context
	i.setApplicationContext(ctx, issuable)

	// Check if the request is for HTML or JSON
	format := ctx.DefaultQuery("format", "html")
	if format == "html" {
		// Render the sidebar
		sidebar, err := i.serializerService.Represent(issuable, "sidebar")
		if err != nil {
			i.handleError(ctx, err)
			return
		}
		ctx.Set("issuable_sidebar", sidebar)
		ctx.HTML(http.StatusOK, "show", nil)
	} else {
		// Get the serializer from the query parameters
		serializer := ctx.DefaultQuery("serializer", "default")
		// Render the JSON response
		representation, err := i.serializerService.Represent(issuable, serializer)
		if err != nil {
			i.handleError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, representation)
	}
}

// Update updates an issuable item
func (i *IssuableActions) Update(ctx *gin.Context) {
	// Get the issuable from the context
	issuable, err := i.getIssuable(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Update the issuable
	updatedIssuable, err := i.issuableService.Update(issuable, ctx.Request.Form)
	if err != nil {
		// Check if it's a stale object error
		if strings.Contains(err.Error(), "stale object") {
			i.renderConflictResponse(ctx, issuable)
			return
		}
		i.handleError(ctx, err)
		return
	}

	// Check if the issuable is a spammable
	if spammable, ok := updatedIssuable.(model.Spammable); ok {
		// Check if the request is for HTML or JSON
		format := ctx.DefaultQuery("format", "html")
		if format == "html" {
			if updatedIssuable.IsValid() {
				// Redirect to the spammable path
				ctx.Redirect(http.StatusSeeOther, i.getSpammablePath(spammable))
			} else {
				// Render the edit form with CAPTCHA check
				i.withCaptchaCheckHTMLFormat(ctx, spammable, func() {
					ctx.HTML(http.StatusOK, "edit", gin.H{
						"issuable": updatedIssuable,
					})
				})
			}
		} else {
			// Render the JSON response with CAPTCHA check
			i.withCaptchaCheckJSONFormat(ctx, spammable, func() {
				i.renderEntityJSON(ctx, updatedIssuable)
			})
		}
	} else {
		// Check if the request is for HTML or JSON
		format := ctx.DefaultQuery("format", "html")
		if format == "html" {
			ctx.HTML(http.StatusOK, "edit", gin.H{
				"issuable": updatedIssuable,
			})
		} else {
			i.renderEntityJSON(ctx, updatedIssuable)
		}
	}
}

// RealtimeChanges returns realtime changes for an issuable item
func (i *IssuableActions) RealtimeChanges(ctx *gin.Context) {
	// Get the issuable from the context
	issuable, err := i.getIssuable(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Set the polling interval header
	ctx.Header("Poll-Interval", "3000")

	// Create the response
	response := map[string]interface{}{
		"title":               i.markdownService.RenderField(issuable, "title"),
		"title_text":          issuable.GetTitle(),
		"description":         i.markdownService.RenderField(issuable, "description"),
		"description_text":    issuable.GetDescription(),
		"task_completion_status": issuable.GetTaskCompletionStatus(),
		"lock_version":        issuable.GetLockVersion(),
	}

	// Check if the issuable has been edited
	if issuable.IsEdited() {
		response["updated_at"] = issuable.GetLastEditedAt().Format(time.RFC3339)
		response["updated_by_name"] = issuable.GetLastEditedBy().GetName()
		response["updated_by_path"] = i.getUserPath(issuable.GetLastEditedBy())
	}

	// Render the JSON response
	ctx.JSON(http.StatusOK, response)
}

// Destroy deletes an issuable item
func (i *IssuableActions) Destroy(ctx *gin.Context) {
	// Get the issuable from the context
	issuable, err := i.getIssuable(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Authorize the destroy action
	if err := i.authorizeDestroyIssuable(ctx, issuable); err != nil {
		i.handleError(ctx, err)
		return
	}

	// Check the destroy confirmation
	if err := i.checkDestroyConfirmation(ctx, issuable); err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the project from the context
	project, err := i.getProject(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Destroy the issuable
	if err := i.issuableService.Destroy(project, currentUser, issuable); err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the human class name
	humanClassName := issuable.GetHumanClassName()

	// Set the flash notice
	ctx.Set("flash_notice", fmt.Sprintf("The %s was successfully deleted.", humanClassName))

	// Get the index path
	indexPath := i.getPolymorphicPath(i.getParent(ctx), issuable)

	// Check if the request is for HTML or JSON
	format := ctx.DefaultQuery("format", "html")
	if format == "html" {
		ctx.Redirect(http.StatusSeeOther, indexPath)
	} else {
		ctx.JSON(http.StatusOK, map[string]string{
			"web_url": indexPath,
		})
	}
}

// CheckDestroyConfirmation checks if the destroy confirmation is provided
func (i *IssuableActions) CheckDestroyConfirmation(ctx *gin.Context) error {
	// Get the issuable from the context
	issuable, err := i.getIssuable(ctx)
	if err != nil {
		return err
	}

	// Check if the destroy confirmation is provided
	if ctx.DefaultQuery("destroy_confirm", "") == "" {
		// Get the human class name
		humanClassName := issuable.GetHumanClassName()

		// Create the error message
		errorMessage := fmt.Sprintf("Destroy confirmation not provided for %s", humanClassName)

		// Get the project from the context
		project, err := i.getProject(ctx)
		if err != nil {
			return err
		}

		// Track the exception
		i.errorTrackingService.TrackException(
			fmt.Errorf(errorMessage),
			map[string]interface{}{
				"project_path":   project.GetFullPath(),
				"issuable_type": issuable.GetType(),
				"issuable_id":   issuable.GetID(),
			},
		)

		// Get the index path
		indexPath := i.getPolymorphicPath(i.getParent(ctx), issuable)

		// Check if the request is for HTML or JSON
		format := ctx.DefaultQuery("format", "html")
		if format == "html" {
			ctx.Set("flash_notice", errorMessage)
			ctx.Redirect(http.StatusSeeOther, indexPath)
		} else {
			ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"errors": errorMessage,
			})
		}

		return fmt.Errorf(errorMessage)
	}

	return nil
}

// BulkUpdate updates multiple issuable items
func (i *IssuableActions) BulkUpdate(ctx *gin.Context) {
	// Get the parent from the context
	parent, err := i.getParent(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Authorize the admin action
	if err := i.authorizeAdminIssuable(ctx, parent); err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the resource name
	resourceName := i.getResourceName(ctx)

	// Get the bulk update parameters
	bulkUpdateParams, err := i.getBulkUpdateParams(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Update the issuables
	result, err := i.issuableService.BulkUpdate(parent, currentUser, bulkUpdateParams, resourceName)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Check if the result is successful
	if result.IsSuccess() {
		// Get the count
		count := result.GetPayload()["count"].(int)

		// Render the JSON response
		ctx.JSON(http.StatusOK, map[string]string{
			"notice": fmt.Sprintf("%d %s updated", count, i.pluralize(resourceName, count)),
		})
	} else if result.IsError() {
		// Render the JSON response
		ctx.JSON(result.GetHTTPStatus(), map[string]interface{}{
			"errors": result.GetMessage(),
		})
	}
}

// Discussions returns the discussions for an issuable item
func (i *IssuableActions) Discussions(ctx *gin.Context) {
	// Get the issuable from the context
	issuable, err := i.getIssuable(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the finder parameters
	finderParams, err := i.getFinderParamsForIssuable(ctx, issuable)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Get the discussions
	discussionNotes, paginator, err := i.discussionService.GetDiscussions(currentUser, issuable, finderParams)
	if err != nil {
		i.handleError(ctx, err)
		return
	}

	// Check if there is a next page
	if paginator != nil && paginator.HasNextPage() {
		ctx.Header("X-Next-Page-Cursor", paginator.GetCursorForNextPage())
	}

	// Check the issuable type
	switch issuable.GetType() {
	case "MergeRequest", "Issue":
		// Check if the response is stale
		if i.isStale(ctx, discussionNotes) {
			// Render the JSON response
			representation, err := i.discussionService.SerializeDiscussions(discussionNotes, ctx)
			if err != nil {
				i.handleError(ctx, err)
				return
			}
			ctx.JSON(http.StatusOK, representation)
		}
	default:
		// Render the JSON response
		representation, err := i.discussionService.SerializeDiscussions(discussionNotes, ctx)
		if err != nil {
			i.handleError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, representation)
	}
}

// RenderConflictResponse renders a conflict response
func (i *IssuableActions) renderConflictResponse(ctx *gin.Context, issuable model.Issuable) {
	// Check if the request is for HTML or JSON
	format := ctx.DefaultQuery("format", "html")
	if format == "html" {
		ctx.HTML(http.StatusOK, "edit", gin.H{
			"conflict": true,
			"issuable": issuable,
		})
	} else {
		// Get the human class name
		humanClassName := issuable.GetHumanClassName()

		// Create the error message
		errorMessage := fmt.Sprintf("Someone edited this %s at the same time you did. Please refresh your browser and make sure your changes will not unintentionally remove theirs.", humanClassName)

		// Render the JSON response
		ctx.JSON(http.StatusConflict, map[string]interface{}{
			"errors": []string{errorMessage},
		})
	}
}

// AuthorizeDestroyIssuable authorizes the destroy action
func (i *IssuableActions) authorizeDestroyIssuable(ctx *gin.Context, issuable model.Issuable) error {
	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	// Check if the user can destroy the issuable
	if !i.authService.Can(currentUser, fmt.Sprintf("destroy_%s", issuable.GetAbilityName()), issuable) {
		return fmt.Errorf("access denied")
	}

	return nil
}

// AuthorizeAdminIssuable authorizes the admin action
func (i *IssuableActions) authorizeAdminIssuable(ctx *gin.Context, parent interface{}) error {
	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	// Get the resource name
	resourceName := i.getResourceName(ctx)

	// Check if the user can admin the parent
	if !i.authService.Can(currentUser, fmt.Sprintf("admin_%s", resourceName), parent) {
		return fmt.Errorf("access denied")
	}

	return nil
}

// SetApplicationContext sets the application context
func (i *IssuableActions) setApplicationContext(ctx *gin.Context, issuable model.Issuable) {
	// This is a no-op in the base implementation
	// The logic is defined in EE module
}

// GetBulkUpdateParams gets the bulk update parameters
func (i *IssuableActions) getBulkUpdateParams(ctx *gin.Context) (map[string]interface{}, error) {
	// Get the update parameters
	updateParams, exists := ctx.GetPostForm("update")
	if !exists {
		return nil, fmt.Errorf("update parameters not provided")
	}

	// Parse the update parameters
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(updateParams), &params); err != nil {
		return nil, err
	}

	// Clean the bulk update parameters
	return i.cleanBulkUpdateParams(params), nil
}

// CleanBulkUpdateParams cleans the bulk update parameters
func (i *IssuableActions) cleanBulkUpdateParams(params map[string]interface{}) map[string]interface{} {
	// Create a new map for the cleaned parameters
	cleanedParams := make(map[string]interface{})

	// Iterate over the parameters
	for k, v := range params {
		// Skip if the key is issuable_ids
		if k == "issuable_ids" {
			cleanedParams[k] = v
			continue
		}

		// Check if the value is an array
		if arr, ok := v.([]interface{}); ok {
			// Check if the array is empty
			if len(arr) == 0 {
				continue
			}
		}

		// Check if the value is empty
		if v == nil || v == "" {
			continue
		}

		// Add the parameter to the cleaned parameters
		cleanedParams[k] = v
	}

	return cleanedParams
}

// GetBulkUpdatePermittedKeys gets the bulk update permitted keys
func (i *IssuableActions) getBulkUpdatePermittedKeys() []string {
	return []string{
		"issuable_ids",
		"assignee_id",
		"milestone_id",
		"state_event",
		"subscription_event",
		"confidential",
		"assignee_ids",
		"add_label_ids",
		"remove_label_ids",
	}
}

// GetResourceName gets the resource name
func (i *IssuableActions) getResourceName(ctx *gin.Context) string {
	// Get the controller name
	controllerName := ctx.Request.URL.Path
	parts := strings.Split(controllerName, "/")
	if len(parts) > 0 {
		controllerName = parts[len(parts)-1]
	}

	// Singularize the controller name
	return i.singularize(controllerName)
}

// RenderEntityJSON renders the entity JSON
func (i *IssuableActions) renderEntityJSON(ctx *gin.Context, issuable model.Issuable) {
	// Check if the issuable is valid
	if issuable.IsValid() {
		// Render the JSON response
		representation, err := i.serializerService.Represent(issuable)
		if err != nil {
			i.handleError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, representation)
	} else {
		// Render the JSON response
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"errors": issuable.GetErrors().FullMessages(),
		})
	}
}

// GetParent gets the parent from the context
func (i *IssuableActions) getParent(ctx *gin.Context) (interface{}, error) {
	// Try to get the project from the context
	project, err := i.getProject(ctx)
	if err == nil {
		return project, nil
	}

	// Try to get the group from the context
	group, err := i.getGroup(ctx)
	if err == nil {
		return group, nil
	}

	return nil, fmt.Errorf("parent not found")
}

// GetFinderParamsForIssuable gets the finder parameters for an issuable
func (i *IssuableActions) getFinderParamsForIssuable(ctx *gin.Context, issuable model.Issuable) (map[string]interface{}, error) {
	// Get the notes filter
	notesFilter, err := i.getNotesFilter(ctx, issuable)
	if err != nil {
		return nil, err
	}

	// Create the finder parameters
	finderParams := map[string]interface{}{
		"notes_filter": notesFilter,
		"cursor":       ctx.DefaultQuery("cursor", ""),
		"per_page":     ctx.DefaultQuery("per_page", "20"),
	}

	return finderParams, nil
}

// GetNotesFilter gets the notes filter
func (i *IssuableActions) getNotesFilter(ctx *gin.Context, issuable model.Issuable) (int, error) {
	// Get the notes filter from the query parameters
	notesFilterStr := ctx.DefaultQuery("notes_filter", "")
	if notesFilterStr == "" {
		// Get the current user from the context
		currentUser, err := i.userService.GetCurrentUser(ctx)
		if err != nil {
			return 0, err
		}

		// Get the notes filter for the issuable
		return currentUser.GetNotesFilterFor(issuable), nil
	}

	// Parse the notes filter
	notesFilter, err := i.parseInt(notesFilterStr)
	if err != nil {
		return 0, err
	}

	// Check if the database is read-only
	if i.isDatabaseReadOnly() {
		return notesFilter, nil
	}

	// Check if the persist filter is false
	if ctx.DefaultQuery("persist_filter", "true") == "false" {
		return notesFilter, nil
	}

	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		return 0, err
	}

	// Set the notes filter for the issuable
	if err := currentUser.SetNotesFilter(notesFilter, issuable); err != nil {
		return 0, err
	}

	return notesFilter, nil
}

// GetDiscussionCacheContext gets the discussion cache context
func (i *IssuableActions) getDiscussionCacheContext(ctx *gin.Context, issuable model.Issuable) string {
	// Get the current user from the context
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		return ""
	}

	// Get the project from the context
	project, err := i.getProject(ctx)
	if err != nil {
		return ""
	}

	// Get the cache key
	cacheKey := currentUser.GetCacheKey()

	// Get the human max access
	humanMaxAccess := project.GetTeam().GetHumanMaxAccess(currentUser.GetID())

	// Create the cache context
	return fmt.Sprintf("%s:%s:v2", cacheKey, humanMaxAccess)
}

// IsStale checks if the response is stale
func (i *IssuableActions) isStale(ctx *gin.Context, discussionNotes []model.DiscussionNote) bool {
	// Get the etag from the request
	etag := ctx.GetHeader("If-None-Match")
	if etag == "" {
		return true
	}

	// Get the discussion cache context
	cacheContext := i.getDiscussionCacheContext(ctx, discussionNotes[0].GetNoteable())

	// Create the etag
	newEtag := fmt.Sprintf("\"%s\"", cacheContext)

	// Check if the etag is stale
	return etag != newEtag
}

// HandleError handles an error
func (i *IssuableActions) handleError(ctx *gin.Context, err error) {
	// Log the error
	i.logger.Error(err)

	// Check if the request is for HTML or JSON
	format := ctx.DefaultQuery("format", "html")
	if format == "html" {
		ctx.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// GetIssuable gets the issuable from the context
func (i *IssuableActions) getIssuable(ctx *gin.Context) (model.Issuable, error) {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return nil, fmt.Errorf("not implemented")
}

// GetProject gets the project from the context
func (i *IssuableActions) getProject(ctx *gin.Context) (*model.Project, error) {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return nil, fmt.Errorf("not implemented")
}

// GetGroup gets the group from the context
func (i *IssuableActions) getGroup(ctx *gin.Context) (*model.Group, error) {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return nil, fmt.Errorf("not implemented")
}

// GetSpammablePath gets the spammable path
func (i *IssuableActions) getSpammablePath(spammable model.Spammable) string {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return ""
}

// GetUserPath gets the user path
func (i *IssuableActions) getUserPath(user *model.User) string {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return ""
}

// GetPolymorphicPath gets the polymorphic path
func (i *IssuableActions) getPolymorphicPath(parent interface{}, issuable model.Issuable) string {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return ""
}

// WithCaptchaCheckHTMLFormat performs a CAPTCHA check for HTML format
func (i *IssuableActions) withCaptchaCheckHTMLFormat(ctx *gin.Context, spammable model.Spammable, callback func()) {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	callback()
}

// WithCaptchaCheckJSONFormat performs a CAPTCHA check for JSON format
func (i *IssuableActions) withCaptchaCheckJSONFormat(ctx *gin.Context, spammable model.Spammable, callback func()) {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	callback()
}

// IsDatabaseReadOnly checks if the database is read-only
func (i *IssuableActions) isDatabaseReadOnly() bool {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return false
}

// ParseInt parses an integer
func (i *IssuableActions) parseInt(s string) (int, error) {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return 0, fmt.Errorf("not implemented")
}

// Singularize singularizes a string
func (i *IssuableActions) singularize(s string) string {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	return s
}

// Pluralize pluralizes a string
func (i *IssuableActions) pluralize(s string, count int) string {
	// This is a placeholder method
	// The actual implementation will depend on the specific controller
	if count == 1 {
		return s
	}
	return s + "s"
}
