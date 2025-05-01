// Ported from: orig_app/app/controllers/groups/boards_controller.rb
// Ported: 2025-05-01

use crate::controllers::concerns::{BoardsActions, RecordUserLastActivity};
use crate::models::{Group, User};
use crate::services::boards::{BoardsFinder, CreateService};
use actix_web::HttpRequest;
use std::sync::{Arc, Mutex};

/// Controller for group boards (port of Groups::BoardsController)
pub struct GroupsBoardsController {
    pub group: Arc<Group>,
    pub current_user: Arc<User>,
    pub request: HttpRequest,
    board_finder: Mutex<Option<BoardsFinder>>, // strong memoize
    board_create_service: Mutex<Option<CreateService>>, // strong memoize
}

impl GroupsBoardsController {
    pub fn new(group: Arc<Group>, current_user: Arc<User>, request: HttpRequest) -> Self {
        Self {
            group,
            current_user,
            request,
            board_finder: Mutex::new(None),
            board_create_service: Mutex::new(None),
        }
    }

    // Mimics before_action for feature flags
    pub fn push_feature_flags(&self) {
        // TODO: Implement push_frontend_feature_flag(:board_multi_select, group)
        // TODO: Implement push_frontend_feature_flag(:issues_list_drawer, group)
        // TODO: Implement push_force_frontend_feature_flag(:work_items_beta, group.work_items_beta_feature_flag_enabled())
    }

    // Memoized board finder
    pub fn board_finder(&self, board_id: i64) -> Arc<BoardsFinder> {
        let mut finder = self.board_finder.lock().unwrap();
        if finder.is_none() {
            *finder = Some(BoardsFinder::new(
                self.group.clone(),
                self.current_user.clone(),
                board_id,
            ));
        }
        Arc::new(finder.as_ref().unwrap().clone())
    }

    // Memoized board create service
    pub fn board_create_service(&self) -> Arc<CreateService> {
        let mut svc = self.board_create_service.lock().unwrap();
        if svc.is_none() {
            *svc = Some(CreateService::new(
                self.group.clone(),
                self.current_user.clone(),
            ));
        }
        Arc::new(svc.as_ref().unwrap().clone())
    }

    // Authorization logic
    pub fn authorize_read_board(&self) -> bool {
        // TODO: Replace with real permission check
        // can?(current_user, :read_issue_board, group)
        true
    }
}

// Implement BoardsActions and RecordUserLastActivity traits for controller
#[async_trait::async_trait]
impl BoardsActions for GroupsBoardsController {
    // ...implement trait methods as needed...
}

impl RecordUserLastActivity for GroupsBoardsController {
    fn current_user(&self) -> Option<&User> {
        Some(&self.current_user)
    }
    fn request(&self) -> &HttpRequest {
        &self.request
    }
    fn group(&self) -> Option<&Group> {
        Some(&self.group)
    }
    fn is_db_read_only(&self) -> bool {
        false // TODO: implement real check
    }
}
