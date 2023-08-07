use serde::{Deserialize, Serialize};

// leaderboard_api
use crate::leaderboard_api::models::LeaderboardEntry;

table! {
    submissions (unix_time_stamp) {
        uuid -> Varchar,
        unix_time_stamp -> Bigint,
        game -> Varchar,
        valid -> Bool,
    }
}

#[derive(Deserialize, Serialize, Queryable, Insertable)]
#[diesel(table_name = submissions)]
pub struct SubmissionRow {
    pub uuid: String,
    pub unix_time_stamp: i64,
    pub game: String,
    pub valid: bool
}

#[derive(Deserialize, Serialize, Queryable, Insertable)]
#[diesel(table_name = leaderboards)]
pub struct LeaderboardRow {
    pub game: String,
    pub player: String,
    pub position: i32,
    pub score: i32,
    pub unix_time_stamp: i64,
}

impl LeaderboardRow {

    pub fn from_entry(entry: &LeaderboardEntry, unix: i64) -> LeaderboardRow {
        LeaderboardRow { game: entry.game.clone(), player: entry.player.clone(), position: entry.position, score: entry.score, unix_time_stamp: unix}
    }

}

table! {
    leaderboards (unix_time_stamp) {
        game -> VarChar,
        player -> VarChar,
        position -> Integer,
        score -> Integer,
        unix_time_stamp -> Bigint,
    }
}

table! {
    ban (uuid) {
        uuid -> Varchar,
    }
}

joinable!(leaderboards -> submissions (unix_time_stamp));
allow_tables_to_appear_in_same_query!(leaderboards, submissions);

// chest_api

table! {
    seasons (season_name) {
        season_name -> VarChar,
        running -> Bool,
    }
}

table! {
    chest_locations (season_name) {
        season_name -> VarChar,
        x -> Integer,
        y -> Integer,
        z -> Integer,
    }
}

#[derive(Deserialize, Serialize, Queryable)]
#[diesel(table_name = chest_locations)]
pub struct ChestLocation {
    pub season_name: String,
    pub x: i32,
    pub y: i32,
    pub z: i32,
}

joinable!(chest_locations -> seasons (season_name));
allow_tables_to_appear_in_same_query!(chest_locations, seasons);