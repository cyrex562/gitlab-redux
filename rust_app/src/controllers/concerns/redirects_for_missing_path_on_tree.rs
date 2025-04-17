use actix_web::{HttpRequest, HttpResponse};
use std::path::Path;

pub trait RedirectsForMissingPathOnTree {
    fn redirect_to_tree_root_for_missing_path(
        &self,
        req: &HttpRequest,
        project_path: &str,
        ref_name: &str,
        path: &str,
    ) -> HttpResponse {
        let truncated_path = self.truncate_path(path);
        let message = self.missing_path_on_ref(project_path, ref_name, &truncated_path);

        HttpResponse::Found()
            .header("Location", format!("/{}/tree/{}", project_path, ref_name))
            .header("X-Flash-Message", message)
            .finish()
    }

    fn missing_path_on_ref(&self, project_path: &str, ref_name: &str, path: &str) -> String {
        format!("The path '{}' does not exist on '{}'", path, ref_name)
    }

    fn truncate_path(&self, path: &str) -> String {
        const MAX_LENGTH: usize = 60;

        if path.len() <= MAX_LENGTH {
            return path.to_string();
        }

        let path_obj = Path::new(path);
        if let Some(file_name) = path_obj.file_name() {
            if let Some(file_name_str) = file_name.to_str() {
                let truncated = format!(".../{}", file_name_str);
                if truncated.len() <= MAX_LENGTH {
                    return truncated;
                }
            }
        }

        // If we can't get a filename or the truncated version is too long,
        // just return the first MAX_LENGTH characters
        path.chars().take(MAX_LENGTH).collect()
    }
}
