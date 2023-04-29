use std::time::Duration;

use actix_extensible_rate_limit::{backend::{memory::InMemoryBackend, SimpleInputFunctionBuilder}, RateLimiter};
use actix_web::{web::Data, App, HttpServer, middleware::Logger};
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};

mod common;

mod routes;
use routes::{hello, submit_leaderboard_entries, get_leaderboards_from_player, get_leaderboard,get_leaderboard_between, get_leaderboard_for_all};
use toml::Value;

pub struct AppState {
    db: Pool<Postgres>
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let config_file = include_str!("../config.toml");
    let config = config_file.parse::<Value>().unwrap();

    let database_url = config["database_url"].as_str().unwrap();
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
        .app_data(Data::new(AppState {db : pool.clone()}))
        .service(hello)
        .service(submit_leaderboard_entries)
        .service(get_leaderboards_from_player)
        .service(get_leaderboard)
        .service(get_leaderboard_between)
        .service(get_leaderboard_for_all)
    })
    .bind((config["address"].as_str().unwrap(), config["port"].as_integer().unwrap() as u16))?
    .run()
    .await
}
