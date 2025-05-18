// Ported from: orig_app/app/controllers/profiles/notifications_controller.rb
// Ported on: 2025-05-07
// This file implements the Profiles::NotificationsController from the Ruby codebase.

use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/profile/notifications")]
pub async fn show() -> impl Responder {
    // TODO: Implement current_user, user_groups, UserGroupNotificationSettingsFinder, project_notifications_with_preloaded_associations, global_notification_setting
    // let user = current_user();
    // let user_groups = user_groups();
    // let group_notifications = UserGroupNotificationSettingsFinder::new(user, user_groups).execute();
    // let project_notifications = project_notifications_with_preloaded_associations();
    // let global_notification_setting = user.global_notification_setting();
    HttpResponse::Ok().body("Notifications page (stub)")
}

#[post("/profile/notifications")]
pub async fn update() -> impl Responder {
    // TODO: Implement Users::UpdateService and user_params merging logic
    let result_status = "success"; // Placeholder
    if result_status == "success" {
        // TODO: Set flash notice: "Notification settings saved"
    } else {
        // TODO: Set flash alert: "Failed to save new settings"
    }
    // TODO: Implement redirect_back_or_default
    HttpResponse::Found()
        .header("Location", "/profile/notifications")
        .finish()
}

// TODO: Implement user_params extraction and validation
// TODO: Implement user_groups, project_notifications_with_preloaded_associations, project_associations helpers
// TODO: Integrate with app state, user session, and flash messages as needed.
