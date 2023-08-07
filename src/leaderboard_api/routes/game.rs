use actix_web::{web::{Data, Path, Query}, Responder, HttpResponse, get};
use serde::Deserialize;
use crate::{API, database::leaderboard::messages::FetchLeaderboardForGame};

#[get("/leaderboard_api/leaderboard/{game}")]
pub async fn get_leaderboard(state: Data<API>, path: Path<String>) -> impl Responder {
    let game = path.into_inner();
    if !state.username_regex.is_match(&game) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &game + ">")
    }
    match state.db.send(FetchLeaderboardForGame{game_name: game, max: 200, min: 1}).await {
        Ok(Ok(leaderboards)) => HttpResponse::Ok().json(leaderboards),
        Ok(Err(err)) => HttpResponse::InternalServerError().body(format!("Unable to retrieve leaderboards: {}", err)),
        _ => HttpResponse::InternalServerError().body("Unable to retrieve leaderboards"),
    }
}

#[derive(Deserialize)]
pub struct BoundedRequest {
    lower: i32,
    upper: i32,
}

#[get("/leaderboard_api/leaderboard/{game}/bounded")]
pub async fn get_leaderboard_between(state: Data<API>, path: Path<String>, info: Query<BoundedRequest>) -> impl Responder {
    let game = path.into_inner();
    if !state.username_regex.is_match(&game) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &game + ">")
    }
    match state.db.send(FetchLeaderboardForGame{game_name: game, max: info.upper, min: info.lower}).await {
        Ok(Ok(leaderboards)) => HttpResponse::Ok().json(leaderboards),
        Ok(Err(err)) => HttpResponse::InternalServerError().body(format!("Unable to retrieve leaderboards: {}", err)),
        _ => HttpResponse::InternalServerError().body("Unable to retrieve leaderboards"),
    }
}