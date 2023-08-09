use actix::Handler;
use diesel::{QueryResult, QueryDsl, RunQueryDsl};
use diesel::{self, prelude::*};

use crate::database::DbActor;
use crate::database::schema::{EggWarsMap, Generator};

use super::EggWarsMapJson;
use super::messages::{FetchEggWarsMaps, FetchEggWarsMap};



impl Handler<FetchEggWarsMaps> for DbActor {
    type Result = QueryResult<Vec<EggWarsMapJson>>;

    fn handle(&mut self, _msg: FetchEggWarsMaps, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::eggwars_maps::dsl::eggwars_maps;
        use crate::database::schema::generators::dsl::{generators, ordering};

        let mut con = self.0.get()
        .expect("Fetch Leaderboard From Player: Unable to establish connection");

        let mut maps = eggwars_maps
            .load::<EggWarsMap>(&mut con)?
            .iter()
            .map(|map| EggWarsMapJson::from_eggwars_map(map))
            .collect::<Vec<EggWarsMapJson>>();

        let gens = generators
            .order(ordering)
            .load::<Generator>(&mut con)?;

        for gen in gens {
            let map_option = maps.iter_mut().find(|map| map.unique_name == gen.unique_name);
            if let Some(map) = map_option {
                map.add_gen_layout(gen);
            }
        }

        QueryResult::Ok(maps)
    }
}

impl Handler<FetchEggWarsMap> for DbActor {
    type Result = QueryResult<EggWarsMapJson>;

    fn handle(&mut self, msg: FetchEggWarsMap, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::eggwars_maps::dsl::{eggwars_maps, unique_name};
        use crate::database::schema::generators::dsl::{generators, unique_name as name, ordering};

        let mut con = self.0.get()
        .expect("Fetch Leaderboard From Player: Unable to establish connection");


        let map: EggWarsMap = eggwars_maps
            .filter(unique_name.eq(&msg.name))
            .first::<EggWarsMap>(&mut con)?;

            let mut map_json = EggWarsMapJson::from_eggwars_map(&map);

        let gens: Vec<Generator> = generators
            .filter(name.eq(&msg.name))
            .order(ordering)
            .load::<Generator>(&mut con)?;

        for gen in gens {
            map_json.add_gen_layout(gen);
        }

        QueryResult::Ok(map_json)
    }
}