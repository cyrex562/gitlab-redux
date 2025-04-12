package auth

import (
	"sync"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// InitializesCurrentUserMode handles initializing and providing access to the current user's mode
type InitializesCurrentUserMode struct {
	userService *service.UserService
	userModeService *service.UserModeService
	mu sync.Mutex
	currentUserMode *model.UserMode
}

// NewInitializesCurrentUserMode creates a new instance of InitializesCurrentUserMode
func NewInitializesCurrentUserMode(
	userService *service.UserService,
	userModeService *service.UserModeService,
) *InitializesCurrentUserMode {
	return &InitializesCurrentUserMode{
		userService: userService,
		userModeService: userModeService,
	}
}

// GetCurrentUserMode returns the current user's mode
func (i *InitializesCurrentUserMode) GetCurrentUserMode(ctx *gin.Context) (*model.UserMode, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.currentUserMode != nil {
		return i.currentUserMode, nil
	}

	// Get the current user
	currentUser, err := i.userService.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize the current user mode
	i.currentUserMode = i.userModeService.NewUserMode(currentUser)
	return i.currentUserMode, nil
}

// ResetCurrentUserMode resets the current user mode
func (i *InitializesCurrentUserMode) ResetCurrentUserMode() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.currentUserMode = nil
}
