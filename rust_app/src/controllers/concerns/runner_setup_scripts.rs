use actix_web::{web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use crate::ci::runner_instructions::RunnerInstructions;

#[derive(Debug, Deserialize)]
pub struct ScriptParams {
    pub os: Option<String>,
    pub arch: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct RunnerSetupOutput {
    pub install: String,
    pub register: String,
}

#[derive(Debug, Serialize)]
pub struct RunnerSetupError {
    pub errors: Vec<String>,
}

pub trait RunnerSetupScripts {
    fn private_runner_setup_scripts(&self, params: web::Query<ScriptParams>) -> impl Responder;
}

pub struct RunnerSetupScriptsHandler;

impl RunnerSetupScriptsHandler {
    pub fn new() -> Self {
        RunnerSetupScriptsHandler
    }
}

impl RunnerSetupScripts for RunnerSetupScriptsHandler {
    fn private_runner_setup_scripts(&self, params: web::Query<ScriptParams>) -> impl Responder {
        let os = params.os.clone().unwrap_or_default();
        let arch = params.arch.clone().unwrap_or_default();
        
        let instructions = RunnerInstructions::new(os, arch);
        
        if !instructions.errors.is_empty() {
            let error_response = RunnerSetupError {
                errors: instructions.errors,
            };
            return HttpResponse::BadRequest().json(error_response);
        }
        
        let output = RunnerSetupOutput {
            install: instructions.install_script,
            register: instructions.register_command,
        };
        
        HttpResponse::Ok().json(output)
    }
}

// This would be implemented in a separate module
pub mod ci {
    pub mod runner_instructions {
        use serde::{Deserialize, Serialize};
        
        #[derive(Debug, Serialize, Deserialize)]
        pub struct RunnerInstructions {
            pub os: String,
            pub arch: String,
            pub install_script: String,
            pub register_command: String,
            pub errors: Vec<String>,
        }
        
        impl RunnerInstructions {
            pub fn new(os: String, arch: String) -> Self {
                // In a real implementation, this would generate the appropriate scripts
                // based on the OS and architecture
                let mut errors = Vec::new();
                
                // Validate OS and arch
                if os.is_empty() {
                    errors.push("OS parameter is required".to_string());
                }
                
                if arch.is_empty() {
                    errors.push("Architecture parameter is required".to_string());
                }
                
                // Generate scripts based on OS and arch
                let (install_script, register_command) = if errors.is_empty() {
                    match (os.as_str(), arch.as_str()) {
                        ("linux", "amd64") => (
                            "#!/bin/bash\n# Linux AMD64 installation script".to_string(),
                            "gitlab-runner register --url https://gitlab.com/ --token YOUR_TOKEN".to_string(),
                        ),
                        ("linux", "arm64") => (
                            "#!/bin/bash\n# Linux ARM64 installation script".to_string(),
                            "gitlab-runner register --url https://gitlab.com/ --token YOUR_TOKEN".to_string(),
                        ),
                        ("windows", "amd64") => (
                            "powershell -Command \"# Windows AMD64 installation script\"".to_string(),
                            "gitlab-runner.exe register --url https://gitlab.com/ --token YOUR_TOKEN".to_string(),
                        ),
                        _ => {
                            errors.push(format!("Unsupported OS/architecture combination: {}/{}", os, arch));
                            (String::new(), String::new())
                        }
                    }
                } else {
                    (String::new(), String::new())
                };
                
                RunnerInstructions {
                    os,
                    arch,
                    install_script,
                    register_command,
                    errors,
                }
            }
        }
    }
} 