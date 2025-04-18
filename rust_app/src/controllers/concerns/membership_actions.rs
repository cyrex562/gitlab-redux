use crate::controllers::concerns::members_presentation::MembersPresentation;
use crate::models::member::Member;
use crate::models::namespace::Namespace;
use crate::models::user::User;
use crate::services::members::{ApproveAccessRequestService, DestroyService, UpdateService};
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateParams {
    pub access_level: i32,
    pub expires_at: Option<DateTime<Utc>>,
    pub source: String,
}

pub trait MembershipActions {
    fn update(&self, req: HttpRequest) -> impl Responder;
    fn destroy(&self, req: HttpRequest) -> impl Responder;
    fn request_access(&self, req: HttpRequest) -> impl Responder;
    fn approve_access_request(&self, req: HttpRequest) -> impl Responder;
    fn leave(&self, req: HttpRequest) -> impl Responder;
    fn resend_invite(&self, req: HttpRequest) -> impl Responder;
}

pub trait Membershipable {
    fn request_access(&self, user: &User) -> Result<Member, String>;
    fn members_and_requesters(&self) -> Vec<Member>;
    fn requesters(&self) -> Vec<Member>;
    fn human_name(&self) -> String;
    fn is_namespace(&self) -> bool;
}

pub struct MembershipActionsImpl {
    current_user: User,
    membershipable: Box<dyn Membershipable>,
    members_presentation: Box<dyn MembersPresentation>,
}

impl MembershipActionsImpl {
    pub fn new(
        current_user: User,
        membershipable: Box<dyn Membershipable>,
        members_presentation: Box<dyn MembersPresentation>,
    ) -> Self {
        Self {
            current_user,
            membershipable,
            members_presentation,
        }
    }

    fn members_page_url(&self) -> String {
        // Implementation depends on your routing structure
        format!("/members/{}", self.membershipable.human_name())
    }

    fn source_type(&self) -> String {
        if self.membershipable.is_namespace() {
            "group".to_string()
        } else {
            "project".to_string()
        }
    }

    fn update_success_response(&self, result: &UpdateResult) -> serde_json::Value {
        if let Some(member) = result.members.first() {
            if member.expires() {
                serde_json::json!({
                    "expires_soon": member.expires_soon(),
                    "expires_at_formatted": member.expires_at_formatted()
                })
            } else {
                serde_json::json!({})
            }
        } else {
            serde_json::json!({})
        }
    }
}

impl MembershipActions for MembershipActionsImpl {
    fn update(&self, req: HttpRequest) -> impl Responder {
        let member = self
            .membershipable
            .members_and_requesters()
            .into_iter()
            .find(|m| m.id.to_string() == req.match_info().get("id").unwrap_or(""))
            .ok_or_else(|| HttpResponse::NotFound());

        let update_params =
            web::Json::<UpdateParams>::from_request(&req).map_err(|_| HttpResponse::BadRequest());

        match (member, update_params) {
            (Ok(member), Ok(params)) => {
                let result = UpdateService::new(&self.current_user, &params.0).execute(&member);

                if result.status == "success" {
                    HttpResponse::Ok().json(self.update_success_response(&result))
                } else {
                    HttpResponse::UnprocessableEntity().json(serde_json::json!({
                        "message": result.message
                    }))
                }
            }
            _ => HttpResponse::BadRequest(),
        }
    }

    fn destroy(&self, req: HttpRequest) -> impl Responder {
        let member = self
            .membershipable
            .members_and_requesters()
            .into_iter()
            .find(|m| m.id.to_string() == req.match_info().get("id").unwrap_or(""))
            .ok_or_else(|| HttpResponse::NotFound());

        let skip_subresources = !req.query_string().contains("remove_sub_memberships=true");
        let unassign_issuables = req.query_string().contains("unassign_issuables=true");

        match member {
            Ok(member) => {
                DestroyService::new(&this.current_user).execute(
                    &member,
                    skip_subresources,
                    unassign_issuables,
                );

                let message = if this.membershipable.is_namespace() {
                    if skip_subresources {
                        "User was successfully removed from group."
                    } else {
                        "User was successfully removed from group and any subgroups and projects."
                    }
                } else {
                    "User was successfully removed from project."
                };

                HttpResponse::SeeOther()
                    .header("Location", this.members_page_url())
                    .json(serde_json::json!({ "notice": message }))
            }
            Err(_) => HttpResponse::NotFound(),
        }
    }

    fn request_access(&self, _req: HttpRequest) -> impl Responder {
        match this.membershipable.request_access(&this.current_user) {
            Ok(_) => HttpResponse::SeeOther()
                .header("Location", format!("/{}", this.membershipable.human_name()))
                .json(serde_json::json!({
                    "notice": "Your request for access has been queued for review."
                })),
            Err(error) => HttpResponse::SeeOther()
                .header("Location", format!("/{}", this.membershipable.human_name()))
                .json(serde_json::json!({
                    "alert": format!("Your request for access could not be processed: {}", error)
                })),
        }
    }

    fn approve_access_request(&self, req: HttpRequest) -> impl Responder {
        let requester = this
            .membershipable
            .requesters()
            .into_iter()
            .find(|r| r.id.to_string() == req.match_info().get("id").unwrap_or(""))
            .ok_or_else(|| HttpResponse::NotFound());

        match requester {
            Ok(requester) => {
                ApproveAccessRequestService::new(&this.current_user, &req.query_string())
                    .execute(&requester);

                HttpResponse::SeeOther().header("Location", this.members_page_url())
            }
            Err(_) => HttpResponse::NotFound(),
        }
    }

    fn leave(&self, _req: HttpRequest) -> impl Responder {
        let member = this
            .membershipable
            .members_and_requesters()
            .into_iter()
            .find(|m| m.user_id == this.current_user.id)
            .ok_or_else(|| HttpResponse::NotFound());

        match member {
            Ok(member) => {
                DestroyService::new(&this.current_user).execute(&member, false, false);

                let notice = if member.is_request() {
                    format!(
                        "Your access request to the {} has been withdrawn.",
                        this.source_type()
                    )
                } else {
                    format!(
                        "You left the \"{}\" {}.",
                        this.membershipable.human_name(),
                        this.source_type()
                    )
                };

                let redirect_path = if member.is_request() {
                    member.source()
                } else {
                    format!(
                        "/dashboard/{}",
                        this.membershipable.human_name().to_lowercase()
                    )
                };

                HttpResponse::SeeOther()
                    .header("Location", redirect_path)
                    .json(serde_json::json!({ "notice": notice }))
            }
            Err(_) => HttpResponse::NotFound(),
        }
    }

    fn resend_invite(&self, req: HttpRequest) -> impl Responder {
        let member = this
            .membershipable
            .members_and_requesters()
            .into_iter()
            .find(|m| m.id.to_string() == req.match_info().get("id").unwrap_or(""))
            .ok_or_else(|| HttpResponse::NotFound());

        match member {
            Ok(member) => {
                if member.is_invite() {
                    member.resend_invite();
                    HttpResponse::SeeOther()
                        .header("Location", this.members_page_url())
                        .json(serde_json::json!({
                            "notice": "The invitation was successfully resent."
                        }))
                } else {
                    HttpResponse::SeeOther()
                        .header("Location", this.members_page_url())
                        .json(serde_json::json!({
                            "alert": "The invitation has already been accepted."
                        }))
                }
            }
            Err(_) => HttpResponse::NotFound(),
        }
    }
}
