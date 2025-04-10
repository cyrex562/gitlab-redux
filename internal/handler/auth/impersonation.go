package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SessionKeysToDelete are the session keys that should be cleared during impersonation
var SessionKeysToDelete = []string{
	"github_access_token",
	"gitea_access_token",
	"gitlab_access_token",
	"bitbucket_token",
	"bitbucket_refresh_token",
	"bitbucket_server_personal_access_token",
	"bulk_import_gitlab_access_token",
	"fogbugz_token",
	"cloud_platform_access_token",
}

// Impersonation handles user impersonation functionality
type Impersonation struct {
	userService    *service.UserService
	configService  *service.ConfigService
	logger         *service.Logger
	sessionManager *service.SessionManager
}

// NewImpersonation creates a new instance of Impersonation
func NewImpersonation(
	userService *service.UserService,
	configService *service.ConfigService,
	logger *service.Logger,
	sessionManager *service.SessionManager,
) *Impersonation {
	return &Impersonation{
		userService:    userService,
		configService:  configService,
		logger:         logger,
		sessionManager: sessionManager,
	}
}

// GetCurrentUser returns the current user with impersonation information
func (i *Impersonation) GetCurrentUser(ctx *gin.Context) (*model.User, error) {
	user, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	impersonator, err := i.GetImpersonator(ctx)
	if err != nil {
		return nil, err
	}

	if impersonator != nil {
		user.Impersonator = impersonator
	}

	return user, nil
}

// CheckImpersonationAvailability checks if impersonation is available and enabled
func (i *Impersonation) CheckImpersonationAvailability(ctx *gin.Context) error {
	if !i.IsImpersonationInProgress(ctx) {
		return nil
	}

	if !i.configService.IsImpersonationEnabled() {
		if err := i.StopImpersonation(ctx); err != nil {
			return err
		}
		return ErrImpersonationDisabled
	}

	return nil
}

// StopImpersonation stops the current impersonation session
func (i *Impersonation) StopImpersonation(ctx *gin.Context) error {
	if err := i.LogImpersonationEvent(ctx); err != nil {
		return err
	}

	impersonator, err := i.GetImpersonator(ctx)
	if err != nil {
		return err
	}

	// Set the user back to the impersonator
	if err := i.sessionManager.SetUser(ctx, impersonator); err != nil {
		return err
	}

	// Clear the impersonator ID from the session
	if err := i.sessionManager.Delete(ctx, "impersonator_id"); err != nil {
		return err
	}

	// Clear access token session keys
	if err := i.ClearAccessTokenSessionKeys(ctx); err != nil {
		return err
	}

	return nil
}

// IsImpersonationInProgress checks if an impersonation session is in progress
func (i *Impersonation) IsImpersonationInProgress(ctx *gin.Context) bool {
	impersonatorID, exists := i.sessionManager.Get(ctx, "impersonator_id")
	return exists && impersonatorID != ""
}

// LogImpersonationEvent logs the impersonation event
func (i *Impersonation) LogImpersonationEvent(ctx *gin.Context) error {
	impersonator, err := i.GetImpersonator(ctx)
	if err != nil {
		return err
	}

	currentUser, err := i.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	i.logger.Info("User %s has stopped impersonating %s", impersonator.Username, currentUser.Username)
	return nil
}

// ClearAccessTokenSessionKeys clears access token session keys
func (i *Impersonation) ClearAccessTokenSessionKeys(ctx *gin.Context) error {
	sessionKeys, err := i.sessionManager.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, key := range sessionKeys {
		if contains(SessionKeysToDelete, key) {
			if err := i.sessionManager.Delete(ctx, key); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetImpersonator returns the impersonator user if one exists
func (i *Impersonation) GetImpersonator(ctx *gin.Context) (*model.User, error) {
	impersonatorID, exists := i.sessionManager.Get(ctx, "impersonator_id")
	if !exists || impersonatorID == "" {
		return nil, nil
	}

	return i.userService.GetUserByID(ctx, impersonatorID)
}

// StartImpersonation starts impersonating a user
func (i *Impersonation) StartImpersonation(ctx *gin.Context, targetUser *model.User) error {
	currentUser, err := i.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	// Store the current user as the impersonator
	if err := i.sessionManager.Set(ctx, "impersonator_id", currentUser.ID); err != nil {
		return err
	}

	// Set the target user as the current user
	if err := i.sessionManager.SetUser(ctx, targetUser); err != nil {
		return err
	}

	i.logger.Info("User %s is now impersonating %s", currentUser.Username, targetUser.Username)
	return nil
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
