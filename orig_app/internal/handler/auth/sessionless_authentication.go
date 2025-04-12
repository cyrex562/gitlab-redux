package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// SessionlessAuthentication handles authentication methods that don't require a session
// such as Personal Access Tokens (PAT), RSS tokens, and static object tokens
type SessionlessAuthentication struct {
	authService *service.AuthService
	userService *service.UserService
	settings    *service.SettingsService
}

// NewSessionlessAuthentication creates a new instance of SessionlessAuthentication
func NewSessionlessAuthentication(
	authService *service.AuthService,
	userService *service.UserService,
	settings *service.SettingsService,
) *SessionlessAuthentication {
	return &SessionlessAuthentication{
		authService: authService,
		userService: userService,
		settings:    settings,
	}
}

// AuthenticateSessionlessUser authenticates a user without creating a session
// This handles personal access tokens, atom requests with RSS tokens, and static object tokens
func (s *SessionlessAuthentication) AuthenticateSessionlessUser(c *gin.Context, requestFormat string) {
	user := s.getRequestAuthenticator(c).FindSessionlessUser(requestFormat)
	if user != nil {
		s.sessionlessSignIn(c, user)
	}
}

// IsSessionlessUser checks if the current user was authenticated without a session
func (s *SessionlessAuthentication) IsSessionlessUser(c *gin.Context) bool {
	currentUser := s.getCurrentUser(c)
	sessionlessSignIn := c.GetBool("sessionless_sign_in")
	return currentUser != nil && sessionlessSignIn
}

// SessionlessSignIn signs in a user without creating a session
func (s *SessionlessAuthentication) SessionlessSignIn(c *gin.Context, user *model.User) {
	// Mark that this is a sessionless sign-in
	c.Set("sessionless_sign_in", true)

	if s.userService.CanLogInWithNonExpiredPassword(user) {
		// Sign in the user without storing in session
		// A token will be needed for every request
		s.authService.SignIn(c, user, false, "sessionless_sign_in", true)
	} else if s.getRequestAuthenticator(c).CanSignInBot(user) {
		// Sign in the bot user without storing in session and without callbacks
		s.authService.SignIn(c, user, false, "sessionless_sign_in", false)
	}
}

// SessionlessBypassAdminMode executes a block of code while bypassing admin mode
func (s *SessionlessAuthentication) SessionlessBypassAdminMode(c *gin.Context, block func() error) error {
	if !s.settings.IsAdminMode() {
		return block()
	}

	currentUser := s.getCurrentUser(c)
	if currentUser == nil {
		return block()
	}

	return s.authService.BypassSession(c, currentUser.ID, block)
}

// Private helper methods

func (s *SessionlessAuthentication) getRequestAuthenticator(c *gin.Context) *service.RequestAuthenticator {
	// Get or create the request authenticator
	authenticator, exists := c.Get("request_authenticator")
	if exists {
		return authenticator.(*service.RequestAuthenticator)
	}

	newAuthenticator := service.NewRequestAuthenticator(c.Request)
	c.Set("request_authenticator", newAuthenticator)
	return newAuthenticator
}

func (s *SessionlessAuthentication) getCurrentUser(c *gin.Context) *model.User {
	// Get the current user from the context
	user, exists := c.Get("current_user")
	if !exists {
		return nil
	}
	return user.(*model.User)
}

func (s *SessionlessAuthentication) sessionlessSignIn(c *gin.Context, user *model.User) {
	s.SessionlessSignIn(c, user)
}
