package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// FeatureCategory represents the feature category for initial setup
const FeatureCategory = "system_access"

// UserParams represents the parameters for updating a user during initial setup
type UserParams struct {
	Email                string    `json:"email" binding:"required,email"`
	Password            string    `json:"password" binding:"required,min=8"`
	PasswordConfirmation string    `json:"password_confirmation" binding:"required,eqfield=Password"`
}

// InitialSetupHandler handles initial setup requests
type InitialSetupHandler struct {
	userService *service.UserService
}

// NewInitialSetupHandler creates a new InitialSetupHandler instance
func NewInitialSetupHandler(userService *service.UserService) *InitialSetupHandler {
	return &InitialSetupHandler{
		userService: userService,
	}
}

// New handles the GET request to show the initial setup form
func (h *InitialSetupHandler) New(c *gin.Context) {
	// Set feature category
	c.Set("feature_category", FeatureCategory)

	// Check if in initial setup state
	if !h.isInInitialSetupState(c) {
		c.Redirect(http.StatusFound, "/")
		c.SetFlash("notice", "Initial setup complete!")
		return
	}

	// Get the last admin user
	user, err := h.userService.GetLastAdminUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get admin user"})
		return
	}

	// Set user in context for template
	c.Set("user", user)

	// Render the initial setup form
	c.HTML(http.StatusOK, "admin/initial_setup/new", gin.H{
		"user": user,
	})
}

// Update handles the POST request to update the initial user account
func (h *InitialSetupHandler) Update(c *gin.Context) {
	// Set feature category
	c.Set("feature_category", FeatureCategory)

	// Check if in initial setup state
	if !h.isInInitialSetupState(c) {
		c.Redirect(http.StatusFound, "/")
		c.SetFlash("notice", "Initial setup complete!")
		return
	}

	var params UserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	// Get the last admin user
	user, err := h.userService.GetLastAdminUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get admin user"})
		return
	}

	// Update user
	result, err := h.userService.UpdateUser(c, user.ID, &model.User{
		Email:     params.Email,
		Password:  params.Password,
		SkipReconfirmation: true,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Clean up non-primary emails
	if err := h.userService.CleanupNonPrimaryEmails(c, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup emails"})
		return
	}

	// Redirect to login page
	c.Redirect(http.StatusFound, "/users/sign_in")
	c.SetFlash("notice", "Initial account configured! Please sign in.")
}

// isInInitialSetupState checks if the system is in initial setup state
func (h *InitialSetupHandler) isInInitialSetupState(c *gin.Context) bool {
	// TODO: Implement proper initial setup check
	// This should check:
	// 1. If the system is in initial setup mode
	// 2. If there are any admin users
	// 3. If the initial setup is required
	return true
}
