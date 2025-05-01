pub mod google_api;
pub mod jira_connect_app_descriptor;

pub mod jira_connect {
    pub mod app_descriptor_controller;
    pub mod branches_controller;
    pub mod events_controller;
    pub mod installations_controller;
    pub mod oauth_application_ids_controller;
    pub mod oauth_callbacks_controller;
    pub mod public_keys_controller; // Added for public keys endpoint
}

pub mod groups {
    pub mod settings {
        pub mod applications;
        pub mod ci_cd;
        pub mod repository_controller;
    }
}

pub mod groups_labels;
pub mod slacks_controller;

pub mod ldap {
    pub mod omniauth_callbacks_controller;
}

pub mod mailgun {
    pub mod webhooks_controller;
}

pub mod oauth {
    pub mod applications_controller;
    pub mod authorizations_controller;
    pub mod token_info_controller;
    pub mod tokens_controller; // Added for OAuth token endpoint
}

pub mod organizations {
    pub mod application_controller;
}
