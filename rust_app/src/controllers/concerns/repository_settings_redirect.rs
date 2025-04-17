use actix_web::HttpResponse;

/// Module for handling repository settings redirects
pub trait RepositorySettingsRedirect {
    /// Redirect to repository settings path
    fn redirect_to_repository_settings(&self) -> HttpResponse {
        HttpResponse::Found()
            .header("Location", self.repository_settings_path())
            .finish()
    }

    /// Get the repository settings path
    fn repository_settings_path(&self) -> String {
        format!("/projects/{}/settings/repository", self.project().id())
    }

    // Required trait method that needs to be implemented by the controller
    fn project(&self) -> &crate::models::project::Project;
}
