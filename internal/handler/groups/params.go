package groups

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// Params provides group parameter validation and processing
type Params struct {
	// Add any dependencies here if needed
}

// NewParams creates a new instance of Params
func NewParams() *Params {
	return &Params{}
}

// GroupParams represents the parameters for creating or updating a group
type GroupParams struct {
	Avatar                              *string `json:"avatar" binding:"omitempty"`
	Description                         *string `json:"description" binding:"omitempty"`
	EmailsDisabled                      *bool   `json:"emails_disabled" binding:"omitempty"`
	EmailsEnabled                       *bool   `json:"emails_enabled" binding:"omitempty"`
	ShowDiffPreviewInEmail             *bool   `json:"show_diff_preview_in_email" binding:"omitempty"`
	MentionsDisabled                    *bool   `json:"mentions_disabled" binding:"omitempty"`
	LfsEnabled                         *bool   `json:"lfs_enabled" binding:"omitempty"`
	Name                               *string `json:"name" binding:"required,min=1,max=255"`
	Path                               *string `json:"path" binding:"required,min=1,max=255,alphanum"`
	Public                             *bool   `json:"public" binding:"omitempty"`
	RequestAccessEnabled               *bool   `json:"request_access_enabled" binding:"omitempty"`
	ShareWithGroupLock                 *bool   `json:"share_with_group_lock" binding:"omitempty"`
	VisibilityLevel                    *int    `json:"visibility_level" binding:"omitempty,oneof=0 10 20"`
	ParentID                           *int64  `json:"parent_id" binding:"omitempty"`
	CreateChatTeam                     *bool   `json:"create_chat_team" binding:"omitempty"`
	ChatTeamName                       *string `json:"chat_team_name" binding:"omitempty"`
	RequireTwoFactorAuthentication     *bool   `json:"require_two_factor_authentication" binding:"omitempty"`
	TwoFactorGracePeriod              *int    `json:"two_factor_grace_period" binding:"omitempty,min=0"`
	EnabledGitAccessProtocol          *string `json:"enabled_git_access_protocol" binding:"omitempty,oneof=ssh http"`
	ProjectCreationLevel              *int    `json:"project_creation_level" binding:"omitempty,oneof=0 1 2"`
	SubgroupCreationLevel             *int    `json:"subgroup_creation_level" binding:"omitempty,oneof=0 1 2"`
	DefaultBranchProtection           *int    `json:"default_branch_protection" binding:"omitempty,oneof=0 1 2 3"`
	DefaultBranchProtectionDefaults   *DefaultBranchProtectionDefaults `json:"default_branch_protection_defaults" binding:"omitempty"`
	DefaultBranchName                 *string `json:"default_branch_name" binding:"omitempty"`
	AllowMfaForSubgroups              *bool   `json:"allow_mfa_for_subgroups" binding:"omitempty"`
	ResourceAccessTokenCreationAllowed *bool   `json:"resource_access_token_creation_allowed" binding:"omitempty"`
	ResourceAccessTokenNotifyInherited *bool   `json:"resource_access_token_notify_inherited" binding:"omitempty"`
	LockResourceAccessTokenNotifyInherited *bool `json:"lock_resource_access_token_notify_inherited" binding:"omitempty"`
	PreventSharingGroupsOutsideHierarchy *bool `json:"prevent_sharing_groups_outside_hierarchy" binding:"omitempty"`
	SetupForCompany                    *bool   `json:"setup_for_company" binding:"omitempty"`
	JobsToBeDone                       *string `json:"jobs_to_be_done" binding:"omitempty"`
	CrmEnabled                         *bool   `json:"crm_enabled" binding:"omitempty"`
	CrmSourceGroupID                   *int64  `json:"crm_source_group_id" binding:"omitempty"`
	ForcePagesAccessControl            *bool   `json:"force_pages_access_control" binding:"omitempty"`
	EnableNamespaceDescendantsCache    *bool   `json:"enable_namespace_descendants_cache" binding:"omitempty"`
}

