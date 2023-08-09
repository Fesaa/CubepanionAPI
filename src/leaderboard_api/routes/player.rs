use actix_web::{get, web::{Data, Path}, Responder, HttpResponse};
use diesel::result::Error::NotFound;

use crate::database::{API, leaderboard::messages::FetchLeaderboardFromPlayer};

/// Get all LeaderboardRow for a player
#[utoipa::path(get,
    responses(
        (status = 200, description = "All leaderboards for the player", body = Vec<LeaderboardRow>),
        (status = 400, description = "Invalid name"),
        (status = 404, description = "No leaderboards found"),
        (status = 500, description = "SQL error", example = json!(HttpResponse::InternalServerError().body("Unable to retrieve leaderboards")))
    ),
    params(
        ("name", description = "Player name")
    )
)]
#[get("/leaderboard_api/player/{name}")]
pub async fn get_leaderboards_from_player(state: Data<API>, path: Path<String>) -> impl Responder {
    let name = path.into_inner();
    if !state.username_regex.is_match(&name) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &name + ">")
    }

    match state.db.send(FetchLeaderboardFromPlayer{player_name: name}).await {
        Ok(Ok(leaderboards)) => HttpResponse::Ok().json(leaderboards),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No leaderboards found"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve leaderboards: {}", err)),
        },
        _ => HttpResponse::InternalServerError().body("Unable to retrieve leaderboards"),
    }
}