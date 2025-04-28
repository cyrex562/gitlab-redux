// Ported from: orig_app/app/controllers/concerns/redirects_for_missing_path_on_tree.rb
// Provides redirect and notice for missing paths in a project tree.
use actix_web::{HttpRequest, HttpResponse, http::header};
use chrono::Utc;
use std::fs::OpenOptions;
use std::io::Write;

/// Truncates a path to a maximum length, preserving the last segments separated by '/'.
pub fn truncate_path(path: &str, max_len: usize) -> String {
    if path.len() <= max_len {
        return path.to_string();
    }
    let rev: String = path.chars().rev().collect();
    let mut parts: Vec<&str> = rev.splitn(2, '/').collect();
    let mut truncated = String::new();
    let mut count = 0;
    for c in rev.chars() {
        if count >= max_len {
            break;
        }
        truncated.push(c);
        count += 1;
    }
    let mut truncated: String = truncated.chars().collect();
    truncated = truncated.chars().rev().collect();
    if truncated.len() < path.len() {
        format!("...{}", &truncated[truncated.len().saturating_sub(max_len)..])
    } else {
        truncated
    }
}

/// Formats the notice message for a missing path on a ref.
pub fn missing_path_on_ref(path: &str, r#ref: &str) -> String {
    let truncated = truncate_path(path, 60);
    format!(\"\"{path}\" did not exist on \"{ref}\"\", path = truncated, ref = r#ref)
}

/// Redirects to the tree root for a missing path, with a notice message.
pub fn redirect_to_tree_root_for_missing_path(
    project: &str,
    r#ref: &str,
    path: &str,
    ref_type: Option<&str>,
    req: &HttpRequest,
) -> HttpResponse {
    // Construct the tree path (stub: replace with actual route logic)
    let tree_path = if let Some(rt) = ref_type {
        format!("/projects/{}/-/tree/{}?ref_type={}", project, r#ref, rt)
    } else {
        format!("/projects/{}/-/tree/{}", project, r#ref)
    };
    let notice = missing_path_on_ref(path, r#ref);
    // Log the porting action
    let log_entry = format!(
        "{} | Ported: redirects_for_missing_path_on_tree.rb -> redirects_for_missing_path_on_tree.rs | Ruby: {} | Rust: {}\n",
        Utc::now().to_rfc3339(),
        "orig_app/app/controllers/concerns/redirects_for_missing_path_on_tree.rb",
        "rust_app/src/controllers/concerns/redirects_for_missing_path_on_tree.rs"
    );
    if let Ok(mut file) = OpenOptions::new().append(true).create(true).open("../../../../porting_log.txt") {
        let _ = file.write_all(log_entry.as_bytes());
    }
    HttpResponse::Found()
        .append_header((header::LOCATION, tree_path))
        .append_header(("X-Notice", notice))
        .finish()
}