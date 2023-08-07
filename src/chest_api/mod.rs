use actix_web::{web::{Data, Path}, Responder, HttpResponse, get};

use crate::database::{API, chests::messages::{FetchChestLocationsForRunningSeason, FetchChestLocationsForSeason, FetchSeasons}};

#[get("/chest_api/current")]
pub async fn get_current_chests(state: Data<API>) -> impl Responder {
    match state.db.send(FetchChestLocationsForRunningSeason{}).await {
        Ok(Ok(chests)) => HttpResponse::Ok().json(chests),
        Ok(Err(err)) => HttpResponse::InternalServerError().body(format!("Unable to retrieve chest locations: {}", err)),
        _ => HttpResponse::InternalServerError().body("Unable to retrieve chest locations"),
    }
}

#[get("/chest_api/season/{season}")]
pub async fn get_season_chests(state: Data<API>, path: Path<String>) -> impl Responder {
    match state.db.send(FetchChestLocationsForSeason{season: path.into_inner()}).await {
        Ok(Ok(chests)) => HttpResponse::Ok().json(chests),
        Ok(Err(err)) => HttpResponse::InternalServerError().body(format!("Unable to retrieve chest locations: {}", err)),
        _ => HttpResponse::InternalServerError().body("Unable to retrieve chest locations"),
    }
}

#[get("/chest_api/seasons/{running}")]
pub async fn get_seasons(state: Data<API>, path: Path<bool>) -> impl Responder {
    match state.db.send(FetchSeasons{running: path.into_inner()}).await {
        Ok(Ok(seasons)) => HttpResponse::Ok().json(seasons),
        Ok(Err(err)) => HttpResponse::InternalServerError().body(format!("Unable to retrieve chest locations: {}", err)),
        _ => HttpResponse::InternalServerError().body("Unable to retrieve chest locations"),
    }
}