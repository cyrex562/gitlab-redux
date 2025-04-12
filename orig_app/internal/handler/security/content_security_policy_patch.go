package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// ContentSecurityPolicyPatch provides a patch for content security policy functionality
type ContentSecurityPolicyPatch struct {
	railsService *service.RailsService
	logger       *util.Logger
}

// NewContentSecurityPolicyPatch creates a new instance of ContentSecurityPolicyPatch
func NewContentSecurityPolicyPatch(
	railsService *service.RailsService,
	logger *util.Logger,
) *ContentSecurityPolicyPatch {
	return &ContentSecurityPolicyPatch{
		railsService: railsService,
		logger:       logger,
	}
}

// ContentSecurityPolicyWithContext makes the caller's context available to the policy block
func (c *ContentSecurityPolicyPatch) ContentSecurityPolicyWithContext(
	ctx *gin.Context,
	enabled bool,
	options map[string]interface{},
	block func(policy *ContentSecurityPolicy) error,
) error {
	// Check if we're using Rails 7.2 or later
	if c.railsService.GemVersion().GreaterThanOrEqual(c.railsService.NewVersion("7.2")) {
		c.logger.Warn(
			"content_security_policy_with_context should only be used with Rails < 7.2. " +
				"Use content_security_policy instead.",
		)
	}

	// Get the current content security policy
	policy := c.getCurrentContentSecurityPolicy(ctx)

	// Execute the block with the policy
	if block != nil {
		if err := block(policy); err != nil {
			return err
		}
	}

	// Set the content security policy on the request
	if enabled {
		ctx.Set("content_security_policy", policy)
	} else {
		ctx.Set("content_security_policy", nil)
	}

	return nil
}

// getCurrentContentSecurityPolicy gets the current content security policy
func (c *ContentSecurityPolicyPatch) getCurrentContentSecurityPolicy(ctx *gin.Context) *ContentSecurityPolicy {
	policy, exists := ctx.Get("content_security_policy")
	if !exists {
		policy = NewContentSecurityPolicy()
	}
	return policy.(*ContentSecurityPolicy)
}

// ContentSecurityPolicy represents a content security policy
type ContentSecurityPolicy struct {
	Directives map[string][]string
}

// NewContentSecurityPolicy creates a new content security policy
func NewContentSecurityPolicy() *ContentSecurityPolicy {
	return &ContentSecurityPolicy{
		Directives: make(map[string][]string),
	}
}

// SetDirective sets a directive on the policy
func (p *ContentSecurityPolicy) SetDirective(name string, values []string) {
	p.Directives[name] = values
}

// GetDirective gets a directive from the policy
func (p *ContentSecurityPolicy) GetDirective(name string) []string {
	return p.Directives[name]
}

// ToHeader converts the policy to a header string
func (p *ContentSecurityPolicy) ToHeader() string {
	header := ""
	for name, values := range p.Directives {
		if len(values) > 0 {
			if len(header) > 0 {
				header += "; "
			}
			header += name + " " + joinValues(values)
		}
	}
	return header
}

// joinValues joins an array of values with spaces
func joinValues(values []string) string {
	result := ""
	for i, value := range values {
		if i > 0 {
			result += " "
		}
		result += value
	}
	return result
}
