package security

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// OneTrustCSP handles OneTrust Content Security Policy configuration
type OneTrustCSP struct {
	configService *service.ConfigService
	helperService *service.HelperService
	logger        *service.Logger
}

// NewOneTrustCSP creates a new instance of OneTrustCSP
func NewOneTrustCSP(
	configService *service.ConfigService,
	helperService *service.HelperService,
	logger *service.Logger,
) *OneTrustCSP {
	return &OneTrustCSP{
		configService: configService,
		helperService: helperService,
		logger:        logger,
	}
}

// SetupMiddleware sets up the middleware for OneTrust CSP
func (o *OneTrustCSP) SetupMiddleware(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		// Check if OneTrust is enabled or if there are existing directives
		oneTrustEnabled, err := o.helperService.IsOneTrustEnabled(c)
		if err != nil {
			o.logger.Error("Failed to check if OneTrust is enabled", err)
			c.Next()
			return
		}

		// Get existing CSP directives
		existingDirectives := o.getExistingDirectives(c)

		// If OneTrust is not enabled and there are no existing directives, skip
		if !oneTrustEnabled && len(existingDirectives) == 0 {
			c.Next()
			return
		}

		// Configure CSP with OneTrust domains
		o.configureCSP(c, existingDirectives)

		c.Next()
	})
}

// getExistingDirectives gets the existing CSP directives from the response headers
func (o *OneTrustCSP) getExistingDirectives(c *gin.Context) map[string][]string {
	directives := make(map[string][]string)

	// Get existing CSP header
	cspHeader := c.GetHeader("Content-Security-Policy")
	if cspHeader == "" {
		return directives
	}

	// Parse CSP header
	parts := strings.Split(cspHeader, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split directive and values
		directiveParts := strings.SplitN(part, " ", 2)
		if len(directiveParts) != 2 {
			continue
		}

		directive := strings.TrimSpace(directiveParts[0])
		values := strings.Split(directiveParts[1], " ")

		// Clean up values
		cleanValues := make([]string, 0, len(values))
		for _, value := range values {
			value = strings.TrimSpace(value)
			if value != "" {
				cleanValues = append(cleanValues, value)
			}
		}

		directives[directive] = cleanValues
	}

	return directives
}

// configureCSP configures the Content Security Policy with OneTrust domains
func (o *OneTrustCSP) configureCSP(c *gin.Context, existingDirectives map[string][]string) {
	// Get default script-src or default-src
	defaultScriptSrc := existingDirectives["script-src"]
	if len(defaultScriptSrc) == 0 {
		defaultScriptSrc = existingDirectives["default-src"]
	}

	// Add OneTrust domains to script-src
	scriptSrcValues := make([]string, 0, len(defaultScriptSrc)+3)
	scriptSrcValues = append(scriptSrcValues, defaultScriptSrc...)
	scriptSrcValues = append(scriptSrcValues, "'unsafe-eval'", "https://cdn.cookielaw.org", "https://*.onetrust.com")

	// Get default connect-src or default-src
	defaultConnectSrc := existingDirectives["connect-src"]
	if len(defaultConnectSrc) == 0 {
		defaultConnectSrc = existingDirectives["default-src"]
	}

	// Add OneTrust domains to connect-src
	connectSrcValues := make([]string, 0, len(defaultConnectSrc)+2)
	connectSrcValues = append(connectSrcValues, defaultConnectSrc...)
	connectSrcValues = append(connectSrcValues, "https://cdn.cookielaw.org", "https://*.onetrust.com")

	// Build CSP header
	var cspBuilder strings.Builder

	// Add script-src directive
	cspBuilder.WriteString("script-src ")
	cspBuilder.WriteString(strings.Join(scriptSrcValues, " "))
	cspBuilder.WriteString("; ")

	// Add connect-src directive
	cspBuilder.WriteString("connect-src ")
	cspBuilder.WriteString(strings.Join(connectSrcValues, " "))
	cspBuilder.WriteString("; ")

	// Add other existing directives
	for directive, values := range existingDirectives {
		if directive != "script-src" && directive != "connect-src" {
			cspBuilder.WriteString(directive)
			cspBuilder.WriteString(" ")
			cspBuilder.WriteString(strings.Join(values, " "))
			cspBuilder.WriteString("; ")
		}
	}

	// Set CSP header
	c.Header("Content-Security-Policy", cspBuilder.String())
}
