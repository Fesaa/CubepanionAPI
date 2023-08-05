use actix_web::{web::{Data, Path, Query}, Responder, HttpResponse, get};
use serde::Deserialize;
use crate::{API, leaderboard_api::models::LeaderboardRow};

#[get("/leaderboard_api/leaderboard/{game}")]
pub async fn get_leaderboard(state: Data<API>, path: Path<String>) -> impl Responder {
    let game = path.into_inner();
    if !state.username_regex.is_match(&game) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &game + ">")
    }
    match sqlx::query_as::<_, LeaderboardRow>("
    SELECT 
        player,position,score,game
    FROM
        leaderboards
    WHERE
        game = $1
    AND
        unix_time_stamp
    = (SELECT
            MAX(unix_time_stamp)
        FROM
            submissions
        WHERE
            valid = TRUE
        AND
            game = $1)
    ORDER BY
        position
    ASC;")
        .bind(game)
        .fetch_all(&state.db)
        .await {
            Ok(leaderboards) => HttpResponse::Ok().json(leaderboards),
            Err(err) => HttpResponse::InternalServerError().body(err.to_string())
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
    match sqlx::query_as::<_, LeaderboardRow>("
    SELECT 
        player,position,score,game
    FROM
        leaderboards
    WHERE
        game = $1
    AND
        unix_time_stamp
    = (SELECT
            MAX(unix_time_stamp)
        FROM
            submissions
        WHERE
            valid = TRUE
        AND
            game = $1)
    AND
        position
    BETWEEN
        $2
    AND
        $3
    ORDER BY
        position
    ASC;")
        .bind(game)
        .bind(info.lower)
        .bind(info.upper)
        .fetch_all(&state.db)
        .await {
            Ok(leaderboards) => HttpResponse::Ok().json(leaderboards),
            Err(err) => HttpResponse::InternalServerError().body(err.to_string())
        }
}