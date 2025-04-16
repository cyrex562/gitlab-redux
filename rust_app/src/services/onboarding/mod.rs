pub mod redirectable;

use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};

// TODO: Implement OnboardingStatusPresenter
pub trait OnboardingStatusPresenter {
    fn single_invite(&self) -> bool;
    fn last_invited_member(&self) -> Option<&dyn InvitedMember>;
    fn last_invited_member_source(&self) -> Option<&dyn MemberSource>;
}

// TODO: Implement InvitedMember trait
pub trait InvitedMember {
    fn id(&self) -> i32;
    // Add other necessary fields
}

// TODO: Implement MemberSource trait
pub trait MemberSource {
    fn id(&self) -> i32;
    fn type_name(&self) -> &str;
    // Add other necessary fields
}

// TODO: Implement User model
pub struct User {
    pub id: i32,
    // Add other necessary fields
}
