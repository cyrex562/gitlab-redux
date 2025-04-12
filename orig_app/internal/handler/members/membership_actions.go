package members

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// MembershipActions handles membership-related actions
type MembershipActions struct {
	membersPresentation *MembersPresentation
	membersService      *service.MembersService
	authService         *service.AuthService
	logger              *service.Logger
}

// NewMembershipActions creates a new instance of MembershipActions
func NewMembershipActions(
	membersPresentation *MembersPresentation,
	membersService *service.MembersService,
	authService *service.AuthService,
	logger *service.Logger,
) *MembershipActions {
	return &MembershipActions{
		membersPresentation: membersPresentation,
		membersService:      membersService,
		authService:         authService,
		logger:              logger,
	}
}

// SetupMiddleware sets up the membership actions middleware
func (m *MembershipActions) SetupMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authenticate user for request_access
		if ctx.Request.URL.Path == "/request_access" {
			if err := m.authenticateUser(ctx); err != nil {
				return
			}

			if err := m.alreadyAMember(ctx); err != nil {
				return
			}
		}

		ctx.Next()
	}
}

// Update updates a member
func (m *MembershipActions) Update(ctx *gin.Context) error {
	// Get member ID from params
	memberID := ctx.Param("id")

	// Find member
	member, err := m.membersAndRequesters(ctx).Find(memberID)
	if err != nil {
		return err
	}

	// Get update params
	updateParams, err := m.updateParams(ctx)
	if err != nil {
		return err
	}

	// Update member
	result, err := m.membersService.Update(ctx, member, updateParams)
	if err != nil {
		return err
	}

	// Check result status
	if result.Status == "success" {
		ctx.JSON(http.StatusOK, m.updateSuccessResponse(result))
	} else {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": result.Message,
		})
	}

	return nil
}

// Destroy removes a member
func (m *MembershipActions) Destroy(ctx *gin.Context) error {
	// Get member ID from params
	memberID := ctx.Param("id")

	// Find member
	member, err := m.membersAndRequesters(ctx).Find(memberID)
	if err != nil {
		return err
	}

	// Get params
	skipSubresources := !ctx.DefaultQuery("remove_sub_memberships", "true")
	unassignIssuables := ctx.DefaultQuery("unassign_issuables", "false") == "true"

	// Destroy member
	err = m.membersService.Destroy(ctx, member, map[string]interface{}{
		"skip_subresources": skipSubresources,
		"unassign_issuables": unassignIssuables,
	})
	if err != nil {
		return err
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "text/html":
		message := m.destroySuccessMessage(skipSubresources)
		ctx.Redirect(http.StatusSeeOther, m.membersPageURL(ctx))
		ctx.Set("notice", message)
	case "application/javascript":
		ctx.Status(http.StatusOK)
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"notice": m.destroySuccessMessage(skipSubresources),
		})
	}

	return nil
}

// RequestAccess requests access to a membershipable
func (m *MembershipActions) RequestAccess(ctx *gin.Context) error {
	// Get current user
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Request access
	accessRequester, err := m.membershipable(ctx).RequestAccess(currentUser)
	if err != nil {
		return err
	}

	// Check if request was successful
	if accessRequester.Persisted() {
		ctx.Redirect(http.StatusSeeOther, m.polymorphicPath(ctx, m.membershipable(ctx)))
		ctx.Set("notice", "Your request for access has been queued for review.")
	} else {
		ctx.Redirect(http.StatusSeeOther, m.polymorphicPath(ctx, m.membershipable(ctx)))
		ctx.Set("alert", "Your request for access could not be processed: " + accessRequester.ErrorMessages())
	}

	return nil
}

// ApproveAccessRequest approves an access request
func (m *MembershipActions) ApproveAccessRequest(ctx *gin.Context) error {
	// Get requester ID from params
	requesterID := ctx.Param("id")

	// Find requester
	accessRequester, err := m.requesters(ctx).Find(requesterID)
	if err != nil {
		return err
	}

	// Approve request
	err = m.membersService.ApproveAccessRequest(ctx, accessRequester, ctx.Request.URL.Query())
	if err != nil {
		return err
	}

	// Redirect to members page
	ctx.Redirect(http.StatusSeeOther, m.membersPageURL(ctx))

	return nil
}

