pub mod crm {
    pub mod contacts;
    pub mod organizations_controller;
}

pub mod dependency_proxy {
    pub mod application_controller;
    pub mod auth_controller; // Ported from Ruby: groups/dependency_proxy_auth_controller.rb
}

pub mod achievements_controller;

pub mod application_controller;

pub mod autocomplete_sources_controller; // Register the new controller

pub mod avatars_controller;

pub mod boards_controller;

pub mod children_controller; // Register the new controller

pub mod clusters_controller; // Ported from Ruby: groups/clusters_controller.rb

pub mod custom_emoji_controller; // Ported from Ruby: groups/custom_emoji_controller.rb

pub mod dependency_proxies_controller; // Ported from Ruby: groups/dependency_proxies_controller.rb

pub mod deploy_tokens_controller; // Ported from Ruby: groups/deploy_tokens_controller.rb

pub mod group_links_controller; // Ported from Ruby: groups/group_links_controller.rb

pub mod group_members_controller; // Ported from Ruby: groups/group_members_controller.rb

pub mod imports_controller; // Ported from Ruby: groups/imports_controller.rb

pub mod infrastructure_registry_controller; // Ported from Ruby: groups/infrastructure_registry_controller.rb

pub mod milestones_controller; // Ported from Ruby: groups/milestones_controller.rb

pub mod packages_controller; // Ported from Ruby: groups/packages_controller.rb

pub mod redirect_controller; // Ported from Ruby: groups/redirect_controller.rb

pub mod releases_controller; // Ported from Ruby: groups/releases_controller.rb

pub mod runners_controller; // Ported from Ruby: groups/runners_controller.rb

pub mod uploads_controller; // Ported from Ruby: groups/uploads_controller.rb

pub mod usage_quotas_controller; // Ported from Ruby: groups/usage_quotas_controller.rb

pub mod variables_controller; // Ported from Ruby: groups/variables_controller.rb

pub mod work_items_controller;
