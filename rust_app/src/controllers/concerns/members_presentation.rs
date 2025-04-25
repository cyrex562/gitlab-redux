// Ported from: orig_app/app/controllers/concerns/members_presentation.rb (Ruby)
// Ported on: 2025-04-25
// This file implements the MembersPresentation concern in Rust.

use crate::models::user::User;
use crate::presenters::members_presenter::MembersPresenter;
use crate::services::members_preloader::MembersPreloader;
use crate::view::presenter::factory::PresenterFactory;

pub trait MembersPresentation {
    fn present_members(&self, members: Vec<User>) -> Vec<MembersPresenter>;
    fn preload_associations(&self, members: &[User]);
}

pub struct MembersPresentationImpl {
    current_user: User,
}

impl MembersPresentationImpl {
    pub fn new(current_user: User) -> Self {
        Self { current_user }
    }
}

impl MembersPresentation for MembersPresentationImpl {
    fn present_members(&self, members: Vec<User>) -> Vec<MembersPresenter> {
        // Preload associations for all members
        self.preload_associations(&members);

        // Create and return presenters for all members
        PresenterFactory::new(members, self.current_user.clone(), MembersPresenter::new).fabricate()
    }

    fn preload_associations(&self, members: &[User]) {
        MembersPreloader::new(members).preload_all();
    }
}
