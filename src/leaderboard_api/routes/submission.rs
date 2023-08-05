use actix_web::{post, web::{Data, Json}, Responder, HttpResponse};
use uuid::Uuid;

use crate::{API, leaderboard_api::models::LeaderboardSubmission};

pub fn is_valid_uuid(uuid_string: &str) -> bool {
    match Uuid::parse_str(uuid_string) {
        Ok(_) => true,
        Err(_) => false,
    }
}


#[post("/leaderboard_api")]
pub async fn submit_leaderboard_entries(state: Data<API>, body: Json<LeaderboardSubmission>) -> impl Responder {
    if body.entries.len() != 200 {
        return HttpResponse::BadRequest().body(String::from("Invalid entries length <") + &body.entries.len().to_string() + ">")
    }
    if !is_valid_uuid(&body.uuid) {
        return HttpResponse::BadRequest().body(String::from("Invalid uuid <") + &body.uuid + ">")
    }

    match sqlx::query(
        "INSERT INTO submissions (uuid, unix_time_stamp, game) VALUES ($1, $2, $3)"
    )
    .bind(body.uuid.to_string())
    .bind(body.unix_time_stamp)
    .bind(body.game.to_string())
    .execute(&state.db)
    .await {
        Err(err) => HttpResponse::InternalServerError().body(err.to_string()),
        Ok(_) => {
            let sql: String = String::from("
                INSERT INTO
                    leaderboards
                    (game, player, position, score, unix_time_stamp) 
                SELECT 
                    *
                FROM
                    UNNEST($1::varchar[], $2::varchar[], $3::int8[], $4::int8[], $5::int8[]);
                    ");

            let mut game_vec: Vec<String> = Vec::new();
            let mut player_vec: Vec<String> = Vec::new();
            let mut position_vec: Vec<i32> = Vec::new();
            let mut score_vec: Vec<i32> = Vec::new();
            let mut unix_vec: Vec<i64> = Vec::new();

            body.entries.iter().for_each(|row| {
                game_vec.push(row.game.to_owned());
                player_vec.push(row.player.to_owned());
                position_vec.push(row.position);
                score_vec.push(row.score);
                unix_vec.push(body.unix_time_stamp);
            });

            match sqlx::query(&sql)
            .bind(game_vec)
            .bind(player_vec)
            .bind(position_vec)
            .bind(score_vec)
            .bind(unix_vec)
            .execute(&state.db)
            .await {
                Err(err) => {
                    remove_submission(body.unix_time_stamp, state).await;
                    HttpResponse::InternalServerError().body(err.to_string())
                },
                Ok(_) => HttpResponse::Accepted().body("Success")
            }
        }
    }
}

async fn remove_submission(unix_time_stamp: i64, state: Data<API>) {
    match sqlx::query("UPDATE submissions SET valid = false WHERE unix_time_stamp = $1;")
        .bind(unix_time_stamp)
        .execute(&state.db)
        .await {
            Ok(_) => (),
            Err(err) => println!("{}", String::from("Could not update submission: ") + &err.to_string()),
        }
}