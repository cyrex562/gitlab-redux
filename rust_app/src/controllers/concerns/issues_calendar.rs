use actix_web::{web, HttpResponse, Responder};
use serde::Serialize;
use crate::models::issuable::Issuable;
use crate::config::Settings;

pub trait IssuesCalendar {
    fn render_issues_calendar(&self, issuables: Vec<Issuable>) -> impl Responder;
}

pub struct IssuesCalendarImpl {
    settings: Settings,
}

impl IssuesCalendarImpl {
    pub fn new(settings: Settings) -> Self {
        Self { settings }
    }
}

impl IssuesCalendar for IssuesCalendarImpl {
    fn render_issues_calendar(&self, mut issuables: Vec<Issuable>) -> impl Responder {
        // Filter and limit issuables
        let filtered_issuables: Vec<Issuable> = issuables
            .into_iter()
            .filter(|i| !i.is_archived() && i.has_due_date())
            .take(100)
            .collect();

        // Create response with appropriate content type
        HttpResponse::Ok()
            .content_type("text/calendar")
            .body(self.generate_ics_content(&filtered_issuables))
    }
}

impl IssuesCalendarImpl {
    fn generate_ics_content(&self, issuables: &[Issuable]) -> String {
        // Basic ICS format implementation
        let mut ics_content = String::from(
            "BEGIN:VCALENDAR\r\n\
             VERSION:2.0\r\n\
             PRODID:-//GitLab//Issues Calendar//EN\r\n"
        );

        for issuable in issuables {
            ics_content.push_str(&format!(
                "BEGIN:VEVENT\r\n\
                 UID:{}\r\n\
                 DTSTAMP:{}\r\n\
                 DTSTART:{}\r\n\
                 SUMMARY:{}\r\n\
                 END:VEVENT\r\n",
                issuable.id(),
                chrono::Utc::now().format("%Y%m%dT%H%M%SZ"),
                issuable.due_date().unwrap().format("%Y%m%d"),
                issuable.title()
            ));
        }

        ics_content.push_str("END:VCALENDAR\r\n");
        ics_content
    }
} 