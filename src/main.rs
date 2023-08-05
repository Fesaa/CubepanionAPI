use std::time::Duration;

use actix_extensible_rate_limit::{backend::{memory::InMemoryBackend, SimpleInputFunctionBuilder}, RateLimiter};
use actix_web::{web::Data, App, HttpServer, middleware::Logger, Responder, HttpResponse, get};
use regex::Regex;
use serde::Deserialize;
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};
use toml::from_str;

mod leaderboard_api;

pub struct API {
    db: Pool<Postgres>,
    username_regex: Regex
}

#[derive(Debug, Deserialize)]
struct APIConfig {
    database_url: String,
    address: String,
    port: u16
}

#[get("/")]
pub async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Ameliah loves you <3")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let config_file = include_str!("../config.toml");
    let config: APIConfig = match from_str(&config_file) {
        Ok(c) => c,
        Err(e) => panic!("Could not parse config: {}", e),
    };

    let database_url = config.database_url;
    let pool = PgPoolOptions::new()
    .max_connections(5)
    .connect(&database_url)
    .await
    .expect("Error building a connection pool");

    let backend = InMemoryBackend::builder().build();

    HttpServer::new(move || {
        let input = SimpleInputFunctionBuilder::new(Duration::from_secs(60), 10)
            .real_ip_key()
            .build();
        let middleware = RateLimiter::builder(backend.clone(), input)
            .add_headers()
            .build();
        App::new()
        .wrap(middleware)
        .wrap(Logger::default())
        .app_data(Data::new(API {db : pool.clone(), username_regex: Regex::new(r"[a-zA-Z0-9_]{2,16}").unwrap()}))
        .service(hello)
        .service(leaderboard_api::routes::submission::submit_leaderboard_entries)
        .service(leaderboard_api::routes::player::get_leaderboards_from_player)
        .service(leaderboard_api::routes::game::get_leaderboard)
        .service(leaderboard_api::routes::game::get_leaderboard_between)
        .service(leaderboard_api::routes::players::get_leaderboard_for_all)
    })
    .bind((config.address, config.port))?
    .run()
    .await
}
