// Ported from: orig_app/app/controllers/concerns/boards_actions.rb
// Date: 2025-04-24
// This module provides board-related controller actions and helpers.

use actix_web::{HttpResponse, Responder};
use async_trait::async_trait;
use std::sync::{Arc, Mutex};

// Placeholder types for parent, board, user, and services
pub struct Board;
pub struct User;
pub struct Parent;

pub trait BoardFinder {
    fn execute(&self) -> Vec<Board>;
}

pub trait BoardCreateService {
    fn execute(&self) -> BoardCreateResult;
}

pub struct BoardCreateResult {
    pub payload: Option<Board>,
}

pub trait BoardVisitService {
    fn execute(&self, board: &Board);
}

pub trait BoardsVisitsFinder {
    fn latest(&self) -> Option<BoardVisit>;
}

pub struct BoardVisit {
    pub board: Board,
}

#[async_trait]
pub trait BoardsActions {
    async fn index(&self) -> HttpResponse;
    async fn show(&self) -> HttpResponse;
    fn redirect_to_recent_board(&self) -> Option<HttpResponse>;
    fn latest_visited_board(&self) -> Option<BoardVisit>;
    fn push_licensed_features(&self) {}
    fn board(&self) -> Option<Board>;
    fn board_visit_service(&self) -> Arc<dyn BoardVisitService + Send + Sync>;
    fn parent(&self) -> Arc<Parent>;
    fn board_path(&self, board: &Board) -> String;
    fn group(&self) -> Option<Arc<Parent>>;
    fn project(&self) -> Option<Arc<Parent>>;
    fn group_(&self) -> bool;
}

pub struct BoardsActionsHandler {
    pub current_user: Arc<User>,
    pub parent: Arc<Parent>,
    pub board_finder: Arc<dyn BoardFinder + Send + Sync>,
    pub board_create_service: Arc<dyn BoardCreateService + Send + Sync>,
    pub board_visit_service: Arc<dyn BoardVisitService + Send + Sync>,
    pub boards_visits_finder: Arc<dyn BoardsVisitsFinder + Send + Sync>,
    pub memo_board: Mutex<Option<Board>>,
    pub memo_latest_visited_board: Mutex<Option<BoardVisit>>,
}

#[async_trait]
impl BoardsActions for BoardsActionsHandler {
    async fn index(&self) -> HttpResponse {
        // If no board exists, create one
        let board = self.board();
        if board.is_none() {
            if let Some(new_board) = self.board_create_service.execute().payload {
                let mut memo = self.memo_board.lock().unwrap();
                *memo = Some(new_board);
            }
        }
        HttpResponse::Ok().finish()
    }

    async fn show(&self) -> HttpResponse {
        let board = self.board();
        if board.is_none() {
            return HttpResponse::NotFound().finish();
        }
        // Add/update the board in the recent visits table
        if let Some(ref b) = board {
            self.board_visit_service.execute(b);
        }
        HttpResponse::Ok().finish()
    }

    fn redirect_to_recent_board(&self) -> Option<HttpResponse> {
        // Placeholder: parent.multiple_issue_boards_available? logic
        let parent_has_multiple = true;
        let latest = self.latest_visited_board();
        if !parent_has_multiple || latest.is_none() {
            return None;
        }
        let board = &latest.unwrap().board;
        Some(
            HttpResponse::Found()
                .header("Location", self.board_path(board))
                .finish(),
        )
    }

    fn latest_visited_board(&self) -> Option<BoardVisit> {
        let mut memo = self.memo_latest_visited_board.lock().unwrap();
        if memo.is_none() {
            *memo = self.boards_visits_finder.latest();
        }
        memo.clone()
    }

    fn push_licensed_features(&self) {}

    fn board(&self) -> Option<Board> {
        let mut memo = self.memo_board.lock().unwrap();
        if memo.is_none() {
            let found = self.board_finder.execute().into_iter().next();
            *memo = found;
        }
        memo.clone()
    }

    fn board_visit_service(&self) -> Arc<dyn BoardVisitService + Send + Sync> {
        self.board_visit_service.clone()
    }

    fn parent(&self) -> Arc<Parent> {
        self.parent.clone()
    }

    fn board_path(&self, _board: &Board) -> String {
        // Placeholder: implement group/project path logic
        "/boards/1".to_string()
    }

    fn group(&self) -> Option<Arc<Parent>> {
        None
    }

    fn project(&self) -> Option<Arc<Parent>> {
        None
    }

    fn group_(&self) -> bool {
        false
    }
}
