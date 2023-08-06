use actix::Message;
use diesel::QueryResult;
use crate::database::schema::LeaderboardRow;

#[derive(Message)]
#[rtype(result = "QueryResult<Vec<LeaderboardRow>>")]
pub struct FetchLeaderboardFromPlayer {
    pub player_name: String
}

#[derive(Message)]
#[rtype(result = "QueryResult<Vec<LeaderboardRow>>")]
pub struct FetchLeaderboardForGame{
    pub game_name: String,
    pub min: i32,
    pub max: i32,
}