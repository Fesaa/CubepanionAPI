use actix_web::{web::{Data, Path}, Responder, HttpResponse, get};
use diesel::result::Error::NotFound;

use crate::database::{API, eggwars_maps::messages::{FetchEggWarsMaps, FetchEggWarsMap}};


#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All EggWars maps", body = Vec<EggWarsMap>),
        (status = 404, description =  "No EggWars maps found"),
        (status = 500, description = "Unable to retrieve EggWars maps")
    )
)]
#[get("/eggwars_map_api")]
pub async fn get_all_eggwars_maps(state: Data<API>) -> impl Responder {
    match state.db.send(FetchEggWarsMaps{}).await {
        Ok(Ok(maps)) => HttpResponse::Ok().json(maps),
        Ok(Err(err)) => HttpResponse::InternalServerError().body(format!("Unable to retrieve eggwars map: {}", err)),
        _ => HttpResponse::InternalServerError().body("Unable to retrieve eggwars maps"),
    }
}

#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All EggWars maps", body = Vec<EggWarsMap>),
        (status = 404, description =  "No EggWars maps found"),
        (status = 500, description = "Unable to retrieve EggWars maps")
    ),
    params(
        ("name", description = "EggWars map name")
    )
)]
#[get("/eggwars_map_api/{name}")]
pub async fn get_eggwars_map(state: Data<API>, path: Path<String>) -> impl Responder {
    match state.db.send(FetchEggWarsMap{name: path.into_inner()}).await {
        Ok(Ok(map)) => HttpResponse::Ok().json(map),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No eggwars maps found"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve eggwars maps: {}", err)),
        },
        _ => HttpResponse::InternalServerError().body("Unable to retrieve eggwars map"),
    }
}