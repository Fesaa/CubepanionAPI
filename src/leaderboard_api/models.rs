use serde::{Deserialize, Serialize};
use utoipa::ToSchema;

#[derive(Deserialize)]
pub struct Players {
    pub players: Vec<String>,
}

#[derive(Deserialize, ToSchema)]
pub struct LeaderboardSubmission {
    pub uuid: String,
    pub unix_time_stamp: i64,
    pub game: String,
    pub entries: Vec<LeaderboardEntry>
}

#[derive(Deserialize, Serialize)]
pub struct LeaderboardEntry {
    pub player: String,
    pub position: i32,
    pub score: i32,
    pub game: String
}

#[derive(Deserialize, Serialize)]
pub struct LeaderboardGame {
    pub game: String,
    pub display_name: String,
    pub aliases: Vec<String>,
    pub active: bool,
    pub score_type: String,
}
