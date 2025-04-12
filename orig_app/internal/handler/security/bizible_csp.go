package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// BizibleCSP handles Content Security Policy configuration for Bizible integration
type BizibleCSP struct {
	configService *service.ConfigService
	helperService *service.HelperService
	logger        *util.Logger
}

// NewBizibleCSP creates a new instance of BizibleCSP
func NewBizibleCSP(
	configService *service.ConfigService,
	helperService *service.HelperService,
	logger *util.Logger,
) *BizibleCSP {
	return &BizibleCSP{
		configService: configService,
		helperService: helperService,
		logger:        logger,
	}
}

// ApplyBizibleCSP applies the Bizible CSP to the response
func (b *BizibleCSP) ApplyBizibleCSP(ctx *gin.Context) {
	// Check if Bizible is enabled or if there are existing directives
	if !b.helperService.BizibleEnabled() && len(ctx.GetHeader("Content-Security-Policy")) == 0 {
		return
	}

	// Get the existing CSP directives
	existingDirectives := b.getExistingDirectives(ctx)

	// Get the default script-src or default-src
	defaultScriptSrc := existingDirectives["script-src"]
	if defaultScriptSrc == nil {
		defaultScriptSrc = existingDirectives["default-src"]
	}

	// Convert the default script-src to an array if it's not already
	scriptSrcValues := b.convertToArray(defaultScriptSrc)

	// Add the Bizible script-src values
	scriptSrcValues = append(scriptSrcValues, "'unsafe-eval'", "https://cdn.bizible.com/scripts/bizible.js")

	// Set the script-src directive
	existingDirectives["script-src"] = scriptSrcValues

	// Set the Content-Security-Policy header
	b.setCSPHeader(ctx, existingDirectives)
}

// getExistingDirectives gets the existing CSP directives from the response
func (b *BizibleCSP) getExistingDirectives(ctx *gin.Context) map[string][]string {
	// Get the existing CSP header
	cspHeader := ctx.GetHeader("Content-Security-Policy")
	if len(cspHeader) == 0 {
		return make(map[string][]string)
	}

	// Parse the CSP header into directives
	directives := make(map[string][]string)
	// This is a simplified implementation
	// In a real implementation, you would parse the CSP header properly
	return directives
}

// convertToArray converts a value to an array if it's not already
func (b *BizibleCSP) convertToArray(value interface{}) []string {
	if value == nil {
		return []string{}
	}

	switch v := value.(type) {
	case []string:
		return v
	case string:
		return []string{v}
	default:
		return []string{}
	}
}

// setCSPHeader sets the Content-Security-Policy header
func (b *BizibleCSP) setCSPHeader(ctx *gin.Context, directives map[string][]string) {
	// Build the CSP header
	cspHeader := ""
	for directive, values := range directives {
		if len(values) > 0 {
			if len(cspHeader) > 0 {
				cspHeader += "; "
			}
			cspHeader += directive + " " + b.joinValues(values)
		}
	}

	// Set the Content-Security-Policy header
	ctx.Header("Content-Security-Policy", cspHeader)
}

// joinValues joins an array of values with spaces
func (b *BizibleCSP) joinValues(values []string) string {
	result := ""
	for i, value := range values {
		if i > 0 {
			result += " "
		}
		result += value
	}
	return result
}
