use chrono::{NaiveDate, TimeZone, Utc};

pub trait ParseCommitDate {
    fn convert_date_to_epoch(&self, date: Option<&str>) -> Option<i64> {
        if let Some(date_str) = date {
            match NaiveDate::parse_from_str(date_str, "%Y-%m-%d") {
                Ok(naive_date) => {
                    let datetime = Utc.from_utc_datetime(&naive_date.and_hms_opt(0, 0, 0).unwrap());
                    Some(datetime.timestamp())
                }
                Err(_) => None,
            }
        } else {
            None
        }
    }
}