// DefaultBranchProtectionDefaults represents the default branch protection settings
type DefaultBranchProtectionDefaults struct {
	AllowForcePush           *bool `json:"allow_force_push" binding:"omitempty"`
	DeveloperCanInitialPush  *bool `json:"developer_can_initial_push" binding:"omitempty"`
	CodeOwnerApprovalRequired *bool `json:"code_owner_approval_required" binding:"omitempty"`
	AllowedToMerge          *AccessLevel `json:"allowed_to_merge" binding:"omitempty"`
	AllowedToPush           *AccessLevel `json:"allowed_to_push" binding:"omitempty"`
}

// AccessLevel represents the access level for branch protection
type AccessLevel struct {
	AccessLevel *int `json:"access_level" binding:"omitempty,oneof=0 10 20 30 40 50"`
}

// GetGroupParams processes and validates group parameters from the request
func (p *Params) GetGroupParams(ctx *gin.Context) (*GroupParams, error) {
	var params GroupParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group parameters", "details": err.Error()})
		return nil, err
	}

	// Validate required fields for creation
	if ctx.Request.Method == http.MethodPost {
		if err := p.validateRequiredFields(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields", "details": err.Error()})
			return nil, err
		}
	}

	// Normalize default branch parameters
	if err := p.normalizeDefaultBranchParams(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid default branch parameters", "details": err.Error()})
		return nil, err
	}

	// Validate visibility level
	if err := p.validateVisibilityLevel(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visibility level", "details": err.Error()})
		return nil, err
	}

	return &params, nil
}

// validateRequiredFields checks that all required fields are present
func (p *Params) validateRequiredFields(params *GroupParams) error {
	if params.Name == nil || *params.Name == "" {
		return fmt.Errorf("name is required")
	}
	if params.Path == nil || *params.Path == "" {
		return fmt.Errorf("path is required")
	}
	return nil
}

// validateVisibilityLevel validates the visibility level
func (p *Params) validateVisibilityLevel(params *GroupParams) error {
	if params.VisibilityLevel == nil {
		return nil
	}

	validLevels := map[int]bool{
		0:  true, // Private
		10: true, // Internal
		20: true, // Public
	}

	if !validLevels[*params.VisibilityLevel] {
		return fmt.Errorf("invalid visibility level: %d", *params.VisibilityLevel)
	}

	return nil
}

// normalizeDefaultBranchParams normalizes the default branch protection parameters
func (p *Params) normalizeDefaultBranchParams(params *GroupParams) error {
	if params.DefaultBranchProtectionDefaults == nil {
		return nil
	}

	// Set default values if not provided
	if params.DefaultBranchProtectionDefaults.AllowedToMerge == nil {
		params.DefaultBranchProtectionDefaults.AllowedToMerge = &AccessLevel{
			AccessLevel: &model.AccessLevelDeveloper,
		}
	}

	if params.DefaultBranchProtectionDefaults.AllowedToPush == nil {
		params.DefaultBranchProtectionDefaults.AllowedToPush = &AccessLevel{
			AccessLevel: &model.AccessLevelDeveloper,
		}
	}

	// Validate access levels
	if err := p.validateAccessLevels(params.DefaultBranchProtectionDefaults); err != nil {
		return err
	}

	return nil
}

// validateAccessLevels validates the access levels for branch protection
func (p *Params) validateAccessLevels(defaults *DefaultBranchProtectionDefaults) error {
	validLevels := map[int]bool{
		0:  true, // No access
		10: true, // Guest
		20: true, // Reporter
		30: true, // Developer
		40: true, // Maintainer
		50: true, // Owner
	}

	if defaults.AllowedToMerge != nil && defaults.AllowedToMerge.AccessLevel != nil {
		if !validLevels[*defaults.AllowedToMerge.AccessLevel] {
			return fmt.Errorf("invalid merge access level: %d", *defaults.AllowedToMerge.AccessLevel)
		}
	}

	if defaults.AllowedToPush != nil && defaults.AllowedToPush.AccessLevel != nil {
		if !validLevels[*defaults.AllowedToPush.AccessLevel] {
			return fmt.Errorf("invalid push access level: %d", *defaults.AllowedToPush.AccessLevel)
		}
	}

	return nil
}

// GetGroupFeatureAttributes returns the list of feature attributes that can be modified
func (p *Params) GetGroupFeatureAttributes() []string {
	// This can be extended in EE
	return []string{}
}
