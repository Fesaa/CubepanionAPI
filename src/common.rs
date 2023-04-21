use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
pub struct LeaderboardEntry {
    pub player: String,
    pub position: i32,
    pub score: i32,
    pub game: String
}

use uuid::Uuid;

pub fn is_valid_uuid(uuid_string: &str) -> bool {
    match Uuid::parse_str(uuid_string) {
        Ok(_) => true,
        Err(_) => false,
    }
}
