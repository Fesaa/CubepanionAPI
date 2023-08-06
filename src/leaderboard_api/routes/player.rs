use actix_web::{get, web::{Data, Path}, Responder, HttpResponse};

use crate::database::{API, leaderboard::messages::FetchLeaderboardFromPlayer};

#[get("/leaderboard_api/player/{name}")]
pub async fn get_leaderboards_from_player(state: Data<API>, path: Path<String>) -> impl Responder {
    let name = path.into_inner();
    if !state.username_regex.is_match(&name) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &name + ">")
    }

    match state.db.send(FetchLeaderboardFromPlayer{player_name: name}).await {
        Ok(Ok(leaderboards)) => HttpResponse::Ok().json(leaderboards),
        Ok(Err(_)) => HttpResponse::NotFound().finish(),
        _ => HttpResponse::InternalServerError().body("Unable to retrieve leaderboards"),
    }
}