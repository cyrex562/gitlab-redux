use actix_web::web;
use std::collections::HashMap;
use std::sync::OnceLock;

pub trait Snippet {
    fn id(&self) -> &str;
}

pub trait SnippetRepository {
    fn find_by_id(&self, id: &str) -> Option<Box<dyn Snippet>>;
    fn inc_relations_for_view(&self) -> &Self;
}

pub trait FindSnippet {
    fn snippet(&self) -> Option<Box<dyn Snippet>> {
        static SNIPPET: OnceLock<Option<Box<dyn Snippet>>> = OnceLock::new();
        SNIPPET
            .get_or_init(|| {
                self.snippet_repository()
                    .inc_relations_for_view()
                    .find_by_id(&self.snippet_id())
            })
            .clone()
    }

    fn snippet_klass(&self) -> &'static str {
        unimplemented!("snippet_klass must be implemented by concrete types")
    }

    fn snippet_id(&self) -> String {
        // In a real implementation, this would get the ID from the request parameters
        "default_id".to_string()
    }

    fn snippet_find_params(&self) -> HashMap<String, String> {
        let mut params = HashMap::new();
        params.insert("id".to_string(), self.snippet_id());
        params
    }

    // Required method to be implemented by concrete types
    fn snippet_repository(&self) -> &dyn SnippetRepository;
}
