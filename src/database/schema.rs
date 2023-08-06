use serde::{Deserialize, Serialize};

table! {
    submissions (unix_time_stamp) {
        uuid -> Varchar,
        unix_time_stamp -> Bigint,
        game -> Varchar,
        valid -> Bool,
    }
}

#[derive(Deserialize, Serialize, Queryable)]
pub struct LeaderboardRow {
    pub game: String,
    pub player: String,
    pub position: i32,
    pub score: i32,
    pub unix_time_stamp: i64,
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
