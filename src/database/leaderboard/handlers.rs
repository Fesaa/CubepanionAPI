use actix::Handler;
use diesel::{QueryResult, QueryDsl};
use crate::database::DbActor;
use crate::database::schema::LeaderboardRow;
use diesel::{self, prelude::*};

use super::messages::*;


impl Handler<FetchLeaderboardFromPlayer> for DbActor {
    type Result = QueryResult<Vec<LeaderboardRow>>;

    fn handle(&mut self, msg: FetchLeaderboardFromPlayer, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::leaderboards::dsl::{leaderboards, player, unix_time_stamp};

        let con = self.0.get()
        .expect("Fetch Leaderboard From Player: Unable to establish connection");

        leaderboards
        .filter(unix_time_stamp.eq_all(self.max_unix(&con)))
        .filter(player.eq_all(msg.player_name))
        .load::<LeaderboardRow>(&con)
    }
}

impl Handler<FetchLeaderboardForGame> for DbActor {
    type Result = QueryResult<Vec<LeaderboardRow>>;

    fn handle(&mut self, msg: FetchLeaderboardForGame, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::leaderboards::dsl::{leaderboards, game, unix_time_stamp, position};

        let con = self.0.get()
        .expect("Fetch Leaderboard For Game: Unable to establish connection");

        leaderboards
        .filter(unix_time_stamp.eq_all(self.max_unix(&con)))
        .filter(game.eq_all(msg.game_name))
        .filter(position.between(msg.min, msg.max))
        .load::<LeaderboardRow>(&con)
    }
}

impl Handler<InsertLeaderboardRows> for DbActor {
    type Result = QueryResult<String>;

    fn handle(&mut self, msg: InsertLeaderboardRows, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::leaderboards::dsl::leaderboards;

        let con = self.0.get()
        .expect("Insert Leaderboard Rows: Unable to establish connection");

        match diesel::insert_into(leaderboards)
        .values(msg.rows)
        .execute(&con) {
            Ok(_) => Ok(String::from("")),
            Err(err) => Err(err),
        }
    }
}

impl Handler<InsertSubmission> for DbActor {
    type Result = QueryResult<()>;

    fn handle(&mut self, msg: InsertSubmission, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::submissions::dsl::submissions;
        
        let con = self.0.get()
        .expect("Insert submission: Unable to establish connection");

    match diesel::insert_into(submissions)
        .values(msg.sub)
        .execute(&con) {
            Ok(_) => Ok(()),
            Err(err) => Err(err),
        }
    }
}
