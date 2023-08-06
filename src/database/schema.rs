table! {
    submissions (unix_time_stamp) {
        uuid -> Varchar,
        unix_time_stamp -> Bigint,
        game -> Varchar,
        valid -> Bool,
    }
}

table! {
    leaderboards (unix_time_stamp) {
        game -> VarChar,
        player -> VarChar,
        position -> Integer,
        score -> Integer,
        unix_time_stamp -> Bigint,
    }
}

table! {
    ban (uuid) {
        uuid -> Varchar,
    }
}

joinable!(leaderboards -> submissions (unix_time_stamp));
allow_tables_to_appear_in_same_query!(leaderboards, submissions);
