use actix_web::{web::{Data, Path}, Responder, HttpResponse, get};
use diesel::result::Error::NotFound;

use crate::database::{API, chests::messages::{FetchChestLocationsForRunningSeason, FetchChestLocationsForSeason, FetchSeasons}};

#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All chests for the current season", body = Vec<ChestLocation>),
        (status = 404, description =  "No chests found for the current season"),
        (status = 500, description = "Unable to retrieve chest locations")
    )
)]
#[get("/chest_api/current")]
pub async fn get_current_chests(state: Data<API>) -> impl Responder {
    match state.db.send(FetchChestLocationsForRunningSeason{}).await {
        Ok(Ok(chests)) => HttpResponse::Ok().json(chests),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No chests found for the current season"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve chest locations: {}", err)),
        }
        _ => HttpResponse::InternalServerError().body("Unable to retrieve chest locations"),
    }
}

#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All chests for the given season", body = Vec<ChestLocation>),
        (status = 404, description =  "No chests found for the current season"),
        (status = 500, description = "Unable to retrieve chest locations")
    ),
    params(
        ("season", description = "Season name")
    )
)]
#[get("/chest_api/season/{season}")]
pub async fn get_season_chests(state: Data<API>, path: Path<String>) -> impl Responder {
    match state.db.send(FetchChestLocationsForSeason{season: path.into_inner()}).await {
        Ok(Ok(chests)) => HttpResponse::Ok().json(chests),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No chests found for that season"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve chest locations: {}", err)),
            
        },
        _ => HttpResponse::InternalServerError().body("Unable to retrieve chest locations"),
    }
}

#[utoipa::path(
    get,
    responses(
        (status = 200, description = "All queried seasons", body = Vec<String>),
        (status = 404, description =  "No seasons found"),
        (status = 500, description = "Unable to retrieve seasons")
    ),
    params(
        ("running", description = "Running seasons only, or all seasons")
    )
)]
#[get("/chest_api/seasons/{running}")]
pub async fn get_seasons(state: Data<API>, path: Path<bool>) -> impl Responder {
    match state.db.send(FetchSeasons{running: path.into_inner()}).await {
        Ok(Ok(seasons)) => HttpResponse::Ok().json(seasons),
        Ok(Err(err)) => match err {
            NotFound => HttpResponse::NotFound().body("No seasons found"),
            _ => HttpResponse::InternalServerError().body(format!("Unable to retrieve seasons: {}", err)),
        },
        _ => HttpResponse::InternalServerError().body("Unable to retrieve seasons"),
    }
}