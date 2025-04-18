use chrono::{DateTime, NaiveDate, TimeZone, Utc};
use std::str::FromStr;

pub trait ParseCommitDate {
    fn convert_date_to_epoch(&self, date: Option<&str>) -> Option<i64>;
}

pub struct ParseCommitDateImpl;

impl ParseCommitDateImpl {
    pub fn new() -> Self {
        Self
    }
}

impl ParseCommitDate for ParseCommitDateImpl {
    fn convert_date_to_epoch(&self, date: Option<&str>) -> Option<i64> {
        date.and_then(|date_str| {
            NaiveDate::parse_from_str(date_str, "%Y-%m-%d")
                .ok()
                .map(|date| {
                    let datetime = Utc.from_utc_datetime(&date.and_hms(0, 0, 0));
                    datetime.timestamp()
                })
        })
    }
} 