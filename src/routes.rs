use std::{collections::HashSet};
use regex::Regex;

use actix_web::{get, post, web::{Json}, HttpResponse, Responder, web::{Data, Path, Query}};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;

use crate::{common::{LeaderboardEntry, is_valid_uuid}, AppState};

#[get("/")]
pub async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Ameliah loves you <3")
}

#[get("/leaderboard_api/leaderboard/{game}")]
pub async fn get_leaderboard(state: Data<AppState>, path: Path<String>) -> impl Responder {
    let game = path.into_inner();
    if !Regex::new(r"[a-zA-Z0-9_]{2,16}").unwrap().is_match(&game) {
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
pub async fn get_leaderboard_between(state: Data<AppState>, path: Path<String>, info: Query<BoundedRequest>) -> impl Responder {
    let game = path.into_inner();
    if !Regex::new(r"[a-zA-Z0-9_]{2,16}").unwrap().is_match(&game) {
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

#[derive(Deserialize, Serialize, FromRow)]
pub struct LeaderboardRow {
    pub player: String,
    pub position: i32,
    pub score: i32,
    pub game: String
}

#[get("/leaderboard_api/player/{name}")]
pub async fn get_leaderboards_from_player(state: Data<AppState>, path: Path<String>) -> impl Responder {
    let name = path.into_inner();
    if !Regex::new(r"[a-zA-Z0-9_]{2,16}").unwrap().is_match(&name) {
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
        player = $1
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

#[derive(Deserialize)]
pub struct LeaderboardSubmission {
    uuid: String,
    unix_time_stamp: i64,
    game: String,
    entries: Vec<LeaderboardEntry>
}


#[post("/leaderboard_api")]
pub async fn submit_leaderboard_entries(state: Data<AppState>, body: Json<LeaderboardSubmission>) -> impl Responder {
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
            let mut sql: String = String::from("INSERT INTO leaderboards (game, player, position, score, unix_time_stamp) VALUES ");
            let mut player_set: HashSet<String> = HashSet::new();
            let mut position_set: HashSet<i32> = HashSet::new();

            for row in body.entries.iter() {
                if player_set.contains(&row.player) || position_set.contains(&row.position) {
                    remove_submission(body.unix_time_stamp, state).await;
                    return HttpResponse::BadRequest().body("Entries contains duplicate name or position.")
                }
                sql = sql + &format!("('{}', '{}', {}, {}, {}),", &body.game, &row.player, &row.position, &row.score, &body.unix_time_stamp);
                player_set.insert(row.player.to_owned());
                position_set.insert(row.position.to_owned());
            }

            match sql.strip_suffix(",") {
                Some(s) => sql = s.to_owned() + ";",
                None => sql = sql + ";",
            }

            match sqlx::query(&sql)
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

async fn remove_submission(unix_time_stamp: i64, state: Data<AppState>) {
    match sqlx::query("UPDATE submissions SET valid = false WHERE unix_time_stamp = $1;")
        .bind(unix_time_stamp)
        .execute(&state.db)
        .await {
            Ok(_) => (),
            Err(err) => println!("{}", String::from("Could not update submission: ") + &err.to_string()),
        }
}