use actix_web::{get, web::{Data, Path}, Responder, HttpResponse};

use crate::{API, leaderboard_api::models::LeaderboardRow};

#[get("/leaderboard_api/player/{name}")]
pub async fn get_leaderboards_from_player(state: Data<API>, path: Path<String>) -> impl Responder {
    let name = path.into_inner();
    if !state.username_regex.is_match(&name) {
        return HttpResponse::BadRequest().body(String::from("Invalid name <") + &name + ">")
    }
    match sqlx::query_as::<_, LeaderboardRow>("
    SELECT 
        player,position,score,game
    FROM
        leaderboards
    WHERE
        (game, unix_time_stamp)
    IN (SELECT
            game, MAX(unix_time_stamp)
        FROM
            submissions
        WHERE
            valid = TRUE
        GROUP BY
            game)
    AND
        UPPER(player) = UPPER($1)
    ORDER BY
        position
    ASC;")
        .bind(name)
        .fetch_all(&state.db)
        .await {
            Ok(leaderboards) => HttpResponse::Ok().json(leaderboards),
            Err(err) => HttpResponse::InternalServerError().body(err.to_string())
        }
}