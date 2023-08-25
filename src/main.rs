use std::time::Duration;

use actix_extensible_rate_limit::{backend::{memory::InMemoryBackend, SimpleInputFunctionBuilder}, RateLimiter};
use actix_web::{web::Data, App, HttpServer, middleware::Logger, Responder, HttpResponse, get};
use config::APIConfig;
use database::API;

use docs::ApiDoc;
use utoipa::OpenApi;
use utoipa_rapidoc::RapiDoc;

#[macro_use] extern crate diesel;

mod leaderboard_api;
mod chest_api;
mod eggwars_map_api;
mod database;
mod config;
mod docs;


#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Amelia loves you very much <3 <3")
}

#[actix_web::main]
async fn main() -> Result<(), std::io::Error> {
    let config = match APIConfig::from_file(String::from("config.toml")) {
        Ok(config) => config,
        Err(e) => panic!("Couldn't make config: {}", e),
    };
    let config_clone = config.clone();

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
        .app_data(Data::new(API::new(&config)))
        .service(RapiDoc::with_openapi("/api-docs/openapi.json", ApiDoc::openapi()).path("/rapidoc"))
        .service(hello)
        .service(leaderboard_api::routes::submission::submit_leaderboard_entries)
        .service(leaderboard_api::routes::player::get_leaderboards_from_player)
        .service(leaderboard_api::routes::game::get_leaderboard)
        .service(leaderboard_api::routes::game::get_leaderboard_between)
        .service(leaderboard_api::routes::game::get_games)
        .service(chest_api::get_current_chests)
        .service(chest_api::get_season_chests)
        .service(chest_api::get_seasons)
        .service(eggwars_map_api::get_all_eggwars_maps)
        .service(eggwars_map_api::get_eggwars_map)
    })
    .bind((config_clone.address, config_clone.port))?
    .run()
    .await
}
