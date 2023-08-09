use actix::Message;
use diesel::QueryResult;
use super::EggWarsMapJson;


#[derive(Message)]
#[rtype(result = "QueryResult<Vec<EggWarsMapJson>>")]
pub struct FetchEggWarsMaps {
}

#[derive(Message)]
#[rtype(result = "QueryResult<EggWarsMapJson>")]
pub struct FetchEggWarsMap {
    pub name: String,
}