// Leave leaves a membershipable
func (m *MembershipActions) Leave(ctx *gin.Context) error {
	// Get current user
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Find member
	member, err := m.membersAndRequesters(ctx).FindByUserID(currentUser.ID)
	if err != nil {
		return err
	}

	// Destroy member
	err = m.membersService.Destroy(ctx, member, nil)
	if err != nil {
		return err
	}

	// Generate notice
	notice := m.leaveNotice(ctx, member)

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "text/html":
		redirectPath := m.leaveRedirectPath(ctx, member)
		ctx.Redirect(http.StatusSeeOther, redirectPath)
		ctx.Set("notice", notice)
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"notice": notice,
		})
	}

	return nil
}

// ResendInvite resends an invitation
func (m *MembershipActions) ResendInvite(ctx *gin.Context) error {
	// Get member ID from params
	memberID := ctx.Param("id")

	// Find member
	member, err := m.membershipableMembers(ctx).Find(memberID)
	if err != nil {
		return err
	}

	// Check if member is invited
	if member.Invite() {
		// Resend invite
		err = member.ResendInvite()
		if err != nil {
			return err
		}

		// Redirect to members page
		ctx.Redirect(http.StatusSeeOther, m.membersPageURL(ctx))
		ctx.Set("notice", "The invitation was successfully resent.")
	} else {
		// Redirect to members page
		ctx.Redirect(http.StatusSeeOther, m.membersPageURL(ctx))
		ctx.Set("alert", "The invitation has already been accepted.")
	}

	return nil
}

// Protected methods

// Membershipable returns the membershipable
func (m *MembershipActions) Membershipable(ctx *gin.Context) interface{} {
	panic("Not implemented")
}

// MembershipableMembers returns the membershipable members
func (m *MembershipActions) MembershipableMembers(ctx *gin.Context) interface{} {
	panic("Not implemented")
}

// RootParamsKey returns the root params key
func (m *MembershipActions) RootParamsKey() string {
	panic("Not implemented")
}

// MembersPageURL returns the members page URL
func (m *MembershipActions) MembersPageURL(ctx *gin.Context) string {
	panic("Not implemented")
}

// SourceType returns the source type
func (m *MembershipActions) SourceType() string {
	panic("Not implemented")
}

// PlainSourceType returns the plain source type
func (m *MembershipActions) PlainSourceType() string {
	panic("Not implemented")
}

// Source returns the source
func (m *MembershipActions) Source() interface{} {
	panic("Not implemented")
}

// MembersAndRequesters returns the members and requesters
func (m *MembershipActions) MembersAndRequesters(ctx *gin.Context) interface{} {
	return m.membershipable(ctx).(interface{ MembersAndRequesters() interface{} }).MembersAndRequesters()
}

// Requesters returns the requesters
func (m *MembershipActions) Requesters(ctx *gin.Context) interface{} {
	return m.membershipable(ctx).(interface{ Requesters() interface{} }).Requesters()
}

// UpdateParams returns the update params
func (m *MembershipActions) UpdateParams(ctx *gin.Context) (map[string]interface{}, error) {
	rootParamsKey := m.RootParamsKey()
	params := ctx.Request.URL.Query()

	updateParams := map[string]interface{}{
		"access_level": params.Get("access_level"),
		"source":       m.Source(),
	}

	// Parse expires_at if present
	if expiresAt := params.Get("expires_at"); expiresAt != "" {
		parsedTime, err := time.Parse("2006-01-02", expiresAt)
		if err != nil {
			return nil, err
		}
		updateParams["expires_at"] = parsedTime
	}

	return updateParams, nil
}

