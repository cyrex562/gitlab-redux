// Ported from: orig_app/app/controllers/concerns/internal_redirect.rb
// Provides safe redirect path utilities for actix-web controllers.
use actix_web::HttpRequest;
use url::Url;

/// Verifies that the path starts with a '/' and a known route character.
pub fn safe_redirect_path(path: &str) -> Option<String> {
    if path.is_empty() {
        return None;
    }
    let re = regex::Regex::new(r"^/[-\w].*").unwrap();
    if !re.is_match(path) {
        return None;
    }
    // Try to parse as a relative URL
    let uri = Url::parse(&format!("http://localhost{}", path)).ok()?;
    Some(full_path_for_uri(&uri))
}

/// Checks if the host and port of the uri match the request's host and port.
pub fn host_allowed(uri: &Url, req: &HttpRequest) -> bool {
    let host = req.connection_info().host();
    let port = uri.port_or_known_default().unwrap_or(80);
    let req_url = format!("http://{}", host);
    if let Ok(req_uri) = Url::parse(&req_url) {
        uri.host_str() == req_uri.host_str()
            && port == req_uri.port_or_known_default().unwrap_or(80)
    } else {
        false
    }
}

/// Returns a safe redirect path for a full URL if the host is allowed.
pub fn safe_redirect_path_for_url(url: &str, req: &HttpRequest) -> Option<String> {
    let uri = Url::parse(url).ok()?;
    if host_allowed(&uri, req) {
        safe_redirect_path(&full_path_for_uri(&uri))
    } else {
        None
    }
}

/// Sanitizes a redirect by checking both path and url forms.
pub fn sanitize_redirect(url_or_path: &str, req: &HttpRequest) -> Option<String> {
    safe_redirect_path(url_or_path).or_else(|| safe_redirect_path_for_url(url_or_path, req))
}

/// Returns the path, query, and fragment from a Url as a single string.
pub fn full_path_for_uri(uri: &Url) -> String {
    let mut path_with_query = uri.path().to_string();
    if let Some(q) = uri.query() {
        path_with_query.push('?');
        path_with_query.push_str(q);
    }
    if let Some(frag) = uri.fragment() {
        path_with_query.push('#');
        path_with_query.push_str(frag);
    }
    path_with_query
}

/// Returns the referer path from the request, if present.
pub fn referer_path(req: &HttpRequest) -> Option<String> {
    req.headers()
        .get("referer")
        .and_then(|v| v.to_str().ok())
        .and_then(|referer| Url::parse(referer).ok())
        .map(|uri| uri.path().to_string())
}
