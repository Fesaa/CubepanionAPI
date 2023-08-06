use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
pub struct Players {
    pub players: Vec<String>,
}

#[derive(Deserialize, Serialize, Queryable)]
pub struct LeaderboardRow {
    pub game: String,
    pub player: String,
    pub position: i32,
    pub score: i32,
    pub unix_time_stamp: i64,
}

#[derive(Deserialize)]
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