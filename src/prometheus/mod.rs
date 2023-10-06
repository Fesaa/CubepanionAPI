use actix_web::{Responder, get, HttpResponse};
use metrics::{register_counter, register_histogram};
use metrics_prometheus::Recorder;

pub async fn setup_metrics() -> Result<Recorder, Box<dyn std::error::Error>> {
    let recorder = metrics_prometheus::try_install()?;

    register_counter!("total_requests", "endpoint" => "");
    register_counter!("success_requests", "endpoint" => "");
    register_histogram!("request_duration", "endpoint" => "");

    return Ok(recorder);
}

#[get("/prometheus")]
pub async fn get_metrics() -> impl Responder {
    let report = prometheus::TextEncoder::new()
        .encode_to_string(&prometheus::default_registry().gather());
    return match report {
        Ok(s) => HttpResponse::Ok().body(s),
        Err(e) => HttpResponse::InternalServerError().body(format!("Error while encoding metrics, {}", e)),
    };
}