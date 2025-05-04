use std::sync::{Arc, Mutex};
use std::collections::HashMap;
use std::hash::Hash;

pub trait StrongMemoize {
    fn memoized<K, V, F>(&self, key: K, f: F) -> V 
    where
        K: Hash + Eq + Clone,
        V: Clone,
        F: FnOnce() -> V;
}

#[derive(Default)]
pub struct StrongMemoizeImpl {
    cache: Arc<Mutex<HashMap<String, Box<dyn std::any::Any + Send + Sync>>>>,
}

impl StrongMemoizeImpl {
    pub fn new() -> Self {
        Self {
            cache: Arc::new(Mutex::new(HashMap::new())),
        }
    }
}

impl StrongMemoize for StrongMemoizeImpl {
    fn memoized<K, V, F>(&self, key: K, f: F) -> V 
    where
        K: Hash + Eq + Clone,
        V: Clone + Send + Sync + 'static,
        F: FnOnce() -> V,
    {
        let key_str = format!("{:?}", key);
        let mut cache = self.cache.lock().unwrap();
        
        if let Some(value) = cache.get(&key_str) {
            if let Some(value) = value.downcast_ref::<V>() {
                return value.clone();
            }
        }
        
        let value = f();
        cache.insert(key_str, Box::new(value.clone()));
        value
    }
}