use actix::Handler;
use diesel::{QueryResult, QueryDsl, dsl::sql, sql_types::Bool};
use crate::database::{schema::LeaderboardRow, DbActor};
use diesel::{self, prelude::*};

use super::messages::*;


impl Handler<FetchLeaderboardFromPlayer> for DbActor {
    type Result = QueryResult<Vec<LeaderboardRow>>;

    fn handle(&mut self, msg: FetchLeaderboardFromPlayer, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::leaderboards::dsl::{leaderboards, player, position};
        use crate::database::schema::submissions::dsl as s;
        use diesel::dsl::max;

        let mut con = self.0.get()
        .expect("Fetch Leaderboard From Player: Unable to establish connection");

        let mut max_unix_time_stamps_per_game: Vec<(String, Option<i64>)> = s::submissions
            .filter(s::valid)
            .group_by(s::game)
            .select((s::game, max(s::unix_time_stamp)))
            .load::<(String, Option<i64>)>(&mut con)?;

        let s = max_unix_time_stamps_per_game
            .iter_mut()
            .map(|row| format!("('{}', {})", row.0, row.1.unwrap_or(0)))
            .collect::<Vec<String>>()
            .join(",");

        leaderboards
        .filter(player.eq(msg.player_name))
        .filter(sql::<Bool>(&format!("({}, {}) IN ({})", "game", "unix_time_stamp", s)))
        .order_by(position)
        .load::<LeaderboardRow>(&mut con)
    }
}

impl Handler<FetchLeaderboardForGame> for DbActor {
    type Result = QueryResult<Vec<LeaderboardRow>>;

    fn handle(&mut self, msg: FetchLeaderboardForGame, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::leaderboards::dsl::{leaderboards, game, position};

        let mut con = self.0.get()
        .expect("Fetch Leaderboard For Game: Unable to establish connection");

        leaderboards
        //.filter(unix_time_stamp.eq(self.max_unix(&con)))
        .filter(game.eq(msg.game_name))
        .filter(position.between(msg.min, msg.max))
        .load::<LeaderboardRow>(&mut con)
    }
}

impl Handler<InsertLeaderboardRows> for DbActor {
    type Result = QueryResult<String>;

    fn handle(&mut self, msg: InsertLeaderboardRows, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::leaderboards::dsl::leaderboards;

        let mut con = self.0.get()
        .expect("Insert Leaderboard Rows: Unable to establish connection");

        match diesel::insert_into(leaderboards)
        .values(msg.rows)
        .execute(&mut con) {
            Ok(_) => Ok(String::from("")),
            Err(err) => Err(err),
        }
    }
}

impl Handler<InsertSubmission> for DbActor {
    type Result = QueryResult<()>;

    fn handle(&mut self, msg: InsertSubmission, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::submissions::dsl::submissions;
        
        let mut con = self.0.get()
        .expect("Insert submission: Unable to establish connection");

    match diesel::insert_into(submissions)
        .values(msg.sub)
        .execute(&mut con) {
            Ok(_) => Ok(()),
            Err(err) => Err(err),
        }
    }
}
