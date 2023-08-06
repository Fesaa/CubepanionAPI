use actix::Handler;
use diesel::{QueryResult, QueryDsl, RunQueryDsl};
use diesel::{self, prelude::*};

use crate::database::{DbActor, schema::ChestLocation};

use super::messages::{FetchChestLocationsForRunningSeason, FetchChestLocationsForSeason, FetchSeasons};

impl Handler<FetchChestLocationsForRunningSeason> for DbActor {
    type Result = QueryResult<Vec<ChestLocation>>;

    fn handle(&mut self, _msg: FetchChestLocationsForRunningSeason, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::seasons::dsl::{seasons, season_name, running};
        use crate::database::schema::chest_locations::dsl::{chest_locations, season_name as season};

        let mut con = self.0.get()
        .expect("Fetch Leaderboard From Player: Unable to establish connection");

        let current_season: String = seasons
            .filter(running)
            .select(season_name)
            .first::<String>(&mut con)?;

        chest_locations
        .filter(season.eq(current_season))
        .load::<ChestLocation>(&mut con)
    }
}

impl Handler<FetchChestLocationsForSeason> for DbActor {
    type Result = QueryResult<Vec<ChestLocation>>;

    fn handle(&mut self, msg: FetchChestLocationsForSeason, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::chest_locations::dsl::{chest_locations, season_name as season};

        let mut con = self.0.get()
        .expect("Fetch Leaderboard From Player: Unable to establish connection");

        chest_locations
        .filter(season.eq(msg.season))
        .load::<ChestLocation>(&mut con)
    }
}

impl Handler<FetchSeasons> for DbActor {
    type Result = QueryResult<Vec<String>>;

    fn handle(&mut self, msg: FetchSeasons, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::seasons::dsl::{seasons, season_name, running};

        let mut con = self.0.get()
        .expect("Fetch Leaderboard From Player: Unable to establish connection");

        if msg.running {
            seasons
            .filter(running)
            .select(season_name)
            .load::<String>(&mut con)
        } else {
            seasons
            .select(season_name)
            .load::<String>(&mut con)
        }
    }
}