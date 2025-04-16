use crate::config::settings::Settings;
use actix_web::{web, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;
use tokio::fs;

/// Module for managing runner setup scripts
pub trait RunnerSetupScripts {
    /// Get the runner installation path
    fn installation_path(&self) -> PathBuf;

    /// Get the runner configuration path
    fn config_path(&self) -> PathBuf;

    /// Get the runner token
    fn runner_token(&self) -> String;

    /// Get the runner registration token
    fn registration_token(&self) -> String;

    /// Get the runner tags
    fn runner_tags(&self) -> Vec<String>;

    /// Generate the installation script
    fn generate_installation_script(&self) -> Result<String, HttpResponse> {
        let script = format!(
            r#"#!/bin/bash
# Install GitLab Runner
curl -L https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh | sudo bash
sudo apt-get install gitlab-runner

# Register the runner
sudo gitlab-runner register \
  --non-interactive \
  --url "{}" \
  --registration-token "{}" \
  --executor "docker" \
  --docker-image alpine:latest \
  --description "{}" \
  --tag-list "{}" \
  --run-untagged="true" \
  --locked="false" \
  --access-level="not_protected"

# Start the runner
sudo gitlab-runner start"#,
            self.config_path().to_str().unwrap_or(""),
            self.registration_token(),
            "GitLab Runner",
            self.runner_tags().join(",")
        );

        Ok(script)
    }

    /// Generate the configuration script
    fn generate_config_script(&self) -> Result<String, HttpResponse> {
        let config = format!(
            r#"concurrent = 1
check_interval = 0

[session_server]
  session_timeout = 1800

[[runners]]
  name = "{}"
  url = "{}"
  token = "{}"
  executor = "docker"
  [runners.docker]
    tls_verify = false
    image = "alpine:latest"
    privileged = false
    disable_entrypoint_overwrite = false
    oom_kill_disable = false
    disable_cache = false
    volumes = ["/cache"]
    shm_size = 0
"#,
            "GitLab Runner",
            self.config_path().to_str().unwrap_or(""),
            self.runner_token()
        );

        Ok(config)
    }

    /// Save the installation script
    async fn save_installation_script(&self) -> Result<(), HttpResponse> {
        let script = self.generate_installation_script()?;
        let path = self.installation_path();

        tokio::fs::write(&path, script)
            .await
            .map_err(|e| HttpResponse::InternalServerError().body(e.to_string()))?;

        Ok(())
    }

    /// Save the configuration script
    async fn save_config_script(&self) -> Result<(), HttpResponse> {
        let config = self.generate_config_script()?;
        let path = self.config_path();

        tokio::fs::write(&path, config)
            .await
            .map_err(|e| HttpResponse::InternalServerError().body(e.to_string()))?;

        Ok(())
    }

    /// Get script status
    async fn get_script_status(&self) -> Result<HashMap<String, bool>, HttpResponse> {
        let mut status = HashMap::new();

        status.insert(
            "installation_script_exists".to_string(),
            self.installation_path().exists(),
        );

        status.insert(
            "config_script_exists".to_string(),
            self.config_path().exists(),
        );

        Ok(status)
    }
}
