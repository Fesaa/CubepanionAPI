use actix::Message;
use diesel::QueryResult;
use crate::database::schema::ChestLocation;

#[derive(Message)]
#[rtype(result = "QueryResult<Vec<ChestLocation>>")]
pub struct FetchChestLocationsForRunningSeason {
}

#[derive(Message)]
#[rtype(result = "QueryResult<Vec<ChestLocation>>")]
pub struct FetchChestLocationsForSeason {
    pub season: String
}

#[derive(Message)]
#[rtype(result = "QueryResult<Vec<String>>")]
pub struct FetchSeasons {
    pub running: bool
}