// RequestedRelations returns the requested relations
func (m *MembershipActions) RequestedRelations(ctx *gin.Context, inheritedPermissions string) []string {
	params := ctx.Request.URL.Query()
	inheritedPermissionsValue := params.Get(inheritedPermissions)

	switch inheritedPermissionsValue {
	case "exclude":
		return []string{"direct"}
	case "only":
		return append([]string{"inherited"}, m.sharedMembersRelations(ctx)...)
	default:
		return append([]string{"inherited", "direct"}, m.sharedMembersRelations(ctx)...)
	}
}

// AuthenticateUser authenticates the user
func (m *MembershipActions) AuthenticateUser(ctx *gin.Context) error {
	// Get current user
	currentUser, err := ctx.Get("current_user")
	if err != nil || currentUser == nil {
		// Store location
		ctx.Set("return_to", ctx.Request.URL.Path)

		// Redirect to login
		ctx.Redirect(http.StatusSeeOther, "/users/sign_in")
		return nil
	}

	return nil
}

// AlreadyAMember checks if the user is already a member
func (m *MembershipActions) AlreadyAMember(ctx *gin.Context) error {
	// Get current user
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Check if user is a member
	member := m.members(ctx).WithUser(currentUser)
	if member != nil {
		ctx.Redirect(http.StatusSeeOther, m.polymorphicPath(ctx, m.membershipable(ctx)))
		ctx.Set("notice", "You already have access.")
		return nil
	}

	// Check if user is a requester
	requester := m.requesters(ctx).WithUser(currentUser)
	if requester != nil {
		ctx.Redirect(http.StatusSeeOther, m.polymorphicPath(ctx, m.membershipable(ctx)))
		ctx.Set("notice", "You have already requested access.")
		return nil
	}

	return nil
}

// Private methods

// UpdateSuccessResponse returns the update success response
func (m *MembershipActions) UpdateSuccessResponse(result *service.MembersUpdateResult) map[string]interface{} {
	member := result.Members[0]
	if member.Expires() {
		return map[string]interface{}{
			"expires_soon":         member.ExpiresSoon(),
			"expires_at_formatted": member.ExpiresAt().Format("Jan 02, 2006 15:04"),
		}
	}
	return map[string]interface{}{}
}

// SharedMembersRelations returns the shared members relations
func (m *MembershipActions) SharedMembersRelations(ctx *gin.Context) []string {
	params := ctx.Request.URL.Query()
	projectID := params.Get("project_id")

	relations := []string{"shared_from_groups"}
	if projectID != "" {
		relations = append(relations, "invited_groups", "shared_into_ancestors")
	}

	return relations
}

// DestroySuccessMessage returns the destroy success message
func (m *MembershipActions) DestroySuccessMessage(skipSubresources bool) string {
	membershipable := m.membershipable(nil)
	sourceType := m.SourceType()

	switch membershipable.(type) {
	case *service.Namespace:
		if skipSubresources {
			return "User was successfully removed from group."
		}
		return "User was successfully removed from group and any subgroups and projects."
	default:
		return "User was successfully removed from project."
	}
}

// LeaveNotice returns the leave notice
func (m *MembershipActions) LeaveNotice(ctx *gin.Context, member *service.Member) string {
	membershipable := m.membershipable(ctx)
	sourceType := m.SourceType()

	if member.Request() {
		return "Your access request to the " + sourceType + " has been withdrawn."
	}

	return "You left the \"" + membershipable.(interface{ HumanName() string }).HumanName() + "\" " + sourceType + "."
}

// LeaveRedirectPath returns the leave redirect path
func (m *MembershipActions) LeaveRedirectPath(ctx *gin.Context, member *service.Member) string {
	membershipable := m.membershipable(ctx)

	if member.Request() {
		return member.Source().(string)
	}

	return "/dashboard/" + membershipable.(interface{ Class() string }).Class()
}

// PolymorphicPath returns the polymorphic path
func (m *MembershipActions) PolymorphicPath(ctx *gin.Context, record interface{}) string {
	// This is a simplified implementation
	// In a real application, this would use a more sophisticated routing system
	switch record.(type) {
	case *service.Project:
		return "/projects/" + record.(*service.Project).ID
	case *service.Group:
		return "/groups/" + record.(*service.Group).ID
	default:
		return "/"
	}
}
