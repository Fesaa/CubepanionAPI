use actix::Message;
use diesel::QueryResult;
use crate::database::schema::{LeaderboardRow, SubmissionRow};
use crate::leaderboard_api::models::LeaderboardGame;

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

#[derive(Message)]
#[rtype(result = "QueryResult<usize>")]
pub struct InsertLeaderboardRows {
    pub rows: Vec<LeaderboardRow>
}

#[derive(Message)]
#[rtype(result = "QueryResult<usize>")]
pub struct InsertSubmission {
    pub sub: SubmissionRow
}


#[derive(Message)]
#[rtype(result = "QueryResult<usize>")]
pub struct DisableSubmission {
    pub unix: i64
}

#[derive(Message)]
#[rtype(result = "QueryResult<Vec<LeaderboardGame>>")]
pub struct GetGames {
    pub must_be_active: bool
}
