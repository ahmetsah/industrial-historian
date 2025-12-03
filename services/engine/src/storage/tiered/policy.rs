#![allow(dead_code)]
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Clone, Debug)]
pub struct TieringPolicy {
    pub max_age_ms: i64,
}

impl TieringPolicy {
    pub fn new(max_age_ms: i64) -> Self {
        Self { max_age_ms }
    }

    pub fn is_eligible(&self, timestamp_ms: i64) -> bool {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_millis() as i64;
        timestamp_ms < (now - self.max_age_ms)
    }
}
