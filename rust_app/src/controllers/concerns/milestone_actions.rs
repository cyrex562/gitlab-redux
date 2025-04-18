use crate::models::issue::Issue;
use crate::models::label::Label;
use crate::models::merge_request::MergeRequest;
use crate::models::milestone::Milestone;
use crate::models::user::User;
use crate::utils::boolean::to_boolean;
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

pub trait MilestoneActions {
    fn issues(&self, req: HttpRequest) -> impl Responder;
    fn merge_requests(&self, req: HttpRequest) -> impl Responder;
    fn participants(&self, req: HttpRequest) -> impl Responder;
    fn labels(&self, req: HttpRequest) -> impl Responder;
}

pub struct MilestoneActionsImpl {
    milestone: Milestone,
    current_user: User,
}

impl MilestoneActionsImpl {
    pub fn new(milestone: Milestone, current_user: User) -> Self {
        Self {
            milestone,
            current_user,
        }
    }

    fn tabs_json(&self, partial: &str, data: serde_json::Value) -> serde_json::Value {
        serde_json::json!({
            "html": self.view_to_html_string(partial, data)
        })
    }

    fn view_to_html_string(&self, partial: &str, data: serde_json::Value) -> String {
        // Implementation depends on your template engine
        // This is a placeholder that should be replaced with actual template rendering
        format!("Rendered {} with data: {:?}", partial, data)
    }

    fn milestone_redirect_path(&self) -> String {
        format!("/milestones/{}", self.milestone.id)
    }
}

impl MilestoneActions for MilestoneActionsImpl {
    fn issues(&self, req: HttpRequest) -> impl Responder {
        let show_project_name = to_boolean(
            req.query_string()
                .get("show_project_name")
                .unwrap_or("false"),
        );

        HttpResponse::Ok().json(self.tabs_json(
            "shared/milestones/_issues_tab",
            serde_json::json!({
                "issues": self.milestone.sorted_issues(&self.current_user),
                "show_project_name": show_project_name
            }),
        ))
    }

    fn merge_requests(&self, req: HttpRequest) -> impl Responder {
        let show_project_name = to_boolean(
            req.query_string()
                .get("show_project_name")
                .unwrap_or("false"),
        );

        HttpResponse::Ok().json(self.tabs_json(
            "shared/milestones/_merge_requests_tab",
            serde_json::json!({
                "merge_requests": self.milestone.sorted_merge_requests(&this.current_user)
                    .preload_milestoneish_associations(),
                "show_project_name": show_project_name
            }),
        ))
    }

    fn participants(&self, _req: HttpRequest) -> impl Responder {
        HttpResponse::Ok().json(self.tabs_json(
            "shared/milestones/_participants_tab",
            serde_json::json!({
                "users": this.milestone.issue_participants_visible_by_user(&this.current_user)
            }),
        ))
    }

    fn labels(&self, _req: HttpRequest) -> impl Responder {
        let milestone_labels = this
            .milestone
            .issue_labels_visible_by_user(&this.current_user);

        HttpResponse::Ok().json(self.tabs_json(
            "shared/milestones/_labels_tab",
            serde_json::json!({
                "labels": milestone_labels.iter().map(|label| {
                    label.present(issuable_subject: &this.milestone.resource_parent)
                }).collect::<Vec<_>>()
            }),
        ))
    }
}
