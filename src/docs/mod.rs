use utoipa::OpenApi;

use crate::database::schema::{ChestLocation, LeaderboardRow, EggWarsMap};
use crate::leaderboard_api::routes::game::BoundedRequest;

use super::eggwars_map_api as EggWarsMaps;
use super::chest_api as Chests;
use super::leaderboard_api::{models::LeaderboardSubmission, routes::submission, routes::player, routes::game};


#[derive(OpenApi)]
#[openapi(
    paths(
        submission::submit_leaderboard_entries,
        player::get_leaderboards_from_player,
        game::get_leaderboard,
        game::get_leaderboard_between,
        game::get_games,
        Chests::get_current_chests,
        Chests::get_season_chests,
        Chests::get_seasons,
        EggWarsMaps::get_all_eggwars_maps,
        EggWarsMaps::get_eggwars_map
    ),
    components(
        schemas(LeaderboardSubmission, LeaderboardRow, BoundedRequest, ChestLocation, EggWarsMap)
    ),
    info(
        title = "Cubepanion API",
        description = "Backend for https://github.com/Fesaa/Cubepanion",
        license(
            name = "GNU General Public License v3.0",
            url = "https://github.com/Fesaa/CubepanionAPI/blob/master/LICENSE"
        ),
        contact(
            name = "Amelia",
            url = "https://github.com/Fesaa/CubepanionAPI"
        ),
        version = "1.0.0"
    )

)]
pub struct ApiDoc;
