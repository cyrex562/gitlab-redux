// Ported from: orig_app/app/controllers/import/history_controller.rb
// Ported on: 2025-05-01
// This file implements the Import::HistoryController from the Ruby codebase.
//
// See porting_log.txt for details.

use actix_web::{HttpResponse, Responder};

/// Controller for import history actions
pub struct HistoryController;

impl HistoryController {
    /// GET /import/history
    pub async fn index() -> impl Responder {
        // No-op: Ruby controller is empty, so just return 200 OK
        HttpResponse::Ok().finish()
    }
}
