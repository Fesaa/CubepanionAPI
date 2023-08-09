use super::schema::{GenLayout, EggWarsMap};
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
    pub gen_layout: Vec<GenLayout>,
}


impl EggWarsMapJson {

    pub fn from_eggwars_map(map: &EggWarsMap) -> EggWarsMapJson {
        EggWarsMapJson {
            unique_name: map.unique_name.clone(),
            map_name: map.map_name.clone(),
            team_size: map.team_size,
            build_limit: map.build_limit,
            colours: map.colours.clone(),
            gen_layout: vec![],
        }
    }

    pub fn add_gen_layout(&mut self, gen_layout: GenLayout) {
        self.gen_layout.push(gen_layout);
    }

}