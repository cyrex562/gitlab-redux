use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::controllers::admin::ApplicationController;

/// Controller for system information
pub struct SystemInfoController {
    /// The admin application controller
    app_controller: ApplicationController,
}

impl SystemInfoController {
    /// Create a new system info controller
    pub fn new() -> Self {
        Self {
            app_controller: ApplicationController::new(),
        }
    }

    /// Handle the show action
    pub async fn show(&self) -> impl Responder {
        // TODO: Implement proper CPU information
        let cpus = json!({
            "user": 10.5,
            "nice": 0.0,
            "system": 5.2,
            "idle": 84.3,
            "iowait": 0.0,
            "irq": 0.0,
            "softirq": 0.0,
            "steal": 0.0,
            "guest": 0.0
        });

        // TODO: Implement proper memory information
        let memory = json!({
            "total": 16777216,
            "used": 8388608,
            "free": 8388608,
            "cached": 4194304,
            "buffers": 1048576
        });

        // TODO: Implement proper disk information
        let disks = vec![json!({
            "bytes_total": 107374182400,
            "bytes_used": 53687091200,
            "disk_name": "/dev/sda1",
            "mount_path": "/"
        })];

        let response = json!({
            "cpus": cpus,
            "memory": memory,
            "disks": disks
        });

        HttpResponse::Ok().json(response)
    }
}
