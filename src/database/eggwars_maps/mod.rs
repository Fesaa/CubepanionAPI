use super::schema::{EggWarsMap, Generator};
use serde::{Serialize, Deserialize};

pub mod handlers;
pub mod messages;

#[derive(Serialize, Deserialize)]
pub struct EggWarsMapJson {
    pub unique_name: String,
    pub map_name: String,
    pub team_size: i32,
    pub build_limit: i32,
    pub colours: String,
    pub layout: String,
    pub generators: Vec<Generator>,
}


impl EggWarsMapJson {

    pub fn from_eggwars_map(map: &EggWarsMap) -> EggWarsMapJson {
        EggWarsMapJson {
            unique_name: map.unique_name.clone(),
            map_name: map.map_name.clone(),
            team_size: map.team_size,
            build_limit: map.build_limit,
            colours: map.colours.clone(),
            layout: map.layout.clone(),
            generators: vec![],
        }
    }

    pub fn add_gen_layout(&mut self, gen: Generator) {
        self.generators.push(gen);
    }

}