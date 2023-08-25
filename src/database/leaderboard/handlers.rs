use actix::Handler;
use diesel::{QueryResult, QueryDsl, dsl::sql, sql_types::Bool};
use crate::{database::{schema::{LeaderboardRow, GameRow}, DbActor}, leaderboard_api::models::LeaderboardGame};
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
        use crate::database::schema::leaderboards::dsl::{leaderboards, game, position, unix_time_stamp};
        use crate::database::schema::submissions::dsl as s;
        use diesel::dsl::max;

        let mut con = self.0.get()
        .expect("Fetch Leaderboard For Game: Unable to establish connection");

        let max_unix: i64 = s::submissions
            .filter(s::game.eq(&msg.game_name))
            .select(max(s::unix_time_stamp))
            .first::<Option<i64>>(&mut con)?.unwrap_or(0);

        leaderboards
        .filter(unix_time_stamp.eq(max_unix))
        .filter(game.eq(&msg.game_name))
        .filter(position.between(msg.min, msg.max))
        .order(position)
        .load::<LeaderboardRow>(&mut con)
    }
}

impl Handler<InsertLeaderboardRows> for DbActor {
    type Result = QueryResult<usize>;

    fn handle(&mut self, msg: InsertLeaderboardRows, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::leaderboards::dsl::leaderboards;

        let mut con = self.0.get()
        .expect("Insert Leaderboard Rows: Unable to establish connection");

        diesel::insert_into(leaderboards)
        .values(msg.rows)
        .execute(&mut con)
    }
}

impl Handler<InsertSubmission> for DbActor {
    type Result = QueryResult<usize>;

    fn handle(&mut self, msg: InsertSubmission, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::submissions::dsl::submissions;
        
        let mut con = self.0.get()
        .expect("Insert submission: Unable to establish connection");

        diesel::insert_into(submissions)
        .values(msg.sub)
        .execute(&mut con)
    }
}

impl Handler<DisableSubmission> for DbActor {
    type Result = QueryResult<usize>;

    fn handle(&mut self, msg: DisableSubmission, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::submissions::dsl::{submissions, unix_time_stamp, valid};
        
        let mut con = self.0.get()
        .expect("Insert submission: Unable to establish connection");

        diesel::update(submissions.filter(unix_time_stamp.eq(msg.unix)))
        .set(valid.eq(false))
        .execute(&mut con)
    }
}


impl Handler<GetGames> for DbActor {
    type Result = QueryResult<Vec<LeaderboardGame>>;

    fn handle(&mut self, msg: GetGames, _ctx: &mut Self::Context) -> Self::Result {
        use crate::database::schema::games::dsl::{games, active};
        
        let mut con = self.0.get()
        .expect("Get Games: Unable to establish connection");


        let rows = if msg.must_be_active {
            games
                .filter(active.eq(true))
                .load::<GameRow>(&mut con)?
        } else {
            games
                .load::<GameRow>(&mut con)?
        };


        rows.iter()
            .map(|row| {
                Ok(LeaderboardGame {
                    game: row.game.clone(),
                    display_name: row.display_name.clone(),
                    aliases: row.aliases.clone().split(",").map(|s| s.to_string()).collect(),
                    active: row.active,
                })
            })
            .collect::<QueryResult<Vec<LeaderboardGame>>>()


    }
}

