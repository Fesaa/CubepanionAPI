use actix_web::{web::{Data, Path, Query}, Responder, HttpResponse, get};
use diesel::result::Error::NotFound;
use serde::Deserialize;
use utoipa::{ToSchema, IntoParams};
use crate::{API, database::leaderboard::messages::{FetchLeaderboardForGame, GetGames}};

const GAME: &'static str = "[GET] LB game";
const BOUNDED_GAME: &'static str = "[GET] LB game - bounded";
const GAMES: &'static str = "[GET] LB games";

/// Get all LeaderboardRow for a game
#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All leaderboards for the game", body = Vec<LeaderboardRow>),
        (status = 400, description = "Invalid game"),
        (status = 404, description = "No leaderboards found"),
        (status = 500, description = "SQL error", example = json!(HttpResponse::InternalServerError().body("Unable to retrieve leaderboards")))
    ),
    params(
        ("game", description = "Game name")
    )
)]
#[get("/leaderboard_api/leaderboard/{game}")]
pub async fn get_leaderboard(state: Data<API>, path: Path<String>) -> impl Responder {
    let game = path.into_inner();
    if !state.username_regex.is_match(&game) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &game + ">")
    }
    match state.db.send(FetchLeaderboardForGame{game_name: game, max: 200, min: 1}, GAME).await {
        Ok(Ok(leaderboards)) => HttpResponse::Ok().json(leaderboards),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No leaderboards found"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve leaderboards: {}", err)),
        },
        _ => HttpResponse::InternalServerError().body("Unable to retrieve leaderboards"),
    }
}

#[derive(Deserialize, ToSchema, IntoParams)]
pub struct BoundedRequest {
    lower: i32,
    upper: i32,
}

/// Get all LeaderboardRow for a game between the bounds	
#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All leaderboards for the game between the bounds", body = Vec<LeaderboardRow>),
        (status = 400, description = "Invalid game"),
        (status = 404, description = "No leaderboards found"),
        (status = 500, description = "SQL error", example = json!(HttpResponse::InternalServerError().body("Unable to retrieve leaderboards")))
    ),
    params(
        ("game", description = "Game name"),
        BoundedRequest
    )
)]
#[get("/leaderboard_api/leaderboard/{game}/bounded")]
pub async fn get_leaderboard_between(state: Data<API>, path: Path<String>, info: Query<BoundedRequest>) -> impl Responder {
    let game = path.into_inner();
    if !state.username_regex.is_match(&game) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &game + ">")
    }
    match state.db.send(FetchLeaderboardForGame{game_name: game, max: info.upper, min: info.lower}, BOUNDED_GAME).await {
        Ok(Ok(leaderboards)) => HttpResponse::Ok().json(leaderboards),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No leaderboards found"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve leaderboards: {}", err)),
        },
        _ => HttpResponse::InternalServerError().body("Unable to retrieve leaderboards"),
    }
}


/// Get all (active) games
#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All games", body = Vec<GameRow>),
        (status = 404, description = "No games found"),
        (status = 500, description = "SQL error", example = json!(HttpResponse::InternalServerError().body("Unable to retrieve games")))
    ),
    params(
        ("active", description = "Whether to only get active games")
    )
)]
#[get("/leaderboard_api/games/{active}")]
pub async fn get_games(state: Data<API>, path: Path<bool>) -> impl Responder {
    let active = path.into_inner();
    match state.db.send(GetGames{must_be_active: active}, GAMES).await {
        Ok(Ok(games)) => HttpResponse::Ok().json(games),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No games found"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve games: {}", err)),
        },
        _ => HttpResponse::InternalServerError().body("Unable to retrieve games"),
    }
}


