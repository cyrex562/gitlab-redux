package branch

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// DefaultBranchProtection handles default branch protection settings
type DefaultBranchProtection struct {
	accessService *service.AccessService
	logger        *util.Logger
}

// NewDefaultBranchProtection creates a new instance of DefaultBranchProtection
func NewDefaultBranchProtection(
	accessService *service.AccessService,
	logger *util.Logger,
) *DefaultBranchProtection {
	return &DefaultBranchProtection{
		accessService: accessService,
		logger:        logger,
	}
}

// NormalizeDefaultBranchParams normalizes the default branch protection parameters
func (d *DefaultBranchProtection) NormalizeDefaultBranchParams(ctx *gin.Context, formKey string) map[string]interface{} {
	// Get the entity settings parameters from the form
	entitySettingsParams, exists := ctx.GetPostFormMap(formKey)
	if !exists {
		return make(map[string]interface{})
	}

	// Convert the map to a more usable format
	params := make(map[string]interface{})
	for key, value := range entitySettingsParams {
		params[key] = value
	}

	// Check if default branch protection is disabled
	defaultBranchProtected, exists := params["default_branch_protected"]
	if exists {
		isProtected, err := strconv.ParseBool(defaultBranchProtected.(string))
		if err == nil && !isProtected {
			// Set the protection to none
			params["default_branch_protection_defaults"] = d.accessService.GetProtectionNone()
			return params
		}
	}

	// Check if default branch protection defaults are present
	defaultBranchProtectionDefaults, exists := params["default_branch_protection_defaults"]
	if !exists {
		return params
	}

	// Convert to a map for easier manipulation
	defaults, ok := defaultBranchProtectionDefaults.(map[string]interface{})
	if !ok {
		return params
	}

	// Remove the protection level
	delete(params, "default_branch_protection_level")

	// Process allowed to push entries
	if allowedToPush, ok := defaults["allowed_to_push"]; ok {
		if entries, ok := allowedToPush.([]interface{}); ok {
			for i, entry := range entries {
				if entryMap, ok := entry.(map[string]interface{}); ok {
					if accessLevel, ok := entryMap["access_level"]; ok {
						// Convert access level to integer
						level, err := strconv.Atoi(accessLevel.(string))
						if err == nil {
							entryMap["access_level"] = level
							entries[i] = entryMap
						}
					}
				}
			}
			defaults["allowed_to_push"] = entries
		}
	}

	// Process allowed to merge entries
	if allowedToMerge, ok := defaults["allowed_to_merge"]; ok {
		if entries, ok := allowedToMerge.([]interface{}); ok {
			for i, entry := range entries {
				if entryMap, ok := entry.(map[string]interface{}); ok {
					if accessLevel, ok := entryMap["access_level"]; ok {
						// Convert access level to integer
						level, err := strconv.Atoi(accessLevel.(string))
						if err == nil {
							entryMap["access_level"] = level
							entries[i] = entryMap
						}
					}
				}
			}
			defaults["allowed_to_merge"] = entries
		}
	}

	// Process boolean fields
	booleanFields := []string{"allow_force_push", "code_owner_approval_required", "developer_can_initial_push"}
	for _, key := range booleanFields {
		if value, ok := defaults[key]; ok {
			// Convert to boolean
			boolValue, err := strconv.ParseBool(value.(string))
			if err == nil {
				defaults[key] = boolValue
			} else {
				// Use the default value from the fully protected settings
				defaults[key] = d.accessService.GetProtectedFully()[key]
			}
		}
	}

	// Update the params with the processed defaults
	params["default_branch_protection_defaults"] = defaults

	return params
}
