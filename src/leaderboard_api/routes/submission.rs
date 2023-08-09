use actix_web::{post, web::{Data, Json}, Responder, HttpResponse};
use uuid::Uuid;

use crate::{API, leaderboard_api::models::LeaderboardSubmission, database::{schema::{LeaderboardRow, SubmissionRow}, leaderboard::messages::{InsertLeaderboardRows, InsertSubmission, DisableSubmission}}};

pub fn is_valid_uuid(uuid_string: &str) -> bool {
    match Uuid::parse_str(uuid_string) {
        Ok(_) => true,
        Err(_) => false,
    }
}

// Submit a leaderboard
//
//  Inserts the data into the database, requests are validated by the server
#[utoipa::path(post,
    request_body = LeaderboardSubmission,
    responses(
        (status = 202, description = "Accepted your submission"),
        (status = 400, description = "Invalid submission data"),
        (status = 500, description = "SQL error", example = json!(HttpResponse::InternalServerError().body("Submission insert didn't go through, won't be queried:")))
    )
)]
#[post("/leaderboard_api")]
pub async fn submit_leaderboard_entries(state: Data<API>, body: Json<LeaderboardSubmission>) -> impl Responder {
    if body.entries.len() != 200 {
        return HttpResponse::BadRequest().body(String::from("Invalid entries length <") + &body.entries.len().to_string() + ">")
    }
    if !is_valid_uuid(&body.uuid) {
        return HttpResponse::BadRequest().body(String::from("Invalid uuid <") + &body.uuid + ">")
    }

    let rows: Vec<LeaderboardRow> = body.entries
        .iter()
        .map(|entry| LeaderboardRow::from_entry(entry, body.unix_time_stamp))
        .collect::<Vec<LeaderboardRow>>();

        match state.db.send(InsertSubmission{sub: SubmissionRow{uuid: body.uuid.clone(), game: body.game.clone(), valid: true, unix_time_stamp: body.unix_time_stamp}}).await {
            Ok(Ok(_)) => (),
            Ok(Err(err)) => return HttpResponse::InternalServerError().body(format!("Submission insert didn't go through, won't be queried: {}", err)),
            _ => return HttpResponse::InternalServerError().body("Cannot process request"),
        };

    match state.db.send(InsertLeaderboardRows {rows: rows}).await {
        Ok(Ok(_)) => HttpResponse::Accepted().body("Success"),
        Ok(Err(err)) => {
            match state.db.send(DisableSubmission{unix: body.unix_time_stamp}).await {
                Ok(Ok(_)) => (),
                Ok(Err(err)) => return HttpResponse::InternalServerError().body(format!("Could not disable submission. We're in a bad sate! Error: {}", err)),
                _ => return HttpResponse::InternalServerError().body("Cannot process request"),
            }
            HttpResponse::InternalServerError().body(err.to_string())
        },
        _ => HttpResponse::InternalServerError().body("Cannot process request"),
    }
}