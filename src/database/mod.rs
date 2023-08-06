use actix::{Actor, SyncContext, Addr, SyncArbiter};
use diesel::{r2d2::{Pool, ConnectionManager}, PgConnection, dsl::max};
use diesel::{self, prelude::*};
use regex::Regex;

use crate::APIConfig;

pub mod leaderboard;
mod schema;

#[derive(Clone)]
pub struct DbActor(pub Pool<ConnectionManager<PgConnection>>);

impl Actor for DbActor {
    type Context = SyncContext<Self>;
}

impl DbActor {

    pub fn max_unix(&self, con: &PgConnection) -> i64 {
        use crate::database::schema::submissions::dsl::{unix_time_stamp, submissions, valid};
        submissions
        .filter(valid.eq_all(true))
        .select(max(unix_time_stamp))
        .first::<Option<i64>>(con)
        .unwrap_or(Some(0))
        .unwrap_or(0)
    }

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