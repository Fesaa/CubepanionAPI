use std::time::Instant;

use actix::{Actor, SyncContext, Addr, SyncArbiter, Message, dev::ToEnvelope, Handler, MailboxError};
use diesel::{r2d2::{Pool, ConnectionManager}, PgConnection};
use diesel;
use metrics::{counter, histogram};
use regex::Regex;

use crate::APIConfig;

pub mod leaderboard;
pub mod chests;
pub mod eggwars_maps;
pub mod schema;

#[derive(Clone)]
pub struct DbActor(pub Pool<ConnectionManager<PgConnection>>);

impl Actor for DbActor {
    type Context = SyncContext<Self>;
}

pub struct Holder<A: Actor> {
    addr: Addr<A>
}

impl <A: Actor> Holder<A> {
    pub async fn send<M>(&self, msg: M, endpoint: &'static str) -> Result<<M as Message>::Result, MailboxError>
    where
        M: Message + Send + 'static,
        M::Result: Send,
        A: Handler<M>,
        A::Context: ToEnvelope<A, M>,
    {
        counter!("total_requests", 1, "endpoint" => endpoint);

        let start = Instant::now();
        let response:Result<M::Result, MailboxError> = self.addr.send(msg).await;
        let delta = start.elapsed();
        histogram!("request_duration", delta, "endpoint" => endpoint);
        
        return match response {
            Ok(r) => {
                counter!("success_requests", 1, "endpoint" => endpoint);
                Ok(r)
            },
            Err(e) => Err(e),
        };
    }
}

pub struct API {
    pub db: Holder<DbActor>,
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
        API { db: Holder { addr: db_addr }, username_regex: Regex::new(r"[a-zA-Z0-9_]{2,16}").unwrap()}
    }    
}