use std::time::Duration;

use actix_extensible_rate_limit::{backend::{memory::InMemoryBackend, SimpleInputFunctionBuilder}, RateLimiter};
use actix_web::{web::Data, App, HttpServer, middleware::Logger};
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};

mod common;

mod routes;
use routes::{hello, submit_leaderboard_entries, get_leaderboards_from_player, get_leaderboard,get_leaderboard_between, get_leaderboard_for_all};

pub struct AppState {
    db: Pool<Postgres>
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let database_url = "postgresql://ameliah:@127.0.0.1:5432/leaderboard_api";
    let pool = PgPoolOptions::new()
    .max_connections(5)
    .connect(&database_url)
    .await
    .expect("Error building a connection pool");

    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

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
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
