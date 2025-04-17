use std::collections::HashMap;

pub struct Group {
    pub id: i32,
    // Add other group fields as needed
}

pub struct User {
    // Add user fields as needed
}

pub trait RendersMemberAccess {
    fn prepare_groups_for_rendering(&self, groups: &[Group]) -> Vec<Group> {
        self.preload_max_member_access_for_collection::<Group>(groups);
        groups.to_vec()
    }

    fn preload_max_member_access_for_collection<T>(&self, collection: &[T])
    where
        T: HasId,
    {
        if let Some(user) = self.get_current_user() {
            if !collection.is_empty() {
                let collection_ids: Vec<i32> =
                    collection.iter().map(|item| item.get_id()).collect();
                let method_name =
                    format!("max_member_access_for_{}_ids", std::any::type_name::<T>());
                self.call_access_method(user, &method_name, &collection_ids);
            }
        }
    }

    // Required methods to be implemented by concrete types
    fn get_current_user(&self) -> Option<&User>;
    fn call_access_method(&self, user: &User, method_name: &str, ids: &[i32]) -> HashMap<i32, i32>;
}

pub trait HasId {
    fn get_id(&self) -> i32;
}

impl HasId for Group {
    fn get_id(&self) -> i32 {
        self.id
    }
}
