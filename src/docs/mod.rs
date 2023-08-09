use utoipa::OpenApi;

use crate::database::schema::{ChestLocation, LeaderboardRow, EggWarsMap};
use crate::leaderboard_api::routes::game::BoundedRequest;

use super::eggwars_map_api;
use super::chest_api;
use super::leaderboard_api::{self, models::LeaderboardSubmission};


#[derive(OpenApi)]
#[openapi(
    paths(
        leaderboard_api::routes::submission::submit_leaderboard_entries,
        leaderboard_api::routes::player::get_leaderboards_from_player,
        leaderboard_api::routes::game::get_leaderboard,
        leaderboard_api::routes::game::get_leaderboard_between,
        chest_api::get_current_chests,
        chest_api::get_season_chests,
        chest_api::get_seasons,
        eggwars_map_api::get_all_eggwars_maps,
        eggwars_map_api::get_eggwars_map
    ),
    components(
        schemas(LeaderboardSubmission, LeaderboardRow, BoundedRequest, ChestLocation, EggWarsMap)
    )

)]
pub struct ApiDoc;