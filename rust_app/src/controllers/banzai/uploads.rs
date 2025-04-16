use actix_web::{web, HttpResponse, Responder};
use serde::Deserialize;
use crate::controllers::base::BaseController;
use crate::models::{Project, Group};

#[derive(Debug, Deserialize)]
pub struct UploadParams {
    model: String,
    model_id: i32,
}

pub struct UploadsController;

impl UploadsController {
    const MODEL_CLASSES: &'static [(&'static str, &'static str)] = &[
        ("project", "Project"),
        ("group", "Group"),
    ];

    pub async fn verify_upload_model_class(params: &UploadParams) -> bool {
        Self::MODEL_CLASSES.iter().any(|(key, _)| key == &params.model)
    }

    pub async fn authorize_access(
        current_user: Option<&crate::models::User>,
        model: &dyn crate::models::Model,
    ) -> bool {
        // TODO: Implement proper authorization checks
        // This should check if the user has read access to the model
        true
    }

    pub fn get_uploader_class(model_type: &str) -> &'static str {
        match model_type {
            "project" => "FileUploader",
            "group" => "NamespaceFileUploader",
            _ => "FileUploader",
        }
    }

    pub async fn find_model(params: &UploadParams) -> Result<Box<dyn crate::models::Model>, actix_web::Error> {
        match params.model.as_str() {
            "project" => {
                let project = Project::find(params.model_id)
                    .await
                    .map_err(|e| actix_web::error::ErrorNotFound(e))?;
                Ok(Box::new(project))
            }
            "group" => {
                let group = Group::find(params.model_id)
                    .await
                    .map_err(|e| actix_web::error::ErrorNotFound(e))?;
                Ok(Box::new(group))
            }
            _ => Err(actix_web::error::ErrorNotFound("Model not found")),
        }
    }

    #[post("/banzai/uploads")]
    pub async fn create(
        params: web::Json<UploadParams>,
        current_user: Option<web::ReqData<crate::models::User>>,
    ) -> impl Responder {
        // Verify model class
        if !Self::verify_upload_model_class(&params).await {
            return HttpResponse::NotFound().finish();
        }

        // Find the model
        let model = match Self::find_model(&params).await {
            Ok(model) => model,
            Err(_) => return HttpResponse::NotFound().finish(),
        };

        // Authorize access
        if !Self::authorize_access(current_user.as_ref().map(|u| u.as_ref()), &*model).await {
            return HttpResponse::Forbidden().finish();
        }

        // TODO: Implement file upload handling
        // This should use the appropriate uploader class based on the model type

        HttpResponse::Ok().json(serde_json::json!({
            "status": "ok",
            "message": "File upload endpoint"
        }))
    }
} 