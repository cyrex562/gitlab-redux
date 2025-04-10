package setup

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// CheckInitialSetup handles checking if the application is in initial setup state
type CheckInitialSetup struct {
	userService *service.UserService
	logger      *util.Logger
}

// NewCheckInitialSetup creates a new instance of CheckInitialSetup
func NewCheckInitialSetup(
	userService *service.UserService,
	logger *util.Logger,
) *CheckInitialSetup {
	return &CheckInitialSetup{
		userService: userService,
		logger:      logger,
	}
}

// InInitialSetupState checks if the application is in initial setup state
func (c *CheckInitialSetup) InInitialSetupState() bool {
	// Count users to check if we have exactly one
	userCount, err := c.userService.Count()
	if err != nil {
		c.logger.Error("Failed to count users", err)
		return false
	}

	// Return false if we don't have exactly one user
	if userCount != 1 {
		return false
	}

	// Get the admin user
	admin, err := c.userService.GetLastAdmin()
	if err != nil {
		c.logger.Error("Failed to get admin user", err)
		return false
	}

	// Return false if no admin user found
	if admin == nil {
		return false
	}

	// Check if the admin user requires password creation for web
	return admin.RequirePasswordCreationForWeb()
}
