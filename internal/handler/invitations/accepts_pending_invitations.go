package invitations

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// AcceptsPendingInvitations handles accepting pending invitations for users
type AcceptsPendingInvitations struct{}

// NewAcceptsPendingInvitations creates a new instance of AcceptsPendingInvitations
func NewAcceptsPendingInvitations() *AcceptsPendingInvitations {
	return &AcceptsPendingInvitations{}
}

// AcceptPendingInvitations accepts pending invitations for a user
func (a *AcceptsPendingInvitations) AcceptPendingInvitations(user model.User) {
	// Return early if the user is not active for authentication
	if !user.IsActiveForAuthentication() {
		return
	}

	// Check if the user has any pending invitations
	pendingInvitations := user.GetPendingInvitations()
	if len(pendingInvitations) > 0 {
		// Accept all pending invitations
		user.AcceptPendingInvitations()

		// Call the hook method
		a.AfterPendingInvitationsHook()
	}
}

// AfterPendingInvitationsHook is called after accepting pending invitations
// This is a no-op by default and can be overridden by embedding this struct
func (a *AcceptsPendingInvitations) AfterPendingInvitationsHook() {
	// no-op
}
