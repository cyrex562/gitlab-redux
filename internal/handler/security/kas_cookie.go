package security

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// KasCookie handles KAS (Kubernetes Agent Server) cookie functionality
type KasCookie struct {
	configService *service.ConfigService
	sessionService *service.SessionService
	kasService    *service.KasService
	logger        *service.Logger
}

// NewKasCookie creates a new instance of KasCookie
func NewKasCookie(
	configService *service.ConfigService,
	sessionService *service.SessionService,
	kasService *service.KasService,
	logger *service.Logger,
) *KasCookie {
	return &KasCookie{
		configService:  configService,
		sessionService: sessionService,
		kasService:     kasService,
		logger:         logger,
	}
}

// SetupMiddleware sets up the middleware for KAS cookie functionality
func (k *KasCookie) SetupMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Configure content security policy for KAS
		k.configureContentSecurityPolicy(ctx)

		// Set the KAS cookie
		k.SetKasCookie(ctx)

		// Continue to the next middleware/handler
		ctx.Next()
	}
}

// configureContentSecurityPolicy configures the content security policy for KAS
func (k *KasCookie) configureContentSecurityPolicy(ctx *gin.Context) {
	// Check if KAS user access is enabled
	if !k.kasService.IsUserAccessEnabled() {
		return
	}

	// Check if content security policy is enabled
	if !k.configService.IsContentSecurityPolicyEnabled() {
		return
	}

	// Get the KAS URL
	kasURL := k.GetKasURL()

	// Parse the KAS URL
	parsedKasURL, err := url.Parse(kasURL)
	if err != nil {
		k.logger.Error("Failed to parse KAS URL", "error", err)
		return
	}

	// Get the GitLab host
	gitlabHost := k.configService.GetGitlabHost()

	// Check if the KAS host is the same as the GitLab host
	if parsedKasURL.Host == gitlabHost {
		// Already allowed, no need for exception
		return
	}

	// Get the KAS WebSocket URL
	kasWsURL := k.GetKasWsURL()

	// Remove trailing slash from URLs if present
	kasURL = strings.TrimRight(kasURL, "/")
	kasWsURL = strings.TrimRight(kasWsURL, "/")

	// Get the current content security policy
	csp := ctx.GetHeader("Content-Security-Policy")
	if csp == "" {
		// If no CSP header is set, create a new one
		csp = "default-src 'self'"
	}

	// Add the KAS URLs to the connect-src directive
	connectSrcRegex := regexp.MustCompile(`connect-src\s+([^;]+)`)
	matches := connectSrcRegex.FindStringSubmatch(csp)

	if len(matches) > 1 {
		// Add the KAS URLs to the existing connect-src directive
		connectSrc := matches[1] + " " + kasURL + " " + kasWsURL
		csp = connectSrcRegex.ReplaceAllString(csp, "connect-src "+connectSrc)
	} else {
		// Add a new connect-src directive
		csp += "; connect-src 'self' " + kasURL + " " + kasWsURL
	}

	// Set the updated content security policy header
	ctx.Header("Content-Security-Policy", csp)
}

// SetKasCookie sets the KAS cookie
func (k *KasCookie) SetKasCookie(ctx *gin.Context) {
	// Check if KAS user access is enabled
	if !k.kasService.IsUserAccessEnabled() {
		return
	}

	// Get the current session
	session, err := k.sessionService.GetCurrentSession(ctx)
	if err != nil {
		k.logger.Error("Failed to get current session", "error", err)
		return
	}

	// Check if the session has a public ID
	if session == nil || session.PublicID == "" {
		return
	}

	// Get the cookie data
	cookieData, err := k.kasService.GetCookieData(session.PublicID)
	if err != nil {
		k.logger.Error("Failed to get KAS cookie data", "error", err)
		return
	}

	// Set the cookie
	cookieKey := k.kasService.GetCookieKey()
	ctx.SetCookie(
		cookieKey,
		cookieData,
		0, // MaxAge: 0 means no 'Max-Age' attribute specified
		"/", // Path
		"", // Domain
		false, // Secure
		true, // HttpOnly
	)
}

// GetKasURL gets the KAS URL
func (k *KasCookie) GetKasURL() string {
	return k.kasService.GetTunnelURL()
}

// GetKasWsURL gets the KAS WebSocket URL
func (k *KasCookie) GetKasWsURL() string {
	return k.kasService.GetTunnelWsURL()
}
