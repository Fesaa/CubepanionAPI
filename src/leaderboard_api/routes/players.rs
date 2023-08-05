use actix_web::{post, web::{Data, Json}, Responder, HttpResponse};

use crate::{API, leaderboard_api::models::{Players, LeaderboardRow}};


#[post("/leaderboard_api/leaderboard/players")]
pub async fn get_leaderboard_for_all(state: Data<API>, body: Json<Players>) -> impl Responder {
    match sqlx::query_as::<_, LeaderboardRow>("
    SELECT l.*
    FROM leaderboards l
    INNER JOIN (
        SELECT game, MAX(unix_time_stamp) AS max_time_stamp
        FROM submissions
        WHERE valid = TRUE
        GROUP BY game
    ) s ON l.game = s.game AND l.unix_time_stamp = s.max_time_stamp
    WHERE UPPER(l.player) = ANY($1)
    ORDER BY
        l.player,l.position
    ASC;")
        .bind(&body.players.iter().map(|p| p.to_uppercase()).collect::<Vec<String>>())
        .fetch_all(&state.db)
        .await {
            Ok(leaderboards) => HttpResponse::Ok().json(leaderboards),
            Err(err) => HttpResponse::InternalServerError().body(err.to_string())
        }
}