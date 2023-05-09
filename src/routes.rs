use regex::Regex;

use actix_web::{get, post, web::{Json}, HttpResponse, Responder, web::{Data, Path, Query}};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;

use crate::{common::{LeaderboardEntry, is_valid_uuid}, AppState};

#[get("/")]
pub async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Ameliah loves you <3")
}

#[derive(Deserialize)]
pub struct Players {
    players: Vec<String>,
}

#[post("/leaderboard_api/leaderboard/players")]
pub async fn get_leaderboard_for_all(state: Data<AppState>, body: Json<Players>) -> impl Responder {
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

async fn remove_submission(unix_time_stamp: i64, state: Data<AppState>) {
    match sqlx::query("UPDATE submissions SET valid = false WHERE unix_time_stamp = $1;")
        .bind(unix_time_stamp)
        .execute(&state.db)
        .await {
            Ok(_) => (),
            Err(err) => println!("{}", String::from("Could not update submission: ") + &err.to_string()),
        }
}