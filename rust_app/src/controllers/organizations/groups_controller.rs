// Ported from: orig_app/app/controllers/organizations/groups_controller.rb
// Ported on: 2025-05-04
// This file implements the Organizations::GroupsController from the Ruby codebase.

use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::controllers::concerns::groups::GroupParams;
use crate::controllers::organizations::application_controller::OrganizationsApplicationController;
use crate::models::group::Group;
use crate::serializers::group_serializer::GroupSerializer;
use crate::services::groups::{CreateService, DestroyService};

#[derive(Debug, Deserialize)]
pub struct PathParams {
    pub id: String,
}

/// Controller for managing groups within organizations
pub struct OrganizationsGroupsController {
    base: OrganizationsApplicationController,
}

impl OrganizationsGroupsController {
    pub fn new(base: OrganizationsApplicationController) -> Self {
        Self { base }
    }

    /// GET /organizations/:organization_id/groups/new
    pub async fn new(&self) -> impl Responder {
        if let Err(e) = self.authorize_create_group() {
            return e;
        }
        HttpResponse::Ok().finish()
    }

    /// GET /organizations/:organization_id/groups/:id/edit
    pub async fn edit(&self) -> impl Responder {
        if let Err(e) = self.base.authorize_read_organization(None) {
            return e;
        }
        if let Err(e) = self.authorize_view_edit_page() {
            return e;
        }
        HttpResponse::Ok().finish()
    }

    /// POST /organizations/:organization_id/groups
    pub async fn create(&self, params: web::Json<GroupParams>) -> impl Responder {
        let response = self.create_group(params.into_inner()).await;

        match response {
            Ok(group) => {
                let serializer = GroupSerializer::new(self.base.current_user_id);
                let json = serializer.represent(&group);
                HttpResponse::Ok().json(json)
            }
            Err(errors) => HttpResponse::UnprocessableEntity().json(json!({
                "message": errors
            })),
        }
    }

    /// DELETE /organizations/:organization_id/groups/:id
    pub async fn destroy(&self, group_id: web::Path<String>) -> impl Responder {
        if let Err(e) = self.authorize_remove_group() {
            return e;
        }

        let group = match self.group(&group_id).await {
            Some(g) => g,
            None => return HttpResponse::NotFound().finish(),
        };

        match DestroyService::new(group, self.base.current_user_id)
            .async_execute()
            .await
        {
            Ok(_) => HttpResponse::Ok().json(json!({
                "message": format!("Group '{}' is being deleted.", group.full_name)
            })),
            Err(e) => HttpResponse::UnprocessableEntity().json(json!({
                "message": e.to_string()
            })),
        }
    }

    // Private methods

    async fn group(&self, id: &str) -> Option<Group> {
        Group::in_organization(self.base.organization(None)?)
            .find_by_full_path(id)
            .await
    }

    async fn create_group(&self, params: GroupParams) -> Result<Group, String> {
        let create_service_params = GroupParams {
            organization_id: Some(self.base.organization(None)?.id),
            ..params
        };

        CreateService::new(self.base.current_user_id, create_service_params)
            .execute()
            .await
    }

    fn authorize_view_edit_page(&self) -> Result<(), HttpResponse> {
        match self.group(params.id).await {
            None => return Err(HttpResponse::NotFound().finish()),
            Some(group) => {
                if !self.base.can("view_edit_page", Some(&group.id.to_string())) {
                    return Err(self.base.access_denied());
                }
            }
        }
        Ok(())
    }

    fn authorize_remove_group(&self) -> Result<(), HttpResponse> {
        match self.group(params.id).await {
            None => return Err(HttpResponse::NotFound().finish()),
            Some(group) => {
                if !self.base.can("remove_group", Some(&group.id.to_string())) {
                    return Err(self.base.access_denied());
                }
            }
        }
        Ok(())
    }

    fn authorize_create_group(&self) -> Result<(), HttpResponse> {
        self.base.authorize_create_group(None)
    }
}
