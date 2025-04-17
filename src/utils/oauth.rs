use std::fmt;

#[derive(Debug)]
pub enum OAuthError {
    ConfigurationMissing,
    InvalidUrl,
    TokenExchangeFailed,
    InvalidState,
    InvalidCode,
    NetworkError(String),
}

impl fmt::Display for OAuthError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            OAuthError::ConfigurationMissing => write!(f, "OAuth configuration is missing"),
            OAuthError::InvalidUrl => write!(f, "Invalid OAuth URL"),
            OAuthError::TokenExchangeFailed => write!(f, "Failed to exchange code for token"),
            OAuthError::InvalidState => write!(f, "Invalid OAuth state"),
            OAuthError::InvalidCode => write!(f, "Invalid OAuth code"),
            OAuthError::NetworkError(msg) => write!(f, "Network error: {}", msg),
        }
    }
}

impl std::error::Error for OAuthError {} 