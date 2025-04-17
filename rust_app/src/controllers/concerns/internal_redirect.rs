use actix_web::{HttpRequest, HttpResponse};
use regex::Regex;
use std::sync::OnceLock;
use url::Url;

pub trait InternalRedirect {
    fn safe_redirect_path(&self, path: Option<&str>) -> Option<String> {
        let path = path?;

        // Verify that the string starts with a '/' and a known route character
        static PATH_REGEX: OnceLock<Regex> = OnceLock::new();
        let path_regex = PATH_REGEX.get_or_init(|| Regex::new(r"^/[-\w].*$").unwrap());

        if !path_regex.is_match(path) {
            return None;
        }

        // Parse the URI and extract the path components
        let uri = match Url::parse(&format!("http://dummy{}", path)) {
            Ok(uri) => uri,
            Err(_) => return None,
        };

        // Build the full path with query and fragment
        let mut full_path = uri.path().to_string();
        if let Some(query) = uri.query() {
            full_path.push('?');
            full_path.push_str(query);
        }
        if let Some(fragment) = uri.fragment() {
            full_path.push('#');
            full_path.push_str(fragment);
        }

        Some(full_path)
    }

    fn safe_redirect_path_for_url(&self, url: Option<&str>, req: &HttpRequest) -> Option<String> {
        let url = url?;

        let uri = match Url::parse(url) {
            Ok(uri) => uri,
            Err(_) => return None,
        };

        if self.host_allowed(&uri, req) {
            self.safe_redirect_path(Some(&format!("http://dummy{}", uri.path())))
        } else {
            None
        }
    }

    fn sanitize_redirect(&self, url_or_path: Option<&str>, req: &HttpRequest) -> Option<String> {
        self.safe_redirect_path(url_or_path)
            .or_else(|| self.safe_redirect_path_for_url(url_or_path, req))
    }

    fn host_allowed(&self, uri: &Url, req: &HttpRequest) -> bool {
        let host = req.connection_info().host();
        let port = req.connection_info().port();

        uri.host_str() == Some(host) && uri.port() == Some(port.parse::<u16>().unwrap_or(80))
    }

    fn referer_path(&self, req: &HttpRequest) -> Option<String> {
        req.headers()
            .get("referer")
            .and_then(|v| v.to_str().ok())
            .and_then(|referer| Url::parse(referer).ok().map(|uri| uri.path().to_string()))
    }
}
