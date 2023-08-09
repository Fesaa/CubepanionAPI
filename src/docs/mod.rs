use utoipa::OpenApi;

use crate::database::schema::LeaderboardRow;
use crate::leaderboard_api::routes::game::BoundedRequest;

use super::leaderboard_api;
use super::leaderboard_api::models::LeaderboardSubmission;


#[derive(OpenApi)]
#[openapi(
    paths(
        leaderboard_api::routes::submission::submit_leaderboard_entries,
        leaderboard_api::routes::player::get_leaderboards_from_player,
        leaderboard_api::routes::game::get_leaderboard,
        leaderboard_api::routes::game::get_leaderboard_between
    ),
    components(
        schemas(LeaderboardSubmission, LeaderboardRow, BoundedRequest)
    )

)]
pub struct ApiDoc;