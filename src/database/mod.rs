use actix::{Actor, SyncContext, Addr, SyncArbiter};
use diesel::{r2d2::{Pool, ConnectionManager}, PgConnection};
use diesel;
use regex::Regex;

use crate::APIConfig;

pub mod leaderboard;
pub mod schema;

#[derive(Clone)]
pub struct DbActor(pub Pool<ConnectionManager<PgConnection>>);

impl Actor for DbActor {
    type Context = SyncContext<Self>;
}

pub struct API {
    pub db: Addr<DbActor>,
    pub username_regex: Regex
}

impl API {

    pub fn new(config: &APIConfig) -> API {
        let manager = ConnectionManager::<PgConnection>::new(&config.database_url);
        let pool = Pool::builder()
        .max_size(5)
        .build(manager)
        .expect("Error building a connection pool");
        let db_addr = SyncArbiter::start(5, move || DbActor(pool.clone()));
        API { db: db_addr, username_regex: Regex::new(r"[a-zA-Z0-9_]{2,16}").unwrap()}
    }

}