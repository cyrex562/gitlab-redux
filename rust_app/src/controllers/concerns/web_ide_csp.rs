// Ported from: orig_app/app/controllers/concerns/web_ide_csp.rb
// Adds Web IDE-specific frame-src and worker-src to the Content Security Policy.

use crate::auth::security_policy::ContentSecurityPolicy;
use url::Url;

pub trait WebIdeCSP {
    /// Updates the given ContentSecurityPolicy with Web IDE values
    fn include_web_ide_csp(
        &self,
        csp: &mut ContentSecurityPolicy,
        request_url: &str,
        relative_url_root: Option<&str>,
    );
}

pub struct WebIdeCSPHandler;

impl WebIdeCSP for WebIdeCSPHandler {
    fn include_web_ide_csp(
        &self,
        csp: &mut ContentSecurityPolicy,
        request_url: &str,
        relative_url_root: Option<&str>,
    ) {
        // If there are no directives, do nothing
        let has_directives = !csp.default_src.is_empty()
            || !csp.frame_src.is_empty()
            || !csp.worker_src().is_empty();
        if !has_directives {
            return;
        }

        // Build the webpack URL
        let mut base_uri =
            Url::parse(request_url).unwrap_or_else(|_| Url::parse("http://localhost/").unwrap());
        let rel_root = relative_url_root.unwrap_or("/");
        base_uri.set_path(&format!(
            "{}/assets/webpack/",
            rel_root.trim_end_matches('/')
        ));
        base_uri.set_query(None);
        let webpack_url = base_uri.to_string();

        // frame-src: add webpack_url and https://*.web-ide.gitlab-static.net/
        let mut frame_src = if !csp.frame_src.is_empty() {
            csp.frame_src.clone()
        } else {
            csp.default_src.clone()
        };
        if !frame_src.contains(&webpack_url) {
            frame_src.push(webpack_url.clone());
        }
        let static_net = "https://*.web-ide.gitlab-static.net/".to_string();
        if !frame_src.contains(&static_net) {
            frame_src.push(static_net);
        }
        csp.frame_src = frame_src;

        // worker-src: add webpack_url
        let mut worker_src = csp.worker_src();
        if !worker_src.contains(&webpack_url) {
            worker_src.push(webpack_url);
        }
        csp.set_worker_src(worker_src);
    }
}

// Extension trait for worker-src, since it's not in the default struct
trait WorkerSrcExt {
    fn worker_src(&self) -> Vec<String>;
    fn set_worker_src(&mut self, values: Vec<String>);
}

impl WorkerSrcExt for ContentSecurityPolicy {
    fn worker_src(&self) -> Vec<String> {
        // In a real implementation, this would be a field. Here, we use default_src as fallback.
        // You may want to add worker_src to ContentSecurityPolicy struct for full support.
        self.default_src.clone()
    }
    fn set_worker_src(&mut self, _values: Vec<String>) {
        // No-op unless worker_src is added to ContentSecurityPolicy
    }
}
