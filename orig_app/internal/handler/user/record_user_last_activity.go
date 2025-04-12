package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// RecordUserLastActivity handles recording user last activity
type RecordUserLastActivity struct {
	userService    *service.UserService
	activityService *service.ActivityService
	eventService   *service.EventService
	cookieService  *service.CookieService
	dbService      *service.DatabaseService
	logger         *service.Logger
}

// NewRecordUserLastActivity creates a new instance of RecordUserLastActivity
func NewRecordUserLastActivity(
	userService *service.UserService,
	activityService *service.ActivityService,
	eventService *service.EventService,
	cookieService *service.CookieService,
	dbService *service.DatabaseService,
	logger *service.Logger,
) *RecordUserLastActivity {
	return &RecordUserLastActivity{
		userService:     userService,
		activityService: activityService,
		eventService:    eventService,
		cookieService:   cookieService,
		dbService:       dbService,
		logger:          logger,
	}
}

// SetUserLastActivity sets the user's last activity
func (r *RecordUserLastActivity) SetUserLastActivity(c *gin.Context) error {
	// Check if request is GET
	if c.Request.Method != http.MethodGet {
		return nil
	}

	// Check if database is read-only
	isReadOnly, err := r.dbService.IsReadOnly()
	if err != nil {
		return err
	}

	if isReadOnly {
		return nil
	}

	// Get current user
	user, err := r.userService.GetCurrentUser(c)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	// Record user activity
	err = r.activityService.RecordUserActivity(user)
	if err != nil {
		return err
	}

	return nil
}

// SetMemberLastActivity sets the member's last activity
func (r *RecordUserLastActivity) SetMemberLastActivity(c *gin.Context) error {
	// Get current user
	user, err := r.userService.GetCurrentUser(c)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	// Get context (group or project)
	context, err := r.getContext(c)
	if err != nil {
		return err
	}

	if context == nil {
		return nil
	}

	// Check if context is persisted
	isPersisted, err := r.isContextPersisted(context)
	if err != nil {
		return err
	}

	if !isPersisted {
		return nil
	}

	// Get root ancestor
	rootAncestor, err := r.getRootAncestor(context)
	if err != nil {
		return err
	}

	// Publish activity event
	err = r.eventService.PublishActivityEvent(map[string]interface{}{
		"user_id":       user.ID,
		"namespace_id":  rootAncestor.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetContext gets the context (group or project)
func (r *RecordUserLastActivity) getContext(c *gin.Context) (interface{}, error) {
	// Try to get group from context
	group, err := r.getGroupFromContext(c)
	if err == nil && group != nil {
		return group, nil
	}

	// Try to get project from context
	project, err := r.getProjectFromContext(c)
	if err == nil && project != nil {
		return project, nil
	}

	return nil, nil
}

// IsContextPersisted checks if the context is persisted
func (r *RecordUserLastActivity) isContextPersisted(context interface{}) (bool, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would check if the context is persisted
	return true, nil
}

// GetRootAncestor gets the root ancestor of the context
func (r *RecordUserLastActivity) getRootAncestor(context interface{}) (interface{}, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the root ancestor of the context
	return context, nil
}

// GetGroupFromContext gets the group from context
func (r *RecordUserLastActivity) getGroupFromContext(c *gin.Context) (interface{}, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the group from the context
	return nil, nil
}

// GetProjectFromContext gets the project from context
func (r *RecordUserLastActivity) getProjectFromContext(c *gin.Context) (interface{}, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the project from the context
	return nil, nil
}
