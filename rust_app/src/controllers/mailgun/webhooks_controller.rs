// Ported from: orig_app/app/controllers/mailgun/webhooks_controller.rb
// Ported: 2025-05-01
// Handles Mailgun webhook events with HMAC signature validation and feature flag check.

use actix_web::{post, web, HttpRequest, HttpResponse, Responder};
use hmac::{Hmac, Mac};
use serde::Deserialize;
use sha2::Sha256;
use std::env;

// --- Settings stub (replace with real settings integration) ---
fn mailgun_events_enabled() -> bool {
    // TODO: Replace with real application settings
    env::var("MAILGUN_EVENTS_ENABLED").unwrap_or_else(|_| "true".into()) == "true"
}
fn mailgun_signing_key() -> Option<String> {
    // TODO: Replace with real application settings
    env::var("MAILGUN_SIGNING_KEY").ok()
}
// ------------------------------------------------------------

#[derive(Deserialize)]
struct Signature {
    timestamp: String,
    token: String,
    signature: String,
}

#[derive(Deserialize)]
struct WebhookPayload {
    signature: Signature,
    #[serde(rename = "event-data")]
    event_data: serde_json::Value,
}

#[post("/mailgun/webhooks")]
pub async fn process_webhook(payload: web::Json<WebhookPayload>) -> impl Responder {
    if !mailgun_events_enabled() {
        // 406 Not Acceptable: feature not enabled
        return HttpResponse::NotAcceptable().finish();
    }
    let signing_key = match mailgun_signing_key() {
        Some(key) if !key.is_empty() => key,
        _ => return HttpResponse::NotFound().finish(),
    };
    if !valid_signature(&payload.signature, &signing_key) {
        // 404 Not Found: invalid signature
        return HttpResponse::NotFound().finish();
    }
    // TODO: Call webhook processors here
    HttpResponse::Ok().finish()
}

fn valid_signature(sig: &Signature, key: &str) -> bool {
    use subtle::ConstantTimeEq;
    let data = format!("{}{}", sig.timestamp, sig.token);
    let mut mac = Hmac::<Sha256>::new_from_slice(key.as_bytes()).unwrap();
    mac.update(data.as_bytes());
    let expected = hex::encode(mac.finalize().into_bytes());
    sig.signature.as_bytes().ct_eq(expected.as_bytes()).into()
}
