// Ported from orig_app/app/controllers/concerns/accepts_pending_invitations.rb
// Provides a trait for accepting pending invitations for a user in a controller context.
// Ported: 2025-04-24

pub trait AcceptsPendingInvitations {
    // The user type must implement these methods:
    // - active_for_authentication() -> bool
    // - pending_invitations() -> &PendingInvitations
    // - accept_pending_invitations(&mut self)
    type User;
    fn resource(&self) -> &Self::User;
    fn resource_mut(&mut self) -> &mut Self::User;

    fn accept_pending_invitations(&mut self, user: Option<&mut Self::User>) {
        let user = match user {
            Some(u) => u,
            None => self.resource_mut(),
        };
        if !user.active_for_authentication() {
            return;
        }
        if user.pending_invitations().load_any() {
            user.accept_pending_invitations();
            self.after_pending_invitations_hook();
        }
    }

    fn after_pending_invitations_hook(&mut self) {
        // no-op by default
    }
}
