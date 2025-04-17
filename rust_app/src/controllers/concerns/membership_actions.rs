use actix_web::{web, HttpResponse, Responder};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct Member {
    pub id: i32,
    pub user_id: i32,
    pub access_level: i32,
    pub expires_at: Option<DateTime<Utc>>,
    pub source_type: String,
    pub source_id: i32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateParams {
    pub access_level: i32,
    pub expires_at: Option<DateTime<Utc>>,
    pub source: String,
}

#[derive(Debug, Serialize)]
pub struct UpdateResponse {
    pub expires_soon: bool,
    pub expires_at_formatted: Option<String>,
}

pub trait MembershipActions {
    fn update(&self, member_id: i32, params: web::Json<UpdateParams>) -> impl Responder {
        let member = self.find_member(member_id);
        match self.update_member(member, params.into_inner()) {
            Ok(result) => {
                let response = self.create_update_response(&result);
                HttpResponse::Ok().json(response)
            }
            Err(e) => HttpResponse::UnprocessableEntity()
                .json(HashMap::from([("message", e.to_string())])),
        }
    }

    fn destroy(
        &self,
        member_id: i32,
        remove_sub_memberships: bool,
        unassign_issuables: bool,
    ) -> impl Responder {
        let member = self.find_member(member_id);
        match self.destroy_member(member, remove_sub_memberships, unassign_issuables) {
            Ok(_) => {
                let message = if self.is_namespace() {
                    if remove_sub_memberships {
                        "User was successfully removed from group."
                    } else {
                        "User was successfully removed from group and any subgroups and projects."
                    }
                } else {
                    "User was successfully removed from project."
                };
                HttpResponse::SeeOther()
                    .append_header(("Location", self.members_page_url()))
                    .append_header(("X-Flash-Message", message))
                    .finish()
            }
            Err(e) => HttpResponse::InternalServerError().body(e.to_string()),
        }
    }

    fn request_access(&self, user_id: i32) -> impl Responder {
        match self.create_access_request(user_id) {
            Ok(_) => HttpResponse::SeeOther()
                .append_header(("Location", self.polymorphic_path()))
                .append_header((
                    "X-Flash-Message",
                    "Your request for access has been queued for review.",
                ))
                .finish(),
            Err(e) => HttpResponse::SeeOther()
                .append_header(("Location", self.polymorphic_path()))
                .append_header((
                    "X-Flash-Message",
                    format!("Your request for access could not be processed: {}", e),
                ))
                .finish(),
        }
    }

    fn approve_access_request(&self, requester_id: i32) -> impl Responder {
        let requester = self.find_requester(requester_id);
        match self.approve_access_request(requester) {
            Ok(_) => HttpResponse::SeeOther()
                .append_header(("Location", self.members_page_url()))
                .finish(),
            Err(e) => HttpResponse::InternalServerError().body(e.to_string()),
        }
    }

    fn leave(&self, user_id: i32) -> impl Responder {
        let member = self.find_member_by_user(user_id);
        match self.destroy_member(member, true, true) {
            Ok(_) => {
                let notice = if member.is_request() {
                    format!(
                        "Your access request to the {} has been withdrawn.",
                        self.source_type()
                    )
                } else {
                    format!(
                        "You left the \"{}\" {}.",
                        self.membershipable_human_name(),
                        self.source_type()
                    )
                };
                HttpResponse::SeeOther()
                    .append_header(("Location", self.leave_redirect_path(&member)))
                    .append_header(("X-Flash-Message", notice))
                    .finish()
            }
            Err(e) => HttpResponse::InternalServerError().body(e.to_string()),
        }
    }

    fn resend_invite(&self, member_id: i32) -> impl Responder {
        let member = self.find_member(member_id);
        if member.is_invite() {
            match self.resend_invite(member) {
                Ok(_) => HttpResponse::SeeOther()
                    .append_header(("Location", self.members_page_url()))
                    .append_header(("X-Flash-Message", "The invitation was successfully resent."))
                    .finish(),
                Err(e) => HttpResponse::InternalServerError().body(e.to_string()),
            }
        } else {
            HttpResponse::SeeOther()
                .append_header(("Location", self.members_page_url()))
                .append_header((
                    "X-Flash-Message",
                    "The invitation has already been accepted.",
                ))
                .finish()
        }
    }

    // Required methods to be implemented by concrete types
    fn find_member(&self, id: i32) -> Member;
    fn update_member(
        &self,
        member: Member,
        params: UpdateParams,
    ) -> Result<Member, Box<dyn std::error::Error>>;
    fn destroy_member(
        &self,
        member: Member,
        skip_subresources: bool,
        unassign_issuables: bool,
    ) -> Result<(), Box<dyn std::error::Error>>;
    fn create_access_request(&self, user_id: i32) -> Result<(), Box<dyn std::error::Error>>;
    fn find_requester(&self, id: i32) -> Member;
    fn approve_access_request(&self, requester: Member) -> Result<(), Box<dyn std::error::Error>>;
    fn find_member_by_user(&self, user_id: i32) -> Member;
    fn resend_invite(&self, member: Member) -> Result<(), Box<dyn std::error::Error>>;
    fn is_namespace(&self) -> bool;
    fn members_page_url(&self) -> String;
    fn polymorphic_path(&self) -> String;
    fn source_type(&self) -> String;
    fn membershipable_human_name(&self) -> String;
    fn leave_redirect_path(&self, member: &Member) -> String;

    // Helper methods
    fn create_update_response(&self, member: &Member) -> UpdateResponse {
        UpdateResponse {
            expires_soon: member
                .expires_at
                .map_or(false, |date| (date - Utc::now()).num_days() <= 7),
            expires_at_formatted: member
                .expires_at
                .map(|date| date.format("%B %d, %Y").to_string()),
        }
    }
}
