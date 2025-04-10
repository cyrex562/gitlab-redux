package security

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

const (
	KnownSignInCookie     = "known_sign_in"
	KnownSignInCookieExpiry = 14 * 24 * time.Hour // 14 days
)

// KnownSignIn handles known sign-in functionality
type KnownSignIn struct {
	configService    *service.ConfigService
	sessionService   *service.SessionService
	notificationService *service.NotificationService
	logger          *service.Logger
}

// NewKnownSignIn creates a new instance of KnownSignIn
func NewKnownSignIn(
	configService *service.ConfigService,
	sessionService *service.SessionService,
	notificationService *service.NotificationService,
	logger *service.Logger,
) *KnownSignIn {
	return &KnownSignIn{
		configService:    configService,
		sessionService:   sessionService,
		notificationService: notificationService,
		logger:          logger,
	}
}

// VerifyKnownSignIn verifies if the current sign-in is from a known device or IP
func (k *KnownSignIn) VerifyKnownSignIn(ctx *gin.Context) error {
	// Check if notifications for unknown sign-ins are enabled
	if !k.configService.NotifyOnUnknownSignIn() {
		return nil
	}

	// Get current user from context
	currentUser, err := k.sessionService.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return nil
	}

	// Check if the sign-in is from a known device or IP
	if k.knownDevice(ctx) || k.knownRemoteIP(ctx) {
		return nil
	}

	// Notify user about unknown sign-in
	if err := k.notifyUser(ctx, currentUser); err != nil {
		return err
	}

	// Update the cookie
	return k.updateCookie(ctx, currentUser)
}

// knownRemoteIP checks if the current IP is known
func (k *KnownSignIn) knownRemoteIP(ctx *gin.Context) bool {
	knownIPs := k.knownIPAddresses(ctx)
	currentIP := ctx.ClientIP()

	for _, ip := range knownIPs {
		if ip == currentIP {
			return true
		}
	}
	return false
}

// knownDevice checks if the current device is known
func (k *KnownSignIn) knownDevice(ctx *gin.Context) bool {
	cookie, err := ctx.Cookie(KnownSignInCookie)
	if err != nil {
		return false
	}

	currentUser, err := k.sessionService.GetCurrentUser(ctx)
	if err != nil {
		return false
	}

	return cookie == currentUser.ID
}

// updateCookie updates the known sign-in cookie
func (k *KnownSignIn) updateCookie(ctx *gin.Context, user *service.User) error {
	return ctx.SetCookie(
		KnownSignInCookie,
		user.ID,
		int(KnownSignInCookieExpiry.Seconds()),
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)
}

// sessions gets the list of active sessions for the current user
func (k *KnownSignIn) sessions(ctx *gin.Context) ([]*service.Session, error) {
	currentUser, err := k.sessionService.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	sessions, err := k.sessionService.ListActiveSessions(currentUser)
	if err != nil {
		return nil, err
	}

	// Filter out impersonated sessions
	var filteredSessions []*service.Session
	for _, session := range sessions {
		if !session.IsImpersonated {
			filteredSessions = append(filteredSessions, session)
		}
	}

	return filteredSessions, nil
}

// knownIPAddresses gets the list of known IP addresses
func (k *KnownSignIn) knownIPAddresses(ctx *gin.Context) []string {
	currentUser, err := k.sessionService.GetCurrentUser(ctx)
	if err != nil {
		return nil
	}

	sessions, err := k.sessions(ctx)
	if err != nil {
		return nil
	}

	// Start with the user's last sign-in IP
	knownIPs := []string{currentUser.LastSignInIP}

	// Add IPs from active sessions
	for _, session := range sessions {
		knownIPs = append(knownIPs, session.IPAddress)
	}

	return knownIPs
}

// notifyUser sends a notification about unknown sign-in
func (k *KnownSignIn) notifyUser(ctx *gin.Context, user *service.User) error {
	requestInfo := &service.VisitorLocation{
		IP: ctx.ClientIP(),
		// Additional visitor location info can be added here
	}

	return k.notificationService.UnknownSignIn(
		user,
		ctx.ClientIP(),
		time.Now(),
		requestInfo,
	)
}
