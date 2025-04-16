pub mod common;
pub mod html_format_actions_support;
pub mod json_format_actions_support;
pub mod rest_api_actions_support;

use crate::spammable::Spammable;
use actix_web::{web, HttpResponse};
use std::future::Future;
